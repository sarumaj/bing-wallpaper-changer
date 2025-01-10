package core

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/tidwall/gjson"
)

//go:embed testdata/*
var embeddedTestData embed.FS

var testData, _ = fs.Sub(embeddedTestData, "testdata")

// getHandler returns an http.HandlerFunc that serves the file with the given filename.
func getHandler(t testing.TB, filename string) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		if ct := mime.TypeByExtension(filepath.Ext(filename)); ct != "" {
			w.Header().Set("Content-Type", ct)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		w.WriteHeader(http.StatusOK)

		f, err := testData.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// FromMock returns true if the test is running in mock mode.
func FromMock(t testing.TB) bool {
	t.Helper()

	switch os.Getenv("TEST_MOCK") {

	case "true", "True", "TRUE", "1", "y", "yes", "YES":
		return true

	}

	return false
}

// MockServers sets up mock servers for the Bing and Hiragana API.
func MockServers(t testing.TB) {
	t.Helper()

	serveMux := http.NewServeMux()
	serveMux.Handle("/HPImageArchive.aspx", getHandler(t, "bing.json"))
	serveMux.Handle("/th", getHandler(t, "bing.jpg"))
	serveMux.Handle("/api/hiragana", getHandler(t, "hiragana.json"))
	serveMux.Handle("/search/", getHandler(t, "jisho.html"))
	server := httptest.NewServer(serveMux)

	backupCfg := crawlerConfig{
		bingUrl:          cfg.bingUrl,
		furiganaApiUrl:   cfg.furiganaApiUrl,
		furiganaApiAppId: cfg.furiganaApiAppId,
		jishoOrgUrl:      cfg.jishoOrgUrl,
	}

	cfg.bingUrl = server.URL
	cfg.furiganaApiUrl = server.URL
	cfg.furiganaApiAppId = "test"
	cfg.jishoOrgUrl = server.URL

	t.Cleanup(func() {
		server.Close()

		cfg.bingUrl = backupCfg.bingUrl
		cfg.furiganaApiUrl = backupCfg.furiganaApiUrl
		cfg.furiganaApiAppId = backupCfg.furiganaApiAppId
		cfg.jishoOrgUrl = backupCfg.jishoOrgUrl
	})
}

// SetupTestImage returns a test image
func SetupTestImage(t testing.TB) *Image {
	t.Helper()

	f, err := testData.Open("bing.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	decoder, err := getDecoder("bing.jpg")
	if err != nil {
		t.Fatal(err)
	}

	img, err := decoder(f)
	if err != nil {
		t.Fatal(err)
	}

	f, err = testData.Open("bing.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	jsonRaw, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	parsedRequestUri, err := url.ParseRequestURI(gjson.GetBytes(jsonRaw, "images.0.url").String())
	if err != nil {
		t.Fatal(err)
	}

	remoteHostUrl, err := url.Parse(cfg.bingUrl)
	if err != nil {
		t.Fatal(err)
	}

	parsedRequestUri.Host = remoteHostUrl.Host
	parsedRequestUri.Scheme = remoteHostUrl.Scheme

	return &Image{
		Description: fmt.Sprintf(
			"%s, %s",
			gjson.GetBytes(jsonRaw, "images.0.title").String(),
			gjson.GetBytes(jsonRaw, "images.0.copyright").String(),
		),
		DownloadURL: parsedRequestUri.String(),
		SearchURL:   gjson.GetBytes(jsonRaw, "images.0.copyrightlink").String(),
		Image:       img,
	}
}
