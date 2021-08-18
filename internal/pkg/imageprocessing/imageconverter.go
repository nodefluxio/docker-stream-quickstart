package imageprocessing

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"

	"golang.org/x/image/bmp"
)

func ConvertToJPG(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/jpeg":
		return imageBytes, nil
	case "image/jpg":
		return imageBytes, nil
	case "image/png":
		img, err := png.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case "image/bmp":
		img, err := bmp.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("image %s format not yet supported", contentType)
	}
}
