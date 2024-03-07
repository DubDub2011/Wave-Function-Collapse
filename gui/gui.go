package gui

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path"
	"wavefunctioncollapse/wfc"

	"github.com/hajimehoshi/ebiten"
)

type ConfigTile struct {
	Name        string         `json:"name"`
	Connections map[int]string `json:"connections"`
}

type tileImage struct {
	img *ebiten.Image // image to output
}

type Simulation struct {
	tileImages                 map[int]*tileImage
	tileSet                    []wfc.Tile
	result                     *[][]int
	width, height              int
	aspectRatioX, aspectRatioY int
	screenWidth, screenHeight  int
}

func RunSimulation(tileDir string, width, height int) {
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowTitle("Wave function collapse")

	configPath := path.Join(tileDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to read file %s with err %v", configPath, err))
	}

	var config []ConfigTile
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data with err %v", err))
	}

	tiles := make(map[int]*tileImage, len(config))
	tileSet := make([]wfc.Tile, 0, len(config))
	for tileIdx, tile := range config {
		id := tileIdx
		conn := tile.Connections
		tileSet = append(tileSet, wfc.Tile{Id: id, Configuration: conn})

		imgPath := path.Join(tileDir, tile.Name)
		imgReader, err := os.Open(imgPath)
		if err != nil {
			panic(fmt.Errorf("failed to open image %s with error %v", imgPath, err))
		}
		defer imgReader.Close()
		img, _, err := image.Decode(imgReader)
		if err != nil {
			panic(fmt.Errorf("failed to decode image %s with error %v", imgPath, err))
		}

		ebitenImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		if err != nil {
			panic(fmt.Errorf("failed to convert to ebiten image, path %s with error %v", imgPath, err))
		}

		tiles[id] = &tileImage{
			ebitenImg,
		}
	}

	res := wfc.Collapse(tileSet, width, height)
	sim := Simulation{
		tileImages:   tiles,
		tileSet:      tileSet,
		width:        width,
		height:       height,
		result:       &res,
		aspectRatioX: 16, aspectRatioY: 9,
		screenWidth: 1280, screenHeight: 720,
	}

	if err := ebiten.RunGame(sim); err != nil {
		panic(err)
	}
}

func (g Simulation) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		res := wfc.Collapse(g.tileSet, g.width, g.height)
		*g.result = res
	}

	return nil
}

func (sim Simulation) Draw(screen *ebiten.Image) {
	imgIds := *sim.result
	for row := range imgIds {
		for col := range imgIds[row] {
			img := sim.tileImages[imgIds[row][col]]
			imgWidth, imgHeight := img.img.Size()

			imgOptions := ebiten.DrawImageOptions{}
			tileLen := float64(sim.screenWidth / sim.width)
			tileWid := float64(sim.screenHeight / sim.height)
			imgOptions.GeoM.Scale(
				tileLen/float64(imgWidth),
				tileWid/float64(imgHeight))
			// due to order of rows returned, need to place them at the bottom
			imgOptions.GeoM.Translate(tileLen*float64(row), tileWid*float64(col))
			screen.DrawImage(img.img, &imgOptions)
		}
	}
}

func (sim Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return sim.screenWidth, sim.screenHeight
}
