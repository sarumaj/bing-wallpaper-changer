package core

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/extras"
	"github.com/sarumaj/bing-wallpaper-changer/pkg/types"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// DrawDescription draws a title onto the given image.
func (img *Image) DrawDescription(position types.Position) error {
	imgBounds := img.Bounds()

	// create a new image with the same dimensions as the original.
	ctx := gg.NewContextForRGBA(image.NewRGBA(imgBounds))

	// copy the original image onto the new image.
	ctx.DrawImage(img.Image, 0, 0)

	// parse font
	parsed, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{Size: 18, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		return err
	}

	// measure text bounding box
	ctx.SetFontFace(face)
	text, lineSpacing := img.Description, 1.2
	textWidth, textHeight := ctx.MeasureString(text)
	if textWidth > float64(imgBounds.Dx()) {
		text = strings.Join(ctx.WordWrap(text, float64(imgBounds.Dx())), "\n")
		textWidth, textHeight = ctx.MeasureMultilineString(text, lineSpacing)
	}

	x_margin, y_margin, r := 10.0, 10.0, textHeight*lineSpacing/5
	var x, y, w, h float64
	switch position {
	case types.TopCenter:
		x, y, w, h = float64(imgBounds.Dx())/2-textWidth/2-x_margin, y_margin, textWidth+x_margin*2, y_margin+textHeight

	case types.BottomCenter:
		x, y, w, h = float64(imgBounds.Dx())/2-textWidth/2-x_margin, float64(imgBounds.Dy())-textHeight-y_margin*2, textWidth+x_margin*2, y_margin+textHeight

	default:
		return fmt.Errorf("unsupported position: %s, expected any of: %s", position, types.Positions{types.TopCenter, types.BottomCenter})

	}

	// draw outline of the text box with rounded corners
	ctx.SetColor(color.White)
	ctx.SetLineWidth(5)
	ctx.DrawRoundedRectangle(x, y, w, h, r)
	ctx.Stroke()

	// fill the text box with a semi-transparent black color (opacity of 50%)
	ctx.SetColor(color.RGBA{R: 0, G: 0, B: 0, A: 128})
	ctx.DrawRoundedRectangle(x, y, w, h, r)
	ctx.Fill()

	// draw the text
	ctx.SetColor(color.White)
	//ctx.DrawString(img.Description, x, y)
	ctx.DrawStringWrapped(text, x+x_margin, y, 0.0, 0.0, float64(imgBounds.Dx()), lineSpacing, gg.AlignLeft)

	img.Image = ctx.Image()
	return nil
}

// DrawQRCode draws a QR code onto the given image.
func (img *Image) DrawQRCode(size int, position types.Position) error {
	x_offset, y_offset := 10, 10
	coder, err := qrcode.New(img.SearchURL, qrcode.Medium)
	if err != nil {
		return err
	}

	imgBounds := img.Bounds()
	ctx := gg.NewContextForRGBA(image.NewRGBA(imgBounds))
	ctx.DrawImage(img.Image, 0, 0)

	switch position {
	case types.TopLeft:
		ctx.DrawImage(coder.Image(size), x_offset, y_offset)

	case types.TopRight:
		ctx.DrawImage(coder.Image(size), imgBounds.Dx()-size-x_offset, y_offset)

	case types.BottomLeft:
		ctx.DrawImage(coder.Image(size), x_offset, imgBounds.Dy()-size-y_offset)

	case types.BottomRight:
		ctx.DrawImage(coder.Image(size), imgBounds.Dx()-size-x_offset, imgBounds.Dy()-size-y_offset)

	default:
		return fmt.Errorf("unsupported position: %s, expected any of: %s", position, types.Positions{types.TopLeft, types.TopRight, types.BottomLeft, types.BottomRight})

	}

	img.Image = ctx.Image()
	return nil
}

// DrawWatermark draws a watermark onto the given image.
func (img *Image) DrawWatermark(watermarkFile string) error {
	var source io.Reader
	var err error
	if r, ok := extras.RegisteredWatermarks[watermarkFile]; ok {
		source = r

	} else {
		source, err = os.OpenFile(watermarkFile, os.O_RDONLY, os.ModePerm)

	}

	if err != nil {
		return err
	}

	decode, err := getDecoder(watermarkFile)
	if err != nil {
		return err
	}

	watermark, err := decode(source)
	if err != nil {
		return err
	}

	// resize watermark to fit the wallpaper dimensions
	resized := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.CatmullRom.Scale(resized, resized.Rect, watermark, watermark.Bounds(), draw.Over, nil)

	ctx := gg.NewContextForRGBA(image.NewRGBA(img.Bounds()))

	// copy the original image onto the new image
	ctx.DrawImage(img.Image, 0, 0)

	// draw the watermark
	ctx.DrawImage(resized, 0, 0)

	img.Image = ctx.Image()
	return nil
}
