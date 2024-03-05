package main

import (
	"flag"
	_ "image/png"
	"log"
	"os"
	"runtime/pprof"
	"wavefunctioncollapse/gui"
	imageprocess "wavefunctioncollapse/imageProcess"
)

var (
	width  = flag.Int("width", 32, "width of grid to collapse")
	height = flag.Int("height", 18, "height of grid to collapse")
	dir    = flag.String("directory", "", "directory of tiles with config to run against")

	process = flag.String("process", "", "directory of tiles to process ")

	cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if dir == nil && process == nil {
		log.Fatalf("Require process or dir flag to be passed")
	}

	if *process != "" {
		imageprocess.ProcessDir(*process)
	}

	if *dir != "" {
		gui.RunSimulation(*dir, *width, *height)
	}
}
