package dijkstra

import (
	"github.com/JesseleDuran/osm-graph-parser/coordinates"
	"github.com/JesseleDuran/osm-graph-parser/graph"
	"github.com/JesseleDuran/osm-graph-parser/graph/shortest_path"
	"github.com/JesseleDuran/osm-graph-parser/graph/shortest_path/dijkstra/heap"
	"math"
)

type DijkstraV2 struct {
	Graph graph.GraphV2
}

const INFINITE = math.MaxInt64

type PreviousV2 map[int64]int64
type PathWeightV2 map[int64]float64

func (d DijkstraV2) FromCoordinates(origin, destiny coordinates.Coordinates) shortest_path.Response {
	originNode := d.Graph.FindNode(origin)
	destinyNode := d.Graph.FindNode(destiny)
	weight, prev := d.FromNodeIDs(originNode, destinyNode)
	return shortest_path.Response{
		//Steps:       StepsV2(pathV2(originCell.Cid, destinyCell.Cid, prev), d.Graph),
		TotalWeight: weight[destinyNode],
		Polyline:    d.pathPolylineV2(d.Graph.Nodes[originNode], d.Graph.Nodes[destinyNode], prev),
	}
}

func (d DijkstraV2) FromNodeIDs(start, end int64) (PathWeightV2, PreviousV2) {
	//maps from each node to the total weight of the total shortest path.
	pathWeight := make(PathWeightV2, 0)

	//maps from each node to the previous node in the "current" shortest path.
	previous := make(PreviousV2, 0)

	remaining := heap.Create()
	// insert first node id the PQ, the start node.
	remaining.Insert(heap.Node{Value: start, Cost: 0})

	// initialize pathWeight all to infinite value.
	for _, v := range d.Graph.Nodes {
		pathWeight[v.ID] = INFINITE
	}
	//start node distance to itself is 0.
	pathWeight[start] = 0

	//the previous node does not exists
	previous[start] = INFINITE

	visit := make(map[int64]bool, 0)

	//while the PQ is not empty.
	for !remaining.IsEmpty() {
		// extract the min value of the PQ.
		min, _ := remaining.Min()
		nei := d.Graph.Nodes[min.Value]
		visit[nei.OsmID] = true
		remaining.DeleteMin()
		if min.Value == end {
			return pathWeight, previous
		}

		// if the node has edges, the loop through it.
		if v, ok := d.Graph.Nodes[min.Value]; ok {

			//change to normal for
			for nodeNeighborID, e := range v.Edges {
				nei := d.Graph.Nodes[nodeNeighborID]

				if visit[nei.OsmID] {
					continue //change to negative condition
				}
				visit[nei.OsmID] = true

				// the value is the one of the current node plus the weight(a, neighbor)
				currentPathValue := pathWeight[min.Value] + e.Weight

				if currentPathValue < pathWeight[nodeNeighborID] {
					pathWeight[nodeNeighborID] = currentPathValue
					previous[nodeNeighborID] = min.Value
				}
				remaining.Insert(heap.Node{Value: nodeNeighborID, Cost: currentPathValue})
			}
		}
	}
	return pathWeight, previous
}

//key : end, value: prev
func (d DijkstraV2) pathPolylineV2(start, end graph.Node, previous PreviousV2) [][3]float64 {
	result := make([][3]float64, 0)
	result = append(result, [3]float64{
		end.Location.LatLng().Lng.Degrees(),
		end.Location.LatLng().Lat.Degrees(),
	})
	var prev int64
	_, startOk := previous[start.ID]
	_, endOk := previous[end.ID]
	if !startOk && !endOk {
		return result
	}

	for prev != start.ID {
		prev = previous[end.ID]
		node := d.Graph.Nodes[prev]
		result = append(result, [3]float64{
			node.Location.LatLng().Lng.Degrees(),
			node.Location.LatLng().Lat.Degrees(),
			float64(node.OsmID),
		})
		end = node
		//log.Println(prev, end)
	}

	resultSorted := make([][3]float64, len(result))
	j := 0
	for i := len(result) - 1; i >= 0; i-- {
		resultSorted[j] = result[i]
		j++
	}
	return resultSorted
}
