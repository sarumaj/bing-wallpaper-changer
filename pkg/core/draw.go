package core

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// DrawDescription draws a title onto the given image.
func (img *Image) DrawDescription(position types.Position, fontName string) error {
	imgBounds := img.Bounds()

	// create a new image with the same dimensions as the original.
	ctx := gg.NewContextForRGBA(image.NewRGBA(imgBounds))

	// copy the original image onto the new image.
	ctx.DrawImage(img.Image, 0, 0)

	fontDataReader, ok := extras.EmbeddedFonts[fontName]
	if !ok {
		return fmt.Errorf("unknown font: %s", fontName)
	}
	defer fontDataReader.Close()

	data, err := io.ReadAll(fontDataReader)
	if err != nil {
		return err
	}

	// parse font
	parsed, err := opentype.Parse(data)
	if err != nil {
		return fmt.Errorf("error parsing font: %v", err)
	}

	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{Size: 20, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		return fmt.Errorf("error creating font face: %v", err)
	}

	// measure text bounding box
	ctx.SetFontFace(face)

	noLines := 1.0
	text, lineSpacing := img.Description, 1.2
	textWidth, textHeight := ctx.MeasureString(text)

	if maxWidth := 0.6 * float64(imgBounds.Dx()); textWidth > maxWidth {
		lines := ctx.WordWrap(text, maxWidth)
		noLines = float64(len(lines))
		text = strings.Join(lines, "\n")
		textWidth, textHeight = ctx.MeasureMultilineString(text, lineSpacing)
	}

	lineHeight := textHeight / noLines
	textHeight = lineHeight*3/2 + lineHeight*(noLines-1)

	y_margin, r := 50.0, textHeight/5
	if r > y_margin {
		r = y_margin
	}

	var x, y, w, h float64
	switch position {
	case types.TopCenter:
		x, y, w, h = float64(imgBounds.Dx())/2-textWidth/2-r, y_margin+r, textWidth+r*2, textHeight+2*r

	case types.BottomCenter:
		x, y, w, h = float64(imgBounds.Dx())/2-textWidth/2-r, float64(imgBounds.Dy())-textHeight*3/2-y_margin-r, textWidth+r*2, textHeight+2*r

	default:
		return fmt.Errorf("unsupported position: %s, expected any of: %s", position, types.Positions{types.TopCenter, types.BottomCenter})

	}

	// draw outline of the text box with rounded corners
	ctx.SetColor(color.White)
	ctx.SetLineWidth(5)
	ctx.DrawRoundedRectangle(x, y, w, h, r)
	ctx.Stroke()

	// fill the text box with a semi-transparent black color (opacity of 64%)
	ctx.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 164})
	ctx.DrawRoundedRectangle(x, y, w, h, r)
	ctx.Fill()

	// draw the text
	ctx.SetColor(color.White)
	ctx.DrawStringWrapped(text, x+r, y+r, 0.0, 0.0, w-2*r, lineSpacing, gg.AlignCenter)

	img.Image = ctx.Image()
	return nil
}

