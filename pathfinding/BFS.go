package pathfinding

import "fmt"

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
func BFS(x int, y int, matrix *[][]string, WALL string, FLOOR string, GOAL string) {
	var visited [][]bool
	for _, row := range *matrix {
		var tmp []bool
		for range row {
			tmp = append(tmp, false)
		}
		visited = append(visited, tmp)
	}
	var parents [][][]int
	for _, row := range *matrix {
		var tmp [][]int
		for range row {
			tmp = append(tmp, []int{-1, -1})
		}
		parents = append(parents, tmp)
	}

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

		var neighbours [][]int
		neighbours = append(neighbours, up)
		neighbours = append(neighbours, down)
		neighbours = append(neighbours, left)
		neighbours = append(neighbours, right)

		//Test for valid up, down, left, right neighbours
		for _, n := range neighbours {
			tmpX, tmpY := n[0], n[1]
			if validSquare(tmpX, tmpY, matrix, &visited, WALL) {
				visited[tmpX][tmpY] = true
				stack = append(stack, n)
				parents[tmpX][tmpY] = []int{i, j}
			}
		}
		stack = stack[1:]

	}
	var path [][]int
	/*
		for _, row := range visited {
			for _, col := range row {
				if col {
					fmt.Print(0)
				} else {
					fmt.Print(1)
				}

			}
			fmt.Println()
		}
	*/
	for goalX != x && goalY != y {
		path = append(path, []int{goalX, goalY})
		goalX, goalY = parents[goalX][goalY][0], parents[goalX][goalY][1]
		fmt.Println(goalX, goalY)
	}
	for i := range parents {
		fmt.Println(parents[i])
	}

	fmt.Println(goalX, goalY, "<- Result")

}
