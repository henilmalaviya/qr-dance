package game

import (
	"fmt"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

type Cell struct {
	X, Y int
}

func NewCellFromCords(x int, y int) *Cell {
	return &Cell{x, y}
}

func NewCell() *Cell {
	return NewCellFromCords(0, 0)
}

func (c *Cell) String() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

// ----

type Grid = map[Cell]bool

func (g *GameState) AddCells(cells []Cell) {
	logger.Trace("Entering AddCells with %d cells", len(cells))
	defer logger.Trace("Exiting AddCells")
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, cell := range cells {
		g.grid[cell] = true
	}
}

func (g *GameState) AddCell(c Cell) {
	logger.Trace("Adding single cell: (%d, %d)", c.X, c.Y)
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.grid[c] = true
}

func (g *GameState) RemoveCell(c Cell) {
	logger.Trace("Removing single cell: (%d, %d)", c.X, c.Y)
	g.mutex.Lock()
	defer g.mutex.Unlock()

	delete(g.grid, c)
}

func (g *GameState) GetMaxX() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	maxX := 0
	for cell := range g.grid {
		if cell.X > maxX {
			maxX = cell.X
		}
	}
	return maxX
}

func (g *GameState) GetMaxY() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	maxY := 0
	for cell := range g.grid {
		if cell.Y > maxY {
			maxY = cell.Y
		}
	}
	return maxY
}

func (g *GameState) GetBounds() (int, int) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	maxX := g.GetMaxX()
	maxY := g.GetMaxY()

	return maxX + 1, maxY + 1
}

func (g *GameState) Bitmap() [][]bool {
	logger.Trace("Entering Bitmap")
	defer logger.Trace("Exiting Bitmap")
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	matrix := make([][]bool, g.height)
	for i := range matrix {
		matrix[i] = make([]bool, g.width)
	}

	for cell := range g.grid {
		if cell.X < 0 || cell.X >= g.width || cell.Y < 0 || cell.Y >= g.height {
			logger.Warn("Cell (%d, %d) is out of bounds and will be ignored", cell.X, cell.Y)
			continue
		}
		matrix[cell.Y][cell.X] = true
	}

	return matrix
}
