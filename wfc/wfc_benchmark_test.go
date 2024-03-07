package wfc

import (
	"fmt"
	"testing"
)

// BenchmarkCollapse benchmarks the Collapse function
func BenchmarkCollapse(b *testing.B) {
	tileset5 := generateTileSet(5)
	tileset10 := generateTileSet(10)

	scenarios := []struct {
		width, height int
		tileset       []Tile
	}{
		{20, 20, tileset5},
		{50, 50, tileset5},
		{100, 100, tileset5},
		{20, 20, tileset10},
		{50, 50, tileset10},
	}

	for _, scenario := range scenarios {
		b.Run(
			fmt.Sprintf("Width%d_Height%d_TileSet%d", scenario.width, scenario.height, len(scenario.tileset)),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					Collapse(scenario.tileset, scenario.width, scenario.height)
				}
			},
		)
	}
}

// generateTileSet generates a tile set with the specified height
func generateTileSet(height int) []Tile {
	tileSet := make([]Tile, height)
	for i := 0; i < height; i++ {
		tileSet[i] = Tile{
			Id:            i,
			Configuration: map[int]string{LEFT: "AAA", UP: "AAA", RIGHT: "AAA", DOWN: "AAA"},
		}
	}
	return tileSet
}
