package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	defaultWidth        = 101 // Largeur par défaut du labyrinthe (doit être impair)
	defaultHeight       = 101 // Hauteur par défaut du labyrinthe (doit être impair)
	defaultCellSize     = 5   // Taille de cellule par défaut pour un petit labyrinthe
	defaultWindowMargin = 150 // Marge pour ajuster la fenêtre pour l'affichage
)

var (
	width        = defaultWidth
	height       = defaultHeight
	cellSize     = defaultCellSize
	windowWidth  int
	windowHeight int
	maze         [defaultHeight][defaultWidth]int
	visited      [defaultHeight][defaultWidth]bool
	startX       = 1
	startY       = 1
	endX         = width - 2
	endY         = height - 2
	currentPath  [][]int
	solvedPaths  map[string][][]int
	bfsTime      time.Duration
	dfsTime      time.Duration
	astarTime    time.Duration
)

type Node struct {
	x, y       int
	cost, heur float64
	index      int // Index dans le heap
	parent     *Node
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost+pq[i].heur < pq[j].cost+pq[j].heur
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

// Fonction pour ajuster la taille des cellules en fonction de la taille du labyrinthe
func adjustCellSize() {
	// Réduire la taille des cellules à mesure que la taille du labyrinthe augmente
	cellSize = defaultCellSize

	// Diminue légèrement la taille des cellules à mesure que la taille du labyrinthe augmente
	// Par exemple, pour chaque 100 unités supplémentaires dans width ou height, diminue de 1 pixel
	totalSize := width + height

	// Réduction très légère de la taille des cellules
	if totalSize > 100 {
		reduction := totalSize / 200 // Réduction d'un pixel tous les 100 de largeur/hauteur
		cellSize = defaultCellSize - reduction
		if cellSize < 1 {
			cellSize = 1 // Ne pas aller en dessous de 1 pixel
		}
	}

	// Recalculer les dimensions de la fenêtre
	windowWidth = width*cellSize + 100
	windowHeight = height*cellSize + defaultWindowMargin
}

// Initialiser le labyrinthe avec des murs
func initMaze() {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			maze[i][j] = 1 // 1 représente un mur
		}
	}
}

// Générer le labyrinthe en utilisant l'algorithme de Prim
func generateMaze(x, y int, rnd *rand.Rand) {
	directions := [][]int{
		{0, 2}, {2, 0}, {0, -2}, {-2, 0},
	}

	rnd.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	maze[y][x] = 0 // 0 représente un chemin

	for _, dir := range directions {
		newX, newY := x+dir[0], y+dir[1]
		if newY > 0 && newY < height-1 && newX > 0 && newX < width-1 && maze[newY][newX] == 1 {
			maze[y+dir[1]/2][x+dir[0]/2] = 0
			generateMaze(newX, newY, rnd)
		}
	}
}

// Fonction BFS pour résoudre le labyrinthe
func bfs(startX, startY, endX, endY int) ([][]int, bool) {
	queue := [][]int{{startX, startY}}
	visited[startY][startX] = true
	parent := make(map[string][]int)
	path := [][]int{}

	directions := [][]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current[0] == endX && current[1] == endY {
			path = append(path, []int{current[0], current[1]}) // Ajouter la destination au chemin
			for {
				if current[0] == startX && current[1] == startY {
					break
				}
				path = append(path, current)
				current = parent[fmt.Sprintf("%d,%d", current[0], current[1])]
			}
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			solvedPaths["BFS"] = path
			return path, true
		}

		for _, dir := range directions {
			newX, newY := current[0]+dir[0], current[1]+dir[1]
			if newX >= 0 && newY >= 0 && newX < width && newY < height && maze[newY][newX] == 0 && !visited[newY][newX] {
				visited[newY][newX] = true
				queue = append(queue, []int{newX, newY})
				parent[fmt.Sprintf("%d,%d", newX, newY)] = current
			}
		}
	}
	return nil, false
}

// Fonction DFS modifiée pour résoudre le labyrinthe
func dfs(x, y, endX, endY int) ([][]int, bool) {
	if x < 0 || x >= width || y < 0 || y >= height || maze[y][x] == 1 || visited[y][x] {
		return nil, false
	}

	visited[y][x] = true
	path := [][]int{{x, y}}

	if x == endX && y == endY {
		return path, true
	}

	directions := [][]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}

	for _, dir := range directions {
		nextPath, found := dfs(x+dir[0], y+dir[1], endX, endY)
		if found {
			solvedPaths["DFS"] = append(path, nextPath...)
			return append(path, nextPath...), true
		}
	}

	return nil, false
}

