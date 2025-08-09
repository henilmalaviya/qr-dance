package engine

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

func PrepareGIFFromImages(imgs []*image.RGBA, delay int, initialFrameDelay int) (*gif.GIF, error) {
	logger.Info("Preparing GIF: frames=%d delay=%dms initialFrameDelay=%dms", len(imgs), delay, initialFrameDelay)
	gifImg := &gif.GIF{
		Image:     make([]*image.Paletted, len(imgs)),
		Delay:     make([]int, len(imgs)),
		LoopCount: 0,
	}

	// Create a simple black and white palette
	palette := []color.Color{
		color.RGBA{255, 255, 255, 255}, // White
		color.RGBA{0, 0, 0, 255},       // Black
	}

	for i, img := range imgs {
		palettedImg := image.NewPaletted(img.Bounds(), palette)
		draw.Draw(palettedImg, palettedImg.Rect, img, image.Point{}, draw.Src)
		gifImg.Image[i] = palettedImg
		gifImg.Delay[i] = delay / 10
		if i < 3 || i == len(imgs)-1 {
			logger.Trace("Prepared frame %d/%d", i+1, len(imgs))
		}
	}

	if initialFrameDelay > 0 {
		gifImg.Delay[0] = initialFrameDelay / 10
	}

	logger.Debug("GIF prepared: palettedFrames=%d", len(gifImg.Image))
	return gifImg, nil
}

func WriteGIFToFile(g *gif.GIF, destPath string) error {
	logger.Info("Writing GIF to file: %s", destPath)
	f, err := os.Create(destPath)
	if err != nil {
		logger.Error("Failed to create file %s: %v", destPath, err)
		return err
	}
	defer f.Close()

	if err := gif.EncodeAll(f, g); err != nil {
		logger.Error("Failed to encode GIF to file: %v", err)
		return err
	}

	logger.Info("GIF successfully written: %s", destPath)
	return nil
}

func GIFToBase64(g *gif.GIF) (string, error) {
	logger.Debug("Encoding GIF to base64")
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		logger.Error("Failed to encode GIF to buffer: %v", err)
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	logger.Debug("Base64 encoding complete: bytes=%d", len(encoded))
	return encoded, nil
}
