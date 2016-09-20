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
		dy = (y1 - y0) / math.Abs(it)
	} else {
		it = y1 - y0
		dx = (x1 - x0) / math.Abs(it)
		dy = it / math.Abs(it)
	}

	for i := 0.; i < math.Abs(it); i++ {
		x := int(x0 + i*dx)
		y := int(y0 + i*dy)
		img.Set(x, y, c)
	}
}

func ImageDrawTriangle(img *image.RGBA, p0, p1, p2 Point, c color.RGBA) {
    // Sort
    if p0.Y > p1.Y {
        tmp := p0
        p0 = p1
        p1 = tmp
    }
    if p0.Y > p2.Y {
        tmp := p0
        p0 = p2
        p2 = tmp
    }
    if p1.Y > p2.Y {
        tmp := p1
        p1 = p2
        p2 = tmp
    }

    // Origin
    plo := p0
    pld := p1

    pro := p0
    prd := p2

    if pld.X > prd.X {
        tmp := pld
        pld = prd
        prd = tmp
    }

    it := p2.Y - p0.Y
    for i := 0.; i < it; i++ {
        y := p0.Y + i
        if y == pld.Y {
            plo = pld
            pld = prd
        } else if y == prd.Y {
            pro = prd
            prd = pld
        }

        // Get line from x0 to x1
        x0 := (i - (plo.Y - p0.Y)) * (pld.X - plo.X) / (pld.Y - plo.Y) + plo.X
        x1 := (i - (pro.Y - p0.Y)) * (prd.X - pro.X) / (prd.Y - pro.Y) + pro.X

        if x0 > x1 {
            tmp := x0
            x0 = x1
            x1 = tmp
        }

        for x := x0; x < x1; x++ {
            img.Set(int(x), int(y), c)
        }
    }
}

func ImageRenderWireframeModel(img *image.RGBA, m *Model) {
	RGBAWhite := color.RGBA{0xff, 0xff, 0xff, 0xff}
    width := img.Rect.Dx()
    height := img.Rect.Dy()

    for f := range m.Faces {
        a := []Point{
            m.Vertices[m.Faces[f].A],
            m.Vertices[m.Faces[f].B],
            m.Vertices[m.Faces[f].C],
        }
        for i := 0; i < 3; i++ {
            v0, v1 := a[i], a[(i + 1) % 3]
            x0 := v0.X*float64(width/2) + float64(img.Rect.Min.X+width/2)
            y0 := v0.Y*float64(height/2) + float64(img.Rect.Min.Y+height/2)
            x1 := v1.X*float64(width/2) + float64(img.Rect.Min.X+width/2)
            y1 := v1.Y*float64(height/2) + float64(img.Rect.Min.Y+height/2)
            ImageDrawLine(img, x0, y0, x1, y1, RGBAWhite)
        }
    }
}

func ImageRenderTriangleModel(img *image.RGBA, m *Model) {
    width := img.Rect.Dx()
    height := img.Rect.Dy()

	RGBARed := color.RGBA{0xff, 0x88, 0x88, 0xff}
    for f := range m.Faces {
        a := []Point{
            m.Vertices[m.Faces[f].A],
            m.Vertices[m.Faces[f].B],
            m.Vertices[m.Faces[f].C],
        }
        ImageDrawTriangle(img,
            Point{
                (a[0].X)*float64(width/2) + float64(img.Rect.Min.X+width/2),
                (a[0].Y)*float64(height/2) + float64(img.Rect.Min.Y+height/2),
                0,
            }, Point{
                (a[0].X + 0.01)*float64(width/2) + float64(img.Rect.Min.X+width/2),
                (a[0].Y)*float64(height/2) + float64(img.Rect.Min.Y+height/2),
                0,
            }, Point{
                (a[0].X + 0.01)*float64(width/2) + float64(img.Rect.Min.X+width/2),
                (a[0].Y + 0.01)*float64(height/2) + float64(img.Rect.Min.Y+height/2),
                0,
            }, RGBARed)
    }
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 800, 800))
	RGBABlack := color.RGBA{0x00, 0x00, 0x00, 0xff}
	// RGBARed := color.RGBA{0xff, 0x88, 0x88, 0xff}

	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, RGBABlack)
		}
	}

	model := Model{}
	model.ReadObj("obj/african_head.obj")
    // ImageRenderWireframeModel(img, &model)
    ImageRenderTriangleModel(img, &model)
	ImageFlipVertical(img)
    // ImageDrawTriangle(img, Point{100, 180, 0}, Point{400, 150, 0}, Point{40, 400, 0}, RGBARed)

	file, err := os.Create(ImageOutFile)
	if err != nil {
		log.Fatal("err")
	}
	defer file.Close()

	tga.Encode(file, img)

	log.Print("Finish!")
}
