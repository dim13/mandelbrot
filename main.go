package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"

	"github.com/dim13/colormap"
)

func mandelbrot(img *image.Paletted, p image.Point) {
	var z complex128
	c := scale(img.Bounds(), p)
	for n := 0; n < len(img.Palette); n++ {
		if z = z*z + c; cmplx.IsInf(z) {
			img.SetColorIndex(p.X, p.Y, uint8(n))
			return
		}
	}
}

func scale(r image.Rectangle, p image.Point) complex128 {
	zoom := float64(r.Max.X) / 3
	re := float64(p.X) - float64(r.Max.X)/2 - float64(r.Max.X)/5
	im := float64(p.Y) - float64(r.Max.Y)/2
	return complex(re/zoom, im/zoom)
}

func generate(r image.Rectangle, p color.Palette) image.Image {
	img := image.NewPaletted(r, p)
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			mandelbrot(img, image.Pt(x, y))
		}
	}
	return img
}

func main() {
	fd, err := os.Create("mandelbrot.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	img := generate(image.Rect(0, 0, 800, 600), colormap.Magma)
	if err := png.Encode(fd, img); err != nil {
		log.Println(err)
	}
}
