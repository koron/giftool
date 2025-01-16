package analyse

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/koron-go/subcmd"
)

var Entropy = subcmd.DefineCommand("entropy", "calc entropy of image", func(ctx context.Context, args []string) error {
	fs := subcmd.FlagSet(ctx)
	fs.Parse(args)

	var highest struct {
		i       int
		name    string
		entropy float64
	}
	highest.entropy = math.Inf(-1)

	for i, name := range fs.Args() {
		cimg, err := loadImage(name)
		if err != nil {
			return fmt.Errorf("failed to load image %q: %w", name, err)
		}
		gray := toGrayImage(cimg)
		hist := calcHist(gray)
		entropy := histToEntropy(hist)
		chist := calcHist(cimg)
		centropy := histToEntropy(chist)
		log.Printf("#%d %s : entropy=%f color_entropy=%f", i, name, entropy, centropy)
		if highest.entropy < entropy {
			highest.i = i
			highest.name = name
			highest.entropy = entropy
		}
	}

	log.Printf("highest entropy: #%d %s : entropy=%f", highest.i, highest.name, highest.entropy)
	return nil
})

func calcHist(img image.Image) map[color.RGBA]int {
	hist := map[color.RGBA]int{}
	rect := img.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			c := color.RGBA{
				R: uint8(r / 0x0101),
				G: uint8(g / 0x0101),
				B: uint8(b / 0x0101),
				A: uint8(a / 0x0101),
			}
			hist[c]++
		}
	}
	return hist
}

func histToEntropy[T comparable](hist map[T]int) float64 {
	var entropy float64
	var sum int
	for _, n := range hist {
		sum += n
	}
	for _, n := range hist {
		p := float64(n) / float64(sum)
		entropy += -p * math.Log2(p)
	}
	return entropy
}

func loadImage(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func toGrayImage(img image.Image) *image.Gray {
	r := img.Bounds()
	gray := image.NewGray(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			c := img.At(x, y)
			g := color.GrayModel.Convert(c)
			gray.Set(x, y, g)
		}
	}
	return gray
}
