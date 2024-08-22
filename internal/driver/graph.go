package driver

type Graph struct {
	adjacency map[string][]string
	vertices  []string
}

// Function to initialize a new graph.
func NewGraph() *Graph {
	return &Graph{
		adjacency: make(map[string][]string),
		vertices:  []string{},
	}
}

// Function to add a vertex to the graph.
func (g *Graph) AddVertex(vertex string) {
	if g.adjacency[vertex] == nil {
		g.vertices = append(g.vertices, vertex)
		g.adjacency[vertex] = []string{}
	}
}

// Function to add a directed edge between two vertices.
func (g *Graph) AddEdge(from, to string) {
	g.AddVertex(from)
	g.AddVertex(to)
	g.adjacency[from] = append(g.adjacency[from], to)
}

// Function to check for cycles using DFS.
func (g *Graph) hasCycle() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, v := range g.vertices {
		if g.detectCycleDFS(v, visited, recStack) {
			return true
		}
	}

	return false
}

func (g *Graph) DFS(startVertex string) []string {
	visited := make(map[string]bool)
	result := []string{}

	g.dfsRecursive(startVertex, visited, &result)

	return result
}

func (g *Graph) dfsRecursive(vertex string, visited map[string]bool, result *[]string) {
	if visited[vertex] {
		return
	}

	visited[vertex] = true

	*result = append(*result, vertex)

	for _, neighbor := range g.adjacency[vertex] {
		g.dfsRecursive(neighbor, visited, result)
	}
}

func (g *Graph) detectCycleDFS(vertex string, visited, recStack map[string]bool) bool {
	if !visited[vertex] {
		visited[vertex] = true
		recStack[vertex] = true

		for _, neighbor := range g.adjacency[vertex] {
			if !visited[neighbor] && g.detectCycleDFS(neighbor, visited, recStack) {
				return true
			} else if recStack[neighbor] {
				return true
			}
		}
	}

	recStack[vertex] = false

	return false
}
