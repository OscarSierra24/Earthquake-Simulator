package game

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	border := getBorder(mapArray)
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
		max := 1200
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
func getFloor(mapData [][]string) [][]chan struct{} {
	var floor [][]chan struct{}
	nTokens := 2
	for _, row := range mapData {
		var tmp []chan struct{}
		for range row {
			tmpChan := make(chan struct{}, nTokens)
			//tmp_chan <- 1
			tmp = append(tmp, tmpChan)

		}
		floor = append(floor, tmp)

	}
	return floor
}

func (p person) run(path [][]int, floor [][]chan struct{}) {
	lastX, lastY := -1, -1
	//var tokenCurrent chan struct{}
	//var tokenNext chan struct{}
	for len(path) > 0 {
		x, y := path[0][0], path[0][1]
		path = path[1:]

		//Ocupy next pos
		floor[x][y] <- struct{}{}
		lastX, lastY = x, y

		//Logic
		p.Position[0] = x
		p.Position[1] = y

		//Release last pos
		<-floor[p.Position[0]][p.Position[1]]

		//Wait for next move
		time.Sleep(time.Duration(p.Speed) * time.Millisecond)
		//fmt.Println("Still here")
	}
	//Take person out of building
	*p.isInside = 0
	//End of path release token
	<-floor[lastX][lastY]
	//
	//<-floor[p.Position[0]][p.Position[1]]

}

//Get border, returns coords of the border to a given matrix
func getBorder(mapArray *[][]string) [][]int {
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
	return border
}

//Check state
//Checks and changes state of person
func checkState(exits [][]int, people *[]person) {
	for _, e := range exits {
		for _, p := range *people {
			if e[0] == p.Position[0] && e[1] == p.Position[1] {
				*p.isInside = 0
			}
		}
	}

}

//Show stats about who is in and out
func showStats(people []person) {
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
}

//Start ...
func Start() {
	//Setup
	mapFile := "game/maps/out.map"

	//Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Number of people (< 300): ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	nPeople, err := strconv.Atoi(text)

	if err != nil {
		fmt.Println("Not a number")
		os.Exit(-1)
	}
	fmt.Print("Number of exits (< 50) : ")
	text, _ = reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	nExits, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("Not a number")
		os.Exit(-1)
	}

	fmt.Print("Time (seconds) for people to run : ")
	text, _ = reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	runningTime, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("Not a number")
		os.Exit(-1)
	}

	//nPeople = 200
	//nExits := 3

	//Texture map
	texture := map[string]string{
		//Wall
		"#": "■",
		//Floor
		".": " ",
		//Door
		"|": "|",
		//Path taken,
		"+": "🔺",
	}

	//People skins
	skins := []string{
		"👮‍", "👩", "👨‍", "👶", "👨",
	}

	//Default
	//texture = map[string]string{
	//	//Wall
	//	"#": "#",
	//	//Floor
	//	".": ".",
	//	//Door
	//	"|": "|",
	//	//Path taken,
	//	"+": "+",
	//}
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
	//for _, p := range people {
	//	floor[p.Position[0]][p.Position[1]] <- 1
	//}

	//Generate the exits for the people
	exits := generateExits(nExits, &mapData)
	//fmt.Println(exits)

	//Tell everyone to start running
	for _, p := range people {
		path := pathfinding.BFS(
			p.Position[0],
			p.Position[1],
			&mapData,
			WALL,
			FLOOR,
			DOOR,
		)
		//time.Sleep(1 * time.Second)
		go p.run(path, floor)
	}

	running := true

	startTime := time.Now()

	//Timer
	go func() {
		t := time.NewTimer(time.Duration(runningTime) * time.Second)

		<-t.C
		running = false
	}()

	//Print render every 250ms while timer is up
	for running {
		clear()
		checkState(exits, &people)
		elapsedTime := time.Now().Sub(startTime).Seconds()
		timeLeft := math.Floor((float64(runningTime)-elapsedTime)*100) / 100

		fmt.Println(timeLeft, "Seconds left for building to collapse")
		renderBuilding(mapData, people, texture, skins)
		showStats(people)

		time.Sleep(250 * time.Millisecond)

	}

	//Final stats
	//Count number of survivors
	var survivors int
	for _, p := range people {
		if *p.isInside == 0 {
			survivors++
		}
	}

	fmt.Println(nPeople, "were inside the building,", survivors, "were able to get out")

	//Debug stuff
	//for _, p := range people {
	//fmt.Println(p.Position, *p.isInside, p.skin)
	//}
}
