package engine

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"strings"
)

func ReadInputFile(filePath string) (*image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return &src, nil
}

func ReadInputBase64(base64Data string) (*image.Image, error) {
	// Remove data URL prefix if present (e.g., "data:image/png;base64,")
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}

	// Create an image from the decoded data
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &src, nil
}

func ReadInput(input string, isBase64 bool) (*image.Image, error) {
	if isBase64 {
		return ReadInputBase64(input)
	}
	return ReadInputFile(input)
}
