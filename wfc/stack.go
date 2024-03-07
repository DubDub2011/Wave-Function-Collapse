package wfc

import "fmt"

// tileStack is responsible for tracking the previous tile positions and their values to support backtracking
type tileStack struct {
	pointer    int
	stackSlice []oldTile
}

// push will add a tile value onto the stack,
func (stack *tileStack) push(stackVal oldTile) {
	if len(stackVal.oldNeighbours) != 4 {
		panic(fmt.Errorf("must have four neighbours (mark edges and collapsed tile as nil)"))
	}

	// iniitalize stack
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

// pop will take a tile value off the stack
func (stack *tileStack) pop() oldTile {
	stack.pointer--

	if stack.pointer == -1 {
		panic(fmt.Errorf("stack empty, attempt to pop when no values left"))
	}

	return stack.stackSlice[stack.pointer]
}

// Tracks an old tile, takes it's position, it's possible tiles, and the neighbours possible tiles
type oldTile struct {
	pos           position       // the position of the tile in the grid
	oldTileConfig []Tile         // the tile config for the given tile
	oldNeighbours map[int][]Tile // the neighbour values BEFORE the current tile was added
}
