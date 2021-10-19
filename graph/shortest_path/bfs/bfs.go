package bfs

import (
	"container/list"
	"github.com/JesseleDuran/osm-graph-parser/graph"
	"log"
)

type BFS struct {
	Graph graph.GraphV2
}

func (g BFS) Path(start, end int64) []int64 {
	result := make([]int64, 0)
	visited := make(map[int64]bool)
	queue := list.New()
	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		qnode := queue.Front()
		queue.Remove(qnode)
		id := qnode.Value.(int64)
		log.Println(id)
		if id == end {
			return result
		}
		for k, _ := range g.Graph.Nodes[id].Edges {
			if _, ok := visited[k]; ok {
				continue
			}
			visited[k] = true
			queue.PushBack(k)
			result = append(result, g.Graph.Nodes[k].ID)
		}
	}
	return result
}
