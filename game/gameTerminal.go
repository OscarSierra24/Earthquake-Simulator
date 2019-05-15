package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const WALL string = "#"
const FLOOR string = "."
const DOOR string = "|"

type point struct {
	X, Y int
}

type person struct {
	//Position in building
	Position []int
	//Walking speed
	Speed float32
	//Reference to map
	mapData *[][]string
}

func loadLevelFromFile(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//var map_data [][]string

	scanner := bufio.NewScanner(file)
	mapData := make([][]string, 0)
	longestRow := 0
	index := 0

	for scanner.Scan() {
		var tmpLine []string
		for _, c := range scanner.Text() {
			tmpLine = append(tmpLine, string(c))
		}
		mapData = append(mapData, tmpLine)
		if len(mapData[index]) > longestRow {
			longestRow = len(mapData[index])
		}
		index++
	}
	return mapData
}

//adds exits to mapArray and returns the coords as an array
func generateExits(nExits int, mapArray *[][]string) [][]int {
	var border [][]int

	for i, row := range *mapArray {
		for j := range row {
			//Checks if coords are in border
			if i == 0 || i == len(*mapArray)-1 || j == 0 || j == len(row)-1 {
				border = append(
					border,
					[]int{i, j},
				)
			}

		}
	}
	var i int
	var exitArray [][]int
	for i < nExits {

		rand.Shuffle(len(border), func(i, j int) {
			border[i], border[j] = border[j], border[i]
		})

		x, y := border[0][0], border[0][1]

		//Tests if exit is valid or not
		isValid := true

		//Upper border
		if x == 0 {
			if (*mapArray)[x+1][y] == WALL {
				isValid = false
			}
		}
		//Lower border
		if x == len(*mapArray)-1 {
			if (*mapArray)[x-1][y] == WALL {
				isValid = false
			}
		}
		//Left border
		if y == 0 {
			if (*mapArray)[x][y+1] == WALL {
				isValid = false
			}
		}
		//Right border
		if y == len((*mapArray)[0])-1 {
			if (*mapArray)[x][y-1] == WALL {
				isValid = false
			}
		}
		if isValid {
			(*mapArray)[x][y] = DOOR
			exitArray = append(exitArray, []int{x, y})
			i++

		}

		border = border[1:]
	}
	return exitArray
}

//Returns an arrray with available positions to walk at
func getPositions(mapArray [][]string) [][]int {
	var positions [][]int

	for i, row := range mapArray {
		for j, column := range row {
			if column == "." {
				var tmpCoord []int
				tmpCoord = append(tmpCoord, i)
				tmpCoord = append(tmpCoord, j)
				positions = append(
					positions,
					tmpCoord)
			}

		}

	}
	return positions
}

func generatePeople(nPeople int, mapArray *[][]string, positions [][]int) []person {
	var people []person

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	positions = positions[:nPeople]

	for _, point := range positions {
		people = append(
			people,
			person{
				point,
				rand.Float32() * 10,
				mapArray,
			},
		)
	}

	return people
}

func clear() {
	fmt.Println("\033[2J")
}

func renderBuilding(mapData [][]string, people []person, textures map[string]string, skins []string) {
	for i, row := range mapData {
		for j, column := range row {
			p := false
			for _, person := range people {
				if i == person.Position[0] && j == person.Position[1] {
					p = true
				}
			}
			if p {
				rand.Shuffle(len(skins), func(i, j int) {
					skins[i], skins[j] = skins[j], skins[i]
				})
				fmt.Print(skins[0])
			} else {
				fmt.Print(

					textures[column],
				)
			}

		}
		fmt.Println()

	}
}

// Returns [][]  for chan(int,1 ) of available floor
func getFloor(mapData [][]string) [][]chan (int) {
	var floor [][]chan (int)
	nTokens := 2
	for _, row := range mapData {
		var tmp []chan (int)
		for range row {
			tmpChan := make(chan int, nTokens)
			//tmp_chan <- 1
			tmp = append(tmp, tmpChan)

		}
		floor = append(floor, tmp)

	}
	return floor
}

func (p person) run(path [][]int, floor [][]chan (int)) {
	for len(path) > 0 {
		time.Sleep(1 * time.Second)
		y, x := path[0][0], path[0][1]
		path = path[1:]
		//fmt.Println("Waiting for token")
		//Ocupy next pos
		floor[x][y] <- 1
		//fmt.Println("Release token")
		//Release last pos
		<-floor[p.Position[0]][p.Position[1]]
		//Logic
		p.Position[0] = x
		p.Position[1] = y

	}
}

//Start ...
func Start() {
	mapFile := "game/maps/map1.map"

	nPeople := 20
	nExits := 5

	//Building data as a 2d array
	mapData := loadLevelFromFile(mapFile)

	//Available positions as an array of [x,y]
	positions := getPositions(mapData)

	//Array of struct of people
	people := generatePeople(nPeople, &mapData, positions)

	//Floor in which one can be
	floor := getFloor(mapData)

	//Populate floor
	for _, p := range people {
		floor[p.Position[0]][p.Position[1]] <- 1
	}

	//Generate the exits for the people
	exits := generateExits(nExits, &mapData)

	//Texture map
	texture := map[string]string{
		//Floor
		"#": "◽️",
		//Wall
		".": " ",
		//Door
		"|": "🚪",
	}

	//People skins
	skins := []string{
		"👮‍", "👩", "👨‍", "👶", "👨",
	}

	fmt.Println(exits)
	renderBuilding(mapData, people, texture, skins)
}
