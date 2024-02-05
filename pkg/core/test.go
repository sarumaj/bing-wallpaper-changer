package core

import (
	"image"
	"os"
	"testing"
)

// SetupTestImage returns a test image
func SetupTestImage() *Image {
	return &Image{
		Description: "Test Wallpaper",
		DownloadURL: "https://www.bing.com/th?id=OHR.SnowyOwl_EN-US2717317224_1920x1080.jpg&rf=LaDigue_1920x1080.jpg&pid=hp",
		SearchURL:   "https://www.bing.com/images/search?view=detailV2&ccid=LaDigue&id=OHR.SnowyOwl_EN-US2717317224_1920x1080.jpg",
		Image:       image.NewRGBA(image.Rect(0, 0, 1920, 1080)),
	}
}

// SkipOAT skips the test if TEST_OAT is not set to true
func SkipOAT(t testing.TB) {
	switch os.Getenv("TEST_OAT") {

	case "true", "True", "TRUE", "1", "y", "yes", "YES":
		return

	}

	t.Skipf("Running only FAT tests, skipping %q, since it requires extensive mock-up", t.Name())
}
