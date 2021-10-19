package edge

import (
	"github.com/JesseleDuran/osm-graph-parser/parser/resources"
	"github.com/JesseleDuran/osm-graph-parser/tag"
	"github.com/golang/geo/s2"
	"github.com/paulmach/osm"
	"log"
	"strconv"
)

type Edges []uint64

type Edges3 struct {
	Nodes  []Node
	Oneway bool
	Layer  int
	ID     int64
}

type Edges2 []Node

type Node struct {
	ID     int64
	CellId uint64
	Name   string
	Layer  int
}

func FromOSMRelation(r osm.Relation, aux Edges) Edges {
	for i := 0; i < len(r.Members); i++ {
		//TODO: make this as a go routine.
		m := r.Members[i]
		if m.Type == "node" {
			if v, ok := resources.Nodes[m.Ref]; ok {
				aux = append(aux, ToToken(v.Lat, v.Lon))
			}
		}
		if m.Type == "way" {
			w := resources.Ways[m.Ref]
			aux = append(aux, FromWay(w)...)
		}
		if m.Type == "relation" {
			if v, ok := resources.Relations[m.Ref]; ok {
				aux = append(aux, FromOSMRelation(v, Edges{})...)
			}
		}
	}
	return aux
}

func FromWay(w osm.Way) Edges {
	r := make(Edges, 0)
	for i := range w.Nodes {
		if v, ok := resources.Nodes[w.Nodes[i].ID.FeatureID().Ref()]; ok {
			r = append(r, ToToken(v.Lat, v.Lon))
		}
	}
	return r
}

func FromOSMRelation2(r osm.Relation, aux Edges2) Edges2 {
	for i := 0; i < len(r.Members); i++ {
		//TODO: make this as a go routine.
		m := r.Members[i]
		if m.Type == "node" {
			if v, ok := resources.Nodes[m.Ref]; ok {
				aux = append(aux, Node{
					ID:     int64(v.ID),
					CellId: ToToken(v.Lat, v.Lon),
				})
			}
		}
		if m.Type == "way" {
			w := resources.Ways[m.Ref]
			aux = append(aux, FromWay2(w)...)
		}
		if m.Type == "relation" {
			if v, ok := resources.Relations[m.Ref]; ok {
				aux = append(aux, FromOSMRelation2(v, Edges2{})...)
			}
		}
	}
	return aux
}

func FromWay2(w osm.Way) Edges2 {
	r := make(Edges2, 0)
	for i := range w.Nodes {
		if v, ok := resources.Nodes[w.Nodes[i].ID.FeatureID().Ref()]; ok {
			tagsNode := tag.FromOSMTags(v.Tags)
			name := ""
			if nameNode, ok := tagsNode["name"]; ok {
				name = nameNode
			} else {
				tagsWay := tag.FromOSMTags(w.Tags)
				if nameWay, ok := tagsWay["name"]; ok {
					name = nameWay
				}
			}
			r = append(r, Node{
				ID:     int64(v.ID),
				CellId: ToToken(v.Lat, v.Lon),
				Name:   name,
			})
		}
	}
	return r
}

func FromWay3(w osm.Way) Edges3 {
	nodes := make([]Node, 0)
	auxTags := tag.FromOSMTags(w.Tags)
	lString, ok := auxTags["layer"]
	layer := 0
	if ok {
		l, err := strconv.Atoi(lString)
		if err != nil {
			log.Println("layer conversion err:", err.Error())
		}
		layer = l
	}
	for i := range w.Nodes {
		if v, ok := resources.Nodes[w.Nodes[i].ID.FeatureID().Ref()]; ok {
			tagsNode := tag.FromOSMTags(v.Tags)
			name := ""
			if nameNode, ok := tagsNode["name"]; ok {
				name = nameNode
			} else {
				tagsWay := tag.FromOSMTags(w.Tags)
				if nameWay, ok := tagsWay["name"]; ok {
					name = nameWay
				}
			}
			nodes = append(nodes, Node{
				ID:     int64(v.ID),
				CellId: ToToken(v.Lat, v.Lon),
				Name:   name,
				Layer:  layer,
			})
		}
	}
	oneway := false
	v, ok := auxTags["oneway"]
	if ok && v == "yes" {
		oneway = true
	}
	//result := make([]Node, 0)
	//if layer != 0 {
	//	first := nodes[0]
	//	result = append(result, Node{
	//		ID:     first.ID,
	//		CellId: first.CellId,
	//		Name:   first.Name,
	//		Layer:  0,
	//	})
	//	result = append(result, nodes...)
	//	last := nodes[len(nodes)-1]
	//	result = append(result, Node{
	//		ID:     last.ID,
	//		CellId: last.CellId,
	//		Name:   last.Name,
	//		Layer:  0,
	//	})
	//	return Edges3{
	//		Nodes:  result,
	//		Oneway: oneway,
	//		Layer:  layer,
	//		ID: w.ID.FeatureID().Ref(),
	//	}
	//}

	return Edges3{
		Nodes:  nodes,
		Oneway: oneway,
		Layer:  layer,
		ID:     w.ID.FeatureID().Ref(),
	}
}

//TODO TO 17
func ToToken(lat, lng float64) uint64 {
	return uint64(s2.CellFromPoint(s2.PointFromLatLng(
		s2.LatLngFromDegrees(lat, lng))).ID().Parent(30)) //19
}
