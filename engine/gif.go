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
	logger.Trace("Entering PrepareGIFFromImages with %d images, delay %d, initialFrameDelay %d", len(imgs), delay, initialFrameDelay)
	defer logger.Trace("Exiting PrepareGIFFromImages")

	gifImg := &gif.GIF{
		Image:     make([]*image.Paletted, len(imgs)),
		Delay:     make([]int, len(imgs)),
		LoopCount: 0, // Loop forever
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
		gifImg.Delay[i] = delay / 10 // delay is in 100ths of a second
	}

	if initialFrameDelay > 0 {
		gifImg.Delay[0] = initialFrameDelay / 10 // delay is in 100ths of a second
	}

	logger.Info("GIF preparation complete with %d frames", len(gifImg.Image))
	return gifImg, nil
}

func WriteGIFToFile(g *gif.GIF, destPath string) error {
	logger.Trace("Entering WriteGIFToFile with destPath %s", destPath)
	defer logger.Trace("Exiting WriteGIFToFile")

	f, err := os.Create(destPath)
	if err != nil {
		logger.Error("Failed to create file %s: %v", destPath, err)
		return err
	}
	defer f.Close()

	if err := gif.EncodeAll(f, g); err != nil {
		logger.Error("Failed to encode GIF: %v", err)
		return err
	}

	return nil
}

func GIFToBase64(g *gif.GIF) (string, error) {
	logger.Trace("Entering GIFToBase64")
	defer logger.Trace("Exiting GIFToBase64")

	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		logger.Error("Failed to encode GIF to buffer: %v", err)
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}
