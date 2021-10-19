package graph

import (
	"fmt"
	"github.com/JesseleDuran/osm-graph-parser/coordinates"
	"github.com/JesseleDuran/osm-graph-parser/edge"
	"github.com/JesseleDuran/osm-graph-parser/file/json"
	"github.com/JesseleDuran/osm-graph-parser/parser/resources"
	"github.com/JesseleDuran/osm-graph-parser/tag"
	"github.com/golang/geo/s2"
	geojson "github.com/paulmach/go.geojson"
	"github.com/paulmach/osm"
	"github.com/umahmood/haversine"
	"math"
	"math/rand"
	"reflect"
)

type Graph struct {
	Edges  []edge.Edges  `json:"edges"`
	Edges2 []edge.Edges2 `json:"edges"`
	Edges3 []edge.Edges3 `json:"edges"`
}

type GraphV2 struct {
	Nodes Nodes
}

//can be a cell id or a osmID
type Nodes map[int64]Node

type Node struct {
	ID       int64
	OsmID    int64
	Location s2.CellID
	Edges    Edges
	IsVia    bool
	WayIDs   map[int64]bool
}

type Edge struct {
	Weight float64
}

type Edges map[int64]Edge

func (edges Edges) Copy() Edges {
	r := make(Edges, 0)
	for k, v := range edges {
		r[k] = Edge{Weight: v.Weight}
	}
	return r
}

func (g *GraphV2) RelateNodesFromWay(way osm.Way) {
	oneway := false
	auxTags := tag.FromOSMTags(way.Tags)
	v, ok := auxTags["oneway"]
	if ok && v == "yes" {
		oneway = true
	}
	v, ok = auxTags["junction"]
	if ok && v == "roundabout" {
		oneway = true
	}
	for i := 0; i < len(way.Nodes)-1; i++ {

		a := resources.Nodes[way.Nodes[i].ID.FeatureID().Ref()]
		b := resources.Nodes[way.Nodes[i+1].ID.FeatureID().Ref()]
		nodeA := g.FindNodeOrCreate(a)
		nodeA.WayIDs[way.FeatureID().Ref()] = true
		nodeB := g.FindNodeOrCreate(b)
		nodeB.WayIDs[way.FeatureID().Ref()] = true
		weight := Distance(nodeA.Location, nodeB.Location)

		nodeA.Edges[nodeB.ID] = Edge{Weight: weight}

		if !oneway {
			nodeB.Edges[nodeA.ID] = Edge{Weight: weight}
		}
	}
}

func (g GraphV2) FindNodeFromOsmID(osmID int64) (Node, error) {
	for _, n := range g.Nodes {
		if n.OsmID == osmID {
			return n, nil
		}
	}
	return Node{}, fmt.Errorf("not found")
}

func (g *GraphV2) Bridge(componentA, componentB []int64) Bridge {
	for _, idA := range componentA {
		node := g.Nodes[idA]
		for n, _ := range node.Edges {
			for _, idB := range componentB {
				if idB == int64(math.Abs(float64(n))) {
					return Bridge{
						From: idA,
						To:   idB,
					}
				}
			}
		}
	}
	for _, idB := range componentB {
		node := g.Nodes[idB]
		for n, _ := range node.Edges {
			for _, idA := range componentA {
				if idA == int64(math.Abs(float64(n))) {
					return Bridge{
						From: idA,
						To:   idB,
					}
				}
			}
		}
	}
	return Bridge{}
}

