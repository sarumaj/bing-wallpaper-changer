package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/core"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/spf13/pflag"
)

// name of remote code repository mirror
const remoteRepository = "sarumaj/bing-wallpaper-changer"

// BuildDate is the date when the binary was built.
var BuildDate = "2024-12-20 21:07:32 UTC"

// Version is the version of the binary.
var Version = "v1.1.5"

func main() {
	var config core.Config
	checkVersionOrUpdate()
	parseArgs(&config, os.Args[1:]...)
	core.Run(execute, &config)
}

// checkVersionOrUpdate checks if there is a new version available and updates the binary if necessary.
func checkVersionOrUpdate() {
	parsed, err := semver.ParseTolerant(Version)
	if err != nil {
		logger.Logger.Printf("Failed to parse version: %s", err)
		return
	}

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{APIToken: ""})
	if err != nil {
		logger.Logger.Printf("Failed to setup source: %s", err)
		return
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{Source: source, Validator: &selfupdate.SHAValidator{}})
	if err != nil {
		logger.Logger.Printf("Failed to setup updater: %s", err)
		return
	}

	repository := selfupdate.ParseSlug(remoteRepository)
	latest, found, err := updater.DetectLatest(context.Background(), repository)
	if err != nil {
		logger.Logger.Printf("Failed to detect latest version: %s", err)
		return
	}

	if !found {
		logger.Logger.Printf("No update found")
		return
	}

	if latest.GreaterThan(parsed.String()) {
		if _, err := updater.UpdateSelf(context.Background(), parsed.String(), repository); err != nil {
			logger.Logger.Printf("Failed to update: %s", err)
			return
		}

		logger.Logger.Printf("Updated to version %s", latest.Version())
		return
	}

	logger.Logger.Printf("Current version %s is the latest", Version)
}

// execute fetches the wallpaper, processes it, and sets it as the desktop wallpaper.
func execute(config *core.Config) *core.Image {
	img, err := core.DownloadAndDecode(
		config.Day.Value(), config.Region.Value(), config.Resolution.Value(),
		core.WithFuriganaApiAppId(config.FuriganaApiAppId),
		core.WithGoogleAppCredentials(config.GoogleAppCredentials),
		core.WithUseGoogleText2SpeechService(config.UseGoogleText2SpeechService),
		core.WithUseGoogleTranslateService(config.UseGoogleTranslateService),
	)
	if err != nil {
		logger.Logger.Println(err)
		return nil
	}

	if config.DimImage > 0.0 {
		if err := img.Dim(config.DimImage); err != nil {
			logger.Logger.Println(err)
			return img
		}
	}

	if config.Watermark != "" {
		if err := img.DrawWatermark(config.Watermark, config.RotateCounterClockwise); err != nil {
			logger.Logger.Println(err)
			return img
		}
	}

	if config.DrawDescription {
		if err := img.DrawDescription(types.PositionTopCenter, extras.DefaultFontName); err != nil {
			logger.Logger.Println(err)
			return img
		}
	}

	if config.DrawQRCode {
		if err := img.DrawQRCode(config.Resolution.Value(), types.PositionTopRight); err != nil {
			logger.Logger.Println(err)
			return img
		}
	}

	path, err := img.EncodeAndDump(config.DownloadDirectory)
	if err != nil {
		logger.Logger.Println(err)
		return img
	}

	logger.Logger.Printf("Wallpaper saved to: %s", path)
	if !config.DownloadOnly {
		if err := core.SetWallpaper(path, config.Mode.Value()); err != nil {
			logger.Logger.Println(err)
			return img
		}

		logger.Logger.Printf("Wallpaper set to: %s", path)
	}

	if img.Audio == nil {
		return img
	}

	logger.Logger.Println("Playing audio description")
	if err := img.Audio.Play(); err != nil {
		logger.Logger.Printf("Failed to play audio: %v", err)
		return img
	}

	logger.Logger.Println("Audio description played")
	return img
}

