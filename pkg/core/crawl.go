package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/sarumaj/bing-wallpaper-changer/pkg/logger"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/tidwall/gjson"
)

// config contains the configuration for the Bing and Goo Labs APIs.
var config = struct {
	AppID       string
	Bing        url.URL
	HiraganaAPI url.URL
}{
	AppID:       "d175fb007ba9ef9931b665291ad75192c62b667360cd6d9ed8e0ef524a2aa442",
	Bing:        url.URL{Scheme: "https", Host: "www.bing.com"},
	HiraganaAPI: url.URL{Scheme: "https", Host: "labs.goo.ne.jp"},
}

// annotateDescription annotates the description in Japanese with Furigana for Kanji sequences.
// It uses the Goo Labs API to convert Kanji to Hiragana.
// It should be replaced through kakasi NLP library in the future (binding similar to https://github.com/Theta-Dev/kakasi).
func annotateDescription(description string) (string, error) {
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

	// request hiragana conversion.
	raw, err := readResponse(http.PostForm(
		(&url.URL{
			Scheme: config.HiraganaAPI.Scheme,
			Host:   config.HiraganaAPI.Host,
			Path:   "/api/hiragana",
		}).String(),
		url.Values{
			"app_id":      {config.AppID},
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

// DownloadAndDecode fetches the Bing wallpaper and decodes it.
func DownloadAndDecode(day types.Day, region types.Region, resolution types.Resolution) (*Image, error) {
	jsonRaw, err := readResponse(http.Get(
		(&url.URL{
			Scheme: config.Bing.Scheme,
			Host:   config.Bing.Host,
			Path:   "/HPImageArchive.aspx",
			RawQuery: url.Values{
				"format": {"js"},
				"idx":    {fmt.Sprintf("%d", day)},
				"n":      {"1"},
				"mkt":    {region.String()},
			}.Encode(),
		}).String(),
	))
	if err != nil {
		return nil, err
	}

	path := gjson.GetBytes(jsonRaw, "images.0.url").String()
	path = regexp.MustCompile(`_\d+x\d+`).ReplaceAllString(path, "_"+resolution.String())

	parsed, err := url.ParseRequestURI(path)
	if err != nil {
		return nil, err
	}

	parsed.Scheme = config.Bing.Scheme
	parsed.Host = config.Bing.Host

	decoder, err := getDecoder(parsed.Query().Get("id"))
	if err != nil {
		return nil, err
	}

	content, err := readResponse(http.Get(parsed.String()))
	if err != nil {
		return nil, err
	}

	img, err := decoder(bytes.NewReader(content))

	imgBounds := img.Bounds()
	if imgBounds.Dx() != resolution.Width || imgBounds.Dy() != resolution.Height {
		return nil, fmt.Errorf("expected resolution: %s, got: %s", resolution, imgBounds.Size())
	}

	description := fmt.Sprintf(
		"%s, %s",
		gjson.GetBytes(jsonRaw, "images.0.title").String(),
		gjson.GetBytes(jsonRaw, "images.0.copyright").String(),
	)
	if region == types.Japan {
		annotated, err := annotateDescription(description)
		if err != nil {
			logger.ErrLogger.Printf("failed to annotate description: %v\n", err)
		} else {
			description = annotated
		}

	}

	return &Image{
		Description: description,
		Image:       img,
		DownloadURL: parsed.String(),
		SearchURL:   gjson.GetBytes(jsonRaw, "images.0.copyrightlink").String(),
	}, err
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
