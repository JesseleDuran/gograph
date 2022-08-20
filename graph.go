package gograph

import (
	"encoding/gob"
	"encoding/json"
	"github.com/JesseleDuran/gograph/nearest_edge"
	"github.com/golang/geo/s2"
	"github.com/umahmood/haversine"
	"log"
	"math"
	"os"
)

// Graph is a collection of nodes and edges between some or all of the nodes.
type Graph struct {
	Nodes         []Node
	IncomingEdges Relations
	OutgoingEdges Relations
	EdgeIndex     nearest_edge.Node
}

// Node also called vertex is the fundamental unit of which graphs are formed.
// Location is a S2 cell ID representing the coordinates of the node.
type Node struct {
	Data       []uint64
	ID         int32
	Location   uint64
	Compressed bool
}

// Edge represents connections between the nodes of a graph.
// The edges can be directed and weighted.
type Edge struct {
	ID     int32
	Weight float32
}

// Relations join the edges of a node, indexed by its ID.
type Relations [][]Edge

type EdgeDirection int

const (
	Bidirectional EdgeDirection = iota
	LeftToRight
	RightToLeft
)

// DegreeNode returns the degree of a directed graph given a node ID. That means
// that returns the sum of the indegree and the outdegree.
func (g Graph) DegreeNode(id int32) int {
	return len(g.IncomingEdges[id]) + len(g.OutgoingEdges[id])
}

// AddNode adds a node to the array of graph nodes, in the position of its id.
func (g *Graph) AddNode(n Node) int32 {
	id := len(g.Nodes)
	n.ID = int32(id)
	g.Nodes = append(g.Nodes, n)
	g.OutgoingEdges = append(g.OutgoingEdges, make([]Edge, 0))
	g.IncomingEdges = append(g.IncomingEdges, make([]Edge, 0))
	return int32(id)
}

func (g *Graph) DeleteRelations(id int32) {

	for _, edgeIn := range g.IncomingEdges[id] {
		result := []Edge{}
		for _, edgeOut := range g.OutgoingEdges[edgeIn.ID] {
			if id != edgeOut.ID {
				result = append(result, edgeOut)
			}
		}
		g.OutgoingEdges[edgeIn.ID] = result
	}

	for _, edgeOut := range g.OutgoingEdges[id] {
		result := []Edge{}
		for _, edgeIn := range g.IncomingEdges[edgeOut.ID] {
			if id != edgeIn.ID {
				result = append(result, edgeIn)
			}
		}
		g.IncomingEdges[edgeOut.ID] = result
	}
	g.IncomingEdges[id] = []Edge{}
	g.OutgoingEdges[id] = []Edge{}
}

// RelateNodes relates two nodes on a given direction.
func (g *Graph) RelateNodes(a, b Node, weight float32, dir EdgeDirection) {
	switch dir {

	case Bidirectional:
		// relate two nodes bidirectionally o<------>o.
		{
			// Left to right relation(relate node n with node x).
			g.addOutgoingEdge(a.ID, b.ID, weight)
			g.addIncomingEdge(b.ID, a.ID, weight)

			// Right to left relation(relate node x with node n).
			g.addOutgoingEdge(b.ID, a.ID, weight)
			g.addIncomingEdge(a.ID, b.ID, weight)
		}

	case LeftToRight:
		// relate two nodes from left to right o------>o.
		{
			g.addOutgoingEdge(a.ID, b.ID, weight)
			g.addIncomingEdge(a.ID, b.ID, weight)
		}

	case RightToLeft:
		// relate two nodes from right to left o<------o.
		{
			g.addOutgoingEdge(b.ID, a.ID, weight)
			g.addIncomingEdge(b.ID, a.ID, weight)
		}
	}
}

// addOutgoingEdge Adds an outgoing edge to the given node.
// An outgoing edge is an edge that leaves a node, for instance:
// o----->
func (g *Graph) addOutgoingEdge(from, to int32, weight float32) {
	if g.OutgoingEdges[from] == nil {
		g.OutgoingEdges[from] = make([]Edge, 0)
	}
	g.OutgoingEdges[from] = append(g.OutgoingEdges[from], Edge{
		ID:     to,
		Weight: weight,
	})
}

// addIncomingEdge Adds an incoming edge to the given node.
// An incoming edge is an edge that enters the node, for instance:
// ----->o
func (g *Graph) addIncomingEdge(from, to int32, weight float32) {
	if g.IncomingEdges[to] == nil {
		g.IncomingEdges[to] = make([]Edge, 0)
	}
	g.IncomingEdges[to] = append(g.IncomingEdges[to], Edge{
		ID:     from,
		Weight: weight,
	})
}

// Edges returns the number of edges of the graph.
func (g Graph) Edges() int {
	result := 0
	for _, in := range g.IncomingEdges {
		result += len(in)
	}

	for _, out := range g.OutgoingEdges {
		result += len(out)
	}
	return result
}

func (g Graph) Serialize(filePath string) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(g)
	}
	file.Close()
	return err
}

// Degree returns the average degree of the graph.
func (g Graph) Degree() float64 {
	nodesDegree := 0.0
	len := 0
	for _, n := range g.Nodes {
		if !n.Compressed {
			nodesDegree += float64(g.DegreeNode(n.ID))
			len++
		}
	}
	return nodesDegree / float64(len)
}

