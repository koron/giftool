package extract

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/koron-go/subcmd"
	"github.com/koron/giftool/internal/gifutil"
	"github.com/koron/giftool/internal/imgutil"
)

func ExtractComposedFrames(ctx context.Context, args []string) error {
	var (
		outdir string
	)
	fs := subcmd.FlagSet(ctx)
	fs.StringVar(&outdir, "outdir", "", "output directory (default: base of input)")
	fs.Parse(args)

	if fs.NArg() != 1 {
		return errors.New("only one GIF can be extracted at a time")
	}
	input := fs.Arg(0)

	if outdir == "" {
		outdir = input[0 : len(input)-len(filepath.Ext(input))]
	}

	return extractComposedFrames(outdir, input)
}

func extractComposedFrames(outdir, input string) error {
	g, err := prepareExpose(outdir, input)
	if err != nil {
		return err
	}
	for i, img := range gifutil.IterateComposed(g) {
		output := filepath.Join(outdir, fmt.Sprintf("%03d_composed.png", i))
		err := imgutil.WritePNGFile(output, img)
		if err != nil {
			return err
		}
	}
	return nil
}
