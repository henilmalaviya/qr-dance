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
	g.mutex.Lock()
	defer g.mutex.Unlock()

	for _, cell := range cells {
		g.grid[cell] = true
	}
	logger.Trace("AddCells: count=%d total=%d", len(cells), len(g.grid))
}

func (g *GameState) AddCell(c Cell) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.grid[c] = true

	logger.Trace("AddCell: %v total=%d", c, len(g.grid))
}

func (g *GameState) RemoveCell(c Cell) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	delete(g.grid, c)
	logger.Trace("RemoveCell: %v total=%d", c, len(g.grid))
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

func (g *GameState) Get2DMatrix() [][]int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	matrix := make([][]int, g.height)
	for i := range matrix {
		matrix[i] = make([]int, g.width)
	}

	for cell := range g.grid {
		if cell.X < 0 || cell.X >= g.width || cell.Y < 0 || cell.Y >= g.height {
			continue
		}
		matrix[cell.Y][cell.X] = 1
	}

	logger.Trace("Get2DMatrix: size=%dx%d", len(matrix[0]), len(matrix))
	return matrix
}