func (g Graph) BuildEdgeIndex() nearest_edge.Node {
	geoSegments := make(nearest_edge.GeoSegments, 0)
	unique := make(map[int32]map[int32]bool)
	for i, e := range g.OutgoingEdges {
		for _, edge := range e {
			nodeA := g.Nodes[i]
			nodeB := g.Nodes[edge.ID]
			A := s2.CellID(g.Nodes[i].Location).LatLng()
			B := s2.CellID(g.Nodes[edge.ID].Location).LatLng()
			_, ok := unique[nodeA.ID][nodeB.ID]
			_, ok1 := unique[nodeB.ID][nodeA.ID]
			if !ok && !ok1 {
				geoSegments = append(geoSegments, nearest_edge.GeoSegment{
					A: nearest_edge.GeoPointFromCoords(A.Lat.Degrees(), A.Lng.Degrees(), g.Nodes[i].ID),
					B: nearest_edge.GeoPointFromCoords(B.Lat.Degrees(), B.Lng.Degrees(), edge.ID),
				})
			}
			if _, ok := unique[nodeA.ID]; !ok {
				unique[nodeA.ID] = make(map[int32]bool)
			}
			if _, ok := unique[nodeB.ID]; !ok {
				unique[nodeB.ID] = make(map[int32]bool)
			}
			unique[nodeA.ID][nodeB.ID] = true
			unique[nodeB.ID][nodeA.ID] = true
		}
	}
	return nearest_edge.FromGeoSegments(geoSegments...)
}

func (g Graph) EdgeDirectionByNodes(a, b int32) (EdgeDirection, float32) {
	toLeft, toRight := false, false
	weight := float32(0.0)
	for _, edge := range g.OutgoingEdges[a] {
		if b == edge.ID {
			weight = edge.Weight
			toRight = true
			break
		}
	}
	for _, edge := range g.OutgoingEdges[b] {
		if a == edge.ID {
			weight = edge.Weight
			toLeft = true
			break
		}
	}
	if toLeft && toRight {
		return Bidirectional, weight
	}
	if toRight {
		return LeftToRight, weight
	}
	if toLeft {
		return RightToLeft, weight
	}

	return -1, weight
}

func Deserialize(filePath string) Graph {
	var g = new(Graph)
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(g)
	}
	file.Close()
	g.EdgeIndex = g.BuildEdgeIndex()
	return *g
}

// Distance returns the haversine distance in meters between two S2 cell IDs.
func Distance(a, b s2.CellID) float32 {
	_, km := haversine.Distance(
		haversine.Coord{Lat: a.LatLng().Lat.Degrees(), Lon: a.LatLng().Lng.Degrees()},
		haversine.Coord{Lat: b.LatLng().Lat.Degrees(), Lon: b.LatLng().Lng.Degrees()},
	)
	return float32(km * 1000)
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func PointsToCoordinates(points []s2.Point) [][]float64 {
	result := make([][]float64, 0)
	for _, p := range points {
		ll := s2.LatLngFromPoint(p)
		coordinates := []float64{
			ll.Lng.Degrees(), ll.Lat.Degrees(),
		}
		result = append(result, coordinates)
	}
	result = append(result, result[0])
	return result
}

func (g *Graph) NodeAsCompressed(id int32) {
	n := g.Nodes[id]
	n.Compressed = true
	g.Nodes[id] = n
}

func (g *Graph) ProjectCoordinate(coords Coordinate) (int32, float32) {
	nearestResult := g.EdgeIndex.GeoQuery(coords.Lat, coords.Lng, []int32{})
	tempNode := s2.CellIDFromLatLng(s2.LatLngFromDegrees(nearestResult.Projection.Coordinates[0], nearestResult.Projection.Coordinates[1]))
	a, b := g.Nodes[nearestResult.Segment.A.ID], g.Nodes[nearestResult.Segment.B.ID]
	dir, _ := g.EdgeDirectionByNodes(a.ID, b.ID)
	distanceA := Distance(tempNode, s2.CellID(a.Location))
	distanceB := Distance(tempNode, s2.CellID(b.Location))
	switch dir {
	case LeftToRight:
		return b.ID, distanceB
	case RightToLeft:
		return a.ID, distanceA
	default:
		if distanceA < distanceB {
			return a.ID, distanceA
		}
		return b.ID, distanceB
	}
}

func Write(name string, content interface{}) string {
	f, err := os.Create(name)
	if err != nil {
		return ""
	}
	d2, _ := json.Marshal(content)
	n2, err := f.Write(d2)
	if err != nil {
		log.Println(err)
		f.Close()
		return ""
	}
	log.Println(n2, "bytes written successfully")
	err = f.Close()
	if err != nil {
		log.Println(err)
		return ""
	}
	return f.Name()
}

type Point struct {
	ID    uint64
	Point Coordinate
}

func (g *Graph) AddData(dataPoints []Point) {
	for _, data := range dataPoints {
		nodeID, _ := g.ProjectCoordinate(data.Point)
		node := g.Nodes[nodeID]
		node.Data = append(node.Data, data.ID)
		g.Nodes[nodeID] = node
	}
}
