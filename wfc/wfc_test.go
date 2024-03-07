package wfc

import (
	"reflect"
	"testing"
)

func Test_Collapse_Simple(t *testing.T) {
	testCases := []struct {
		name          string
		tileSet       []Tile
		width, height int
	}{
		{
			"Small grid, two tile",
			[]Tile{
				{1, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
				{2, map[int]string{LEFT: "BBB", UP: "AAA", RIGHT: "BBB", DOWN: "AAA"}},
				{3, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "BBB", DOWN: "AAA"}},
				{4, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "BBB"}},
			},
			10, 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function from your package
			Collapse(tc.tileSet, tc.width, tc.height)
		})
	}

}

func Test_tileGrid_tileWithLowestEntropy(t *testing.T) {
	tg := newTileGrid(2, 2, []Tile{
		{1, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
		{2, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
	})

	tg.tileConfigurations[1][1] = []Tile{}

	expected := position{1, 1}
	pos := tg.tileWithLowestEntropy()

	if expected != *pos {
		t.Errorf("Failed, expected %v, got %v", expected, pos)
	}
}

func Test_tileGrid_collapseTile_success(t *testing.T) {
	tg := newTileGrid(2, 2, []Tile{
		{1, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
		{2, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
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

	expected := Tile{1, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}}
	if !reflect.DeepEqual(expected, tg.tileConfigurations[pos.x][pos.y][0]) {
		t.Errorf("Failed, expected %v, got %v", expected, tg.tileConfigurations[pos.x][pos.y])
	}
}

func Test_tileGrid_collapseTile_neighboursUpdate(t *testing.T) {
	tile1 := Tile{1, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "BBB"}}
	tile2 := Tile{1, map[int]string{LEFT: "AAA", UP: "BBB", RIGHT: "AAA", DOWN: "AAA"}}
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

	collapsedTile := tg.tileConfigurations[pos.x][pos.y][0]
	neighbourBelow := tg.tileConfigurations[pos.x][pos.y+1][0]

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
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{1, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			true,
		},
		{
			"Tile with up and tile with matching down, should match",
			UP,
			Tile{0, map[int]string{LEFT: "AAA", UP: "BBB", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "BBB"}},
			true,
		},
		{
			"Tile with up and tile without matching down, should not match",
			UP,
			Tile{0, map[int]string{LEFT: "AAA", UP: "BBB", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			false,
		},
		{
			"Tile with up and tile without matching down, should not match",
			UP,
			Tile{0, map[int]string{LEFT: "AAA", UP: "BBB", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "CCC"}},
			false,
		},
		{
			"Full tiles, up, should match",
			UP,
			Tile{0, map[int]string{LEFT: "BBB", UP: "BBB", RIGHT: "BBB", DOWN: "BBB"}},
			Tile{0, map[int]string{LEFT: "BBB", UP: "BBB", RIGHT: "BBB", DOWN: "BBB"}},
			true,
		},
		{
			"Full tiles, left, should match",
			LEFT,
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			true,
		},
		{
			"Full tiles, right, should match",
			RIGHT,
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			true,
		},
		{
			"Full tiles, down, should match",
			DOWN,
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"}},
			true,
		},
		{
			"Assymetric tiles, down, should match",
			DOWN,
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAB"}},
			Tile{0, map[int]string{LEFT: "AAA", UP: "AAB", RIGHT: "AAA", DOWN: "AAA"}},
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
