package wfc

import (
	"reflect"
	"testing"
)

func Test_Collapse_Simple(t *testing.T) {
	testCases := []struct {
		name          string
		tileSet       []Tile
		width, length int
	}{
		{
			"Small grid, two tile",
			[]Tile{
				{1, []int{0, 0, 0, 0}},
				{2, []int{1, 0, 1, 0}},
				{3, []int{0, 0, 1, 0}},
				{4, []int{0, 0, 0, 1}},
			},
			10, 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function from your package
			result, err := Collapse(tc.tileSet, tc.width, tc.length)
			if err != nil {
				// Check if the result matches the expected value
				t.Errorf("Failed, expected no error, got %v", err)
				t.Logf("Result set: %v", result)
			}
		})
	}

}

func Test_tileGrid_tileWithLowestEntropy(t *testing.T) {
	tg := newTileGrid(2, 2, []Tile{
		{1, []int{0, 0, 0, 0}},
		{2, []int{0, 0, 0, 0}},
	})

	tg.tiles[1][1] = []Tile{}

	expected := position{1, 1}
	pos := tg.tileWithLowestEntropy()

	if expected != *pos {
		t.Errorf("Failed, expected %v, got %v", expected, pos)
	}
}

func Test_tileGrid_collapseTile_success(t *testing.T) {
	tg := newTileGrid(2, 2, []Tile{
		{1, []int{0, 0, 0, 0}},
		{1, []int{0, 0, 0, 0}},
	})

	pos := position{0, 0}
	success := tg.collapseTile(pos)
	if !success {
		t.Errorf("Failed, expected %v, got %v", true, success)
	}

	tileMarkedAsCollapsed := tg.positionsCollapsed[pos.x][pos.y]
	if !tileMarkedAsCollapsed {
		t.Errorf("Failed, expected %v, got %v", true, success)
	}

	expected := Tile{1, []int{0, 0, 0, 0}}
	if !reflect.DeepEqual(expected, tg.tiles[pos.x][pos.y][0]) {
		t.Errorf("Failed, expected %v, got %v", expected, tg.tiles[pos.x][pos.y])
	}
}

func Test_tileGrid_collapseTile_neighboursUpdate(t *testing.T) {
	tile1 := Tile{1, []int{0, 0, 0, 1}}
	tile2 := Tile{1, []int{0, 1, 0, 0}}
	tg := newTileGrid(2, 2, []Tile{
		tile1,
		tile2,
	})

	pos := position{0, 0}
	success := tg.collapseTile(pos)
	if !success {
		t.Errorf("Failed, expected %v, got %v", true, success)
	}

	tileMarkedAsCollapsed := tg.positionsCollapsed[pos.x][pos.y]
	if !tileMarkedAsCollapsed {
		t.Errorf("Failed, expected %v, got %v", true, success)
	}

	collapsedTile := tg.tiles[pos.x][pos.y][0]
	neighbourBelow := tg.tiles[pos.x][pos.y+1][0]

	var expectedBelow Tile
	if reflect.DeepEqual(collapsedTile, tile1) {
		expectedBelow = tile2
	} else {
		expectedBelow = tile1
	}

	if !reflect.DeepEqual(expectedBelow, neighbourBelow) {
		t.Errorf("Failed, expected %v, got %v", expectedBelow, neighbourBelow)
	}
}

func Test_match(t *testing.T) {
	testCases := []struct {
		name      string
		direction int
		tile1     Tile
		tile2     Tile
		expected  bool
	}{
		{
			"Empty tiles, should match",
			UP,
			Tile{0, []int{0, 0, 0, 0}},
			Tile{0, []int{0, 0, 0, 0}},
			true,
		},
		{
			"Tile with up and tile with matching down, should match",
			UP,
			Tile{0, []int{0, 1, 0, 0}},
			Tile{0, []int{0, 0, 0, 1}},
			true,
		},
		{
			"Tile with up and tile without matching down, should not match",
			UP,
			Tile{0, []int{0, 1, 0, 0}},
			Tile{0, []int{0, 0, 0, 0}},
			false,
		},
		{
			"Tile with up and tile without matching down, should not match",
			UP,
			Tile{0, []int{0, 1, 0, 0}},
			Tile{0, []int{0, 0, 0, 2}},
			false,
		},
		{
			"Full tiles, up, should match",
			UP,
			Tile{0, []int{1, 1, 1, 1}},
			Tile{0, []int{1, 1, 1, 1}},
			true,
		},
		{
			"Full tiles, left, should match",
			LEFT,
			Tile{0, []int{1, 1, 1, 1}},
			Tile{0, []int{1, 1, 1, 1}},
			true,
		},
		{
			"Full tiles, right, should match",
			RIGHT,
			Tile{0, []int{1, 1, 1, 1}},
			Tile{0, []int{1, 1, 1, 1}},
			true,
		},
		{
			"Full tiles, down, should match",
			DOWN,
			Tile{0, []int{1, 1, 1, 1}},
			Tile{0, []int{1, 1, 1, 1}},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := match(tc.direction, tc.tile1, tc.tile2)
			if tc.expected != res {
				t.Errorf("Failed, expected %v, got %v", tc.expected, res)
			}
		})
	}
}
