// Package extract provides "extract" sub command: extract all images from a GIF image.
package extract

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"path/filepath"
	"strings"

	"github.com/koron-go/subcmd"
	"github.com/koron/giftool/internal/imgutil"
)

type exposeTarget int

const (
	targetFrames exposeTarget = iota
	targetComposed
	targetRepresentative
)

var ExtractSet = subcmd.DefineSet("extract", "extract frames from animation GIF",
	subcmd.DefineCommand("one", "extract a representative frame from GIF", ExtractRepFrame),
	subcmd.DefineCommand("representative", "extract a representative frame from GIF", ExtractRepFrame),
	subcmd.DefineCommand("rep", "extract a representative frame from GIF", ExtractRepFrame),
	subcmd.DefineCommand("cframes", "extract all composed frames from GIF", ExtractComposedFrames),
	//subcmd.DefineCommand("all", "extract all frames from GIF", ExtractAll),
)

var Extract = subcmd.DefineCommand("extract", "extract each frames from GIF", ExtractAll)

func ExtractAll(ctx context.Context, args []string) error {
	fs := subcmd.FlagSet(ctx)
	var (
		outdir string
		target string
	)
	fs.StringVar(&outdir, "outdir", "", "output directory (default: base of input)")
	fs.StringVar(&target, "target", "frames", "what to extract")
	fs.Parse(args)

	var exposeMode exposeTarget
	switch strings.ToLower(target) {
	case "frames":
		exposeMode = targetFrames
	case "composed":
		exposeMode = targetComposed
	case "representative":
		exposeMode = targetRepresentative
	default:
		return fmt.Errorf("unknown expose target: %q", target)
	}

	if fs.NArg() != 1 {
		return errors.New("only one GIF can be extracted at a time")
	}
	input := fs.Arg(0)

	if outdir == "" {
		outdir = input[0 : len(input)-len(filepath.Ext(input))]
	}

	switch exposeMode {
	case targetFrames:
		return extractFrames(outdir, input)
	case targetComposed:
		return extractComposedFrames(outdir, input)
	case targetRepresentative:
		return extractRepresentative(outdir, input)
	default:
		panic("invalid expose mode")
	}
}

func extractRepresentative(outdir, input string) error {
	g, err := prepareExpose(outdir, input)
	if err != nil {
		return err
	}

	images := make([]*image.RGBA, len(g.Image))
	err = forComposedFrames(g, func(i int, img *image.RGBA) error {
		copy := image.NewRGBA(image.Rect(0, 0, g.Config.Width, g.Config.Height))
		draw.Draw(copy, img.Rect, img, img.Rect.Min, draw.Over)
		images[i] = copy
		return nil
	})
	if err != nil {
		return err
	}

	avg := averagingImages(images, g.Delay)
	output := filepath.Join(outdir, "000avg.png")
	err = imgutil.WritePNGFile(output, avg)
	if err != nil {
		return err
	}

	for i, img := range images {
		g16 := diffImages(avg, img)
		p := measureGray16(g16)
		fmt.Printf("#%-3d %f\n", i, p)
		output := filepath.Join(outdir, fmt.Sprintf("%03dd.png", i))
		err := imgutil.WritePNGFile(output, g16)
		if err != nil {
			return err
		}
	}

	return nil
}

func composeImages(dst *image.RGBA, src *image.Paletted, disposal, bgIndex byte) error {
	op := draw.Over
	switch disposal {
	case 0:
		op = draw.Src
	}
	draw.Draw(dst, src.Rect, src, src.Rect.Min, op)
	return nil
}

func extractFrames(outdir, input string) error {
	g, err := prepareExpose(outdir, input)
	if err != nil {
		return err
	}

	for i, img := range g.Image {
		output := filepath.Join(outdir, fmt.Sprintf("%03d.png", i))
		err := imgutil.WritePNGFile(output, img)
		if err != nil {
			return err
		}
	}
	return nil
}
