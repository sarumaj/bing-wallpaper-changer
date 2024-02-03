package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/core"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/spf13/pflag"
)

var config struct {
	Day               types.Day
	Region            types.Region
	Resolution        types.Resolution
	DrawDescription   bool
	DrawQRCode        bool
	Watermark         string
	DownloadOnly      bool
	DownloadDirectory string
}

var logger = log.New(os.Stderr, "bing-wall: ", 0)

func init() {
	var day int
	var region string
	var resolution string

	opts := pflag.NewFlagSet("bing-wallpaper-changer", pflag.ContinueOnError)
	opts.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of bing-wallpaper-changer:\n\n")
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
	if err := opts.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, pflag.ErrHelp) {
			logger.Println(err)
		}
		os.Exit(0)
	}

	var err error
	config.Region, err = types.ParseLocale(region)
	if err != nil {
		opts.Usage()
		logger.Fatalln(err)
	}

	config.Day = types.Day(day)
	if err := config.Day.IsValid(); err != nil {
		opts.Usage()
		logger.Fatalln(err)
	}

	config.Resolution, err = types.ParseResolution(resolution)
	if err != nil {
		opts.Usage()
		logger.Fatalln(err)
	}
}

func main() {
	img, err := core.DownloadAndDecode(config.Day, config.Region, config.Resolution)
	if err != nil {
		pflag.Usage()
		logger.Fatalln(err)
	}

	if config.Watermark != "" {
		if err := img.DrawWatermark(config.Watermark); err != nil {
			pflag.Usage()
			logger.Fatalln(err)
		}
	}

	if config.DrawDescription {
		if err := img.DrawDescription(types.TopCenter, extras.DefaultFontName); err != nil {
			pflag.Usage()
			logger.Fatalln(err)
		}
	}

	if config.DrawQRCode {
		if err := img.DrawQRCode(config.Resolution, types.BottomLeft); err != nil {
			pflag.Usage()
			logger.Fatalln(err)
		}
	}

	path, err := img.EncodeAndDump(config.DownloadDirectory)
	if err != nil {
		pflag.Usage()
		logger.Fatalln(err)
	}

	logger.Printf("Wallpaper saved to: %s", path)
	if !config.DownloadOnly {
		if err := core.SetWallpaper(path); err != nil {
			pflag.Usage()
			logger.Fatalln(err)
		}

		logger.Printf("Wallpaper set to: %s", path)
	}
}
