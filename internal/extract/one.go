package extract

import (
	"context"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"iter"
	"log"

	"github.com/koron-go/subcmd"
)

func ExtractOne(ctx context.Context, args []string) error {
	fs := subcmd.FlagSet(ctx)
	fs.Parse(args)

	if fs.NArg() != 1 {
		return errors.New("only one GIF can be extracted at a time")
	}
	input := fs.Arg(0)

	output := appendFilename(input, "_one.png")
	return extractRepresentativeOne(output, input)
}

type frameInfo struct {
	i       int
	img     image.Image
	entropy float64
}

func extractRepresentativeOne(output, input string) error {
	g, err := loadGIF(input)
	if err != nil {
		return err
	}

	highest := frameInfo{i: -1, entropy: -1}
	for i, img := range iterateComposedRGBA(g) {
		entropy := measureImageEntropy(img)
		if entropy > highest.entropy {
			highest = frameInfo{i: i, img: img, entropy: entropy}
		}
	}
	log.Printf("extracted #%d from %s: entropy=%f", highest.i, input, highest.entropy)
	if highest.img != nil && output != "" {
		writeImage(output, highest.img)
	}

	return nil
}

func iterateComposedPaletted(g *gif.GIF) iter.Seq[*image.Paletted] {
	rect := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	buf := image.NewPaletted(rect, g.Config.ColorModel.(color.Palette))
	return func(yield func(*image.Paletted) bool) {
		for i, src := range g.Image {
			// composed accumulated image
			op := draw.Over
			switch g.Disposal[i] {
			case 0:
				op = draw.Src
				buf.Palette = src.Palette
			}
			draw.Draw(buf, src.Rect, src, src.Rect.Min, op)
			if !yield(buf) {
				break
			}
		}
	}
}

func duplicatePaletted(src *image.Paletted) *image.Paletted {
	dup := image.NewPaletted(src.Rect, src.Palette)
	draw.Draw(dup, src.Rect, src, src.Rect.Min, draw.Over)
	return dup
}

func iterateComposedRGBA(g *gif.GIF) iter.Seq2[int, *image.RGBA] {
	rect := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	buf := image.NewRGBA(rect)
	return func(yield func(int, *image.RGBA) bool) {
		for i, src := range g.Image {
			// composed accumulated image
			op := draw.Over
			switch g.Disposal[i] {
			case 0:
				op = draw.Src
			}
			draw.Draw(buf, src.Rect, src, src.Rect.Min, op)
			if !yield(i, duplicateRGBA(buf)) {
				break
			}
		}
	}
}

func duplicateRGBA(src *image.RGBA) *image.RGBA {
	dup := image.NewRGBA(src.Rect)
	draw.Draw(dup, src.Rect, src, src.Rect.Min, draw.Over)
	return dup
}
