package main

import (
	"github.com/ftrvxmtrx/tga"
	"image"
	"image/color"
	"log"
	"os"
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

	file, err := os.Create("simple.tga")
	if err != nil {
		log.Fatal("err")
	}
	defer file.Close()

	tga.Encode(file, img)

	log.Print("Finish!")
}
