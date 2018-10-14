package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	cellSize = 8
	w, h     = 200, 200
)

var (
	colorBG    = pixel.RGB(0.2, 0.2, 0.2)
	colorCells = pixel.RGB(0.1, 0.1, 0.1)
)

func newMap() [][]int {
	m := make([][]int, h)
	for i := 0; i < h; i++ {
		m[i] = make([]int, w)
	}
	return m
}

func getSizes(fig [][]int) (int, int) {
	var (
		h = len(fig)
		w int
	)
	if h > 0 {
		w = len(fig[0])
	}
	return w, h
}

func placeFigure(dest, fig [][]int, x, y int) ([][]int, error) {
	figW, figH := getSizes(fig)
	if figW*figH == 0 {
		return nil, errors.New("width or height of figure is zero")
	}

	result := make([][]int, len(dest))
	for i := 0; i < len(dest); i++ {
		result[i] = make([]int, len(dest[0]))
		for j := 0; j < len(dest[0]); j++ {
			result[i][j] = dest[i][j]
		}
	}

	for i := y; i < y+figH; i++ {
		for j := x; j < x+figW; j++ {
			result[i][j] = fig[i-y][j-x]
		}
	}
	return result, nil
}

func drawWorld(target pixel.Target, world [][]int) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Aquamarine
	for y := 0; y < len(world); y++ {
		for x := 0; x < len(world[0]); x++ {
			if world[y][x] > 0 {
				var (
					x = float64(x+1) * cellSize
					y = float64(y+1) * cellSize
				)

				imd.Push(
					pixel.V(x, y),
					pixel.V(x+cellSize-1, y+cellSize-1),
				)
				imd.Rectangle(0)
			}
		}
	}
	imd.Draw(target)
}

func countNeighbours(world [][]int, x, y int) int {
	var (
		minX  = x - 1
		minY  = y - 1
		maxX  = x + 1
		maxY  = y + 1
		count = 0
	)

	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}

	w, h := getSizes(world)
	if maxX > w-1 {
		maxX = w - 1
	}
	if maxY > h-1 {
		maxY = h - 1
	}

	for j := minY; j <= maxY; j++ {
		for i := minX; i <= maxX; i++ {
			if i == x && j == y {
				continue
			}
			if world[j][i] > 0 {
				count++
			}
		}
	}
	return count
}

func updateWorld(world [][]int) [][]int {
	newWorld := make([][]int, h)
	for y := 0; y < len(world); y++ {
		newWorld[y] = make([]int, w)
		for x := 0; x < len(world[0]); x++ {
			count := countNeighbours(world, x, y)
			lifeness := world[y][x]
			if world[y][x] > 0 {
				if count < 2 || count > 3 {
					lifeness = 0
				}
			} else if count == 3 {
				lifeness = 1
			}
			newWorld[y][x] = lifeness
		}
	}
	return newWorld
}

func randomPopulation(world [][]int) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if rand.Float64() > 0.5 {
				world[j][i] = 1
			}
		}
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, w*cellSize, h*cellSize),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var (
		frames int
		second = time.NewTicker(time.Second)
	)

	world := newMap()
	randomPopulation(world)
	glider := [][]int{
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
	}
	world, err = placeFigure(world, glider, 50, 50)
	world, err = placeFigure(world, glider, 10, 50)
	world, err = placeFigure(world, glider, 50, 10)
	if err != nil {
		log.Fatal(err)
	}

	imd := imdraw.New(nil)
	imd.Color = colorCells
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			var (
				x = float64(i) * cellSize
				y = float64(j) * cellSize
			)
			imd.Push(
				pixel.V(x, y),
				pixel.V(x+cellSize, y+cellSize),
			)
			imd.Rectangle(1)
		}
	}

	for !win.Closed() {
		world = updateWorld(world)
		win.Clear(colorBG)
		imd.Draw(win)
		drawWorld(win, world)
		win.Update()

		time.Sleep(time.Millisecond * 10)
		select {
		case <-second.C:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
			frames++
		}
	}
}

func main() {
	pixelgl.Run(run)
}
