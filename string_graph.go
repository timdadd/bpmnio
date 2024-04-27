package bpmnio

type stringGraph struct {
	edges    map[string][]string
	vertices []string
}

func (g *stringGraph) addEdge(u, v string) {
	g.edges[u] = append(g.edges[u], v)
}

func (g *stringGraph) topologicalSortUtil(v string, visited map[string]bool, stack *[]string) {
	visited[v] = true

	for _, u := range g.edges[v] {
		if !visited[u] {
			g.topologicalSortUtil(u, visited, stack)
		}
	}

	*stack = append([]string{v}, *stack...)
}

func (g *stringGraph) topologicalSort() []string {
	var stack []string
	visited := make(map[string]bool)

	for _, v := range g.vertices {
		if !visited[v] {
			g.topologicalSortUtil(v, visited, &stack)
		}
	}

	return stack
}

// DAG --> Directed Acyclic Graph is a directed graph with no directed cycles
//                 (A)----------\
//                 /|\           \
//         (B)<---+ | +-->(C)     \
//            \     |     / |      \
//             \    v    /  |      /
//              +->(D)<-+   |     /
//                    \     |    /
//                     +-->(E)<-+
// Nodes are known as vertices whilst links are known as edges
// A vertex (v) is said to be reachable from another vertex (u) where a path starts at u and ends at v