// Fonction A* pour résoudre le labyrinthe
func astar(startX, startY, endX, endY int) ([][]int, bool) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	startNode := &Node{x: startX, y: startY, cost: 0, heur: heuristic(startX, startY, endX, endY), parent: nil}
	heap.Push(&pq, startNode)
	visited[startY][startX] = true
	parent := make(map[string]*Node)

	directions := [][]int{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Node)

		if current.x == endX && current.y == endY {
			path := buildPath(current)
			solvedPaths["A*"] = path
			return path, true
		}

		for _, dir := range directions {
			newX, newY := current.x+dir[0], current.y+dir[1]
			if newX >= 0 && newY >= 0 && newX < width && newY < height && maze[newY][newX] == 0 && !visited[newY][newX] {
				visited[newY][newX] = true
				newCost := current.cost + 1
				newNode := &Node{x: newX, y: newY, cost: newCost, heur: heuristic(newX, newY, endX, endY), parent: current}
				parent[fmt.Sprintf("%d,%d", newX, newY)] = newNode
				heap.Push(&pq, newNode)
			}
		}
	}
	return nil, false
}

// Heuristique pour A* (distance de Manhattan)
func heuristic(x1, y1, x2, y2 int) float64 {
	return float64(abs(x1-x2) + abs(y1-y2))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Construire le chemin à partir du nœud final
func buildPath(node *Node) [][]int {
	var path [][]int
	for node != nil {
		path = append(path, []int{node.x, node.y})
		node = node.parent
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// Dessiner le labyrinthe
func drawMaze(offsetX int32, offsetY int32, paths [][]int, color rl.Color) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if maze[y][x] == 1 {
				rl.DrawRectangle(int32(x*cellSize)+offsetX, int32(y*cellSize)+offsetY, int32(cellSize), int32(cellSize), rl.Black)
			}
		}
	}
	for _, p := range paths {
		rl.DrawRectangle(int32(p[0]*cellSize)+offsetX, int32(p[1]*cellSize)+offsetY, int32(cellSize), int32(cellSize), color)
	}
}

// Réinitialiser le labyrinthe et les chemins
func reset() {
	initMaze()
	for i := range visited {
		for j := range visited[i] {
			visited[i][j] = false
		}
	}
	currentPath = nil
	solvedPaths = make(map[string][][]int) // Réinitialise les chemins mémorisés
}

func resetPath() {
	for i := range visited {
		for j := range visited[i] {
			rl.DrawRectangle(int32(i*cellSize), int32(j*cellSize), int32(cellSize), int32(cellSize), rl.White)
			visited[i][j] = false
		}
	}
	currentPath = nil
}

func main() {
	adjustCellSize() // Appel de la fonction pour ajuster la taille des cellules

	rl.InitWindow(int32(windowWidth), int32(windowHeight), "Générateur de labyrinthe")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	reset()
	generateMaze(startX, startY, rnd)
	chosen_path := 0

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			reset()
			generateMaze(startX, startY, rnd)
			bfsTime, dfsTime, astarTime = 0, 0, 0
		}

		if rl.IsKeyPressed(rl.KeyOne) {
			chosen_path = 1
			if savedPath, ok := solvedPaths["BFS"]; ok {
				currentPath = savedPath
			} else {
				resetPath()
				startTime := time.Now()
				currentPath, _ = bfs(startX, startY, endX, endY)
				bfsTime = time.Since(startTime)
			}
		}
		if rl.IsKeyPressed(rl.KeyTwo) {
			chosen_path = 2
			if savedPath, ok := solvedPaths["DFS"]; ok {
				currentPath = savedPath
			} else {
				resetPath()
				startTime := time.Now()
				currentPath, _ = dfs(startX, startY, endX, endY)
				dfsTime = time.Since(startTime)
			}
		}
		if rl.IsKeyPressed(rl.KeyThree) {
			chosen_path = 3
			if savedPath, ok := solvedPaths["A*"]; ok {
				currentPath = savedPath
			} else {
				resetPath()
				startTime := time.Now()
				currentPath, _ = astar(startX, startY, endX, endY)
				astarTime = time.Since(startTime)
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Affichage du labyrinthe
		switch chosen_path {
		case 1:
			drawMaze(0, 0, currentPath, rl.Green)
		case 2:
			drawMaze(0, 0, currentPath, rl.Blue)
		case 3:
			drawMaze(0, 0, currentPath, rl.Red)
		default:
			drawMaze(0, 0, nil, rl.White)
		}

		// Afficher les instructions et le temps
		rl.DrawText("Appuyez sur 1 pour BFS, 2 pour DFS, 3 pour A*", 10, int32(height*cellSize+10), 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Temps BFS: %v", bfsTime), 10, int32(height*cellSize+40), 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Temps DFS: %v", dfsTime), 10, int32(height*cellSize+70), 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Temps A*: %v", astarTime), 10, int32(height*cellSize+100), 20, rl.Black)

		rl.EndDrawing()
	}
}
