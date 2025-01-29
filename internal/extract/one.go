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

// ExtractRepFrame extracts a representative frame from an animation GIF.
// "representative" means most "informational": highest entropy of image.
func ExtractRepFrame(ctx context.Context, args []string) error {
	var (
		grayedEntropy bool
		output        string
	)

	fs := subcmd.FlagSet(ctx)
	fs.BoolVar(&grayedEntropy, "grayedentropy", false, "use gray scaled image to measure entropy (default: original image)")
	fs.StringVar(&output, "output", "", `output PNG image filename (default: {orignal name} + "_rep.png")`)
	fs.Parse(args)

	if fs.NArg() != 1 {
		return errors.New("only a GIF can be extracted at a time")
	}
	input := fs.Arg(0)

	if output == "" {
		output = appendFilename(input, "_rep.png")
	}

	entropyFunc := colorEntropy
	if grayedEntropy {
		entropyFunc = grayEntropy
	}

	return extractRepresentativeFrame(output, input, entropyFunc)
}

type frameInfo struct {
	i       int
	img     image.Image
	entropy float64
}

func grayEntropy(img image.Image) float64 {
	return imgutil.MeasureEntropy(imgutil.ToGray(img))
}

func colorEntropy(img image.Image) float64 {
	return imgutil.MeasureEntropy(img)
}

func extractRepresentativeFrame(output, input string, measureEntropy func(image.Image) float64) error {
	g, err := gifutil.LoadFile(input)
	if err != nil {
		return err
	}

	highest := frameInfo{i: -1, entropy: -1}
	for i, img := range gifutil.IterateComposed(g) {
		entropy := measureEntropy(img)
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
