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
	return level
}

func Run(ui GameUI) {
	level := LoadLevelFromFile("front/maps/map1.map")
	ui.Draw(level)
}