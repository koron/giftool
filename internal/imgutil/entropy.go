package imgutil

import (
	"image"
	"image/color"
	"math"
)

func calcHist(img image.Image) map[color.RGBA]int {
	hist := map[color.RGBA]int{}
	rect := img.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			c := color.RGBA{
				R: uint8(r / 0x0101),
				G: uint8(g / 0x0101),
				B: uint8(b / 0x0101),
				A: uint8(a / 0x0101),
			}
			hist[c]++
		}
	}
	return hist
}

func histToEntropy[T comparable](hist map[T]int, sum int) float64 {
	var entropy float64
	for _, n := range hist {
		p := float64(n) / float64(sum)
		entropy += -p * math.Log2(p)
	}
	return entropy
}

func ToGray(img image.Image) *image.Gray {
	r := img.Bounds()
	gray := image.NewGray(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			c := img.At(x, y)
			g := color.GrayModel.Convert(c)
			gray.Set(x, y, g)
		}
	}
	return gray
}

func MeasureEntropy(img image.Image) float64 {
	b := img.Bounds()
	return histToEntropy(calcHist(img), b.Dx()*b.Dy())
}
