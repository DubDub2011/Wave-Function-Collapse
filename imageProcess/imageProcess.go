package imageprocess

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path"
	"reflect"
	"wavefunctioncollapse/wfc"

	"github.com/disintegration/imaging"
)

type ConfigTile struct {
	Name        string         `json:"name"`
	Connections map[int]string `json:"connections"`
}

var directory string

func ProcessDir(dirPath string) {
	directory = dirPath
	configPath := path.Join(dirPath, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to read file %s with err %v", configPath, err))
	}

	var config []ConfigTile
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Errorf("failed to marshal data with err %v", err))
	}

	for _, tile := range config {
		config = appendIfNonDuplicate(config, mutateImage(tile, "R"))
		config = appendIfNonDuplicate(config, mutateImage(tile, "RR"))
		config = appendIfNonDuplicate(config, mutateImage(tile, "RRR"))
		config = appendIfNonDuplicate(config, mutateImage(tile, "F"))
		config = appendIfNonDuplicate(config, mutateImage(tile, "FR"))
		config = appendIfNonDuplicate(config, mutateImage(tile, "FRR"))
		config = appendIfNonDuplicate(config, mutateImage(tile, "FRRR"))
	}

	data, err = json.Marshal(config)
	if err != nil {
		panic(err)
	}
	os.WriteFile(configPath, data, 0777)
}

func appendIfNonDuplicate(tiles []ConfigTile, toAppend ConfigTile) []ConfigTile {
	valid := true
	for _, tile := range tiles {
		if reflect.DeepEqual(tile.Connections, toAppend.Connections) {
			valid = false
		}
	}

	if !valid {
		// remove file if not valid
		os.Remove(path.Join(directory, toAppend.Name))
		return tiles
	}

	tiles = append(tiles, toAppend)
	return tiles
}

func mutateImage(conf ConfigTile, op string) ConfigTile {
	imgPath := path.Join(directory, conf.Name)
	imgReader, err := os.Open(imgPath)
	if err != nil {
		panic(fmt.Errorf("failed to open image %s with error %v", imgPath, err))
	}
	defer imgReader.Close()
	img, _, err := image.Decode(imgReader)
	if err != nil {
		panic(fmt.Errorf("failed to decode image %s with error %v", imgPath, err))
	}

	for _, char := range op {
		switch char {
		case 'R':
			// god knows why, but the rotation is counter-clockwise
			img = imaging.Rotate270(img)

			conf.Connections = map[int]string{
				wfc.LEFT:  conf.Connections[wfc.DOWN],
				wfc.UP:    conf.Connections[wfc.LEFT],
				wfc.RIGHT: conf.Connections[wfc.UP],
				wfc.DOWN:  conf.Connections[wfc.RIGHT],
			}
		case 'F':
			img = imaging.FlipH(img)
			img = imaging.FlipV(img)

			conf.Connections = map[int]string{
				wfc.LEFT:  conf.Connections[wfc.RIGHT],
				wfc.UP:    conf.Connections[wfc.DOWN],
				wfc.RIGHT: conf.Connections[wfc.LEFT],
				wfc.DOWN:  conf.Connections[wfc.UP],
			}
		default:
			panic(fmt.Errorf("unsupported char %c", char))
		}
	}

	newName := conf.Name[:len(conf.Name)-len(path.Ext(conf.Name))]
	newName = newName + "-" + op + ".png"
	newPath := path.Join(directory, newName)

	imaging.Save(img, newPath)
	conf.Name = newName

	return conf
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
