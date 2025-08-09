package engine

import (
	"image"

	"github.com/henilmalaviya/qr-dance/game"
	"github.com/henilmalaviya/qr-dance/util/logger"
)

func NewGameFromBitmap(matrix [][]bool) *game.GameState {
	logger.Trace("Entering NewGameFromBitmap")
	defer logger.Trace("Exiting NewGameFromBitmap")

	g := game.NewGameState(len(matrix[0]), len(matrix))

	// move all cell's x and y into single array
	cells := make([]game.Cell, 0, len(matrix)*len(matrix[0]))
	for y := range matrix {
		for x := range matrix[y] {
			if matrix[y][x] {
				cells = append(cells, *game.NewCellFromCords(x, y))
			}
		}
	}
	g.AddCells(cells)

	return g
}

func RunTicks(g *game.GameState, totalTicks int, scaleFactor int) []*image.RGBA {
	logger.Trace("Entering RunTicks with totalTicks=%d, scaleFactor=%d", totalTicks, scaleFactor)
	defer logger.Trace("Exiting RunTicks")

	images := make([]*image.RGBA, totalTicks)

	for i := 0; i < totalTicks; i++ {
		logger.Debug("Running tick %d of %d", i+1, totalTicks)
		// get matrix -> draw it to png -> upscale it
		images[i] = UpscalePNGImage(DrawBitmapToImage(g.Bitmap()), scaleFactor)

		// Update game state
		g.Update()
	}

	return images
}
