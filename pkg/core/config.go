package core

import (
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
)

const AppName = "Bing Wallpaper Changer"

type Config struct {
	ApiPort                     int                                             `json:"apiPort"`
	AutoPlayAudio               bool                                            `json:"autoPlayAudio"`
	Day                         types.Enum[types.Day, types.Days]               `json:"day"`
	Mode                        types.Enum[Mode, Modes]                         `json:"mode"`
	Region                      types.Enum[types.Region, types.Regions]         `json:"region"`
	Resolution                  types.Enum[types.Resolution, types.Resolutions] `json:"resolution"`
	DrawDescription             bool                                            `json:"drawDescription"`
	DrawQRCode                  bool                                            `json:"drawQRCode"`
	Watermark                   string                                          `json:"watermark"`
	DownloadOnly                bool                                            `json:"downloadOnly"`
	DownloadDirectory           string                                          `json:"downloadDirectory"`
	RotateCounterClockwise      bool                                            `json:"rotateCounterClockwise"`
	GoogleAppCredentials        string                                          `json:"googleAppCredentials"`
	FuriganaApiAppId            string                                          `json:"furiganaApiAppId"`
	UseGoogleText2SpeechService bool                                            `json:"useGoogleText2SpeechService"`
	UseGoogleTranslateService   bool                                            `json:"useGoogleTranslateService"`
	Daemon                      bool                                            `json:"daemon"`
	Debug                       bool                                            `json:"debug"`
	DimImage                    types.Percent                                   `json:"dimImage"`
}
