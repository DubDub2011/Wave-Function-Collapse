package wfc

import (
	"fmt"
	"math/rand"
)

const (
	LEFT  = iota
	UP    = iota
	RIGHT = iota
	DOWN  = iota
)

type position struct {
	x, y int
}

type Tile struct {
	Id            int
	Configuration []int
}

func match(dir int, tile1, tile2 Tile) bool {
	// four cardinal directions, adding 2 gets to the opposite and then remainder of 4 to prevent overflow
	return tile1.Configuration[dir] == tile2.Configuration[(dir+2)%4]
}

type tileGrid struct {
	tiles [][][]Tile // tracks the possible tiles in a given position
	grid  [][]bool   // tracks the positions that have been collapsed
}

func newTileGrid(width, length int, tileset []Tile) tileGrid {
	if width <= 0 || length <= 0 {
		panic(fmt.Errorf("error creating tile grid, width or length 0"))
	}

	if len(tileset) <= 1 {
		panic(fmt.Errorf("error creating tile grid, one or less tiles defined in the tile set"))
	}

	tiles := make([][][]Tile, length)
	grid := make([][]bool, length)
	for row := range tiles {
		tiles[row] = make([][]Tile, width)
		grid[row] = make([]bool, width)
		for col := range tiles[row] {
			tiles[row][col] = tileset
		}
	}

	return tileGrid{tiles: tiles, grid: grid}
}

// returns if the tile successfully collapsed
func (tg tileGrid) collapseTile(pos position) bool {
	tileCollapsed := tg.grid[pos.x][pos.y]
	if tileCollapsed {
		panic(fmt.Errorf("somehow tile already collapsed at pos %v", pos))
	}

	// collapse
	possibleTiles := tg.tiles[pos.x][pos.y]
	for len(possibleTiles) > 0 {
		selectedTileIdx := rand.Intn(len(possibleTiles))
		selectedTile := possibleTiles[selectedTileIdx]

		// declare func to check a given neighbour
		checkNeighbourFunc := func(pos position, dir int) (bool, func()) {
			newTiles := make([]Tile, 0)
			for _, tile := range tg.tiles[pos.x][pos.y] {
				if match(dir, selectedTile, tile) {
					newTiles = append(newTiles, tile)
				}
			}

			// no new valid tiles, the selected tile is invalid
			if len(newTiles) == 0 {
				return false, nil
			}
			return true, func() { tg.tiles[pos.x][pos.y] = newTiles }
		}

		updateNeighbourFuncs := make([]func(), 0)

		// check each side is valid
		leftAvailable := pos.x != 0 && !tg.grid[pos.x-1][pos.y]
		if leftAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x - 1, pos.y}, LEFT)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles = append(possibleTiles[selectedTileIdx:], possibleTiles[:selectedTileIdx+1]...)
				continue
			}

			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		upAvailable := pos.y != 0 && !tg.grid[pos.x][pos.y-1]
		if upAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x, pos.y - 1}, UP)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles = append(possibleTiles[selectedTileIdx:], possibleTiles[:selectedTileIdx+1]...)
				continue
			}
			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		rightAvailable := pos.x != len(tg.grid)-1 && !tg.grid[pos.x+1][pos.y]
		if rightAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x + 1, pos.y}, RIGHT)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles = append(possibleTiles[selectedTileIdx:], possibleTiles[:selectedTileIdx+1]...)
				continue
			}
			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		downAvailable := pos.y != len(tg.grid[pos.x])-1 && !tg.grid[pos.x][pos.y+1]
		if downAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x, pos.y + 1}, DOWN)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles = append(possibleTiles[selectedTileIdx:], possibleTiles[:selectedTileIdx+1]...)
				continue
			}
			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		// Checked all neighbours are valid, now update the grid to reflect new neighbours
		for _, updateNeighbour := range updateNeighbourFuncs {
			updateNeighbour()
		}

		tg.tiles[pos.x][pos.y] = []Tile{selectedTile}
		tg.grid[pos.x][pos.y] = true
		return true

	}

	return false
}

func (tg tileGrid) getTileIds() [][]int {
	tileIds := make([][]int, len(tg.tiles))
	for row := range tg.tiles {
		tileIds[row] = make([]int, len(tg.tiles[row]))
		for col, tile := range tg.tiles[row] {
			if !tg.grid[row][col] {
				panic(fmt.Errorf("not finished resolving, tile with pos %v not resolved", position{row, col}))
			}

			tileIds[row][col] = tile[0].Id
		}
	}
	return tileIds
}

func (tg tileGrid) tileWithLowestEntropy() (bool, position) {
	// TODO: Figure out better way than hardcoded 100, as it's dangerous
	tileValue := 100
	possiblePos := make([]position, 0)
	for row := range tg.tiles {
		for col := range tg.tiles[row] {
			tileCollapsed := tg.grid[row][col]
			if tileCollapsed {
				continue
			}

			if len(tg.tiles[row][col]) < tileValue {
				tileValue = len(tg.tiles[row][col])
				possiblePos = []position{{row, col}}
			} else if len(tg.tiles[row][col]) == tileValue {
				possiblePos = append(possiblePos, position{row, col})
			}
		}
	}

	// Means we're done, no tiles to be collapsed left
	if len(possiblePos) == 0 {
		return true, position{}
	}

	res := possiblePos[rand.Intn(len(possiblePos))]
	return false, res
}

func Collapse(tiles []Tile, width int, length int) ([][]int, error) {
	tg := newTileGrid(width, length, tiles)

	pos := getRandPos(width, length)
	finished := false
	for !finished {
		collapsedTile := tg.collapseTile(pos)
		if !collapsedTile {
			return nil, fmt.Errorf("collapse failed at pos %v", pos)
		}
		finished, pos = tg.tileWithLowestEntropy()
	}

	return tg.getTileIds(), nil
}

func getRandPos(width, length int) position {
	return position{
		x: rand.Intn(width),
		y: rand.Intn(length),
	}
}
