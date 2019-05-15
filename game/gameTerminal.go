package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

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

func generateExits(mapArray *[][]string) {

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

func renderBuilding(mapData [][]string, people []person, salidas [][]int) {
	for i, row := range mapData {
		for j, column := range row {
			p := false
			s := false
			for _, person := range people {
				if i == person.Position[0] && j == person.Position[1] {
					p = true
				}
			}

			for pos := 0; pos < len(salidas); pos++ {
				if i == salidas[pos][0] && j == salidas[pos][1] {
					s = true
				}
			}

			if p {
				fmt.Print("🕴")
			} else if s {
				fmt.Print("|")
			} else {
				fmt.Print(column)
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

	nPeople := 100

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

	renderBuilding(mapData, people, nil)
}
