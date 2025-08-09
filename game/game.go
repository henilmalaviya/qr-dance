package game

import (
	"sync"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

type GameState struct {
	width  int
	height int
	grid   Grid
	mutex  sync.RWMutex
}

func (g *GameState) GetGrid() *Grid {
	return &g.grid
}

func (g *GameState) countNeighbors(c Cell) int {
	count := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbor := Cell{c.X + dx, c.Y + dy}
			if g.grid[neighbor] {
				count++
			}
		}
	}
	return count
}

func (g *GameState) Update() (added []Cell, removed []Cell) {
	logger.Trace("Entering GameState.Update")
	defer logger.Trace("Exiting GameState.Update")

	g.mutex.Lock()
	defer g.mutex.Unlock()

	cellsToCheck := make(Grid)
	for cell := range g.grid {
		cellsToCheck[cell] = true
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				neighbor := Cell{cell.X + dx, cell.Y + dy}
				cellsToCheck[neighbor] = true
			}
		}
	}

	newGrid := make(Grid)
	for cell := range cellsToCheck {
		neighbors := g.countNeighbors(cell)
		alive := g.grid[cell]
		if (alive && (neighbors == 2 || neighbors == 3)) || (!alive && neighbors == 3) {
			newGrid[cell] = true
			if !alive {
				added = append(added, cell)
			}
		}
	}

	for cell := range g.grid {
		if !newGrid[cell] {
			removed = append(removed, cell)
		}
	}
	logger.Trace("%d cells added, %d cells removed", len(added), len(removed))

	// make sure the grid size is consistent
	// if the new cells are outside the bounds
	// ignore them
	for cell := range newGrid {
		if cell.X < 0 || cell.X >= g.width || cell.Y < 0 || cell.Y >= g.height {
			delete(newGrid, cell)
			// remove from added and removed if they are outside bounds
			for i, a := range added {
				if a.X == cell.X && a.Y == cell.Y {
					added = append(added[:i], added[i+1:]...)
					break
				}
			}
			for i, r := range removed {
				if r.X == cell.X && r.Y == cell.Y {
					removed = append(removed[:i], removed[i+1:]...)
					break
				}
			}
		}
	}

	g.grid = newGrid

	return added, removed
}

func (g *GameState) PrintGrid() {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.grid[Cell{x, y}] {
				print("1 ")
			} else {
				print("0 ")
			}
		}
		println()
	}

}

func NewGameState(width, height int) *GameState {
	logger.Trace("Creating new game state with width %d and height %d", width, height)
	gameState := &GameState{
		width:  width,
		height: height,
		grid:   make(Grid),
	}

	return gameState
}
