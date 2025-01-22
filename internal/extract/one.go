package extract

import (
	"context"
	"errors"
	"image"
	"log"

	"github.com/koron-go/subcmd"
	"github.com/koron/giftool/internal/gifutil"
	"github.com/koron/giftool/internal/imgutil"
)

func ExtractRepFrame(ctx context.Context, args []string) error {
	fs := subcmd.FlagSet(ctx)
	fs.Parse(args)

	if fs.NArg() != 1 {
		return errors.New("only a GIF can be extracted at a time")
	}
	input := fs.Arg(0)

	output := appendFilename(input, "_rep.png")
	return extractRepresentativeFrame(output, input)
}

type frameInfo struct {
	i       int
	img     image.Image
	entropy float64
}

func extractRepresentativeFrame(output, input string) error {
	g, err := gifutil.LoadFile(input)
	if err != nil {
		return err
	}

	highest := frameInfo{i: -1, entropy: -1}
	for i, img := range gifutil.IterateComposed(g) {
		gray := imgutil.ToGray(img)
		entropy := imgutil.MeasureEntropy(gray)
		if entropy > highest.entropy {
			highest = frameInfo{i: i, img: img, entropy: entropy}
		}
	}
	log.Printf("extracted #%d from %s: entropy=%f", highest.i, input, highest.entropy)
	if highest.img != nil && output != "" {
		imgutil.WritePNGFile(output, highest.img)
	}

	return nil
}
