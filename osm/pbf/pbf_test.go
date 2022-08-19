package osm

import (
	"testing"
)

func TestMakeGraphFromFile(t *testing.T) {
	graph := MakeGraphFromFile(Filter{
		Path:                  "colombia.osm.pbf",
		InCoverageGeoJSONPath: "chico.json",
		Mode:                  0,
	})
	graph.Serialize("chico.gob")
}
