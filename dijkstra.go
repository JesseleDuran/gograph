package gograph

import (
	"github.com/JesseleDuran/gograph/bitset"
	"github.com/JesseleDuran/gograph/heap"
	"github.com/golang/geo/s2"
	geojson "github.com/paulmach/go.geojson"
	"math"
)

const INFINITE = math.MaxFloat32

type Previous map[int32]int32

// Dijkstra is a traditional dijkstra with some conditionals. It can ignore nodes, so it would not pass through those nodes
// in the search. It also has a max distance to reach, so if it exceeds that value, the search will stop.
func (g Graph) Dijkstra(s ShortestPathCriteria) float32 {
	source, target, pMax := s.From, s.To, s.MaxCost
	dist := make(Distances, 0)
	visited := make(map[int32]bool, 0)

	// Source node distance to itself is 0.
	dist[source] = 0

	pq := heap.Create()
	// Insert first node id in the PQ, the source node.
	pq.Insert(heap.Node{Value: source, Cost: 0, Depth: 0})

	for !pq.IsEmpty() {
		min, _ := pq.Min()
		visited[min.Value] = true
		pq.DeleteMin()

		// The max path value was found, so we return the last highest distance found.
		if pMax > 0 && dist.Cost(min.Value) > pMax {
			return dist.Cost(min.Value)
		}

		if min.Value == target {
			return dist.Cost(target)
		}

		for _, e := range g.OutgoingEdges[min.Value] {
			// Validate if we can relax the edge related to the possible ignored node ID.
			if !(g.Nodes[e.ID].Compressed) && !visited[e.ID] {
				// Relax edge.
				currentPathValue := dist.Cost(min.Value) + e.Weight
				if currentPathValue < dist.Cost(e.ID) {
					dist[e.ID] = currentPathValue
					pq.Insert(heap.Node{Value: e.ID, Cost: currentPathValue, Depth: min.Depth + 1})
				}
			}
		}
	}
	return dist.Cost(target)
}

func (g Graph) DijkstraPathCoord(source, target Coordinate) (float32, geojson.FeatureCollection, []int) {
	from, initialCost := g.ProjectCoordinate(source)
	to, _ := g.ProjectCoordinate(target)
	d, path, data := g.DijkstraPath(ShortestPathCriteria{
		From:        from,
		InitialCost: initialCost,
		To:          to,
	})
	fc := geojson.NewFeatureCollection()
	fc.AddFeature(geojson.NewLineStringFeature(path))
	return d, *fc, data
}

func (g Graph) DijkstraPath(s ShortestPathCriteria) (float32, [][]float64, []int) {
	source, target := s.From, s.To
	dist := make(Distances, 0)
	if source < 0 || target < 0 {
		return dist.Cost(target), [][]float64{}, []int{}
	}
	visited := bitset.NewBigInt()
	dataResult := make([]int, 0)
	previous := make(Previous, 0)

	// Source node distance to itself is 0.
	dist[source] = 0

	pq := heap.Create()
	// Insert first node id in the PQ, the source node.
	pq.Insert(heap.Node{Value: source, Cost: 0, Depth: 0})
	//the previous node does not exists
	previous[source] = math.MaxInt32
	var last int32
	for !pq.IsEmpty() {
		min, _ := pq.Min()
		last = min.Value
		if !visited.Exists(min.Value) {
			visited.Set(min.Value, true)
			dataResult = append(dataResult, g.Nodes[min.Value].Data...)
		}

		visited.Set(min.Value, true)
		pq.DeleteMin()

		if min.Value == target {
			return dist.Cost(target), g.pathPolyline(source, target, previous), dataResult
		}

		for _, e := range g.OutgoingEdges[min.Value] {
			// Validate if we can relax the edge related to the possible ignored node ID.
			if !(g.Nodes[e.ID].Compressed) && !visited.Exists(e.ID) {
				// Relax edge.
				currentPathValue := dist.Cost(min.Value) + e.Weight
				if currentPathValue < dist.Cost(e.ID) {
					dist[e.ID] = currentPathValue
					previous[e.ID] = min.Value
					pq.Insert(heap.Node{Value: e.ID, Cost: currentPathValue, Depth: min.Depth + 1})
				}
			}
		}
	}
	return dist.Cost(target), g.pathPolyline(source, last, previous), dataResult
}

func (g Graph) pathPolyline(start, end int32, previous Previous) [][]float64 {
	result := make([][]float64, 0)
	pathval := end
	result = append(result, []float64{
		s2.CellID(g.Nodes[end].Location).LatLng().Lng.Degrees(),
		s2.CellID(g.Nodes[end].Location).LatLng().Lat.Degrees(),
	})
	for pathval != start {
		result = append(result, []float64{
			s2.CellID(g.Nodes[pathval].Location).LatLng().Lng.Degrees(),
			s2.CellID(g.Nodes[pathval].Location).LatLng().Lat.Degrees(),
		})
		pathval = previous[pathval]
	}
	result = append(result, []float64{
		s2.CellID(g.Nodes[pathval].Location).LatLng().Lng.Degrees(),
		s2.CellID(g.Nodes[pathval].Location).LatLng().Lat.Degrees(),
	})
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
