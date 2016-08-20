package main

import (
    "github.com/ftrvxmtrx/tga"
    "image"
    "log"
    "os"
)

func main() {
    img := image.NewRGBA(image.Rect(0, 0, 10, 10))

    file, err := os.Create("simple.tga")
    if err != nil {
        log.Fatal("err")
    }
    defer file.Close()

    tga.Encode(file, img)

    log.Print("Finish!")
}
