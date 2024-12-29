//go:build (!darwin && !linux) || cgo

package core

import (
	"bytes"
	"embed"
	"image/png"
	"os"
	"path/filepath"
	"runtime"

	"github.com/energye/systray"
	"github.com/pkg/browser"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"golang.design/x/clipboard"
)

//go:embed icons/*.png icons/*.ico
var icons embed.FS

var clipboardErr = clipboard.Init()

// modify modifies the given menu items with the given operation.
func modify(op func(*systray.MenuItem), items ...*systray.MenuItem) {
	for _, item := range items {
		op(item)
	}
}

// makeConfigInfo creates a read-only menu item
func makeConfigInfo(option *systray.MenuItem, sensitive bool, cfg *Config, lookup func(*Config) string, editor func(*Config, string)) {
	value := lookup(cfg)
	if value == "" {
		value = "not set"
	} else if sensitive {
		value = "[redacted]"
	}

	if editor != nil {
		option.AddSubMenuItem(value, "Open").Click(func() {
			editor(cfg, lookup(cfg))
		})
		return
	}
	option.AddSubMenuItem(value, "Not editable").Disable()
}

// makeConfigOption creates a menu item with a checkbox
func makeConfigOption(option *systray.MenuItem, cfg *Config, lookup func(*Config) bool, editor func(*Config, bool)) {
	// initialize the menu item
	if lookup(cfg) {
		option.Check()
	} else {
		option.Uncheck()
	}

	// define the click event
	option.Click(func() {
		editor(cfg, !lookup(cfg))
		if lookup(cfg) {
			option.Check()
		} else {
			option.Uncheck()
		}
	})
}

// makeConfigSection creates a menu item with a sub-menu of checkboxes
func makeConfigSection[K comparable](section map[K]*systray.MenuItem, cfg *Config, lookup func(*Config) K, editor func(*Config, K)) {
	for j, item := range section {
		// initialize the menu items
		if j == lookup(cfg) {
			item.Check()
		} else {
			item.Uncheck()
		}

		// define the click event
		item.Click(func() {
			for k, subItem := range section {
				if subItem == item {
					editor(cfg, k)
					subItem.Check()
				} else {
					subItem.Uncheck()
				}
			}
		})
	}
}

// makePropertyCopyAction creates a menu item with a copy action
func makePropertyCopyAction(item *systray.MenuItem, img *Image, format clipboard.Format, getter func(*Image) []byte) {
	item.SetIcon(readIcon("copy"))
	item.Click(func() {
		if img == nil {
			return
		}

		if clipboardErr != nil {
			logger.ErrLogger.Printf("Failed to initialize clipboard: %v", clipboardErr)
			return
		}

		_ = clipboard.Write(format, getter(img))
	})
}

// makePropertyOpenAction creates a menu item with an open URL action
func makePropertyOpenAction(item *systray.MenuItem, img *Image, getter func(*Image) string) {
	item.SetIcon(readIcon("open"))
	item.Click(func() {
		if img == nil {
			return
		}

		if runtime.GOOS == "linux" {
			if value, ok := os.LookupEnv("GTK_PATH"); ok {
				_ = os.Unsetenv("GTK_PATH")
				defer os.Setenv("GTK_PATH", value)
			}
		}

		if err := browser.OpenURL(getter(img)); err != nil {
			logger.ErrLogger.Printf("Failed to open %s: %v", item, err)
		}
	})
}

