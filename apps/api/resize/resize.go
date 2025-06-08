package resize

import (
	"bytes"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
)

func PNG(input []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	srcBounds := img.Bounds()
	width := srcBounds.Dx()
	height := srcBounds.Dy()

	var newImg *image.NRGBA
	if width >= height {
		newImg = imaging.Resize(img, 512, 0, imaging.Lanczos)
	} else {
		newImg = imaging.Resize(img, 0, 512, imaging.Lanczos)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, newImg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
