package imageprocessing

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"strings"

	"github.com/fogleman/gg"
	"github.com/oliamb/cutter"
)

type BoundingBox struct {
	Top    float64 `json:"top"`
	Left   float64 `json:"left"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Base64toCroppedJpg Given a base64 string of a JPEG, encodes it into an cropped JPEG image
func Base64toCroppedJpg(data string, boundingBox *BoundingBox) (*bytes.Buffer, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	ctx := gg.NewContextForImage(m)
	croppedImg, err := cutter.Crop(m, cutter.Config{
		Width:  int(boundingBox.Width * float64(ctx.Width())),
		Height: int(boundingBox.Height * float64(ctx.Height())),
		Anchor: image.Point{int(boundingBox.Left * float64(ctx.Width())), int(boundingBox.Top * float64(ctx.Height()))},
		Mode:   cutter.TopLeft,
	})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, croppedImg, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, err
	}
	return buf, nil
}
