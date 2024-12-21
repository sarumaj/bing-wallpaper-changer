package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"unicode"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/sarumaj/go-kakasi"
	"github.com/tidwall/gjson"
	"google.golang.org/api/option"
)

const (
	defaultBingUrl        = "https://www.bing.com"
	defaultFuriganaApiUrl = "https://labs.goo.ne.jp"
)

type (
	crawlerConfig struct {
		bingUrl                     string
		googleAppCredentials        string
		furiganaApiAppId            string
		furiganaApiUrl              string
		useGoogleText2SpeechService bool
		useGoogleTranslateService   bool
	}

	crawlerConfigOption func(*crawlerConfig)
)

// configuration for the Bing, Goo Labs APIs and Google Cloud Translation Service.
var cfg = crawlerConfig{
	bingUrl:        defaultBingUrl,
	furiganaApiUrl: defaultFuriganaApiUrl,
}

// furiganizeGooLabsApi annotates the description in Japanese with Furigana for Kanji sequences.
// It uses the Goo Labs API to convert Kanji to Furigana.
func furiganizeGooLabsApi(description string) (string, error) {
	// select kanji sequences from the description.
	var kanji []rune
	for _, r := range description {
		if unicode.IsOneOf([]*unicode.RangeTable{unicode.Ideographic}, r) {
			if len(kanji) == 0 || kanji[len(kanji)-1] == ']' {
				kanji = append(kanji, '[')
			}
			kanji = append(kanji, r)

		} else if len(kanji) > 0 && kanji[len(kanji)-1] != ']' {
			kanji = append(kanji, ']')

		}
	}

	// request furigana conversion.
	raw, err := readResponse(retryablehttp.PostForm(
		cfg.furiganaApiUrl+"/api/hiragana",
		url.Values{
			"app_id":      {cfg.furiganaApiAppId},
			"sentence":    {string(kanji)},
			"output_type": {"hiragana"},
		},
	))
	if err != nil {
		return "", err
	}

	// tokens are kanji sequences enclosed in square brackets.
	tokens := strings.Split(strings.Trim(string(kanji), "[]"), "][")

	// converted is the hiragana sequence of the kanji sequence.
	converted := strings.ReplaceAll(gjson.GetBytes(raw, "converted").String(), " ", "")

	// annotations are the hiragana sequences enclosed in square brackets.
	annotations := strings.SplitAfterN(converted, "]", len(tokens))

	var replacements []string
	for i, t := range tokens {
		replacements = append(replacements, t, t+annotations[i])
	}

	// replace the kanji sequences with the hiragana sequences avoiding collisions.
	return strings.NewReplacer(replacements...).Replace(description), nil
}

// furiganizeKakasi annotates the description in Japanese with Furigana for Kanji sequences.
// It uses the kakasi NLP library to convert Kanji to Furigana.
func furiganizeKakasi(description string) (string, error) {
	k, err := kakasi.NewKakasi()
	if err != nil {
		return "", err
	}

	converted, err := k.Convert(description)
	if err != nil {
		return "", err
	}

	return k.Normalize(converted.Furiganize())
}

// readResponse reads the response body and returns the content.
func readResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// speakDescription generates an audio stream reader using Google Cloud Text-to-Speech Service.
func speakDescription(description string, languageCode string) (*Audio, error) {
	contents, err := os.ReadFile(cfg.googleAppCredentials)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsJSON(contents))
	if err != nil {
		return nil, err
	}

	defer client.Close()

	voices_result, err := client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{LanguageCode: languageCode})
	if err != nil {
		return nil, err
	}

	if len(voices_result.Voices) == 0 {
		return nil, fmt.Errorf("no voices found")
	}

	var voice *texttospeechpb.Voice
	for _, v := range voices_result.Voices {
		if v.SsmlGender == texttospeechpb.SsmlVoiceGender_FEMALE {
			voice = v
			break
		}

	}

	if voice == nil {
		voice = voices_result.Voices[0]
	}

	result, err := client.SynthesizeSpeech(ctx, &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: description},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: languageCode,
			SsmlGender:   voice.SsmlGender,
			Name:         voice.Name,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Audio{
		Encoding:   texttospeechpb.AudioEncoding_MP3.String(),
		Source:     bytes.NewReader(result.AudioContent),
		SampleRate: int(voice.NaturalSampleRateHertz),
	}, nil
}

// translateDescription translates the description from the source language to the target language.
// It uses the Google Cloud Translation Service to translate the description.
func translateDescription(description string, source, target string) (string, error) {
	contents, err := os.ReadFile(cfg.googleAppCredentials)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsJSON(contents))
	if err != nil {
		return "", err
	}
	defer client.Close()

	location := fmt.Sprintf("projects/%s/locations/global", gjson.GetBytes(contents, "project_id").String())
	logger.InfoLogger.Println("Using Google Translation service location:", location)

	result, err := client.TranslateText(ctx, &translatepb.TranslateTextRequest{
		Contents:           []string{description},
		MimeType:           "text/plain",
		SourceLanguageCode: source,
		TargetLanguageCode: target,
		Parent:             location,
		Labels:             map[string]string{"requestor": "bing-wallpaper-changer"},
	})
	if err != nil {
		return "", err
	}

	if len(result.Translations) == 0 {
		return "", fmt.Errorf("no translations found")
	}

	return result.Translations[0].TranslatedText, nil
}

