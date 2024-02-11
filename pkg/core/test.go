package core

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
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
		switch filepath.Ext(filename) {
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		case ".jpg":
			w.Header().Set("Content-Type", "image/jpeg")
		default:
			t.Fatalf("unknown file extension: %s", filepath.Ext(filename))
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

	bingServerMux := http.NewServeMux()
	bingServerMux.Handle("/HPImageArchive.aspx", getHandler(t, "bing.json"))
	bingServerMux.Handle("/th", getHandler(t, "bing.jpg"))

	bingServer := httptest.NewServer(bingServerMux)

	parsed, err := url.Parse(bingServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	oldBing := config.Bing
	config.Bing = url.URL{Host: parsed.Host, Scheme: parsed.Scheme}

	hiraganaApiServerMux := http.NewServeMux()
	hiraganaApiServerMux.Handle("/api/hiragana", getHandler(t, "hiragana.json"))
	hiraganaApiServer := httptest.NewServer(hiraganaApiServerMux)

	parsed, err = url.Parse(hiraganaApiServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	oldHiraganaApi := config.HiraganaAPI
	config.HiraganaAPI = url.URL{Host: parsed.Host, Scheme: parsed.Scheme}

	t.Cleanup(func() {
		bingServer.Close()
		hiraganaApiServer.Close()
		config.Bing = oldBing
		config.HiraganaAPI = oldHiraganaApi
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

	parsed, err := url.ParseRequestURI(gjson.GetBytes(jsonRaw, "images.0.url").String())
	if err != nil {
		t.Fatal(err)
	}

	parsed.Scheme = config.Bing.Scheme
	parsed.Host = config.Bing.Host

	return &Image{
		Description: fmt.Sprintf(
			"%s, %s",
			gjson.GetBytes(jsonRaw, "images.0.title").String(),
			gjson.GetBytes(jsonRaw, "images.0.copyright").String(),
		),
		DownloadURL: parsed.String(),
		SearchURL:   gjson.GetBytes(jsonRaw, "images.0.copyrightlink").String(),
		Image:       img,
	}
}