// openDirectory opens the directory in the file explorer
func openDirectory(path string) {
	if path == "" {
		return
	}

	if extras.EmbeddedWatermarks[path] != nil {
		var err error
		path, err = extras.EmbeddedWatermarks.ToFiles("watermarks")
		if err != nil {
			logger.ErrLogger.Printf("Failed to extract watermarks: %v", err)
			return
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		logger.ErrLogger.Printf("Failed to open directory %s: %v", path, err)
		return
	}

	if !info.IsDir() {
		path = filepath.Dir(path)
	}

	if err := browser.OpenURL("file://" + path); err != nil {
		logger.ErrLogger.Printf("Failed to open directory %s: %v", path, err)
	}
}

// readIcon reads the icon file
func readIcon(name string) []byte {
	ext := ".png"
	if runtime.GOOS == "windows" {
		ext = ".ico"
	}

	data, err := icons.ReadFile("icons/" + name + ext)
	if err != nil {
		logger.ErrLogger.Printf("Failed to read icon %s: %s", name+ext, err)
		return nil
	}

	return data
}

// Run executes the given function with the given configuration.
func Run(execute func(*Config) *Image, cfg *Config) {
	if !cfg.Daemon {
		_ = execute(cfg)
		return
	}

	img := &Image{}
	onReady := func() {
		systray.SetTitle(AppName)
		systray.SetTooltip(AppName)
		systray.SetIcon(readIcon("wallpaper"))

		var mRefresh, mSpeak, mQuit, mPropertiesAudio *systray.MenuItem

		// Main section
		mRefresh = systray.AddMenuItem("Refresh", "Refresh the wallpaper")
		mRefresh.SetIcon(readIcon("refresh"))
		mRefresh.Click(func() {
			modify(func(mi *systray.MenuItem) { mi.Disable() }, mRefresh, mSpeak, mQuit)
			mPropertiesAudio.Hide()
			img.Update(execute(cfg))
			modify(func(mi *systray.MenuItem) { mi.Enable() }, mRefresh, mQuit)
			if img != nil && img.Audio != nil {
				mSpeak.Enable()
				mPropertiesAudio.Show()
			}
		})

		mSpeak = systray.AddMenuItem("Speak", "Speak the wallpaper description")
		mSpeak.SetIcon(readIcon("play"))
		mSpeak.Click(func() {
			modify(func(mi *systray.MenuItem) { mi.Disable() }, mRefresh, mSpeak, mQuit)
			if img != nil && img.Audio != nil {
				if err := img.Audio.Play(); err != nil {
					logger.ErrLogger.Printf("Failed to play audio: %v", err)
				}
			}
			modify(func(mi *systray.MenuItem) { mi.Enable() }, mRefresh, mSpeak, mQuit)
		})

		systray.AddSeparator()

		// Config section
		mConfig := systray.AddMenuItem("Configure", "Configure the wallpaper settings")
		mConfigDay := mConfig.AddSubMenuItem("Day", "Day of the Bing wallpaper")
		makeConfigSection(map[types.Day]*systray.MenuItem{
			types.DayToday:                 mConfigDay.AddSubMenuItemCheckbox("Today", "Today's wallpaper", false),
			types.DayYesterday:             mConfigDay.AddSubMenuItemCheckbox("Yesterday", "Yesterday's wallpaper", false),
			types.DayTheDayBeforeYesterday: mConfigDay.AddSubMenuItemCheckbox("The day before yesterday", "The day before yesterday's wallpaper", false),
			types.DayThreeDaysAgo:          mConfigDay.AddSubMenuItemCheckbox("Three days ago", "Three days ago's wallpaper", false),
			types.DayFourDaysAgo:           mConfigDay.AddSubMenuItemCheckbox("Four days ago", "Four days ago's wallpaper", false),
			types.DayFiveDaysAgo:           mConfigDay.AddSubMenuItemCheckbox("Five days ago", "Five days ago's wallpaper", false),
			types.DaySixDaysAgo:            mConfigDay.AddSubMenuItemCheckbox("Six days ago", "Six days ago's wallpaper", false),
			types.DaySevenDaysAgo:          mConfigDay.AddSubMenuItemCheckbox("Seven days ago", "Seven days ago's wallpaper", false),
		}, cfg, func(c *Config) types.Day { return c.Day }, func(c *Config, d types.Day) {
			logger.InfoLogger.Printf("Setting Day: %v", d)
			c.Day = d
		})

		mConfigRegion := mConfig.AddSubMenuItem("Region", "Region of the Bing wallpaper")
		makeConfigSection(map[types.Region]*systray.MenuItem{
			types.RegionBrazil:        mConfigRegion.AddSubMenuItemCheckbox("Brazil Portuguese", "Brazilian region", false),
			types.RegionCanadaEnglish: mConfigRegion.AddSubMenuItemCheckbox("Canada English", "Canadian region (English)", false),
			types.RegionCanadaFrench:  mConfigRegion.AddSubMenuItemCheckbox("Canada French", "Canadian region (French)", false),
			types.RegionChina:         mConfigRegion.AddSubMenuItemCheckbox("China", "Chinese region", false),
			types.RegionFrance:        mConfigRegion.AddSubMenuItemCheckbox("France", "French region", false),
			types.RegionGermany:       mConfigRegion.AddSubMenuItemCheckbox("Germany", "German region", false),
			types.RegionIndia:         mConfigRegion.AddSubMenuItemCheckbox("India", "Indian region", false),
			types.RegionItaly:         mConfigRegion.AddSubMenuItemCheckbox("Italy", "Italian region", false),
			types.RegionJapan:         mConfigRegion.AddSubMenuItemCheckbox("Japan", "Japanese region", false),
			types.RegionNewZealand:    mConfigRegion.AddSubMenuItemCheckbox("New Zealand", "New Zealand's region", false),
			types.RegionOther:         mConfigRegion.AddSubMenuItemCheckbox("Other", "Other regions", false),
			types.RegionSpain:         mConfigRegion.AddSubMenuItemCheckbox("Spain", "Spanish region", false),
			types.RegionUnitedKingdom: mConfigRegion.AddSubMenuItemCheckbox("United Kingdom", "British region", false),
			types.RegionUnitedStates:  mConfigRegion.AddSubMenuItemCheckbox("United States", "US region", false),
		}, cfg, func(c *Config) types.Region { return c.Region }, func(c *Config, r types.Region) {
			logger.InfoLogger.Printf("Setting Region: %v", r)
			c.Region = r
		})

		mConfigResolution := mConfig.AddSubMenuItem("Resolution", "Resolution of the wallpaper")
		makeConfigSection(map[types.Resolution]*systray.MenuItem{
			types.LowDefinition:       mConfigResolution.AddSubMenuItemCheckbox("Low Definition", "Low Definition resolution", false),
			types.HighDefinition:      mConfigResolution.AddSubMenuItemCheckbox("High Definition", "High Definition resolution", false),
			types.UltraHighDefinition: mConfigResolution.AddSubMenuItemCheckbox("Ultra High Definition", "Ultra High Definition resolution", false),
		}, cfg, func(c *Config) types.Resolution { return c.Resolution }, func(c *Config, r types.Resolution) {
			logger.InfoLogger.Printf("Setting Resolution: %v", r)
			c.Resolution = r
		})

		makeConfigOption(mConfig.AddSubMenuItemCheckbox("Draw Description", "Draw the wallpaper description", false), cfg,
			func(c *Config) bool { return c.DrawDescription },
			func(c *Config, b bool) {
				logger.InfoLogger.Printf("Setting DrawDescription: %v", b)
				c.DrawDescription = b
			})

		makeConfigOption(mConfig.AddSubMenuItemCheckbox("Draw QR Code", "Draw the QR code", false), cfg,
			func(c *Config) bool { return c.DrawQRCode },
			func(c *Config, b bool) {
				logger.InfoLogger.Printf("Setting DrawQRCode: %v", b)
				c.DrawQRCode = b
			})

		makeConfigOption(mConfig.AddSubMenuItemCheckbox("Download Only", "Download the wallpaper only", false), cfg,
			func(c *Config) bool { return c.DownloadOnly },
			func(c *Config, b bool) {
				logger.InfoLogger.Printf("Setting DownloadOnly: %v", b)
				c.DownloadOnly = b
			})

		makeConfigOption(mConfig.AddSubMenuItemCheckbox("Rotate Wallpaper counter clockwise", "Rotate the wallpaper counter clockwise", false), cfg,
			func(c *Config) bool { return c.RotateCounterClockwise },
			func(c *Config, b bool) {
				logger.InfoLogger.Printf("Setting RotateWallpaper: %v", b)
				c.RotateCounterClockwise = b
			})

		makeConfigOption(mConfig.AddSubMenuItemCheckbox("Use Google Text2Speech Service", "Use Google Text2Speech Service", false), cfg,
			func(c *Config) bool { return c.UseGoogleText2SpeechService },
			func(c *Config, b bool) {
				logger.InfoLogger.Printf("Setting UseGoogleText2SpeechService: %v", b)
				c.UseGoogleText2SpeechService = b
			})

		makeConfigOption(mConfig.AddSubMenuItemCheckbox("Use Google Translate Service", "Use Google Translate Service", false), cfg,
			func(c *Config) bool { return c.UseGoogleTranslateService },
			func(c *Config, b bool) {
				logger.InfoLogger.Printf("Setting UseGoogleTranslateService: %v", b)
				c.UseGoogleTranslateService = b
			})

		makeConfigInfo(mConfig.AddSubMenuItem("Watermark", "Watermark to be drawn on the wallpaper"), false, cfg,
			func(c *Config) string { return c.Watermark }, func(_ *Config, s string) { openDirectory(s) })

		makeConfigInfo(mConfig.AddSubMenuItem("Google App Credentials", "Google App Credentials"), false, cfg,
			func(c *Config) string { return c.GoogleAppCredentials }, func(_ *Config, s string) { openDirectory(s) })

		makeConfigInfo(mConfig.AddSubMenuItem("Furigana API AppId", "Furigana API AppId"), true, cfg,
			func(c *Config) string { return c.FuriganaApiAppId }, nil)

		makeConfigInfo(mConfig.AddSubMenuItem("Download Directory", "Download Directory"), false, cfg,
			func(c *Config) string { return c.DownloadDirectory }, func(_ *Config, s string) { openDirectory(s) })

		systray.AddSeparator()

		// Property section
		mProperties := systray.AddMenuItem("Properties", "Properties of the wallpaper")

		mPropertiesImage := mProperties.AddSubMenuItem("Image", "Image data of the wallpaper")
		makePropertyCopyAction(mPropertiesImage.
			AddSubMenuItem("Copy", "Copy the image data to the clipboard"),
			img, clipboard.FmtImage,
			func(i *Image) []byte {
				buf := bytes.NewBuffer(nil)
				_ = png.Encode(buf, i)
				return buf.Bytes()
			})
		makePropertyOpenAction(mPropertiesImage.AddSubMenuItem("Open", "Open the image in the browser"), img,
			func(i *Image) string { return "file://" + i.Location })

		mPropertiesAudio = mProperties.
			AddSubMenuItem("Audio", "Audio data of the wallpaper")
		makePropertyOpenAction(mPropertiesAudio.
			AddSubMenuItem("Open", "Open the audio in the browser"), img,
			func(i *Image) string {
				if i.Audio != nil {
					return "file://" + i.Audio.Location
				}

				return ""
			})

		makePropertyCopyAction(mProperties.
			AddSubMenuItem("Description", "Description of the wallpaper").
			AddSubMenuItem("Copy", "Copy the description to the clipboard"),
			img, clipboard.FmtText, func(i *Image) []byte { return []byte(i.Description) })

		mPropertiesSearchUrl := mProperties.AddSubMenuItem("Search URL", "Search URL of the wallpaper")
		makePropertyCopyAction(mPropertiesSearchUrl.AddSubMenuItem("Copy", "Copy the search URL to the clipboard"),
			img, clipboard.FmtText, func(i *Image) []byte { return []byte(i.SearchURL) })
		makePropertyOpenAction(mPropertiesSearchUrl.AddSubMenuItem("Open", "Open the search URL in the browser"), img,
			func(i *Image) string { return i.SearchURL })

		mPropertiesDownloadUrl := mProperties.AddSubMenuItem("Download URL", "Download URL of the wallpaper")
		makePropertyCopyAction(mPropertiesDownloadUrl.AddSubMenuItem("Copy", "Copy the download URL to the clipboard"),
			img, clipboard.FmtText, func(i *Image) []byte { return []byte(i.DownloadURL) })
		makePropertyOpenAction(mPropertiesDownloadUrl.AddSubMenuItem("Open", "Open the download URL in the browser"), img,
			func(i *Image) string { return i.DownloadURL })

		systray.AddSeparator()

		// Quit section
		mQuit = systray.AddMenuItem("Quit", "Quit the whole app")
		mQuit.SetIcon(readIcon("quit"))
		mQuit.Click(systray.Quit)

		// initial execution
		modify(func(mi *systray.MenuItem) { mi.Disable() }, mRefresh, mSpeak, mQuit)
		mPropertiesAudio.Hide()
		img.Update(execute(cfg))
		modify(func(mi *systray.MenuItem) { mi.Enable() }, mRefresh, mQuit)
		if img != nil && img.Audio != nil {
			mSpeak.Enable()
			mPropertiesAudio.Show()
		}
	}

	onExit := func() {
		// close the audio stream
		if img != nil && img.Audio != nil {
			_ = img.Audio.Close()
		}
	}

	systray.Run(onReady, onExit)
}