// DownloadAndDecode fetches the Bing wallpaper and decodes it.
func DownloadAndDecode(day types.Day, region types.Region, resolution types.Resolution, opts ...crawlerConfigOption) (*Image, error) {
	for _, opt := range opts {
		opt(&cfg)
	}

	jsonRaw, err := readResponse(retryablehttp.Get(
		cfg.bingUrl + "/HPImageArchive.aspx?" + url.Values{
			"format": {"js"},
			"idx":    {fmt.Sprintf("%d", day)},
			"n":      {"1"},
			"mkt":    {region.String()},
		}.Encode(),
	))
	if err != nil {
		return nil, err
	}

	path := gjson.GetBytes(jsonRaw, "images.0.url").String()
	path = regexp.MustCompile(`_\d+x\d+`).ReplaceAllString(path, "_"+resolution.String())

	parsedRequestUri, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, err
	}

	remoteHostUrl, err := url.Parse(cfg.bingUrl)
	if err != nil {
		return nil, err
	}

	parsedRequestUri.Host = remoteHostUrl.Host
	parsedRequestUri.Scheme = remoteHostUrl.Scheme

	decoder, err := getDecoder(parsedRequestUri.Query().Get("id"))
	if err != nil {
		return nil, err
	}

	content, err := readResponse(retryablehttp.Get(parsedRequestUri.String()))
	if err != nil {
		return nil, err
	}

	img, err := decoder(bytes.NewReader(content))

	imgBounds := img.Bounds()
	if imgBounds.Dx() != resolution.Width || imgBounds.Dy() != resolution.Height {
		return nil, fmt.Errorf("expected resolution: %s, got: %s", resolution, imgBounds.Size())
	}

	title := gjson.GetBytes(jsonRaw, "images.0.title").String()
	copyright := gjson.GetBytes(jsonRaw, "images.0.copyright").String()
	description := title + ", " + copyright

	var translated string
	if region.IsAny(types.NonEnglishRegions...) && cfg.useGoogleTranslateService && cfg.googleAppCredentials != "" {
		logger.InfoLogger.Println("Using Google Cloud Translation Service for description translation from", region.String(), "to", types.UnitedStates.String())
		translated, err = translateDescription(description, region.String(), types.UnitedStates.String())
		if err != nil {
			logger.ErrLogger.Printf("failed to translate description: %v\n", err)
		}
	}

	if region == types.Japan {
		var fn func(string) (string, error)
		if cfg.furiganaApiAppId != "" {
			logger.InfoLogger.Println("Using Goo Labs API for Furigana conversion")
			fn = furiganizeGooLabsApi
		} else {
			logger.InfoLogger.Println("Using go-kakasi for Furigana conversion")
			fn = furiganizeKakasi
		}

		annotated, err := fn(description)
		if err != nil {
			logger.ErrLogger.Printf("failed to annotate description: %v\n", err)
		} else {
			description = annotated
		}
	}

	lines := []string{description}
	if translated != "" {
		lines = append(lines, translated)
	}

	var audio *Audio
	if cfg.useGoogleText2SpeechService {
		logger.InfoLogger.Println("Using Google Cloud Text-to-Speech Service for audio generation")
		audio, err = speakDescription(title+", "+copyright, types.Map[types.Region, types.Region]{
			types.Canada:     types.UnitedStates,
			types.NewZealand: types.UnitedKingdom,
		}.Get(region, region).String())
		if err != nil {
			logger.ErrLogger.Printf("failed to generate audio stream: %v\n", err)
		}
	}

	return &Image{
		Audio:       audio,
		Description: strings.Join(lines, "\n"),
		Image:       img,
		DownloadURL: parsedRequestUri.String(),
		SearchURL:   gjson.GetBytes(jsonRaw, "images.0.copyrightlink").String(),
	}, err
}

func WithGoogleAppCredentials(credentials string) crawlerConfigOption {
	return func(cfg *crawlerConfig) {
		cfg.googleAppCredentials = credentials
	}
}

func WithFuriganaApiAppId(appId string) crawlerConfigOption {
	return func(cfg *crawlerConfig) {
		cfg.furiganaApiAppId = appId
	}
}

func WithUseGoogleText2SpeechService(use bool) crawlerConfigOption {
	return func(cfg *crawlerConfig) {
		cfg.useGoogleText2SpeechService = use
	}
}

func WithUseGoogleTranslateService(use bool) crawlerConfigOption {
	return func(cfg *crawlerConfig) {
		cfg.useGoogleTranslateService = use
	}
}
