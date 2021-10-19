package bfs

import (
	"github.com/JesseleDuran/osm-graph-parser/parser"
	"log"
	"testing"
)

func TestBFS_Path(t *testing.T) {
	g, _ := parser.FromOSMFileV2("diff_from_multiple_restrictions.osm")
	//restrictions, err := graph.RestrictionsFromFile("testdata/Bogota.osm", g)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//log.Println("restrictions", len(restrictions))
	//g.AddRestrictions(restrictions)
	log.Println("starting")
	d := BFS{Graph: g}
	r := d.Path(1, 9)
	log.Println(r)
}
