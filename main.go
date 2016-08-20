package main

import (
	"github.com/ftrvxmtrx/tga"
	"image"
	"image/color"
	"log"
	"os"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
    RGBAWhite := color.RGBA{0xff, 0xff, 0xff, 0xff}
    RGBARed   := color.RGBA{0xff, 0x88, 0x88, 0xff}

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, RGBAWhite)
		}
	}
    img.Set(30, 40, RGBARed)

	file, err := os.Create("simple.tga")
	if err != nil {
		log.Fatal("err")
	}
	defer file.Close()

	tga.Encode(file, img)

	log.Print("Finish!")
}
