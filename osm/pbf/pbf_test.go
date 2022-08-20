package osm

import (
	graph "github.com/JesseleDuran/gograph"
	"log"
	"testing"
)

func TestMakeGraphFromFile(t *testing.T) {
	g := MakeGraphFromFile(Filter{
		Path: "medellin-def.osm.pbf",
		Mode: Driving,
	})
	node1 := g.BuildEdgeIndex()
	g.AddData([]graph.Point{{
		ID: 2,
		Point: graph.Coordinate{
			6.239964494988126,
			-75.58070182800293,
		},
	}}, node1)
	log.Println(g.Nodes[18553].Data)
	err := g.Serialize("medellin-def.gob")
	log.Println(err)

	g1 := graph.Deserialize("medellin-def.gob")
	log.Println(g1.Nodes[18553].Data)
	log.Println(len(g1.Nodes), g1.Edges())
	node := g1.BuildEdgeIndex()
	d, line, data := g1.DijkstraPathCoord(graph.Coordinate{
		Lat: 6.27797700000,
		Lng: -75.55372100000,
	}, graph.Coordinate{
		Lat: 6.29244400000,
		Lng: -75.58109300000,
	}, node)

	log.Println(d, line, data)
}
