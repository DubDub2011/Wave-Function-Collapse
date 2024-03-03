package wfc

import (
	"fmt"
	"math/rand"
)

type tileGrid struct {
	tiles              [][][]Tile // tracks the possible tiles in a given position
	positionsCollapsed [][]bool   // tracks the positions that have been collapsed
}

func newTileGrid(width, height int, tileset []Tile) tileGrid {
	if width <= 0 || height <= 0 {
		panic(fmt.Errorf("error creating tile grid, width or height 0"))
	}

	if len(tileset) <= 1 {
		panic(fmt.Errorf("error creating tile grid, one or less tiles defined in the tile set"))
	}

	tiles := make([][][]Tile, width)
	grid := make([][]bool, width)
	for row := range tiles {
		tiles[row] = make([][]Tile, height)
		grid[row] = make([]bool, height)
		for col := range tiles[row] {
			tiles[row][col] = tileset
		}
	}

	return tileGrid{tiles: tiles, positionsCollapsed: grid}
}

// returns if the tile successfully collapsed
func (tg tileGrid) collapseTile(pos position) bool {
	tileCollapsed := tg.positionsCollapsed[pos.x][pos.y]
	if tileCollapsed {
		panic(fmt.Errorf("attempt to collapse already collapsed tile at pos %v", pos))
	}

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

			// return a func to update the neighbour, in case other neighbours are invalid
			return true, func() {
				tg.tiles[pos.x][pos.y] = newTiles
			}
		}

		updateNeighbourFuncs := make([]func(), 0)

		// TODO: Sort of repetitive below
		leftAvailable := pos.x != 0 && !tg.positionsCollapsed[pos.x-1][pos.y]
		if leftAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x - 1, pos.y}, LEFT)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles[selectedTileIdx] = possibleTiles[len(possibleTiles)-1]
				possibleTiles = possibleTiles[:len(possibleTiles)-1]
				continue
			}

			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		upAvailable := pos.y != 0 && !tg.positionsCollapsed[pos.x][pos.y-1]
		if upAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x, pos.y - 1}, UP)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles[selectedTileIdx] = possibleTiles[len(possibleTiles)-1]
				possibleTiles = possibleTiles[:len(possibleTiles)-1]
				continue
			}
			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		rightAvailable := pos.x != len(tg.positionsCollapsed)-1 && !tg.positionsCollapsed[pos.x+1][pos.y]
		if rightAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x + 1, pos.y}, RIGHT)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles[selectedTileIdx] = possibleTiles[len(possibleTiles)-1]
				possibleTiles = possibleTiles[:len(possibleTiles)-1]
				continue
			}
			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		downAvailable := pos.y != len(tg.positionsCollapsed[pos.x])-1 && !tg.positionsCollapsed[pos.x][pos.y+1]
		if downAvailable {
			success, updateNeighbourFunc := checkNeighbourFunc(position{pos.x, pos.y + 1}, DOWN)
			if !success {
				// selected tile incompatible, choose new tile
				possibleTiles[selectedTileIdx] = possibleTiles[len(possibleTiles)-1]
				possibleTiles = possibleTiles[:len(possibleTiles)-1]
				continue
			}
			updateNeighbourFuncs = append(updateNeighbourFuncs, updateNeighbourFunc)
		}

		// Checked all neighbours are valid, now update the grid to reflect new neighbours
		for _, updateNeighbour := range updateNeighbourFuncs {
			updateNeighbour()
		}

		tg.tiles[pos.x][pos.y] = []Tile{selectedTile}
		tg.positionsCollapsed[pos.x][pos.y] = true
		return true

	}

	return false
}

// returns the ids of a completely collapsed grid for rendering
func (tg tileGrid) getTileIds() [][]int {
	tileIds := make([][]int, len(tg.tiles))
	for row := range tg.tiles {
		tileIds[row] = make([]int, len(tg.tiles[row]))
		for col, tile := range tg.tiles[row] {
			if !tg.positionsCollapsed[row][col] {
				panic(fmt.Errorf("not finished resolving, tile with pos %v not resolved", position{row, col}))
			}

			tileIds[row][col] = tile[0].Id
		}
	}
	return tileIds
}

// returns position with lowest entropy, and nil if none can be found
func (tg tileGrid) tileWithLowestEntropy() *position {
	lowestOptions := -1
	possiblePos := make([]position, 0)

	// First find the lowest value of a tile
	for row := range tg.tiles {
		for col := range tg.tiles[row] {
			tileCollapsed := tg.positionsCollapsed[row][col]
			if tileCollapsed {
				continue
			}

			if len(tg.tiles[row][col]) < lowestOptions || lowestOptions == -1 {
				lowestOptions = len(tg.tiles[row][col])
			}
		}
	}

	// Means we're done, no tiles to be collapsed left
	if lowestOptions == -1 {
		return nil
	}

	// Now append options that have the same lowest tile
	for row := range tg.tiles {
		for col := range tg.tiles[row] {
			tileCollapsed := tg.positionsCollapsed[row][col]
			if tileCollapsed {
				continue
			}

			if len(tg.tiles[row][col]) == lowestOptions {
				possiblePos = append(possiblePos, position{row, col})
			}
		}
	}

	res := possiblePos[rand.Intn(len(possiblePos))]
	return &res
}

// returns tile config for given position
// returns nil if position out of range, or the tile has already been collapsed (as there is no tile config now)
func (tg tileGrid) getTileConfig(pos position) []Tile {
	if pos.x < 0 || pos.x >= len(tg.tiles) {
		return nil
	}

	if pos.y < 0 || pos.y >= len(tg.tiles[0]) {
		return nil
	}

	if tg.positionsCollapsed[pos.x][pos.y] {
		return nil
	}

	res := make([]Tile, len(tg.tiles[pos.x][pos.y]))
	copy(res, tg.tiles[pos.x][pos.y])
	return res
}

// return tile at a given position
// returns nil if position out of range, or the tile hasn't been collapsed yet
func (tg tileGrid) getTileId(pos position) *Tile {
	if pos.x < 0 || pos.x >= len(tg.tiles) {
		return nil
	}

	if pos.y < 0 || pos.y >= len(tg.tiles[0]) {
		return nil
	}

	if !tg.positionsCollapsed[pos.x][pos.y] {
		return nil
	}

	tile := tg.tiles[pos.x][pos.y][0]
	res := Tile{
		tile.Id,
		tile.Configuration,
	}

	return &res
}

// need to update a position directly, in the scenario of backtracking
func (tg tileGrid) revertTileConfig(pos position, tileConfig []Tile) {
	if pos.x < 0 || pos.x >= len(tg.tiles) {
		panic(fmt.Errorf("attempt to update out of bounds position %v,", pos))
	}

	if pos.y < 0 || pos.y >= len(tg.tiles[0]) {
		panic(fmt.Errorf("attempt to update out of bounds position %v,", pos))
	}

	tg.positionsCollapsed[pos.x][pos.y] = false
	tg.tiles[pos.x][pos.y] = tileConfig
}
