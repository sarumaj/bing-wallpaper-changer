[![test_and_report](https://github.com/sarumaj/bing-wallpaper-changer/actions/workflows/test_and_report.yml/badge.svg)](https://github.com/sarumaj/bing-wallpaper-changer/actions/workflows/test_and_report.yml)
[![build_and_release](https://github.com/sarumaj/bing-wallpaper-changer/actions/workflows/build_and_release.yml/badge.svg)](https://github.com/sarumaj/bing-wallpaper-changer/actions/workflows/build_and_release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sarumaj/bing-wallpaper-changer)](https://goreportcard.com/report/github.com/sarumaj/bing-wallpaper-changer)
[![Maintainability](https://img.shields.io/codeclimate/maintainability-percentage/sarumaj/bing-wallpaper-changer.svg)](https://codeclimate.com/github/sarumaj/bing-wallpaper-changer/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/147f265284b27931c2d2/test_coverage)](https://codeclimate.com/github/sarumaj/bing-wallpaper-changer/test_coverage)
[![Go Reference](https://pkg.go.dev/badge/github.com/sarumaj/bing-wallpaper-changer.svg)](https://pkg.go.dev/github.com/sarumaj/bing-wallpaper-changer)
[![Go version](https://img.shields.io/github/go-mod/go-version/sarumaj/bing-wallpaper-changer?logo=go&label=&labelColor=gray)](https://go.dev)
[![Release](https://img.shields.io/github/v/release/sarumaj/bing-wallpaper-changer?logo=github)](https://github.com/sarumaj/bing-wallpaper-changer/releases/latest)
[![Release Date](https://img.shields.io/github/release-date/sarumaj/bing-wallpaper-changer?logo=github)](https://github.com/sarumaj/bing-wallpaper-changer/releases/latest)
[![Commits since latest release](https://img.shields.io/github/commits-since/sarumaj/bing-wallpaper-changer/latest?logo=github)](https://github.com/sarumaj/bing-wallpaper-changer/releases/latest)
[![Downloads (all assets, all releases)](https://img.shields.io/github/downloads/sarumaj/bing-wallpaper-changer/total?logo=github)](https://github.com/sarumaj/bing-wallpaper-changer/releases)
[![Downloads (all assets, latest release)](https://img.shields.io/github/downloads/sarumaj/bing-wallpaper-changer/latest/total?logo=github)](https://github.com/sarumaj/bing-wallpaper-changer/releases/latest)

---

# bing-wallpaper-changer

**bing-wallpaper-changer** is a cross-platform compatible wallpaper-changer (CLI).

It fetches the newest Bing wallpaper and sets it as a desktop background image.
Custom watermark can be used on the downloaded image.
Done just for fun 😄

## Features

- [x] Crawl and fetch newest Bind wallpaper
  - [x] Support multiple regions
  - [x] Support multiple screen resolutions (😡 UltraHD is broken on the Bing side)
  - [x] Download wallpapers up to seven days in the past
- [x] Draw title on wallpapers
- [x] Place QR code for the copyright links
- [x] Draw watermarks
  - [x] Scale down/up to match the resolution of the wallpaper
  - [x] Rotate if necessary (only clockwise rotation by 90° supported)

## Usage

```console
$ bing-wallpaper-changer -h
>
> Usage: bing-wallpaper-changer [flags]
>       --day int                     the day to fetch the wallpaper for, 0 is today, 1 is yesterday, and so on, 7 is the highest value, which is seven days ago
>       --description                 draw the description on the wallpaper (default true)
>       --download-directory string   the directory to download the wallpaper to (default "~/Pictures/BingWallpapers")
>       --download-only               download the wallpaper only
>       --qrcode                      draw the QR code on the wallpaper (default true)
>       --region string               the region to fetch the wallpaper for, allowed values are: en-CA, zh-CN, de-DE, ja-JP, en-NZ, en-GB, en-US (default "de-DE")
>       --resolution string           the resolution of the wallpaper, allowed values are: 1366x768, 1920x1080, 3840x2160 (default "1920x1080")
>       --rotate-counter-clockwise    rotate the watermark counter-clockwise if necessary (default is clockwise)
>       --watermark string            draw the watermark on the wallpaper (default "sarumaj.png")

```

## Examples

### Default

Using default parameters:

![Bing Wallpaper of the day with QR code, default watermark and title](demo/default.png)

### Resized watermark

Using small PNG watermark: [red-dot.png](pkg/extras/watermarks/red-dot.png)

![Bing Wallpaper of the day with QR code, red-dot watermark and title](demo/red-dot.png)

### Rotated watermark

Using vertical (portrait-mode) PNG watermark: [car.png](pkg/extras/watermarks/car.png)

![Bing Wallpaper of the day with QR code, car watermark and title](demo/car.png)

### Fetching Bing wallpaper for the ja-JP region

Using default parameters with region set to `ja-JP`:

![Bing Wallpaper of the day for ja-JP region with QR code, default watermark and title](demo/unicode.png)
