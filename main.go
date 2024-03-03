package main

import (
	"flag"
	_ "image/png"
	"log"
	"os"
	"runtime/pprof"
	"wavefunctioncollapse/gui"
)

var (
	width  = flag.Int("width", 32, "width of grid to collapse")
	height = flag.Int("height", 18, "height of grid to collapse")
	dir    = flag.String("directory", "assets", "directory of tiles with config to run against")

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

	gui.RunSimulation(*dir, *width, *height)
}
