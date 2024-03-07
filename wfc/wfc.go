package wfc

import (
	"math/rand"
)

// Exposed function to run the collapse algorithm against a tileset
// Mainly responsible for orchestrating interal structures to run the algorithm
func Collapse(tiles []Tile, width int, height int) [][]int {
	tileGrid := newTileGrid(width, height, tiles)
	positionTracker := tileStack{}

	pos := position{
		x: rand.Intn(width),
		y: rand.Intn(height),
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
					prevTile.oldTileConfig[idx] = prevTile.oldTileConfig[len(prevTile.oldTileConfig)-1]
					prevTile.oldTileConfig = prevTile.oldTileConfig[:len(prevTile.oldTileConfig)-1]
					break
				}
			}

			// Now update grid to state prior to collapse
			tileGrid.updateTileConfig(prevTile.pos, prevTile.oldTileConfig)
			revertNeighbourFunc := func(tileConf []Tile, pos position) {
				if tileConf != nil {
					tileGrid.updateTileConfig(pos, tileConf)
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

	return tileGrid.getTileIds()
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
	Configuration map[int]string // left, up, right, down
}

func match(dir int, tile1, tile2 Tile) bool {
	// four cardinal directions, adding 2 gets to the opposite and then remainder of 4 to prevent out of range
	return tile1.Configuration[dir] == reverse(tile2.Configuration[(dir+2)%4])
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
