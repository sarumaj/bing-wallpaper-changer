//go:build (!darwin && !linux) || cgo

package core

import (
	"embed"
	"runtime"

	"github.com/energye/systray"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

//go:embed icons/*.png icons/*.ico
var icons embed.FS

// disable the menu items
func disable(items ...*systray.MenuItem) {
	for _, item := range items {
		item.Disable()
	}
}

// enable the menu items
func enable(items ...*systray.MenuItem) {
	for _, item := range items {
		item.Enable()
	}
}

// makeConfigInfo creates a read-only menu item
func makeConfigInfo(option *systray.MenuItem, cfg *Config, lookup func(*Config) string) {
	option.AddSubMenuItem(lookup(cfg), "Not editable").Disable()
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

// setIcon sets the icon of the menu item
func setIcon(name string, setter func([]byte)) {
	data, err := icons.ReadFile("icons/" + name)
	if err != nil {
		logger.ErrLogger.Printf("Failed to read icon %s: %s", name, err)
		return
	}

	setter(data)
}

// ShowTray shows the tray icon and menu
func ShowTray(execute func(*Config) *Image, cfg *Config) {
	ext := ".png"
	if runtime.GOOS == "windows" {
		ext = ".ico"
	}

	var img *Image
	onReady := func() {
		systray.SetTitle(AppName)
		systray.SetTooltip(AppName)
		setIcon("wallpaper"+ext, systray.SetIcon)

		var mRefresh, mSpeak, mQuit *systray.MenuItem

		// Main section
		mRefresh = systray.AddMenuItem("Refresh", "Refresh the wallpaper")
		setIcon("refresh"+ext, mRefresh.SetIcon)
		mRefresh.Click(func() {
			disable(mRefresh, mSpeak, mQuit)
			img = execute(cfg)
			enable(mRefresh, mSpeak, mQuit)
		})

		mSpeak = systray.AddMenuItem("Speak", "Speak the wallpaper description")
		setIcon("play"+ext, mSpeak.SetIcon)
		mSpeak.Click(func() {
			disable(mRefresh, mSpeak, mQuit)
			if img.Audio != nil {
				if err := img.Audio.Play(); err != nil {
					logger.ErrLogger.Printf("Failed to play audio: %v", err)
				}
			}
			enable(mRefresh, mSpeak, mQuit)
		})

		systray.AddSeparator()
		// Config section

		mConfig := systray.AddMenuItem("Configure", "Configure the wallpaper settings")
		mConfigDay := mConfig.AddSubMenuItem("Day", "Day of the Bing wallpaper")
		makeConfigSection(map[types.Day]*systray.MenuItem{
			types.Today:                 mConfigDay.AddSubMenuItemCheckbox("Today", "Today's wallpaper", false),
			types.Yesterday:             mConfigDay.AddSubMenuItemCheckbox("Yesterday", "Yesterday's wallpaper", false),
			types.TheDayBeforeYesterday: mConfigDay.AddSubMenuItemCheckbox("The day before yesterday", "The day before yesterday's wallpaper", false),
			types.ThreeDaysAgo:          mConfigDay.AddSubMenuItemCheckbox("Three days ago", "Three days ago's wallpaper", false),
			types.FourDaysAgo:           mConfigDay.AddSubMenuItemCheckbox("Four days ago", "Four days ago's wallpaper", false),
			types.FiveDaysAgo:           mConfigDay.AddSubMenuItemCheckbox("Five days ago", "Five days ago's wallpaper", false),
			types.SixDaysAgo:            mConfigDay.AddSubMenuItemCheckbox("Six days ago", "Six days ago's wallpaper", false),
			types.SevenDaysAgo:          mConfigDay.AddSubMenuItemCheckbox("Seven days ago", "Seven days ago's wallpaper", false),
		}, cfg, func(c *Config) types.Day { return c.Day }, func(c *Config, d types.Day) {
			logger.InfoLogger.Printf("Setting Day: %v", d)
			c.Day = d
		})

		mConfigRegion := mConfig.AddSubMenuItem("Region", "Region of the Bing wallpaper")
		makeConfigSection(map[types.Region]*systray.MenuItem{
			types.Canada:        mConfigRegion.AddSubMenuItemCheckbox("Canada", "Canada region", false),
			types.China:         mConfigRegion.AddSubMenuItemCheckbox("China", "China region", false),
			types.Germany:       mConfigRegion.AddSubMenuItemCheckbox("Germany", "Germany region", false),
			types.Japan:         mConfigRegion.AddSubMenuItemCheckbox("Japan", "Japan region", false),
			types.NewZealand:    mConfigRegion.AddSubMenuItemCheckbox("New Zealand", "New Zealand region", false),
			types.UnitedKingdom: mConfigRegion.AddSubMenuItemCheckbox("United Kingdom", "United Kingdom region", false),
			types.UnitedStates:  mConfigRegion.AddSubMenuItemCheckbox("United States", "United States region", false),
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

		makeConfigInfo(mConfig.AddSubMenuItem("Watermark", "Watermark to be drawn on the wallpaper"), cfg,
			func(c *Config) string { return c.Watermark })

		makeConfigInfo(mConfig.AddSubMenuItem("Google App Credentials", "Google App Credentials"), cfg,
			func(c *Config) string { return c.GoogleAppCredentials })

		makeConfigInfo(mConfig.AddSubMenuItem("Furigana API AppId", "Furigana API AppId"), cfg,
			func(c *Config) string {
				if c.FuriganaApiAppId != "" {
					return "[redacted]"
				}

				return "not set"
			})

		makeConfigInfo(mConfig.AddSubMenuItem("Download Directory", "Download Directory"), cfg,
			func(c *Config) string { return c.DownloadDirectory })

		systray.AddSeparator()
		// Quit section

		mQuit = systray.AddMenuItem("Quit", "Quit the whole app")
		setIcon("quit"+ext, mQuit.SetIcon)
		mQuit.Click(systray.Quit)

		// initial execution
		disable(mRefresh, mSpeak, mQuit)
		img = execute(cfg)
		enable(mRefresh, mSpeak, mQuit)
	}

	onExit := func() {
		// close the audio stream
		if img.Audio != nil {
			_ = img.Audio.Close()
		}
	}

	systray.Run(onReady, onExit)
}
