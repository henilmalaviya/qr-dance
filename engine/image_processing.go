package engine

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

func UpscalePNGImage(img *image.RGBA, scale int) *image.RGBA {
	logger.Trace("Entering UpscalePNGImage with scale %d", scale)
	defer logger.Trace("Exiting UpscalePNGImage")

	originalBounds := img.Bounds()
	width := originalBounds.Dx() * scale
	height := originalBounds.Dy() * scale
	upscaledImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origX := x / scale
			origY := y / scale
			upscaledImg.Set(x, y, img.At(origX, origY))
		}
	}

	return upscaledImg
}

func DrawBitmapToImage(matrix [][]bool) *image.RGBA {
	logger.Trace("Entering DrawBitmapToImage")
	defer logger.Trace("Exiting DrawBitmapToImage")

	height := len(matrix)
	if height == 0 {
		logger.Warn("DrawBitmapToImage called with an empty matrix")
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	width := len(matrix[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y, row := range matrix {
		for x, cell := range row {
			if cell {
				img.Set(x, y, image.Black)
			} else {
				img.Set(x, y, image.White)
			}
		}
	}

	return img
}

func SaveImageToFile(img *image.RGBA, path string) error {
	logger.Trace("Entering SaveImageToFile with path %s", path)
	defer logger.Trace("Exiting SaveImageToFile")

	f, err := os.Create(path)
	if err != nil {
		logger.Error("Failed to create file %s: %v", path, err)
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		logger.Error("Failed to encode image to file %s: %v", path, err)
		return fmt.Errorf("failed to encode image to file %s: %w", path, err)
	}

	return nil
}
