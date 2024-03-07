# Wave Function Collapse

Solution for https://www.youtube.com/watch?v=rI_y2GAlQFM

Enjoyed this one a lot, and that's due to the output, really fun seeing a grid of tiles organized randomly but in the correct orientation, feels a bit like magic.

Key takeaways:
- Good organized, annotated code is ESSENTIAL for implementing complex algorithms
- As suspected this algorithm can be pretty slow. Go is reliable as always, I've made some small improvements where I can to speed things up using the profiling tools which gave big gains. Further improvement can be made, but I'm happy with the current performance. Additionally I'm not currently using go routines, I can probably get them added in without too much fuss, but slightly nervous on if they'll help or not, as the WFC algorithm is sequential in nature, would need further investigation.
- Not a massive fan of ebiten, very forced way of handling things, poor ease of use, and lacking some functionality, but gets the job done as a 2D engine.


## Setup
You'll need Go to build the app, see installation instructions [here]("https://go.dev/") 

Go will install missing dependencies when you run `go build` or `go test`, see dependencies below:
-  [hajimehoshi/ebiten](https://github.com/hajimehoshi/ebiten)
- [disintegration/imaging](https://github.com/disintegration/imaging)  

To get started, run command `go run main.go -width=48 -height=27 -directory="assets"`
- `-width=<width>`, width of the tile grid
- `-height=<height>`, height of the tile grid
- `-directory="<path>"`, path of the directory containing the tileset

Click on the screen to regenerate a new tileset

Custom tilesets are supported, these need to be defined with a config file, see inside of `/assets/config.json` for an example
- `/assets/circuit` already exists but without rotated tiles, adding tilesets manually is a slow process. By passing the flag `-process=<path>` on the main command, it'll run the image processor against it. This will create rotated assets and update the config to reflect the new assets.
- We don't want to run this flag against the directory twice however, will start to panic, but as this is a helper app, I've not gone deeper into a fix.

## Future improvements

Happy with what I've got done, understand WFC a lot better now, but definitely more to delve into around the theory behind it. This example is really amazing:
https://github.com/mxgmn/WaveFunctionCollapse

A clearer interface would be nice, could do with a menu that allows you to configure tiles, select an area, then choose the possible orientations.

Performance improvements, I notice a lot of time is spent on finding the position with the lowest entropy, not too sure why as that part is quite simple, I could probably look the better way up, but it's one I'd like to solve myself, feels like a classic leetcode problem.

