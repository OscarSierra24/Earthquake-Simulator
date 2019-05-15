package pathfinding

//BFS searches in a matrix for a path
//x,y are starting coords
//wall is obstacle
//floor is were the path can be drawn
//goal is were the path can end
func DFS(x int, y int, matrix *[][]string, wall string, floor string, goal string) {
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
}
