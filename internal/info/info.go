// Package info provides "info" sub command: to show GIF image information.
package info

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"image/gif"
	"io"
	"os"

	"github.com/koron-go/subcmd"
)

var Info = subcmd.DefineCommand("info", "show GIF info", func(ctx context.Context, args []string) error {
	var (
		palette = false
	)
	fs := subcmd.FlagSet(ctx)
	fs.BoolVar(&palette, "palette", false, "check all pallets are same")
	fs.Parse(args)
	if fs.NArg() == 0 {
		return errors.New("require one or more files")
	}
	out := os.Stdout
	for _, name := range fs.Args() {
		g, err := dumpAsGIF(out, name)
		if err != nil {
			return err
		}
		if palette {
			if err := checkPalette(out, g); err != nil {
				return err
			}
		}
	}
	return nil
})

func checkPalette(w io.Writer, g *gif.GIF) error {
	globalPalette, ok := g.Config.ColorModel.(color.Palette)
	if !ok {
		return fmt.Errorf("color mode is not palette: actual type is %T", g.Config.ColorModel)
	}
	for i, img := range g.Image {
		if err := comparePalette(globalPalette, img.Palette); err != nil {
			return fmt.Errorf("different palette at #%d: %s", i+1, err)
		}
	}
	fmt.Fprintf(w, "all palettes matched")
	return nil
}

func comparePalette(a, b color.Palette) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch: len(a)=%d len(b)=%d", len(a), len(b))
	}
	for i := range len(a) {
		ca := toRGBA(a[i])
		cb := toRGBA(b[i])
		if ca != cb {
			return fmt.Errorf("color #%[1]d mismatch: a[%[1]d]=%+[2]v b[%[1]d]=%+[3]v", i, ca, cb)
		}
	}
	return nil
}

func toRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	r = colorClip(r)
	g = colorClip(g)
	b = colorClip(b)
	a = colorClip(a)
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

func colorClip(v uint32) uint32 {
	return min(0xff, v/0x0101)
}

// dumpAsGIF dumps GIF information to the writer.
func dumpAsGIF(w io.Writer, name string) (*gif.GIF, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(w, "file: %s\n", name)
	return g, dumpGIF(w, g)
}

func dumpGIF(w io.Writer, g *gif.GIF) error {
	fmt.Fprintf(w, "len(Image)=%d\n", len(g.Image))
	for i, img := range g.Image {
		fmt.Fprintf(w, "  #%-2d Stride=%d Rect=%+v len(Palette)=%d\n", i, img.Stride, img.Rect, len(img.Palette))
	}
	fmt.Fprintf(w, "len(Delay)=%d\n", len(g.Delay))
	fmt.Fprintf(w, "  Delay=%+v\n", g.Delay)
	fmt.Fprintf(w, "LoopCount=%d\n", g.LoopCount)
	fmt.Fprintf(w, "len(Disposal)=%d\n", len(g.Disposal))
	fmt.Fprintf(w, "  Disposal=%+v\n", g.Disposal)
	//pp.Fprintf(w, "Config=%s\n", g.Config)
	fmt.Fprintf(w, "ColorModel=%T\n", g.Config.ColorModel)
	var palette color.Palette
	if p, ok := g.Config.ColorModel.(color.Palette); ok {
		fmt.Fprintf(w, "  len(Palette)=%d\n", len(p))
		palette = p
	}
	fmt.Fprintf(w, "Dimension=(%d, %d)\n", g.Config.Width, g.Config.Height)
	fmt.Fprintf(w, "BackgroundIndex=%d\n", g.BackgroundIndex)
	if palette != nil {
		fmt.Fprintf(w, "  %+v\n", palette[g.BackgroundIndex])
	}
	fmt.Fprintln(w)
	return nil
}
