package wfc

import (
	"fmt"
	"math/rand"
)

// tileGrid is responsible for tracking the tiles selected
type tileGrid struct {
	tileConfigurations [][][]Tile // tracks the possible tiles in a given position
	positionsCollapsed [][]bool   // tracks the positions that have been collapsed
}

// Returns a new tileGrid to the given width, height and use the tileset's IDs to track the tiles
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

	return tileGrid{tileConfigurations: tiles, positionsCollapsed: grid}
}

// Selects a random valid tile at the given position and update relevant neighbours
// Returns if the position was successfully collapsed
func (tg tileGrid) collapseTile(pos position) bool {
	// Check position hasn't already been collapsed
	tileCollapsed := tg.positionsCollapsed[pos.x][pos.y]
	if tileCollapsed {
		panic(fmt.Errorf("attempt to collapse already collapsed tile at pos %v", pos))
	}

	possibleTiles := tg.tileConfigurations[pos.x][pos.y]
	for len(possibleTiles) > 0 {
		// Select random tile from the possible tiles
		selectedTileIdx := rand.Intn(len(possibleTiles))
		selectedTile := possibleTiles[selectedTileIdx]

		// Tile is invalid if it makes any of its neighbours invalid, so need to check neighbours
		updatedNeighbourTilesFunc := func(pos position, dir int) []Tile {
			if pos.x < 0 || pos.x >= len(tg.positionsCollapsed) {
				// out of bounds, so don't need to worry about this pos
				return nil
			}

			if pos.y < 0 || pos.y >= len(tg.positionsCollapsed[0]) {
				// out of bounds, so don't need to worry about this pos
				return nil
			}

			if tg.positionsCollapsed[pos.x][pos.y] {
				// already collapsed, so don't need to worry about this pos
				return nil
			}

			// Update current tiles based on selected tile and return them
			newTiles := make([]Tile, 0)
			for _, tile := range tg.tileConfigurations[pos.x][pos.y] {
				if match(dir, selectedTile, tile) {
					newTiles = append(newTiles, tile)
				}
			}

			return newTiles
		}

		// Get the potential new tiles for the neighbours
		neighbours := make(map[int][]Tile)
		neighbours[LEFT] = updatedNeighbourTilesFunc(position{pos.x - 1, pos.y}, LEFT)
		neighbours[UP] = updatedNeighbourTilesFunc(position{pos.x, pos.y - 1}, UP)
		neighbours[RIGHT] = updatedNeighbourTilesFunc(position{pos.x + 1, pos.y}, RIGHT)
		neighbours[DOWN] = updatedNeighbourTilesFunc(position{pos.x, pos.y + 1}, DOWN)

		// Now check if all the neighbours have at least one option (nil is fine, as it means the neighbour is safe)
		tileValid := (neighbours[LEFT] == nil || len(neighbours[LEFT]) > 0) &&
			(neighbours[UP] == nil || len(neighbours[UP]) > 0) &&
			(neighbours[RIGHT] == nil || len(neighbours[RIGHT]) > 0) &&
			(neighbours[DOWN] == nil || len(neighbours[DOWN]) > 0)

		if !tileValid {
			// Selected tile didn't work, remove from array and try again
			possibleTiles[selectedTileIdx] = possibleTiles[len(possibleTiles)-1]
			possibleTiles = possibleTiles[:len(possibleTiles)-1]
			continue
		}

		// Tile is valid, save neighbours and save position
		if neighbours[LEFT] != nil {
			tg.tileConfigurations[pos.x-1][pos.y] = neighbours[LEFT]
		}
		if neighbours[UP] != nil {
			tg.tileConfigurations[pos.x][pos.y-1] = neighbours[UP]
		}
		if neighbours[RIGHT] != nil {
			tg.tileConfigurations[pos.x+1][pos.y] = neighbours[RIGHT]
		}
		if neighbours[DOWN] != nil {
			tg.tileConfigurations[pos.x][pos.y+1] = neighbours[DOWN]
		}

		tg.tileConfigurations[pos.x][pos.y] = []Tile{selectedTile}
		tg.positionsCollapsed[pos.x][pos.y] = true
		return true

	}

	return false
}

