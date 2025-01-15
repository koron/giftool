package gifutil

import (
	"image/gif"
	"os"
)

func LoadFile(name string) (*gif.GIF, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	return g, err
}
