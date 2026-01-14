package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"slices"

	"github.com/dim13/colormap"
)

func mandelbrot(img *image.Paletted, p image.Point) {
	var z complex128
	c := scale(img.Bounds(), p)
	for n := 0; n < len(img.Palette); n++ {
		if z = cmplx.Pow(z, 2) + c; cmplx.IsInf(z) {
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

var palette = map[string]color.Palette{
	"inferno": colormap.Inferno,
	"magma":   colormap.Magma,
	"plasma":  colormap.Plasma,
	"viridis": colormap.Viridis,
	"parula":  colormap.Parula,
}

func main() {
	fname := flag.String("f", "mandelbrot.png", "file name")
	width := flag.Int("w", 800, "width")
	height := flag.Int("h", 600, "height")
	pal := flag.String("p", "magma", "palette (inferno, magma, plasma, viridis)")
	inverse := flag.Bool("i", false, "inverse palette")
	flag.Parse()

	fd, err := os.Create(*fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	p, ok := palette[*pal]
	if !ok {
		log.Printf("no such palette %s, fallback to magma", *pal)
		p = colormap.Magma
	}
	if *inverse {
		slices.Reverse(p)
	}

	img := generate(image.Rect(0, 0, *width, *height), p)
	if err := png.Encode(fd, img); err != nil {
		log.Println(err)
	}
}
