package game

import (
	"os"
	"bufio"
)

type GameUI interface {
	Draw(*Level)
}

type Tile rune

const (
	StoneHall Tile = '#'
	DirtFloor Tile = '.'
	Door Tile = '|'
	Blank Tile = ' ' 	
)

type Level struct {
	Map [][]Tile
}

func LoadLevelFromFile(filename string) *Level {
	file, err := os.Open("game/maps/map1.map")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mapLines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		mapLines = append(mapLines, scanner.Text())
		if len(mapLines[index]) > longestRow {
			longestRow = len(mapLines[index])
		} 
		index++
	}

	level := &Level{}
	level.Map = make([][]Tile, len(mapLines))
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := mapLines[y]
		for x,c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneHall 
			case '|':
				t = Door
			case '.':
				t = DirtFloor
			default: 
				panic("invalid character in map")
			} 
			level.Map[y][x] = t
		}
	}

	return level
}

func Run(ui GameUI) {
	level := LoadLevelFromFile("front/maps/map1.map")
	ui.Draw(level)
}