/*
Package gifutil provides utility functions to manipulate GIF image.
*/
package gifutil

import (
	"image"
	"image/color"
	"image/draw"
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
	dup := image.NewPaletted(src.Rect, src.Palette)
	draw.Draw(dup, src.Rect, src, src.Rect.Min, draw.Over)
	return dup
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
				draw.Over.Draw(last, src.Rect, src, src.Rect.Min)
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
