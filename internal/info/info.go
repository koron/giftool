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
	if len(args) == 0 {
		return errors.New("require one or more files")
	}
	for _, name := range args {
		err := dumpAsGIF(os.Stdout, name)
		if err != nil {
			return err
		}
	}
	return nil
})

// dumpAsGIF dumps GIF information to the writer.
func dumpAsGIF(w io.Writer, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "file: %s\n", name)
	return dumpGIF(w, g)
}

func dumpGIF(w io.Writer, g *gif.GIF) error {
	fmt.Fprintf(w, "len(Image)=%d\n", len(g.Image))
	for i, img := range g.Image {
		fmt.Fprintf(w, "  #%-2d Stride=%d Rect=%+v len(Palette)=%d\n", i, img.Stride, img.Rect, len(img.Palette))
	}
	fmt.Fprintf(w, "len(Delay)=%d\n", len(g.Delay))
	fmt.Fprintf(w, "  Delay=%+v\n", g.Delay)
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
