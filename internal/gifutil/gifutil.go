/*
Package gifutil provides utility functions to manipulate GIF image.
*/
package gifutil

import (
	"image"
	"image/color"
	"image/gif"
	"iter"
	"os"
)

func LoadFile(name string) (*gif.GIF, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	return g, err
}

func duplicatePaletted(src *image.Paletted) *image.Paletted {
	pix := make([]uint8, len(src.Pix))
	copy(pix, src.Pix)
	palette := make(color.Palette, len(src.Palette))
	copy(palette, src.Palette)
	return &image.Paletted{
		Pix:     pix,
		Stride:  src.Stride,
		Rect:    src.Rect,
		Palette: palette,
	}
}

func drawOver(dst, src *image.Paletted) {
	for y := src.Rect.Min.Y; y < src.Rect.Max.Y; y++ {
		for x := src.Rect.Min.X; x < src.Rect.Max.X; x++ {
			index := src.ColorIndexAt(x, y)
			if c := src.Palette[index].(color.RGBA); c.A == 0 {
				continue
			}
			dst.SetColorIndex(x, y, index)
		}
	}
}

// IterateComposed iterates composed frames.
func IterateComposed(g *gif.GIF) iter.Seq2[int, *image.Paletted] {
	rect := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	last := image.NewPaletted(rect, g.Config.ColorModel.(color.Palette))
	return func(yield func(int, *image.Paletted) bool) {
		for i, src := range g.Image {
			var curr *image.Paletted
			// composed accumulated image
			switch g.Disposal[i] {
			case gif.DisposalNone:
				drawOver(last, src)
				curr = duplicatePaletted(last)
			case gif.DisposalBackground:
				curr = src
			case gif.DisposalPrevious:
				// FIXME: support DisposalPrevious
				curr = src
			default:
				curr = src
			}
			if !yield(i, curr) {
				break
			}
		}
	}
}
