package extract

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"path/filepath"
)

func appendFilename(name, suffix string) string {
	ext := filepath.Ext(name)
	return name[0:len(name)-len(ext)] + suffix
}

func loadGIF(name string) (*gif.GIF, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	return g, err
}

func prepareExpose(outdir, input string) (*gif.GIF, error) {
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	g, err := gif.DecodeAll(f)
	f.Close()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(outdir, 0777)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func forComposedFrames(g *gif.GIF, fn func(int, *image.RGBA) error) error {
	buf := image.NewRGBA(image.Rect(0, 0, g.Config.Width, g.Config.Height))
	for i, img := range g.Image {
		err := composeImages(buf, img, g.Disposal[i], g.BackgroundIndex)
		if err != nil {
			return err
		}
		err = fn(i, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func averagingImages(images []*image.RGBA, delay []int) *image.RGBA {
	first := images[0]
	width, height := first.Rect.Dx(), first.Rect.Dy()

	buf := make([]uint64, width*height*4)
	var wsum uint64
	for i, img := range images {
		weight := uint64(delay[i])
		if weight == 0 {
			weight = 1
		}
		wsum += weight
		for y := range height {
			for x := range width {
				c := img.RGBAAt(x, y)
				j := (x + y*width) * 4
				buf[j+0] += uint64(c.R) * weight
				buf[j+1] += uint64(c.G) * weight
				buf[j+2] += uint64(c.B) * weight
				buf[j+3] += uint64(c.A) * weight
			}
		}
	}

	avg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			idx := (x + y*width) * 4
			c := color.RGBA{
				R: uint8(buf[idx+0] / wsum),
				G: uint8(buf[idx+1] / wsum),
				B: uint8(buf[idx+2] / wsum),
				A: uint8(buf[idx+3] / wsum),
			}
			avg.Set(x, y, c)
		}
	}

	return avg
}

type Lab struct {
	L float64
	A float64
	B float64
}

func gammaCorrection(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func labHelper(v float64) float64 {
	if v > 0.008856 {
		return math.Pow(v, 1.0/3.0)
	}
	return (903.3*v + 16) / 116
}

func rgbaToLab(c color.RGBA) Lab {
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	r = gammaCorrection(r)
	g = gammaCorrection(g)
	b = gammaCorrection(b)

	x := r*0.4124 + g*0.3576 + b*0.1805
	y := r*0.2126 + g*0.7152 + b*0.0722
	z := r*0.0193 + g*0.1192 + b*0.9505

	// XYZ -> Lab
	x = labHelper(x / 0.95047)
	y = labHelper(y / 1.00000)
	z = labHelper(z / 1.08883)

	return Lab{
		L: math.Max(0, 116*y-16),
		A: 500 * (x - y),
		B: 200 * (y - z),
	}
}

func diffColors(a, b color.RGBA) color.Gray16 {
	la := rgbaToLab(a)
	lb := rgbaToLab(b)
	d := math.Sqrt(math.Pow(la.L-lb.L, 2) + math.Pow(la.A-lb.A, 2) + math.Pow(la.B-lb.B, 2))
	return color.Gray16{Y: uint16(d * 0xffff)}
}

func diffImages(a, b *image.RGBA) *image.Gray16 {
	width, height := a.Rect.Dx(), a.Rect.Dy()
	diff := image.NewGray16(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			ca := a.RGBAAt(x, y)
			cb := b.RGBAAt(x, y)
			d16 := diffColors(ca, cb)
			diff.SetGray16(x, y, d16)
		}
	}
	return diff
}

func measureGray16(g *image.Gray16) float64 {
	width, height := g.Rect.Dx(), g.Rect.Dy()
	var sum float64
	for y := range height {
		for x := range width {
			c := g.Gray16At(x, y)
			sum += float64(c.Y) / float64(0xffff)
		}
	}
	return sum / (float64(width) * float64(height))
}
