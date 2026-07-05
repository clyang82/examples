package leetcode

func updateMatrix(mat [][]int) [][]int {
	m := len(mat)
	if m == 0 {
		return mat
	}
	n := len(mat[0])

	dist := make([][]int, m)
	for i := range dist {
		dist[i] = make([]int, n)
	}

	type point struct {
		r, c int
	}
	queue := []point{}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if mat[r][c] == 0 {
				dist[r][c] = 0
				queue = append(queue, point{r, c})
			} else {
				dist[r][c] = -1 // Use -1 to indicate unvisited
			}
		}
	}

	dirs := []point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nr, nc := curr.r+d.r, curr.c+d.c
			if nr >= 0 && nr < m && nc >= 0 && nc < n && dist[nr][nc] == -1 {
				dist[nr][nc] = dist[curr.r][curr.c] + 1
				queue = append(queue, point{nr, nc})
			}
		}
	}

	return dist
}