// Returns the position with the lowest entropy (lease possible positions)
// Chooses randomly if multiple tiles have the same least possible positions
// Nil if all positions have been decided, and so no positions to be collapsed
func (tg tileGrid) tileWithLowestEntropy() *position {
	lowestOptions := -1
	possiblePos := make([]position, 0)

	// First find the lowest value of a tile
	for row := range tg.tileConfigurations {
		for col := range tg.tileConfigurations[row] {
			tileCollapsed := tg.positionsCollapsed[row][col]
			if tileCollapsed {
				continue
			}

			if len(tg.tileConfigurations[row][col]) < lowestOptions || lowestOptions == -1 {
				lowestOptions = len(tg.tileConfigurations[row][col])
			}
		}
	}

	// No tiles to be collapsed left
	if lowestOptions == -1 {
		return nil
	}

	// Now append options that have the same lowest tile
	for row := range tg.tileConfigurations {
		for col := range tg.tileConfigurations[row] {
			tileCollapsed := tg.positionsCollapsed[row][col]
			if tileCollapsed {
				continue
			}

			if len(tg.tileConfigurations[row][col]) == lowestOptions {
				possiblePos = append(possiblePos, position{row, col})
			}
		}
	}

	res := possiblePos[rand.Intn(len(possiblePos))]
	return &res
}

// Returns a grid of IDs corresponding to the initial tileset
// IDs can then be used to render the given tile in the correct position
func (tg tileGrid) getTileIds() [][]int {
	tileIds := make([][]int, len(tg.tileConfigurations))
	for row := range tg.tileConfigurations {
		tileIds[row] = make([]int, len(tg.tileConfigurations[row]))
		for col, tile := range tg.tileConfigurations[row] {
			if !tg.positionsCollapsed[row][col] {
				panic(fmt.Errorf("not finished resolving, tile with pos %v not resolved", position{row, col}))
			}

			tileIds[row][col] = tile[0].Id
		}
	}
	return tileIds
}

// Returns the possible tile configurations at a position
// Nil if position is out of range, or position has already been collapsed
func (tg tileGrid) getTileConfig(pos position) []Tile {
	if pos.x < 0 || pos.x >= len(tg.tileConfigurations) {
		return nil
	}

	if pos.y < 0 || pos.y >= len(tg.tileConfigurations[0]) {
		return nil
	}

	if tg.positionsCollapsed[pos.x][pos.y] {
		return nil
	}

	res := make([]Tile, len(tg.tileConfigurations[pos.x][pos.y]))
	copy(res, tg.tileConfigurations[pos.x][pos.y])
	return res
}

// Returns tile at a given position
// Nil if position out of range, or the position hasn't been decided yet
func (tg tileGrid) getTileId(pos position) *Tile {
	if pos.x < 0 || pos.x >= len(tg.tileConfigurations) {
		return nil
	}

	if pos.y < 0 || pos.y >= len(tg.tileConfigurations[0]) {
		return nil
	}

	if !tg.positionsCollapsed[pos.x][pos.y] {
		return nil
	}

	tile := tg.tileConfigurations[pos.x][pos.y][0]
	res := Tile{
		tile.Id,
		tile.Configuration,
	}

	return &res
}

// Updates the tileConfig for a position
// Marks the position as undecided, meaning it still needs to be collapsed
func (tg tileGrid) updateTileConfig(pos position, tileConfig []Tile) {
	if pos.x < 0 || pos.x >= len(tg.tileConfigurations) {
		panic(fmt.Errorf("attempt to update out of bounds position %v,", pos))
	}

	if pos.y < 0 || pos.y >= len(tg.tileConfigurations[0]) {
		panic(fmt.Errorf("attempt to update out of bounds position %v,", pos))
	}

	tg.positionsCollapsed[pos.x][pos.y] = false
	tg.tileConfigurations[pos.x][pos.y] = tileConfig
}