// DrawQRCode draws a QR code onto the given image.
func (img *Image) DrawQRCode(resolution types.Resolution, position types.Position) error {
	var size int
	switch resolution {
	case types.LowDefinition:
		size = 128

	case types.HighDefinition:
		size = 164

	case types.UltraHighDefinition:
		size = 192

	default:
		return fmt.Errorf("unsupported resolution: %s, expected any of: %s", resolution, types.AllowedResolutions)

	}

	coder, err := qrcode.New(img.SearchURL, qrcode.Medium)
	if err != nil {
		return err
	}

	imgBounds := img.Bounds()
	ctx := gg.NewContextForRGBA(image.NewRGBA(imgBounds))

	// copy the original image onto the new image.
	ctx.DrawImage(img.Image, 0, 0)

	// generate QR code.
	x_offset, y_offset, qrCodeImg := 50, 50, coder.Image(size)

	// blur edges of QR code image
	qrCodeImgTransparent := image.NewRGBA(image.Rect(0, 0, size, size))
	center := float64(size) / 2                     // max distance from a pixel at the border to the center
	margin := 0.95                                  // margin, at which blur transition begins
	smooth := (margin - 1) * center / math.Log(0.1) // level of blur when the distance between the margin and center is maximal (1 - 0.1/255 = 96%)

	for x := 0.0; x < float64(size); x++ {
		for y := 0.0; y < float64(size); y++ {
			r, g, b, _ := qrCodeImg.At(int(x), int(y)).RGBA()
			a := uint32(196) // make image semi-transparent per default

			// calculate distance difference and calculate blur level (transparency) besides margin
			if d := math.Max(math.Abs(x-center), math.Abs(y-center)); d > margin*center {
				a = uint32(float64(a) * math.Exp((margin*center-d)/smooth))
			}

			// apply
			qrCodeImgTransparent.Set(int(x), int(y), color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	qrCodeImg = qrCodeImgTransparent

	x, y := 0, 0
	switch position {
	case types.TopLeft:
		x, y = x_offset, y_offset

	case types.TopRight:
		x, y = imgBounds.Dx()-size-x_offset, y_offset

	case types.BottomLeft:
		x, y = x_offset, imgBounds.Dy()-size-y_offset

	case types.BottomRight:
		x, y = imgBounds.Dx()-size-x_offset, imgBounds.Dy()-size-y_offset

	default:
		return fmt.Errorf("unsupported position: %s, expected any of: %s", position, types.Positions{types.TopLeft, types.TopRight, types.BottomLeft, types.BottomRight})

	}

	// draw QR code image.
	ctx.DrawImage(qrCodeImg, x, y)

	img.Image = ctx.Image()
	return nil
}

// DrawWatermark draws a watermark onto the given image.
func (img *Image) DrawWatermark(watermarkFile string, rotateCounterClockwise bool) error {
	var source io.ReadCloser
	var err error
	if r, ok := extras.EmbeddedWatermarks[watermarkFile]; ok {
		source = r

	} else {
		source, err = os.OpenFile(watermarkFile, os.O_RDONLY, os.ModePerm)

	}

	if err != nil {
		return err
	}
	defer source.Close()

	decode, err := getDecoder(watermarkFile)
	if err != nil {
		return err
	}

	watermark, err := decode(source)
	if err != nil {
		return err
	}

	watermarkBounds := watermark.Bounds()
	if watermarkBounds.Dx() < watermarkBounds.Dy() {
		// rotate the image 90 degrees clockwise or counter-clockwise
		rotated := image.NewRGBA(image.Rect(0, 0, watermarkBounds.Dy(), watermarkBounds.Dx()))
		for y := watermarkBounds.Min.Y; y < watermarkBounds.Max.Y; y++ {
			for x := watermarkBounds.Min.X; x < watermarkBounds.Max.X; x++ {
				// set each pixel to the corresponding pixel in the original image
				if rotateCounterClockwise {
					rotated.Set(y, watermarkBounds.Bounds().Max.X-x-1, watermark.At(x, y))
				} else {
					rotated.Set(watermarkBounds.Bounds().Max.Y-y-1, x, watermark.At(x, y))
				}
			}
		}

		watermark = rotated
		watermarkBounds = rotated.Bounds()
	}

	imgBounds := img.Bounds()

	// resize watermark to fit the wallpaper dimensions
	resized := image.NewRGBA(image.Rect(0, 0, imgBounds.Dx(), imgBounds.Dy()))
	draw.CatmullRom.Scale(resized, resized.Rect, watermark, watermarkBounds, draw.Over, nil)

	ctx := gg.NewContextForRGBA(image.NewRGBA(imgBounds))

	// copy the original image onto the new image
	ctx.DrawImage(img.Image, 0, 0)

	// draw the watermark
	ctx.DrawImage(resized, 0, 0)

	img.Image = ctx.Image()
	return nil
}
