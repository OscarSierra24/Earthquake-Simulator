# Earthquake-Simulator

## Setup

clone this repository under your gopath

    mkdir $GOPATH/src/github.com/OscarSierra24/
    cd $GOPATH/src/github.com/OscarSierra24/
    git clone [this repo]


## Run

    go run main.go

Under the source folder 

     $GOPATH/src/github.com/OscarSierra24/

### How it works

It takes a file .map which contains an MxN matrix which is represented like this:

    "#" represents a wall
    "." represents movable path ie, where one can stand

Once the matrix is loaded, the program asks user for number of people / number of exits / time to run

The people is randomly placed in the map and so are the exits

Each person is represented by a letter [A-Za-z] and the program will create a visualization of a top down view of a 1 floor building, each person moves at a different speed so not everyone can make it if the time is tight 


Authors:

- OscarSierra24
- Dcrdn
- QuirinoC
