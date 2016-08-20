package main

import (
	"bufio"
	"github.com/ftrvxmtrx/tga"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	ImageOutFile = "output.tga"
)

type Point struct {
	X, Y, Z float64
}

type FaceVertexIndices struct {
	A, B, C int
}

type Model struct {
	Vertices []Point
	Faces    []FaceVertexIndices
}

func (m *Model) ReadObj(name string) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

    // Add a single point to make m.Vertices one-indexed
    m.Vertices = []Point{Point{}}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line, " ")
		if s[0] == "v" {
			p := Point{}
			p.X, _ = strconv.ParseFloat(s[1], 64)
			p.Y, _ = strconv.ParseFloat(s[2], 64)
			p.Z, _ = strconv.ParseFloat(s[3], 64)

			m.Vertices = append(m.Vertices, p)
		} else if s[0] == "f" {
			f := FaceVertexIndices{}
			as := strings.Split(s[1], "/")
			f.A, _ = strconv.Atoi(as[0])
			bs := strings.Split(s[2], "/")
			f.B, _ = strconv.Atoi(bs[0])
			cs := strings.Split(s[3], "/")
			f.C, _ = strconv.Atoi(cs[0])

			m.Faces = append(m.Faces, f)
		}
	}
}

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

	model := Model{}
	model.ReadObj("obj/african_head.obj")

	log.Print("Finish!")
}
