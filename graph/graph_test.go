package graph

import (
	"log"
	"testing"
)

func TestBridge(t *testing.T) {
	g := createGraph()
	b := g.Bridge([]int64{1, 2, 3}, []int64{3, 4, 5, 6})
	log.Println(b)
}

func createGraph() GraphV2 {
	node1 := Node{
		ID:    1,
		Edges: make(Edges, 0),
	}
	node1.Edges[2] = Edge{Weight: 1}
	node1.Edges[3] = Edge{Weight: 1}

	node2 := Node{
		ID:    2,
		Edges: make(Edges, 0),
	}
	node2.Edges[1] = Edge{Weight: 1}
	node2.Edges[3] = Edge{Weight: 1}

	node3 := Node{
		ID:    3,
		Edges: make(Edges, 0),
	}
	node3.Edges[1] = Edge{Weight: 1}
	node3.Edges[2] = Edge{Weight: 1}
	node3.Edges[4] = Edge{Weight: 1}

	node4 := Node{
		ID:    4,
		Edges: make(Edges, 0),
	}
	node4.Edges[3] = Edge{Weight: 1}
	node4.Edges[5] = Edge{Weight: 1}
	node4.Edges[6] = Edge{Weight: 1}

	node5 := Node{
		ID:    5,
		Edges: make(Edges, 0),
	}
	node5.Edges[4] = Edge{Weight: 1}
	node5.Edges[6] = Edge{Weight: 1}

	node6 := Node{
		ID:    6,
		Edges: make(Edges, 0),
	}
	node6.Edges[4] = Edge{Weight: 1}
	node6.Edges[5] = Edge{Weight: 1}
	nodes := make(Nodes, 0)
	nodes[1] = node1
	nodes[2] = node2
	nodes[3] = node3
	nodes[4] = node4
	nodes[5] = node5
	nodes[6] = node6
	return GraphV2{Nodes: nodes}
}
