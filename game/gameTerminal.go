package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/OscarSierra24/Earthquake-Simulator/pathfinding"
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
	Speed int
	//Reference to map
	mapData *[][]string
	//Skin for printing
	skin string
	//Is inside building
	isInside *int
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
		rand.Seed(time.Now().UTC().UnixNano())
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

func generatePeople(nPeople int, mapArray *[][]string, positions [][]int, skins []string) []person {
	var people []person

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	positions = positions[:nPeople]

	for _, point := range positions {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(skins), func(i, j int) {
			skins[i], skins[j] = skins[j], skins[i]
		})

		//Walking speed aleatorization
		max := 1500
		min := 100

		inside := 1

		people = append(
			people,
			person{
				point,
				(rand.Intn(max-min) + min),
				mapArray,
				skins[0],
				&inside,
			},
		)
	}

	return people
}

func clear() {
	fmt.Println("\033[2J")
	//os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
}

func renderBuilding(mapData [][]string, people []person, textures map[string]string, skins []string) {
	//for _, person := range people {
	//fmt.Print(person.isInside)
	//}
	fmt.Println()
	for i, row := range mapData {
		for j, column := range row {
			p := false
			skin := "p"
			for _, person := range people {
				if i == person.Position[0] && j == person.Position[1] {
					//fmt.Println(person.isInside)
					skin = person.skin
					p = true && *person.isInside == 1
					break
				}
			}
			if p {
				fmt.Print(skin)
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
	nTokens := 3
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
	lastX, lastY := -1, -1
	for len(path) > 0 {
		//p.isInside = false
		time.Sleep(time.Duration(p.Speed) * time.Millisecond)
		x, y := path[0][0], path[0][1]
		path = path[1:]
		//fmt.Println("Waiting for token")

		//Ocupy next pos
		floor[x][y] <- 1
		lastX, lastY = x, y

		//Logic
		p.Position[0] = x
		p.Position[1] = y

		//Release last pos
		<-floor[p.Position[0]][p.Position[1]]

	}
	//Take person out of building
	*p.isInside = 0
	//End of path release token
	<-floor[lastX][lastY]

}

//Start ...
func Start() {
	//Setup
	mapFile := "game/maps/map1.map"

	nPeople := 50
	nExits := 3

	//Texture map
	texture := map[string]string{
		//Wall
		"#": "#",
		//Floor
		".": ".",
		//Door
		"|": "🚪",
		//Path taken,
		"+": "🔺",
	}

	//People skins
	skins := []string{
		"👮‍", "👩", "👨‍", "👶", "👨",
	}

	//Default
	texture = map[string]string{
		//Wall
		"#": "#",
		//Floor
		".": ".",
		//Door
		"|": "|",
		//Path taken,
		"+": "+",
	}
	skins = []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
	//Setup

	//Building data as a 2d array
	mapData := loadLevelFromFile(mapFile)

	//Available positions as an array of [x,y]
	positions := getPositions(mapData)

	//Array of struct of people
	people := generatePeople(nPeople, &mapData, positions, skins)
	//time.Sleep(5 * time.Second)
	//Floor in which one can be
	floor := getFloor(mapData)

	//Populate floor
	for _, p := range people {
		floor[p.Position[0]][p.Position[1]] <- 1
	}

	//Generate the exits for the people
	generateExits(nExits, &mapData)

	for _, p := range people {
		path := pathfinding.BFS(
			p.Position[0],
			p.Position[1],
			&mapData,
			WALL,
			FLOOR,
			DOOR,
		)
		fmt.Println(path)
		//time.Sleep(1 * time.Second)
		go p.run(path, floor)
	}

	var i int
	//Print render every 250ms
	for {
		clear()
		renderBuilding(mapData, people, texture, skins)
		fmt.Print("People inside building: ")
		for _, p := range people {
			if *p.isInside == 1 {
				fmt.Print(p.skin, " ")
			}
		}
		fmt.Println()
		fmt.Print("People outisde building:")
		for _, p := range people {
			if *p.isInside == 0 {
				fmt.Print(p.skin, " ")
			}
		}
		fmt.Println()
		go func() {
			for {

				fmt.Printf("\rOn %d/10", i)
				i++
				time.Sleep(50 * time.Millisecond)
			}
		}()

		time.Sleep(250 * time.Millisecond)

	}

	//fmt.Println(exits, "<- Exits location")
	//path := pathfinding.BFS(1, 20, &mapData, WALL, FLOOR, DOOR)
	//for _, p := range path {
	//	(mapData)[p[0]][p[1]] = "+"
	//}
	//fmt.Println(path)
	//renderBuilding(mapData, people, texture, skins)

}
