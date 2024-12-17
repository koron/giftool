// Package extract provides "extract" sub command: extract all images from a GIF image.
package extract

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"

	"github.com/koron-go/subcmd"
)

var Extract = subcmd.DefineCommand("extract", "extract each frames from GIF", func(ctx context.Context, args []string) error {
	fs := subcmd.FlagSet(ctx)
	var outdir string
	fs.StringVar(&outdir, "outdir", "", "output directory (default: base of input)")
	fs.Parse(args)
	if fs.NArg() != 1 {
		return errors.New("only one GIF can be extracted at a time")
	}
	input := fs.Arg(0)
	if outdir == "" {
		outdir = input[0 : len(input)-len(filepath.Ext(input))]
	}
	return extractFrames(outdir, input)
})

func extractFrames(outdir, input string) error {
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return err
	}
	os.MkdirAll(outdir, 0777)

	for i, img := range g.Image {
		output := filepath.Join(outdir, fmt.Sprintf("%03d.png", i))
		err := writeImage(output, img)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeImage(output string, img image.Image) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
