package gui

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path"
	"wavefunctioncollapse/wfc"

	"github.com/hajimehoshi/ebiten"
)

type ConfigTile struct {
	Name        string `json:"name"`
	RotateBy    int    `json:"rotateBy"`
	Connections []int  `json:"connections"`
}

type tileImage struct {
	img           *ebiten.Image // image to output
	rotateDegrees int           // degrees to rotate by
}

type Simulation struct {
	tileImages                 map[int]*tileImage
	tileSet                    []wfc.Tile
	width, length              int
	result                     *[][]int
	aspectRatioX, aspectRatioY int
	screenWidth, screenHeight  int
}

func NewSimulation(tileDir string, width, length int) Simulation {
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
			tile.RotateBy,
		}
	}

	res, err := wfc.Collapse(tileSet, width, length)
	if err != nil {
		panic(fmt.Errorf("failed to run first collapse with error %v", err))
	}

	return Simulation{
		tileImages:   tiles,
		tileSet:      tileSet,
		width:        width,
		length:       length,
		result:       &res,
		aspectRatioX: 16, aspectRatioY: 9,
		screenWidth: 1280, screenHeight: 720,
	}
}

func (g Simulation) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		res, err := wfc.Collapse(g.tileSet, g.width, g.length)
		if err != nil {
			panic(fmt.Errorf("failed to run collapse with error %v", err))
		}

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
			tileLen := float64(sim.screenWidth / sim.aspectRatioX)
			tileWid := float64(sim.screenHeight / sim.aspectRatioY)
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

func RunSimulation() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Wave function collapse")
	sim := NewSimulation("assets", 16, 9)
	if err := ebiten.RunGame(sim); err != nil {
		log.Fatal(err)
	}
}
