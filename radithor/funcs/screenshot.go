package funcs

import (
	"bytes"
	"image/png"

	"github.com/kbinani/screenshot"
)

func CaptureScreen(index int) (*bytes.Buffer, error) {
	bounds := screenshot.GetDisplayBounds(index)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	err = png.Encode(buffer, img)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
