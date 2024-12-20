package core

import (
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

const AppName = "Bing Wallpaper Changer"

type Config struct {
	Day                         types.Day        `json:"day"`
	Region                      types.Region     `json:"region"`
	Resolution                  types.Resolution `json:"resolution"`
	DrawDescription             bool             `json:"drawDescription"`
	DrawQRCode                  bool             `json:"drawQRCode"`
	Watermark                   string           `json:"watermark"`
	DownloadOnly                bool             `json:"downloadOnly"`
	DownloadDirectory           string           `json:"downloadDirectory"`
	RotateCounterClockwise      bool             `json:"rotateCounterClockwise"`
	GoogleAppCredentials        string           `json:"googleAppCredentials"`
	FuriganaApiAppId            string           `json:"furiganaApiAppId"`
	UseGoogleText2SpeechService bool             `json:"useGoogleText2SpeechService"`
}
