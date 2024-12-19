package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/core"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/spf13/pflag"
)

var config struct {
	Day                         types.Day
	Region                      types.Region
	Resolution                  types.Resolution
	DrawDescription             bool
	DrawQRCode                  bool
	Watermark                   string
	DownloadOnly                bool
	DownloadDirectory           string
	RotateCounterClockwise      bool
	GoogleAppCredentials        string
	FuriganaApiAppId            string
	UseGoogleText2SpeechService bool
}

// name of remote code repository mirror
const remoteRepository = "sarumaj/bing-wallpaper-changer"

// BuildDate is the date when the binary was built.
var BuildDate = "2024-12-19 09:34:58 UTC"

// Version is the version of the binary.
var Version = "v1.0.13"

func main() {
	checkVersionOrUpdate()
	parseArgs(os.Args[1:]...)

	img, err := core.DownloadAndDecode(
		config.Day, config.Region, config.Resolution,
		core.WithFuriganaApiAppId(config.FuriganaApiAppId),
		core.WithGoogleAppCredentials(config.GoogleAppCredentials),
		core.WithUseGoogleText2SpeechService(config.UseGoogleText2SpeechService),
	)
	if err != nil {
		logger.ErrLogger.Fatalln(err)
	}

	if config.Watermark != "" {
		if err := img.DrawWatermark(config.Watermark, config.RotateCounterClockwise); err != nil {
			logger.ErrLogger.Fatalln(err)
		}
	}

	if config.DrawDescription {
		if err := img.DrawDescription(types.TopCenter, extras.DefaultFontName); err != nil {
			logger.ErrLogger.Fatalln(err)
		}
	}

	if config.DrawQRCode {
		if err := img.DrawQRCode(config.Resolution, types.TopRight); err != nil {
			logger.ErrLogger.Fatalln(err)
		}
	}

	path, err := img.EncodeAndDump(config.DownloadDirectory)
	if err != nil {
		logger.ErrLogger.Fatalln(err)
	}

	logger.InfoLogger.Printf("Wallpaper saved to: %s", path)
	if !config.DownloadOnly {
		if err := core.SetWallpaper(path, core.ModeStretch); err != nil {
			logger.ErrLogger.Fatalln(err)
		}

		logger.InfoLogger.Printf("Wallpaper set to: %s", path)
	}

	logger.InfoLogger.Println("Playing audio description")
	if err := img.Audio.Play(); err != nil {
		logger.ErrLogger.Printf("Failed to play audio: %v", err)

	}
}

// checkVersionOrUpdate checks if there is a new version available and updates the binary if necessary.
func checkVersionOrUpdate() {
	parsed, err := semver.ParseTolerant(Version)
	if err != nil {
		logger.ErrLogger.Printf("Failed to parse version: %s", err)
		return
	}

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{APIToken: ""})
	if err != nil {
		logger.ErrLogger.Printf("Failed to setup source: %s", err)
		return
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{Source: source, Validator: &selfupdate.SHAValidator{}})
	if err != nil {
		logger.ErrLogger.Printf("Failed to setup updater: %s", err)
		return
	}

	repository := selfupdate.ParseSlug(remoteRepository)
	latest, found, err := updater.DetectLatest(context.Background(), repository)
	if err != nil {
		logger.ErrLogger.Printf("Failed to detect latest version: %s", err)
		return
	}

	if !found {
		logger.ErrLogger.Printf("No update found")
		return
	}

	if latest.GreaterThan(parsed.String()) {
		if _, err := updater.UpdateSelf(context.Background(), parsed.String(), repository); err != nil {
			logger.ErrLogger.Printf("Failed to update: %s", err)
			return
		}

		logger.InfoLogger.Printf("Updated to version %s", latest.Version())
		return
	}

	logger.InfoLogger.Printf("Current version %s is the latest", Version)
}

// parseArgs parses the command line arguments and sets the configuration accordingly.
func parseArgs(args ...string) {
	var day int
	var region string
	var resolution string

	opts := pflag.NewFlagSet("bing-wallpaper-changer", pflag.ContinueOnError)
	opts.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of bing-wallpaper-changer [Version: %s, BuildDate: %s]:\n\n", Version, BuildDate)
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n\n")
		opts.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "")
	}

	defaultDownloadDirectory, _ := os.UserHomeDir()
	defaultDownloadDirectory += "/Pictures/BingWallpapers"

	opts.IntVar(&day, "day", int(types.Today), "the day to fetch the wallpaper for, 0 is today, 1 is yesterday, and so on, 7 is the highest value, which is seven days ago")
	opts.StringVar(&region, "region", types.Germany.String(), fmt.Sprintf("the region to fetch the wallpaper for, allowed values are: %s", types.AllowedRegions))
	opts.StringVar(&resolution, "resolution", types.HighDefinition.String(), fmt.Sprintf("the resolution of the wallpaper, allowed values are: %s", types.AllowedResolutions))
	opts.BoolVar(&config.DrawDescription, "description", true, "draw the description on the wallpaper")
	opts.BoolVar(&config.DrawQRCode, "qrcode", true, "draw the QR code on the wallpaper")
	opts.StringVar(&config.Watermark, "watermark", extras.DefaultWatermarkName, "draw the watermark on the wallpaper")
	opts.BoolVar(&config.DownloadOnly, "download-only", false, "download the wallpaper only")
	opts.StringVar(&config.DownloadDirectory, "download-directory", defaultDownloadDirectory, "the directory to download the wallpaper to")
	opts.BoolVar(&config.RotateCounterClockwise, "rotate-counter-clockwise", false, "rotate the watermark counter-clockwise if necessary (default is clockwise)")
	opts.StringVar(&config.GoogleAppCredentials, "google-app-credentials", "", fmt.Sprintf("the path to the Google App credentials file for the translation service for %s to %s,\nif not provided, the translation service will not be used", types.NonEnglishRegions, types.UnitedStates))
	opts.StringVar(&config.FuriganaApiAppId, "furigana-api-app-id", "", "the Goo Labs API App ID (labs.goo.ne.jp) for the furigana service, if not provided, github.com/sarumaj/go-kakasi will be used")
	opts.BoolVar(&config.UseGoogleText2SpeechService, "use-google-text2speech-service", false, "use the Google Text2Speech service to record and play the audio description")

	if err := opts.Parse(args); err != nil {
		if !errors.Is(err, pflag.ErrHelp) {
			logger.ErrLogger.Println(err)
		}
		os.Exit(0)
	}

	var err error
	config.Region, err = types.ParseLocale(region)
	if err != nil {
		opts.Usage()
		logger.ErrLogger.Fatalln(err)
	}

	config.Day = types.Day(day)
	if err := config.Day.IsValid(); err != nil {
		opts.Usage()
		logger.ErrLogger.Fatalln(err)
	}

	config.Resolution, err = types.ParseResolution(resolution)
	if err != nil {
		opts.Usage()
		logger.ErrLogger.Fatalln(err)
	}
}
