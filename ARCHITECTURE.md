# ARCHITECTURE

Earthquake-Simulator.go

External libraries

## METHODS
### main.go
#### main() ()
executes the game package start method.

### gameTErminal.go
#### LoadLevelFromFile(string) ([][] string)
Opens the file, read the map and returns a bidimensional array that contains tokens that represent objects

#### generateSalidas(int) ([][]int) 
Creates random exit doors, it receives the amount of exit doors that it should generate and returns a bidimensional
array that contains x and y positions of those doors

#### get_positions([][]string) ([][]int)
Receives the map and returns an arrray with available positions to walk at

#### generatePeople(int, *[][]string, [][]int) ([]Person)
Creates people represeted in the map as P, it receives an int, which is the amount of people that 
has to be created, the map and the available positions

#### clear() ()
Clears the terminal

#### render_building([][]string, []Person, [][]int) ()
As its names says, it renders the building, by printing the tokens,
it receives the map data, the people existing in the moment, and the exit doors

#### get_floor() ()
Returns [][]  for chan(int,1 ) of available floor

#### run() ()
is a method for the struct Person, it moves the person towards the next position in which he should be

#### Start() ()
setups everything