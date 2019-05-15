package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Point struct {
	X, Y int
}

type Level struct {
	Map    [][]string
	Player [][]Person
}

type Person struct {
	//Position in building
	Position []int
	//Walking speed
	Speed float32
	//Reference to map
	Map_data *[][]string
}

func LoadLevelFromFile(filename string) [][]string {
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
		var tmp_line []string
		for _, c := range scanner.Text() {
			tmp_line = append(tmp_line, string(c))
		}
		mapData = append(mapData, tmp_line)
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
func get_positions(mapArray [][]string) [][]int {
	var positions [][]int

	for i, row := range mapArray {
		for j, column := range row {
			if column == "." {
				var tmp_coord []int
				tmp_coord = append(tmp_coord, i)
				tmp_coord = append(tmp_coord, j)
				positions = append(
					positions,
					tmp_coord)
			}

		}

	}
	return positions
}

func generatePeople(nPeople int, mapArray *[][]string, positions [][]int) []Person {
	var people []Person

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(positions), func(i, j int) { positions[i], positions[j] = positions[j], positions[i] })

	positions = positions[:nPeople]

	for _, point := range positions {
		people = append(
			people,
			Person{
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

func render_building(map_data [][]string, people []Person, salidas [][]int) {
	for i, row := range map_data {
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
func get_floor(map_data [][]string) [][]chan (int) {
	var floor [][]chan (int)
	n_tokens := 2
	for _, row := range map_data {
		var tmp []chan (int)
		for range row {
			tmp_chan := make(chan int, n_tokens)
			//tmp_chan <- 1
			tmp = append(tmp, tmp_chan)

		}
		floor = append(floor, tmp)

	}
	return floor
}

func (p Person) run(path [][]int, floor [][]chan (int)) {
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

func Start() {
	map_file := "game/maps/map1.map"

	nPeople := 100

	//Building data as a 2d array
	map_data := LoadLevelFromFile(map_file)

	//Available positions as an array of [x,y]
	positions := get_positions(map_data)

	//Array of struct of people
	people := generatePeople(nPeople, &map_data, positions)
	//salidas := generateSalidas(20)

	//Floor in which one can be
	floor := get_floor(map_data)

	//Populate floor
	for _, p := range people {
		floor[p.Position[0]][p.Position[1]] <- 1
	}

}

/*
package game

import (
	"os"
	"bufio"
	"fmt"
	"math/rand"
	"time"
	"github.com/OscarSierra24/Earthquake-Simulator/goastar"
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
	Pending Tile = -1
)

type Entity struct {
	X, Y int
}

type Player struct {
	Entity
}

func (p *Player) MoveUp() {
	p.Y--
}

func (p *Player) MoveDown() {
	p.Y++
}

func (p *Player) MoveLeft() {
	p.X--
}

func (p *Player) MoveRight() {
	p.X++
}

type Level struct {
	Map [][]Tile
	Player Player
}

func esPersona(personas [][]int, x int, y int) int{
	for i :=0;i<len(personas);i++{
		if(personas[i][0]==x && personas[i][1]==y){
			fmt.Printf("persona: %d %d", x,y)
			return 1
		}
	}
	return 0
}

func esSalida(salidas [][]int, x int, y int) int{
	for i :=0;i<len(salidas);i++{
		if(salidas[i][0]==x && salidas[i][1]==y){
			fmt.Printf("salida: %d %d", x,y)
			return 1
		}
	}
	return 0
}


func LoadLevelFromFile(filename string, salidas [][]int, personas [][]int) *Level {
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
			if(esSalida(salidas,y,x)==1){
				fmt.Printf("es salida %d %d   ",y,x)
				fmt.Println(" ")
				level.Map[y][x] = Door
				continue
			}
			if(esPersona(personas,y,x)==1){

				fmt.Printf("es persona %d %d   ", y,x)
				fmt.Println(" ")
				level.Map[y][x] = Pending

				goastar.GetPath(y,x,salidas[0][0],salidas[1][1])

				level.Player.X = x
				level.Player.Y = y

				continue
			}
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
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				t = Pending

			default:
				panic("invalid character in map")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
				SearchLoop:
				for searchX := x-1; searchX <= x+1; searchX++ {
					for searchY := y-1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchY][searchY]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level
}

func generateSalidas(salidas int) [][]int{
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	pared:=make([]int, 0)
	pared = append(pared, 1,2,3,4) // 1 izq, 2 der, 3 arriba, 4 abajo

	fila := make([]int, 0)
	fila = append(fila,
		1,2,3,4,6,7,8,9,11,12,13,14,15,17,18,19,20,21)

	fila1 := make([]int, 0)
	fila1 = append(fila1,
		1,2,4,5,6,8,9,10,11,12,13,14,15,16,17,18,19,20,21)

	columna := make([]int, 0)
	columna = append(columna,
		1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38)

	columna1 := make([]int, 0)
	columna1 = append(columna1,
		1,2,3,4,5,6,7,8,10,11,12,13,15,16,17,18,19,20,21,22,23,24,25,27,28,29,30,31,32,33,34,35,36,37,38)


	groups := [][]int{}

	for i :=0;i<salidas;i++{
		rnd:=r.Intn(len(pared))
		wall:=pared[rnd]

		if(wall==1){
			pos:=r.Intn(len(fila))
			f:=fila[pos]
			arr := []int{f,0}
			groups=append(groups, arr)
		}
		if(wall==2){
			pos:=r.Intn(len(fila1))
			f:=fila1[pos]
			arr := []int{f,39}
			groups=append(groups, arr)
		}
		if(wall==3){
			pos:=r.Intn(len(columna))
			f:=columna[pos]
			arr := []int{0,f}
			groups=append(groups, arr)
		}
		if(wall==4){
			pos:=r.Intn(len(columna1))
			f:=columna1[pos]
			arr := []int{22,f}
			groups=append(groups, arr)
		}
	}
	return groups
}

func generatePeople(people int) [][]int {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	posiciones:=[][]int{[]int{1,1},[]int{1,2},[]int{1,3},[]int{1,4},[]int{1,5},[]int{1,6},[]int{1,7},[]int{1,8},[]int{1,9},[]int{1,10},[]int{1,11},[]int{1,12},[]int{1,13},[]int{1,14},[]int{1,15},[]int{1,16},[]int{1,17},[]int{1,18},[]int{1,19},[]int{1,20},[]int{1,21},[]int{1,23},[]int{1,24},[]int{1,25},[]int{1,26},[]int{1,27},[]int{1,28},[]int{1,29},[]int{1,30},[]int{1,31},[]int{1,32},[]int{1,33},[]int{1,34},[]int{1,35},[]int{1,36},[]int{1,37},[]int{1,38},[]int{2,1},[]int{2,2},[]int{2,3},[]int{2,4},[]int{2,5},[]int{2,6},[]int{2,7},[]int{2,8},[]int{2,9},[]int{2,10},[]int{2,11},[]int{2,12},[]int{2,13},[]int{2,14},[]int{2,15},[]int{2,16},[]int{2,17},[]int{2,18},[]int{2,19},[]int{2,20},[]int{2,21},[]int{2,23},[]int{2,24},[]int{2,25},[]int{2,26},[]int{2,28},[]int{2,29},[]int{2,30},[]int{2,31},[]int{2,32},[]int{2,33},[]int{2,35},[]int{2,36},[]int{2,37},[]int{2,38},[]int{3,1},[]int{3,2},[]int{3,3},[]int{3,4},[]int{3,5},[]int{3,6},[]int{3,7},[]int{3,8},[]int{3,9},[]int{3,10},[]int{3,11},[]int{3,12},[]int{3,13},[]int{3,14},[]int{3,15},[]int{3,16},[]int{3,17},[]int{3,18},[]int{3,19},[]int{3,20},[]int{3,21},[]int{3,28},[]int{3,29},[]int{3,30},[]int{3,31},[]int{3,32},[]int{3,33},[]int{4,1},[]int{4,2},[]int{4,3},[]int{4,4},[]int{4,5},[]int{4,6},[]int{4,7},[]int{4,8},[]int{4,9},[]int{4,10},[]int{4,11},[]int{4,12},[]int{4,13},[]int{4,14},[]int{4,15},[]int{4,16},[]int{4,17},[]int{4,18},[]int{4,19},[]int{4,20},[]int{4,21},[]int{4,23},[]int{4,24},[]int{4,25},[]int{4,26},[]int{4,28},[]int{4,29},[]int{4,30},[]int{4,31},[]int{4,32},[]int{4,33},[]int{4,35},[]int{4,36},[]int{4,37},[]int{4,38},[]int{5,10},[]int{5,11},[]int{5,12},[]int{5,23},[]int{5,24},[]int{5,25},[]int{5,26},[]int{5,27},[]int{5,28},[]int{5,29},[]int{5,30},[]int{5,31},[]int{5,32},[]int{5,33},[]int{5,34},[]int{5,35},[]int{5,36},[]int{5,37},[]int{5,38},[]int{6,1},[]int{6,2},[]int{6,3},[]int{6,4},[]int{6,5},[]int{6,6},[]int{6,7},[]int{6,8},[]int{6,9},[]int{6,10},[]int{6,11},[]int{6,12},[]int{6,13},[]int{6,14},[]int{6,15},[]int{6,16},[]int{6,17},[]int{6,18},[]int{6,19},[]int{6,20},[]int{6,21},[]int{6,23},[]int{6,24},[]int{6,25},[]int{6,26},[]int{6,28},[]int{6,29},[]int{6,30},[]int{6,31},[]int{6,32},[]int{6,33},[]int{6,35},[]int{6,36},[]int{6,37},[]int{6,38},[]int{7,1},[]int{7,2},[]int{7,3},[]int{7,4},[]int{7,5},[]int{7,6},[]int{7,7},[]int{7,8},[]int{7,9},[]int{7,10},[]int{7,11},[]int{7,12},[]int{7,13},[]int{7,14},[]int{7,15},[]int{7,16},[]int{7,17},[]int{7,18},[]int{7,19},[]int{7,20},[]int{7,21},[]int{7,28},[]int{7,29},[]int{7,30},[]int{7,31},[]int{7,32},[]int{7,33},[]int{8,1},[]int{8,2},[]int{8,3},[]int{8,4},[]int{8,5},[]int{8,6},[]int{8,7},[]int{8,8},[]int{8,9},[]int{8,10},[]int{8,11},[]int{8,12},[]int{8,13},[]int{8,14},[]int{8,15},[]int{8,16},[]int{8,17},[]int{8,18},[]int{8,19},[]int{8,20},[]int{8,21},[]int{8,23},[]int{8,24},[]int{8,25},[]int{8,26},[]int{8,27},[]int{8,28},[]int{8,29},[]int{8,30},[]int{8,31},[]int{8,32},[]int{8,33},[]int{8,34},[]int{8,35},[]int{8,36},[]int{8,37},[]int{8,38},[]int{9,1},[]int{9,2},[]int{9,3},[]int{9,4},[]int{9,5},[]int{9,6},[]int{9,7},[]int{9,8},[]int{9,9},[]int{9,10},[]int{9,11},[]int{9,12},[]int{9,13},[]int{9,14},[]int{9,15},[]int{9,16},[]int{9,17},[]int{9,18},[]int{9,19},[]int{9,20},[]int{9,21},[]int{9,23},[]int{9,24},[]int{9,25},[]int{9,26},[]int{9,28},[]int{9,29},[]int{9,30},[]int{9,31},[]int{9,32},[]int{9,33},[]int{9,35},[]int{9,36},[]int{9,37},[]int{9,38},[]int{10,10},[]int{10,11},[]int{10,12},[]int{10,29},[]int{10,32},[]int{10,38},[]int{11,1},[]int{11,2},[]int{11,3},[]int{11,4},[]int{11,5},[]int{11,6},[]int{11,7},[]int{11,8},[]int{11,9},[]int{11,10},[]int{11,11},[]int{11,12},[]int{11,13},[]int{11,14},[]int{11,15},[]int{11,17},[]int{11,18},[]int{11,19},[]int{11,20},[]int{11,21},[]int{11,22},[]int{11,23},[]int{11,25},[]int{11,26},[]int{11,27},[]int{11,28},[]int{11,29},[]int{11,30},[]int{11,31},[]int{11,32},[]int{11,33},[]int{11,34},[]int{11,35},[]int{11,36},[]int{11,37},[]int{11,38},[]int{12,1},[]int{12,2},[]int{12,3},[]int{12,4},[]int{12,5},[]int{12,6},[]int{12,7},[]int{12,8},[]int{12,9},[]int{12,10},[]int{12,11},[]int{12,12},[]int{12,13},[]int{12,14},[]int{12,15},[]int{12,16},[]int{12,17},[]int{12,18},[]int{12,19},[]int{12,20},[]int{12,21},[]int{12,22},[]int{12,23},[]int{12,24},[]int{12,25},[]int{12,26},[]int{12,27},[]int{12,28},[]int{12,29},[]int{12,30},[]int{12,31},[]int{12,32},[]int{12,33},[]int{12,34},[]int{12,35},[]int{12,36},[]int{12,37},[]int{12,38},[]int{13,1},[]int{13,2},[]int{13,3},[]int{13,4},[]int{13,5},[]int{13,6},[]int{13,7},[]int{13,8},[]int{13,9},[]int{13,10},[]int{13,11},[]int{13,12},[]int{13,13},[]int{13,14},[]int{13,15},[]int{13,16},[]int{13,17},[]int{13,18},[]int{13,19},[]int{13,20},[]int{13,22},[]int{13,23},[]int{13,24},[]int{13,25},[]int{13,26},[]int{13,27},[]int{13,28},[]int{13,29},[]int{13,30},[]int{13,31},[]int{13,32},[]int{13,33},[]int{13,34},[]int{13,35},[]int{13,36},[]int{13,37},[]int{13,38},[]int{14,1},[]int{14,2},[]int{14,3},[]int{14,4},[]int{14,5},[]int{14,6},[]int{14,7},[]int{14,8},[]int{14,9},[]int{14,10},[]int{14,11},[]int{14,12},[]int{14,13},[]int{14,14},[]int{14,15},[]int{14,16},[]int{14,17},[]int{14,18},[]int{14,19},[]int{14,20},[]int{14,21},[]int{14,22},[]int{14,23},[]int{14,24},[]int{14,25},[]int{14,26},[]int{14,27},[]int{14,28},[]int{14,29},[]int{14,30},[]int{14,31},[]int{14,32},[]int{14,33},[]int{14,34},[]int{14,35},[]int{14,36},[]int{14,37},[]int{14,38},[]int{15,1},[]int{15,2},[]int{15,3},[]int{15,4},[]int{15,5},[]int{15,6},[]int{15,7},[]int{15,8},[]int{15,9},[]int{15,10},[]int{15,11},[]int{15,12},[]int{15,13},[]int{15,14},[]int{15,15},[]int{15,16},[]int{15,17},[]int{15,18},[]int{15,19},[]int{15,20},[]int{15,21},[]int{15,22},[]int{15,23},[]int{15,24},[]int{15,25},[]int{15,26},[]int{15,27},[]int{15,28},[]int{15,29},[]int{15,30},[]int{15,31},[]int{15,32},[]int{15,33},[]int{15,34},[]int{15,35},[]int{15,36},[]int{15,37},[]int{15,38},[]int{16,2},[]int{16,3},[]int{16,15},[]int{16,16},[]int{16,17},[]int{16,18},[]int{16,19},[]int{16,20},[]int{16,21},[]int{16,22},[]int{16,23},[]int{16,24},[]int{16,25},[]int{16,36},[]int{16,37},[]int{17,1},[]int{17,2},[]int{17,3},[]int{17,4},[]int{17,5},[]int{17,6},[]int{17,7},[]int{17,8},[]int{17,10},[]int{17,11},[]int{17,12},[]int{17,13},[]int{17,15},[]int{17,16},[]int{17,17},[]int{17,18},[]int{17,19},[]int{17,20},[]int{17,21},[]int{17,22},[]int{17,23},[]int{17,24},[]int{17,25},[]int{17,27},[]int{17,28},[]int{17,29},[]int{17,30},[]int{17,31},[]int{17,32},[]int{17,33},[]int{17,34},[]int{17,35},[]int{17,36},[]int{17,37},[]int{17,38},[]int{18,1},[]int{18,2},[]int{18,3},[]int{18,4},[]int{18,5},[]int{18,6},[]int{18,7},[]int{18,8},[]int{18,9},[]int{18,10},[]int{18,11},[]int{18,12},[]int{18,13},[]int{18,15},[]int{18,16},[]int{18,17},[]int{18,18},[]int{18,19},[]int{18,20},[]int{18,21},[]int{18,22},[]int{18,23},[]int{18,24},[]int{18,25},[]int{18,27},[]int{18,28},[]int{18,29},[]int{18,30},[]int{18,31},[]int{18,32},[]int{18,33},[]int{18,34},[]int{18,35},[]int{18,36},[]int{18,37},[]int{18,38},[]int{19,1},[]int{19,2},[]int{19,3},[]int{19,4},[]int{19,5},[]int{19,6},[]int{19,7},[]int{19,8},[]int{19,15},[]int{19,16},[]int{19,17},[]int{19,18},[]int{19,19},[]int{19,20},[]int{19,21},[]int{19,22},[]int{19,23},[]int{19,24},[]int{19,25},[]int{19,27},[]int{19,28},[]int{19,29},[]int{19,30},[]int{19,31},[]int{19,32},[]int{19,33},[]int{19,34},[]int{19,35},[]int{19,36},[]int{19,37},[]int{19,38},[]int{20,1},[]int{20,2},[]int{20,3},[]int{20,4},[]int{20,5},[]int{20,6},[]int{20,7},[]int{20,8},[]int{20,9},[]int{20,10},[]int{20,11},[]int{20,12},[]int{20,13},[]int{20,15},[]int{20,16},[]int{20,17},[]int{20,18},[]int{20,19},[]int{20,20},[]int{20,21},[]int{20,22},[]int{20,23},[]int{20,24},[]int{20,25},[]int{20,27},[]int{20,28},[]int{20,29},[]int{20,30},[]int{20,31},[]int{20,32},[]int{20,33},[]int{20,34},[]int{20,35},[]int{20,36},[]int{20,37},[]int{20,38},[]int{21,1},[]int{21,2},[]int{21,3},[]int{21,4},[]int{21,5},[]int{21,6},[]int{21,7},[]int{21,8},[]int{21,10},[]int{21,11},[]int{21,12},[]int{21,13},[]int{21,15},[]int{21,16},[]int{21,17},[]int{21,18},[]int{21,19},[]int{21,20},[]int{21,21},[]int{21,22},[]int{21,23},[]int{21,24},[]int{21,25},[]int{21,27},[]int{21,28},[]int{21,29},[]int{21,30},[]int{21,31},[]int{21,32},[]int{21,33},[]int{21,34},[]int{21,35},[]int{21,36},[]int{21,37},[]int{21,38}}

	fila := make([]int, 0)
	fila = append(fila,
		1,2,3,4,5,6,7,8,9,10)

	columna := make([]int, 0)
	columna = append(columna,
		1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21)

	groups := [][]int{}

	for i :=0;i<people;i++{
		pos:=r.Intn(len(posiciones))
		f:=posiciones[pos]
		groups=append(groups, f)
	}
	return groups
}


func Run(ui GameUI, salidas int, personas int) {
	res:=generatePeople(personas)
	res2:=generateSalidas(salidas)
	fmt.Printf("pasando personas: %v    salidas: %v ", res, res2)
	level := LoadLevelFromFile("front/maps/map1.map", res2,res)
	ui.Draw(level)
}
*/
