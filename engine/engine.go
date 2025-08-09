package engine

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/henilmalaviya/qr-dance/game"
	"github.com/henilmalaviya/qr-dance/util/logger"
)

func NewGameFromQRMatrix(matrix [][]int) *game.GameState {
	logger.Debug("Creating game from matrix: w=%d h=%d", len(matrix[0]), len(matrix))
	g := game.NewGameState(len(matrix[0]), len(matrix))

	// add cells to the game grid based on the matrix
	added := 0
	for y := range matrix {
		for x := range matrix[y] {
			if matrix[y][x] == 1 {
				g.AddCell(*game.NewCellFromCords(x, y))
				added++
			}
		}
	}
	logger.Debug("Initial live cells added: %d", added)

	return g
}

func RunTicks(g *game.GameState, totalTicks int, scaleFactor int) []*image.RGBA {
	logger.Debug("RunTicks start: totalTicks=%d scale=%d", totalTicks, scaleFactor)
	images := make([]*image.RGBA, totalTicks)

	for i := 0; i < totalTicks; i++ {
		// get matrix -> draw it to png -> upscale it
		images[i] = UpscalePNGImage(DrawMatrixToPNG(g.Get2DMatrix()), scaleFactor)

		added, removed := g.Update()
		if i < 5 || (i%50 == 0) || i == totalTicks { // sample logs
			logger.Trace("tick=%d added=%d removed=%d", i, len(added), len(removed))
		}
	}

	logger.Debug("RunTicks end: frames=%d", len(images))
	return images
}

func SaveImageToFile(img *image.RGBA, path string) error {
	logger.Info("Saving image to file: %s", path)
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

	logger.Info("Image saved to file: %s", path)
	return nil
}
