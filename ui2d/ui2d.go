package ui2d

import (
	"fmt"
	"github.com/OscarSierra24/Earthquake-Simulator/game"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"image/png"
	"bufio"
	"strings"
	"strconv"
) 

const winWidth, winHeight = 1280, 720

var renderer *sdl.Renderer
var textureAtlas *sdl.Texture

var textureIndex map[game.Tile]sdl.Rect

func loadTextureIndex() {
	textureIndex = make(map[game.Tile]sdl.Rect)

	infile,err := os.Open("ui2d/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := game.Tile(line[0])
		xy := line[1:]
		splitXy := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXy[0]), 10, 64)
		if err != nil {
			panic(err)
		}

		y, err := strconv.ParseInt(strings.TrimSpace(splitXy[1]), 10, 64)
		if err != nil {
			panic(err)
		}	
	
		rect := sdl.Rect{int32(x * 32), int32(y * 32), 32, 32}			
		textureIndex[tileRune] = rect

	}
}

type UI2d struct { 
}

func imgFileToTexture(filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y	

	pixels := make([]byte, w*h*4)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
			
		}
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	return tex
}

func init() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	window, err := sdl.CreateWindow("game", 200, 200, 
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	textureAtlas = imgFileToTexture("ui2d/assets/tiles.png")
	loadTextureIndex()
} 	

func (ui *UI2d) Draw(level *game.Level) {

	for y, row := range level.Map {
		for x, tile := range row {
			srcRect := textureIndex[tile]
			dstRect := sdl.Rect{int32(x*32), int32(y*32), int32(32), int32(32)}
			renderer.Copy(textureAtlas, &srcRect, &dstRect)

		} 
	}

	renderer.Present()
	for{}
}
