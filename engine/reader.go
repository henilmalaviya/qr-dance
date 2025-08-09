package engine

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

func ReadInputFile(filePath string) (*image.Image, error) {
	logger.Info("Reading input file: %s", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open file: %s: %v", filePath, err)
		return nil, err
	}
	defer f.Close()

	src, err := png.Decode(f)
	if err != nil {
		logger.Error("Failed to decode PNG: %v", err)
		return nil, err
	}
	logger.Debug("Decoded PNG from file: bounds=%v", src.Bounds())
	return &src, nil
}

func ReadInputBase64(base64Data string) (*image.Image, error) {
	logger.Info("Reading input from base64 string")

	// Remove data URL prefix if present (e.g., "data:image/png;base64,")
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			logger.Debug("Stripping data URL prefix from base64 input")
			base64Data = parts[1]
		}
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		logger.Error("Failed to decode base64 input: %v", err)
		return nil, err
	}
	logger.Debug("Base64 decoded: bytes=%d", len(data))

	// Create an image from the decoded data
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		logger.Error("Failed to decode PNG from base64: %v", err)
		return nil, err
	}
	logger.Debug("Decoded PNG from base64: bounds=%v", src.Bounds())
	return &src, nil
}

func ReadInput(input string, isBase64 bool) (*image.Image, error) {
	if isBase64 {
		return ReadInputBase64(input)
	}
	return ReadInputFile(input)
}