func (g *GraphV2) AddRestriction(r Restriction) {
	bridge := r.FromBridge
	from := g.Nodes[bridge.From]
	from = from.RemoveRelation(bridge.To)
	to := g.Nodes[bridge.To]
	toCopy := to.Copy()
	weight := Distance(from.Location, toCopy.Location)
	from.Edges[toCopy.ID] = Edge{Weight: weight}
	g.AddNodes(from, toCopy)

	if len(r.Via) == 1 {
		b2 := r.ToBridge
		if r.IsNegate() {
			toCopy = toCopy.RemoveRelation(b2.To)
			g.AddNodes(toCopy)
		} else if r.IsExclusive() {
			for k, _ := range toCopy.Edges {
				if k != b2.To {
					toCopy = toCopy.RemoveRelation(k)
				}
			}
			g.AddNodes(toCopy)
		}
		return
	}
	b2 := r.ToBridge
	idFromTo := int64(0)
	copyA, copyB := Node{}, Node{}
	for i := 0; i < len(r.Via)-1; i++ {
		a := g.Nodes[r.Via[i]]
		b := g.Nodes[r.Via[i+1]]
		if a.ID == bridge.To {
			copyA = toCopy
		} else {
			copyA = a.Copy()
		}

		if b.ID == bridge.To {
			copyB = toCopy
		} else {
			copyB = b.Copy()
		}
		copyA = copyA.RemoveRelation(b.ID)
		weight := Distance(copyA.Location, copyB.Location)
		copyA.Edges[copyB.ID] = Edge{Weight: weight}
		g.AddNodes(copyA, copyB)
		if a.ID == b2.From {
			idFromTo = copyA.ID
		}
		if b.ID == b2.From {
			idFromTo = copyB.ID
		}
	}
	fromToCopy := g.Nodes[idFromTo]
	if r.IsNegate() {
		fromToCopy = fromToCopy.RemoveRelation(b2.To)
		g.AddNodes(fromToCopy)
	} else if r.IsExclusive() {
		for k, _ := range fromToCopy.Edges {
			if k != b2.To {
				fromToCopy = fromToCopy.RemoveRelation(k)
			}
		}
		g.AddNodes(fromToCopy)
	}
}

func (g *GraphV2) AddRestrictions(restrictions Restrictions) {
	for _, r := range restrictions {
		g.AddRestriction(r)
	}
}

func randomID() int64 {
	return int64(rand.Intn(math.MaxInt64-1+1) + 1)
}

func (n Node) Copy() Node {
	return Node{
		ID:       randomID(),
		OsmID:    n.OsmID,
		Location: n.Location,
		Edges:    n.Edges.Copy(),
		IsVia:    n.IsVia,
		WayIDs:   n.WayIDs,
	}
}

func (n Node) RemoveRelation(id int64) Node {
	edges := n.Edges
	delete(edges, id)
	return Node{
		ID:       n.ID,
		OsmID:    n.OsmID,
		Location: n.Location,
		Edges:    edges,
		IsVia:    n.IsVia,
		WayIDs:   n.WayIDs,
	}
}

//func (g *GraphV2) CreateRestriction(r restriction.Restriction) {
//	_, ok := g.Restrictions[r.From]
//	if !ok {
//		g.Restrictions[r.From] = make(map[int64]restriction.Restriction)
//	}
//	g.Restrictions[r.From][r.Via] = r
//}

func (g *GraphV2) FindNodeOrCreate(n osm.Node) *Node {
	node, ok := g.Nodes[n.ID.FeatureID().Ref()]
	if !ok {
		node = Node{
			ID:       n.ID.FeatureID().Ref(),
			OsmID:    n.ID.FeatureID().Ref(),
			Location: s2.CellID(edge.ToToken(n.Lat, n.Lon)),
			Edges:    make(map[int64]Edge),
			WayIDs:   make(map[int64]bool),
		}
		g.AddNodes(node)
	}
	return &node
}

func (g *GraphV2) AddNodes(nodes ...Node) {
	for _, n := range nodes {
		g.Nodes[n.ID] = n
	}
}

func Distance(a, b s2.CellID) float64 {
	_, km := haversine.Distance(
		haversine.Coord{Lat: a.LatLng().Lat.Degrees(), Lon: a.LatLng().Lng.Degrees()},
		haversine.Coord{Lat: b.LatLng().Lat.Degrees(), Lon: b.LatLng().Lng.Degrees()},
	)
	return km * 1000
}

func (g *GraphV2) FindNode(c coordinates.Coordinates) int64 {
	keys := reflect.ValueOf(g.Nodes).MapKeys()
	id := keys[0].Interface().(int64)
	node := g.Nodes[id]
	lowest := coordinates.Distance(c, coordinates.FromS2LatLng(node.Location.LatLng()))

	for _, n := range g.Nodes {
		point := coordinates.FromS2LatLng(n.Location.LatLng())
		distance := coordinates.Distance(c, point)
		if distance < lowest {
			lowest = distance
			id = n.ID
		}
	}
	return id
}

func (nodes Nodes) ToGeoJSON() {
	fc := geojson.NewFeatureCollection()
	for k, n := range nodes {
		f := geojson.NewPointFeature([]float64{n.Location.LatLng().Lng.Degrees(), n.Location.LatLng().Lat.Degrees()})
		f.ID = k
		f.Properties = map[string]interface{}{
			"id": n.ID,
		}
		fc.AddFeature(f)
	}
	json.Write("nodes.json", fc)
}
