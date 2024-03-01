package wfc

import (
	"fmt"
	"math/rand"
)

func Collapse(tiles []Tile, width int, length int) ([][]int, error) {
	tileGrid := newTileGrid(width, length, tiles)
	positionTracker := tileStack{}

	pos := position{
		x: rand.Intn(length),
		y: rand.Intn(width),
	}
	finished := false
	for !finished {
		currTileConf := tileGrid.getTileConfig(pos)
		// Track neighbours before tiles are collapsed incase we need to backtrack
		neighbours := make(map[int][]Tile)
		neighbours[UP] = tileGrid.getTileConfig(position{pos.x, pos.y + 1})
		neighbours[RIGHT] = tileGrid.getTileConfig(position{pos.x + 1, pos.y})
		neighbours[DOWN] = tileGrid.getTileConfig(position{pos.x, pos.y - 1})
		neighbours[LEFT] = tileGrid.getTileConfig(position{pos.x - 1, pos.y})

		collapsedTile := tileGrid.collapseTile(pos)
		if !collapsedTile {
			// Tile at position could not be collapsed, need to backtrack
			prevTile := positionTracker.pop()
			// We now know the ID for the previous tile was invalid, so we'll remove it as an option
			prevTileId := tileGrid.getTileId(prevTile.pos)
			for idx, otc := range prevTile.oldTileConfig {
				if otc.Id == prevTileId.Id {
					prevTile.oldTileConfig = append(prevTile.oldTileConfig[:idx], prevTile.oldTileConfig[idx+1:]...)
				}
			}

			// Now update grid to state prior to collapse
			tileGrid.revertTileConfig(prevTile.pos, prevTile.oldTileConfig)
			revertNeighbourFunc := func(tileConf []Tile, pos position) {
				if tileConf != nil {
					tileGrid.revertTileConfig(pos, tileConf)
				}
			}
			revertNeighbourFunc(prevTile.oldNeighbours[UP], position{prevTile.pos.x, prevTile.pos.y + 1})
			revertNeighbourFunc(prevTile.oldNeighbours[RIGHT], position{prevTile.pos.x + 1, prevTile.pos.y})
			revertNeighbourFunc(prevTile.oldNeighbours[DOWN], position{prevTile.pos.x, prevTile.pos.y - 1})
			revertNeighbourFunc(prevTile.oldNeighbours[LEFT], position{prevTile.pos.x - 1, prevTile.pos.y})

			// Now collapse the tile that was reverted
			pos = prevTile.pos
			continue
		}

		trackedTile := oldTile{
			pos,
			currTileConf,
			neighbours,
		}
		positionTracker.push(trackedTile)

		// If returns nil, means no tiles left to collapse, so we're done
		nextPos := tileGrid.tileWithLowestEntropy()
		if nextPos == nil {
			finished = true
		} else {
			pos = *nextPos
		}
	}

	return tileGrid.getTileIds(), nil
}

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
	// four cardinal directions, adding 2 gets to the opposite and then remainder of 4 to prevent out of range
	return tile1.Configuration[dir] == tile2.Configuration[(dir+2)%4]
}

type tileGrid struct {
	tiles              [][][]Tile // tracks the possible tiles in a given position
	positionsCollapsed [][]bool   // tracks the positions that have been collapsed
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
			return true, func() { tg.tiles[pos.x][pos.y] = newTiles }
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

// returns position with lowest entropy, and if none can be found
func (tg tileGrid) tileWithLowestEntropy() *position {
	tileValue := -1
	possiblePos := make([]position, 0)
	for row := range tg.tiles {
		for col := range tg.tiles[row] {
			tileCollapsed := tg.positionsCollapsed[row][col]
			if tileCollapsed {
				continue
			}

			if len(tg.tiles[row][col]) < tileValue || tileValue == -1 {
				tileValue = len(tg.tiles[row][col])
				possiblePos = []position{{row, col}}
			} else if len(tg.tiles[row][col]) == tileValue {
				possiblePos = append(possiblePos, position{row, col})
			}
		}
	}

	// Means we're done, no tiles to be collapsed left
	if len(possiblePos) == 0 {
		return nil
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

	if tg.positionsCollapsed[pos.x][pos.y] {
		panic(fmt.Errorf("attempt to update already collapsed position at %v which has tileData: %v ", pos, tg.tiles[pos.x][pos.y]))
	}

	tg.positionsCollapsed[pos.x][pos.y] = false
	tg.tiles[pos.x][pos.y] = tileConfig
}

type tileStack struct {
	pointer    int
	stackSlice []oldTile
}

func (stack *tileStack) push(stackVal oldTile) {
	if len(stackVal.oldNeighbours) != 4 {
		panic(fmt.Errorf("must have four neighbours (mark edges and collapsed tile as nil)"))
	}

	if stack.stackSlice == nil {
		stack.stackSlice = make([]oldTile, 0, 1)
	}

	if len(stack.stackSlice) <= stack.pointer {
		stack.stackSlice = append(stack.stackSlice, stackVal)
	} else {
		stack.stackSlice[stack.pointer] = stackVal
	}

	stack.pointer++
}

func (stack *tileStack) pop() oldTile {
	stack.pointer--

	if stack.pointer == -1 {
		panic(fmt.Errorf("stack empty, attempt to pop when no values left"))
	}

	return stack.stackSlice[stack.pointer]
}

type oldTile struct {
	pos           position       // the position of the tile in the grid
	oldTileConfig []Tile         // the tile config for the given tile
	oldNeighbours map[int][]Tile // the neighbour values BEFORE the current tile was added
}
