package main

import (
	"github.com/ftrvxmtrx/tga"
	"image"
	"image/color"
	"log"
	"math"
	"os"
)

const (
	ImageOutFile = "output.tga"
)

func ImageFlipVertical(img *image.RGBA) {
	h := img.Rect.Dy()

	for y := img.Rect.Min.Y; y < img.Rect.Min.Y+h/2; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			tmp := img.At(x, img.Rect.Max.Y-y-1)
			img.Set(x, img.Rect.Max.Y-y-1, img.At(x, y))
			img.Set(x, y, tmp)
		}
	}
}

func ImageDrawLine(img *image.RGBA, x0, y0, x1, y1 float64, c color.RGBA) {
	it, dx, dy := 0., 0., 0.
	if math.Abs(x1-x0) > math.Abs(y1-y0) {
		it = x1 - x0
		dx = it / math.Abs(it)
		dy = (y1 - y0) / it
	} else {
		it = y1 - y0
		dx = (x1 - x0) / it
		dy = it / math.Abs(it)
	}

	for i := 0.; i < math.Abs(it); i++ {
		x := int(x0 + i*dx)
		y := int(y0 + i*dy)
		img.Set(x, y, c)
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	RGBAWhite := color.RGBA{0xff, 0xff, 0xff, 0xff}
	RGBARed := color.RGBA{0xff, 0x88, 0x88, 0xff}

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, RGBAWhite)
		}
	}
	img.Set(30, 40, RGBARed)

	ImageFlipVertical(img)
	ImageDrawLine(img, 0., 0., 80., 100., RGBARed)
	ImageDrawLine(img, 0., 0., 100., 80., RGBARed)

	file, err := os.Create(ImageOutFile)
	if err != nil {
		log.Fatal("err")
	}
	defer file.Close()

	tga.Encode(file, img)

	log.Print("Finish!")
}
