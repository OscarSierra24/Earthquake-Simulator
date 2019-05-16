package pathfinding

func validSquare(x int, y int, matrix *[][]string, visited *[][]bool, WALL string) bool {
	if x < 0 || x >= len(*matrix) || y < 0 || y >= len((*matrix)[0]) {
		return false
	}
	return !(*visited)[x][y] && (*matrix)[x][y] != WALL
}

//BFS searches in a matrix for a path
//x,y are starting coords
//wall is obstacle
//floor is were the path can be drawn
//goal is were the path can end
func BFS(x int, y int, matrix *[][]string, WALL string, FLOOR string, GOAL string) [][]int {
	var visited [][]bool
	for _, row := range *matrix {
		var tmp []bool
		for range row {
			tmp = append(tmp, false)
		}
		visited = append(visited, tmp)
	}
	parents := make(map[int][]int, len(*matrix)*len((*matrix)[0]))
	//var parents map[int][]int

	stack := [][]int{[]int{x, y}}

	goalX, goalY := -1, -1

	for len(stack) > 0 {
		//Extract the coords
		i, j := stack[0][0], stack[0][1]

		if (*matrix)[i][j] == GOAL {
			goalX, goalY = i, j
		}

		//neighbours
		up := []int{i + 1, j}
		down := []int{i + -1, j}
		left := []int{i, j - 1}
		right := []int{i, j + 1}
		//ADVANCED
		upright := []int{i + 1, j + 1}
		upleft := []int{i + 1, j - 1}
		downleft := []int{i - 1, j - 1}
		downright := []int{i - 1, j + 2}

		var neighbours [][]int
		neighbours = append(neighbours, right)
		neighbours = append(neighbours, left)
		neighbours = append(neighbours, up)
		neighbours = append(neighbours, down)
		neighbours = append(neighbours, upright)
		neighbours = append(neighbours, upleft)
		neighbours = append(neighbours, downleft)
		neighbours = append(neighbours, downright)
		//Shuffle
		//rand.Seed(time.Now().UTC().UnixNano())
		//rand.Shuffle(len(neighbours), func(i, j int) {
		//	neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
		//})

		//Test for valid up, down, left, right neighbours
		for _, n := range neighbours {
			tmpX, tmpY := n[0], n[1]
			if validSquare(tmpX, tmpY, matrix, &visited, WALL) {
				visited[tmpX][tmpY] = true
				stack = append(stack, n)

				//Generate unique key for map
				key := n[0]*len((*matrix)[0]) + n[1]
				parents[key] = []int{i, j}
			}
		}
		stack = stack[1:]

	}
	var path [][]int
	for {

		path = append(path, []int{goalX, goalY})
		key := goalX*len((*matrix)[0]) + goalY
		goalX, goalY = parents[key][0], parents[key][1]
		//fmt.Println(goalX, goalY)
		if goalX == x {
			if goalY == y {
				break
			}
		}
	}

	//The other one is reversed
	var actualPath [][]int
	for i := range path {
		actualPath = append(actualPath, path[len(path)-i-1])
	}
	return actualPath
}
