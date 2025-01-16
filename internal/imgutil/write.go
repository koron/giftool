package imgutil

import (
	"image"
	"image/gif"
	"image/png"
	"os"
)

// WritePNGFile writes an image to a file in PNG format.
func WritePNGFile(output string, img image.Image) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// WriteGIFFile writes an image to a file in GIF format.
func WriteGIFFile(output string, img image.Image) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	opts := gif.Options{}
	return gif.Encode(f, img, &opts)
}