// parseArgs parses the command line arguments and sets the configuration accordingly.
func parseArgs(config *core.Config, args ...string) {
	config.Day.SetDefault(types.DayToday)
	config.Day.SetValues(types.AllowedDays...)

	config.Mode.SetDefault(core.ModeFit)
	config.Mode.SetValues(core.AllowedModes...)

	config.Region.SetDefault(types.RegionGermany)
	config.Region.SetValues(types.AllowedRegions...)

	config.Resolution.SetAlias(func(r types.Resolution) string { return r.Alias })
	config.Resolution.SetDefault(types.HighDefinition)
	config.Resolution.SetValues(types.AllowedResolutions...)

	opts := pflag.NewFlagSet("bing-wallpaper-changer", pflag.ContinueOnError)
	opts.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of bing-wallpaper-changer [Version: %s, BuildDate: %s]:\n\n", Version, BuildDate)
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n\n")
		opts.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "")
	}

	defaultDownloadDirectory, _ := os.UserHomeDir()
	defaultDownloadDirectory = filepath.Join(defaultDownloadDirectory, "Pictures", "BingWallpapers")

	opts.Var(&config.Day, "day", fmt.Sprintf("the day to fetch the wallpaper for, allowed values are: %s", config.Day.Values()))
	opts.Var(&config.Mode, "mode", fmt.Sprintf("the mode of the wallpaper, allowed values are: %s", config.Mode.Values()))
	opts.Var(&config.Region, "region", fmt.Sprintf("the region to fetch the wallpaper for, allowed values are: %s", config.Region.Values()))
	opts.Var(&config.Resolution, "resolution", fmt.Sprintf("the resolution of the wallpaper, allowed values are: %s", config.Resolution.Values()))
	opts.BoolVar(&config.DrawDescription, "description", true, "draw the description on the wallpaper")
	opts.BoolVar(&config.DrawQRCode, "qrcode", true, "draw the QR code on the wallpaper")
	opts.StringVar(&config.Watermark, "watermark", extras.DefaultWatermarkName, "draw the watermark on the wallpaper")
	opts.BoolVar(&config.DownloadOnly, "download-only", false, "download the wallpaper only")
	opts.StringVar(&config.DownloadDirectory, "download-directory", defaultDownloadDirectory, "the directory to download the wallpaper to")
	opts.BoolVar(&config.RotateCounterClockwise, "rotate-counter-clockwise", false, "rotate the watermark counter-clockwise if necessary (default is clockwise)")
	opts.StringVar(&config.GoogleAppCredentials, "google-app-credentials", "", fmt.Sprintf("the path to the Google App credentials file for the translation service for %s to %s,\nif not provided, the translation service will not be used", types.NonEnglishRegions, types.RegionUnitedStates))
	opts.StringVar(&config.FuriganaApiAppId, "furigana-api-app-id", "", "the Goo Labs API App ID (labs.goo.ne.jp) for the furigana service, if not provided, Jisho.org (if available) or github.com/sarumaj/go-kakasi will be used")
	opts.BoolVar(&config.UseGoogleText2SpeechService, "use-google-text2speech-service", false, "use the Google Text2Speech service to record and play the audio description (not supported on darwin, and linux unless compiled with cgo)")
	opts.BoolVar(&config.UseGoogleTranslateService, "use-google-translate-service", false, "use the Google Translate service to translate the description to English")
	opts.BoolVar(&config.Daemon, "daemon", false, "run the application as a daemon process")
	opts.BoolVar(&config.Debug, "debug", false, "enable debug mode")
	opts.Var(&config.DimImage, "dim-image", "dim the image by the given percentage (0.0 to 100.0)")
	opts.IntVar(&config.ApiPort, "api-port", 44244, "the port number of the API server")

	if err := opts.Parse(args); err != nil {
		if !errors.Is(err, pflag.ErrHelp) {
			logger.Logger.Fatalln(err)
		}
		os.Exit(0)
	}

	if config.Debug {
		logger.Logger.SetLevel(logger.LogLevelDebug)
		logger.Logger.SetLevel(logger.LogLevelDebug)
	}
}
