//go:build (!darwin && !linux) || cgo

package core

import (
	"bytes"
	"embed"
	"fmt"
	"image/png"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"

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

// Controller is the controller of the application.
type Controller struct {
	img         *Image
	cfg         *Config
	execute     func(*Config) *Image
	refreshLock sync.Mutex
}

// OnReady initializes the application.
func (c *Controller) OnReady() {
	c.refreshLock.Lock()
	defer c.refreshLock.Unlock()

	// reset the menu
	systray.ResetMenu()

	// set the title and tooltip
	systray.SetTitle(AppName)
	systray.SetTooltip(AppName)

	// set the icon
	systray.SetIcon(readIcon("wallpaper"))

	var mRefresh, mSpeak, mQuit, mPropertiesAudio *systray.MenuItem

	// Main section
	mRefresh = systray.AddMenuItem("Refresh", "Refresh the wallpaper")
	mRefresh.SetIcon(readIcon("refresh"))
	mRefresh.Click(func() {
		modify(func(mi *systray.MenuItem) { mi.Disable() }, mRefresh, mSpeak, mQuit)
		mPropertiesAudio.Hide()
		autoPlayAudio := c.cfg.AutoPlayAudio
		c.cfg.AutoPlayAudio = false
		c.img.Update(c.execute(c.cfg))
		c.cfg.AutoPlayAudio = autoPlayAudio
		modify(func(mi *systray.MenuItem) { mi.Enable() }, mRefresh, mQuit)
		if c.img != nil && c.img.Audio != nil {
			mSpeak.Enable()
			mPropertiesAudio.Show()
		}
	})

	mSpeak = systray.AddMenuItem("Speak", "Speak the wallpaper description")
	mSpeak.SetIcon(readIcon("play"))
	mSpeak.Click(func() {
		modify(func(mi *systray.MenuItem) { mi.Disable() }, mRefresh, mSpeak, mQuit)
		if c.img != nil && c.img.Audio != nil {
			if err := c.img.Audio.Play(); err != nil {
				logger.Logger.Printf("Failed to play audio: %v", err)
			}
		}
		modify(func(mi *systray.MenuItem) { mi.Enable() }, mRefresh, mSpeak, mQuit)
	})

	systray.AddSeparator()

	apiMenu := systray.AddMenuItem("API", "API of the wallpaper")
	serverAddr := fmt.Sprintf("localhost:%d", c.cfg.ApiPort)
	apiMenuRetrieveConfig := apiMenu.AddSubMenuItem("Retrieve Config", "Retrieve the config")
	makeApiCommand(apiMenuRetrieveConfig, serverAddr, http.MethodGet, "/config", "", "")
	apiMenuUpdateConfig := apiMenu.AddSubMenuItem("Update Config", "Update the config")
	makeApiCommand(apiMenuUpdateConfig, serverAddr, http.MethodPatch, "/config", "refresh=true", "{...}")

	systray.AddSeparator()

	// Config section
	mConfig := systray.AddMenuItem("Configure", "Configure the wallpaper settings")
	mConfigDay := mConfig.AddSubMenuItem("Day", "Day of the Bing wallpaper")
	makeConfigSection(map[types.Day]*systray.MenuItem{
		types.DayToday: mConfigDay.AddSubMenuItemCheckbox("Today", "Today's wallpaper", false),
		types.Day1Ago:  mConfigDay.AddSubMenuItemCheckbox("Yesterday", "Yesterday's wallpaper", false),
		types.Day2Ago:  mConfigDay.AddSubMenuItemCheckbox("The day before yesterday", "The day before yesterday's wallpaper", false),
		types.Day3Ago:  mConfigDay.AddSubMenuItemCheckbox("Three days ago", "Three days ago's wallpaper", false),
		types.Day4Ago:  mConfigDay.AddSubMenuItemCheckbox("Four days ago", "Four days ago's wallpaper", false),
		types.Day5Ago:  mConfigDay.AddSubMenuItemCheckbox("Five days ago", "Five days ago's wallpaper", false),
		types.Day6Ago:  mConfigDay.AddSubMenuItemCheckbox("Six days ago", "Six days ago's wallpaper", false),
		types.Day7Ago:  mConfigDay.AddSubMenuItemCheckbox("Seven days ago", "Seven days ago's wallpaper", false),
	}, c.cfg, func(c *Config) types.Day { return c.Day.Value() }, func(c *Config, d types.Day) {
		logger.Logger.Printf("Setting Day: %v", d)
		c.Day.SetDefault(d)
	})

	mConfigMode := mConfig.AddSubMenuItem("Mode", "Define how the wallpaper is set")
	makeConfigSection(map[Mode]*systray.MenuItem{
		ModeCenter:  mConfigMode.AddSubMenuItemCheckbox("Center", "Center mode", false),
		ModeCrop:    mConfigMode.AddSubMenuItemCheckbox("Crop", "Crop mode", false),
		ModeFit:     mConfigMode.AddSubMenuItemCheckbox("Fit", "Fit mode", false),
		ModeSpan:    mConfigMode.AddSubMenuItemCheckbox("Span", "Span mode", false),
		ModeStretch: mConfigMode.AddSubMenuItemCheckbox("Stretch", "Stretch mode", false),
		ModeTile:    mConfigMode.AddSubMenuItemCheckbox("Tile", "Tile mode", false),
	}, c.cfg, func(c *Config) Mode { return c.Mode.Value() }, func(c *Config, m Mode) {
		logger.Logger.Printf("Setting Mode: %v", m)
		c.Mode.SetDefault(m)
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
	}, c.cfg, func(c *Config) types.Region { return c.Region.Value() }, func(c *Config, r types.Region) {
		logger.Logger.Printf("Setting Region: %v", r)
		c.Region.SetDefault(r)
	})

	mConfigResolution := mConfig.AddSubMenuItem("Resolution", "Resolution of the wallpaper")
	makeConfigSection(map[types.Resolution]*systray.MenuItem{
		types.LowDefinition:       mConfigResolution.AddSubMenuItemCheckbox("Low Definition", "Low Definition resolution", false),
		types.HighDefinition:      mConfigResolution.AddSubMenuItemCheckbox("High Definition", "High Definition resolution", false),
		types.UltraHighDefinition: mConfigResolution.AddSubMenuItemCheckbox("Ultra High Definition", "Ultra High Definition resolution", false),
	}, c.cfg, func(c *Config) types.Resolution { return c.Resolution.Value() }, func(c *Config, r types.Resolution) {
		logger.Logger.Printf("Setting Resolution: %v", r)
		c.Resolution.SetDefault(r)
	})

	mConfigDimImage := mConfig.AddSubMenuItem("Dim Image", "Dim the image")
	mConfigDimImageMap := make(map[types.Percent]*systray.MenuItem)
	for i := 0; i <= 100; i += 10 {
		mConfigDimImageMap[types.Percent(i)] = mConfigDimImage.AddSubMenuItemCheckbox(fmt.Sprintf("%d%%", i), fmt.Sprintf("%d%% dim", i), false)
	}
	makeConfigSection(mConfigDimImageMap, c.cfg, func(c *Config) types.Percent {
		// round to the nearest 10% to find the closest matching value
		return types.Percent(math.Round(float64(c.DimImage)/10) * 10)
	}, func(c *Config, p types.Percent) {
		logger.Logger.Printf("Setting DimImage: %v", p)
		c.DimImage = p
	})

	makeConfigOption(mConfig.AddSubMenuItemCheckbox("Draw Description", "Draw the wallpaper description", false), c.cfg,
		func(c *Config) bool { return c.DrawDescription },
		func(c *Config, b bool) {
			logger.Logger.Printf("Setting DrawDescription: %v", b)
			c.DrawDescription = b
		})

	makeConfigOption(mConfig.AddSubMenuItemCheckbox("Draw QR Code", "Draw the QR code", false), c.cfg,
		func(c *Config) bool { return c.DrawQRCode },
		func(c *Config, b bool) {
			logger.Logger.Printf("Setting DrawQRCode: %v", b)
			c.DrawQRCode = b
		})

	makeConfigOption(mConfig.AddSubMenuItemCheckbox("Download Only", "Download the wallpaper only", false), c.cfg,
		func(c *Config) bool { return c.DownloadOnly },
		func(c *Config, b bool) {
			logger.Logger.Printf("Setting DownloadOnly: %v", b)
			c.DownloadOnly = b
		})

	makeConfigOption(mConfig.AddSubMenuItemCheckbox("Rotate Wallpaper counter clockwise", "Rotate the wallpaper counter clockwise", false), c.cfg,
		func(c *Config) bool { return c.RotateCounterClockwise },
		func(c *Config, b bool) {
			logger.Logger.Printf("Setting RotateWallpaper: %v", b)
			c.RotateCounterClockwise = b
		})

	makeConfigOption(mConfig.AddSubMenuItemCheckbox("Use Google Text2Speech Service", "Use Google Text2Speech Service", false), c.cfg,
		func(c *Config) bool { return c.UseGoogleText2SpeechService },
		func(c *Config, b bool) {
			logger.Logger.Printf("Setting UseGoogleText2SpeechService: %v", b)
			c.UseGoogleText2SpeechService = b
		})

	makeConfigOption(mConfig.AddSubMenuItemCheckbox("Use Google Translate Service", "Use Google Translate Service", false), c.cfg,
		func(c *Config) bool { return c.UseGoogleTranslateService },
		func(c *Config, b bool) {
			logger.Logger.Printf("Setting UseGoogleTranslateService: %v", b)
			c.UseGoogleTranslateService = b
		})

	makeConfigInfo(mConfig.AddSubMenuItem("Watermark", "Watermark to be drawn on the wallpaper"), false, c.cfg,
		func(c *Config) string { return c.Watermark }, func(_ *Config, s string) { openDirectory(s) })

	makeConfigInfo(mConfig.AddSubMenuItem("Google App Credentials", "Google App Credentials"), false, c.cfg,
		func(c *Config) string { return c.GoogleAppCredentials }, func(_ *Config, s string) { openDirectory(s) })

	makeConfigInfo(mConfig.AddSubMenuItem("Furigana API AppId", "Furigana API AppId"), true, c.cfg,
		func(c *Config) string { return c.FuriganaApiAppId }, nil)

	makeConfigInfo(mConfig.AddSubMenuItem("Download Directory", "Download Directory"), false, c.cfg,
		func(c *Config) string { return c.DownloadDirectory }, func(_ *Config, s string) { openDirectory(s) })

	systray.AddSeparator()

	// Property section
	mProperties := systray.AddMenuItem("Properties", "Properties of the wallpaper")

	mPropertiesImage := mProperties.AddSubMenuItem("Image", "Image data of the wallpaper")
	makePropertyCopyAction(mPropertiesImage.
		AddSubMenuItem("Copy", "Copy the image data to the clipboard"),
		c.img, clipboard.FmtImage,
		func(i *Image) []byte {
			buf := bytes.NewBuffer(nil)
			_ = png.Encode(buf, i)
			return buf.Bytes()
		})
	makePropertyOpenAction(mPropertiesImage.AddSubMenuItem("Open", "Open the image in the browser"), c.img,
		func(i *Image) string { return "file://" + i.Location })

	mPropertiesAudio = mProperties.
		AddSubMenuItem("Audio", "Audio data of the wallpaper")
	makePropertyOpenAction(mPropertiesAudio.
		AddSubMenuItem("Open", "Open the audio in the browser"), c.img,
		func(i *Image) string {
			if i.Audio != nil {
				return "file://" + i.Audio.Location
			}

			return ""
		})

	makePropertyCopyAction(mProperties.
		AddSubMenuItem("Description", "Description of the wallpaper").
		AddSubMenuItem("Copy", "Copy the description to the clipboard"),
		c.img, clipboard.FmtText, func(i *Image) []byte { return []byte(i.Description) })

	mPropertiesSearchUrl := mProperties.AddSubMenuItem("Search URL", "Search URL of the wallpaper")
	makePropertyCopyAction(mPropertiesSearchUrl.AddSubMenuItem("Copy", "Copy the search URL to the clipboard"),
		c.img, clipboard.FmtText, func(i *Image) []byte { return []byte(i.SearchURL) })
	makePropertyOpenAction(mPropertiesSearchUrl.AddSubMenuItem("Open", "Open the search URL in the browser"), c.img,
		func(i *Image) string { return i.SearchURL })

	mPropertiesDownloadUrl := mProperties.AddSubMenuItem("Download URL", "Download URL of the wallpaper")
	makePropertyCopyAction(mPropertiesDownloadUrl.AddSubMenuItem("Copy", "Copy the download URL to the clipboard"),
		c.img, clipboard.FmtText, func(i *Image) []byte { return []byte(i.DownloadURL) })
	makePropertyOpenAction(mPropertiesDownloadUrl.AddSubMenuItem("Open", "Open the download URL in the browser"), c.img,
		func(i *Image) string { return i.DownloadURL })

	systray.AddSeparator()

	// Quit section
	mQuit = systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.SetIcon(readIcon("quit"))
	mQuit.Click(systray.Quit)

	// initial execution
	modify(func(mi *systray.MenuItem) { mi.Disable() }, mRefresh, mSpeak, mQuit)
	mPropertiesAudio.Hide()
	c.img.Update(c.execute(c.cfg))
	modify(func(mi *systray.MenuItem) { mi.Enable() }, mRefresh, mQuit)
	if c.img != nil && c.img.Audio != nil {
		mSpeak.Enable()
		mPropertiesAudio.Show()
	}
}

// OnExit is called when the application is closed.
func (c *Controller) OnExit() {
	// close the audio stream
	if c.img != nil && c.img.Audio != nil {
		_ = c.img.Audio.Close()
	}
}

// modify modifies the given menu items with the given operation.
func modify(op func(*systray.MenuItem), items ...*systray.MenuItem) {
	for _, item := range items {
		op(item)
	}
}

// makeApiCommand creates a menu item with a copy action
func makeApiCommand(item *systray.MenuItem, addr, verb, path, query, payload string) {
	item.SetIcon(readIcon("copy"))
	item.Click(func() {
		if clipboardErr != nil {
			logger.Logger.Printf("Failed to initialize clipboard: %v", clipboardErr)
			return
		}
		uri := &url.URL{
			Scheme:   "http",
			Host:     addr,
			Path:     path,
			RawQuery: query,
		}
		cmd := fmt.Sprintf("curl -X %s %s", verb, uri.String())
		if payload != "" {
			cmd += " -d '" + payload + "'"
		}
		_ = clipboard.Write(clipboard.FmtText, []byte(cmd))
	})
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
			logger.Logger.Printf("Failed to initialize clipboard: %v", clipboardErr)
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
			logger.Logger.Printf("Failed to open %s: %v", item, err)
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
			logger.Logger.Printf("Failed to extract watermarks: %v", err)
			return
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		logger.Logger.Printf("Failed to open directory %s: %v", path, err)
		return
	}

	if !info.IsDir() {
		path = filepath.Dir(path)
	}

	if err := browser.OpenURL("file://" + path); err != nil {
		logger.Logger.Printf("Failed to open directory %s: %v", path, err)
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
		logger.Logger.Printf("Failed to read icon %s: %s", name+ext, err)
		return nil
	}

	return data
}

// Run executes the given function with the given configuration.
func Run(execute func(*Config) *Image, cfg *Config) {
	img := &Image{}
	if !cfg.Daemon {
		img.Update(execute(cfg))
		return
	}

	controller := &Controller{img: img, cfg: cfg, execute: execute}
	server := NewServer(cfg, controller)
	defer func() {
		if err := server.Stop(); err != nil {
			logger.Logger.Printf("Failed to stop API server: %v", err)
		}
	}()

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatalf("Failed to start API server: %v", err)
		}
	}()

	systray.Run(controller.OnReady, controller.OnExit)
}
