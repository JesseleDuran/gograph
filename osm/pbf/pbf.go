package osm

import (
	graph "github.com/JesseleDuran/gograph"
	"github.com/golang/geo/s2"
	"github.com/qedus/osmpbf"
	"io"
	"log"
	"os"
	"runtime"
)

const CellLevel = 30

type Mode int

const (
	Driving Mode = iota
	Cycling
)

func (m Mode) ToString() string {
	if m == Driving {
		return "drive"
	}
	return "bike"
}

type Filter struct {
	Path      string
	Mode      Mode
	Coverage  s2.Loop
	SetWeight SetWeight
}

type SetWeight func(graph.Coordinate, graph.Coordinate) float32

func MakeGraphFromFile(filter Filter) graph.Graph {
	return createGraph(filter)
}

// createGraph make a graph from an osm file.
// the nodes map represents valid nodes if the file contains a node
// that is not present on this map then the node will be ignored.
func createGraph(filter Filter) graph.Graph {
	nodes := determineValidNodesFromFile(filter.Path, filter.Mode)
	log.Println("nodes", len(nodes))
	nodes = getCoverageNodes(filter.Path, filter.Coverage, nodes)

	f, err := os.Open(filter.Path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("nodes", len(nodes))
	g := graph.Graph{Nodes: make([]graph.Node, 0, len(nodes))}
	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}
	for {
		if o, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch o := o.(type) {

			case *osmpbf.Node:
				osmID := o.ID
				if _, ok := nodes[osmID]; ok {
					id := g.AddNode(graph.Node{
						Location: CoordinatesToCellID(o.Lat, o.Lon),
					})
					nodes[osmID] = id
				}

			case *osmpbf.Way:
				w := o
				if validWay(*w, filter.Mode) {
					for i := 0; i < len(w.NodeIDs)-1; i++ {
						nodeA := graph.Node{}
						nodeB := graph.Node{}
						if idA, ok1 := nodes[w.NodeIDs[i]]; ok1 {
							nodeA = g.Nodes[idA]
							if idB, ok2 := nodes[w.NodeIDs[i+1]]; ok2 {
								nodeB = g.Nodes[idB]
								weight := float32(0.0)
								if filter.SetWeight == nil {
									weight = graph.Distance(s2.CellID(nodeA.Location), s2.CellID(nodeB.Location))
								} else {
									weight = filter.SetWeight(graph.Coordinate{
										Lat: s2.CellID(nodeA.Location).LatLng().Lat.Degrees(),
										Lng: s2.CellID(nodeA.Location).LatLng().Lng.Degrees(),
									}, graph.Coordinate{
										Lat: s2.CellID(nodeB.Location).LatLng().Lat.Degrees(),
										Lng: s2.CellID(nodeB.Location).LatLng().Lng.Degrees(),
									})
								}
								g.RelateNodes(nodeA, nodeB, weight, edgeDirectionFromWay(*w, filter.Mode))
							}
						}
					}
				}
			}
		}
	}
	nodes = nil
	return g
}

// determineValidNodes creates a map of the node of interest.
func determineValidNodesFromFile(path string, mode Mode) map[int64]int32 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)
	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[int64]int32)
	i := 0
	for {
		if o, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch o := o.(type) {
			case *osmpbf.Way:
				w := *o
				if validWay(w, mode) {
					for _, n := range w.NodeIDs {
						if _, ok := result[n]; !ok {
							result[n] = int32(i)
							i++
						}
					}
				}
			}
		}
	}
	_ = f.Close()
	return result
}

func getCoverageNodes(path string, loop s2.Loop, nodes map[int64]int32) map[int64]int32 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	d := osmpbf.NewDecoder(f)
	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)
	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[int64]int32)
	i, j := 0, 0
	for {
		if o, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch o := o.(type) {
			case *osmpbf.Node:
				if _, ok := nodes[o.ID]; ok {
					if loop.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(o.Lat, o.Lon))) {
						result[o.ID] = int32(i)
						i++
					}
				}
				j++
			}
		}
	}
	_ = f.Close()
	return result
}

// CoordinatesToCellID tranform a coordinate into a S2 Cell ID of level 30.
func CoordinatesToCellID(lat, lng float64) uint64 {
	return uint64(s2.CellFromPoint(s2.PointFromLatLng(
		s2.LatLngFromDegrees(lat, lng))).ID().Parent(CellLevel))
}

// validWay determine if the given way is valid or not.
// a valid way is a road segment of interest to build the graph.
func validWay(w osmpbf.Way, mode Mode) bool {
	// valid way tags.
	tags := map[string]struct{}{
		"motorway": {}, "motorway_link": {}, "trunk": {},
		"trunk_link": {}, "primary": {}, "primary_link": {},
		"secondary": {}, "secondary_link": {}, "tertiary": {},
		"tertiary_link": {}, "residential": {},
		"unclassified": {}, "living_street": {},
	}
	_, ok := tags[(w.Tags)["highway"]]
	if mode == Driving {
		return ok
	}
	tags["road"] = struct{}{}
	tags["track"] = struct{}{}
	tags["path"] = struct{}{}
	tags["footway"] = struct{}{}
	tags["pedestrian"] = struct{}{}
	tags["steps"] = struct{}{}
	tags["cycleway"] = struct{}{}
	_, okB := w.Tags["bicycle"]
	_, ok = tags[(w.Tags)["highway"]]
	biciOk := okB || ok
	return biciOk
}

func edgeDirectionFromWay(w osmpbf.Way, mode Mode) graph.EdgeDirection {
	tags := w.Tags
	if mode == Cycling {
		if cycleway, ok := tags["cycleway"]; ok && cycleway == "opposite" || cycleway == "opposite_track" || cycleway == "opposite_lane" {
			return graph.Bidirectional
		}
		if oneWay, ok := tags["oneway:bicycle"]; ok && oneWay == "no" {
			return graph.Bidirectional
		}
	}
	if oneWay, ok := tags["oneway"]; ok && oneWay == "yes" {
		return graph.LeftToRight
	}
	if junction, ok := tags["junction"]; ok && junction == "roundabout" {
		return graph.LeftToRight
	}
	return graph.Bidirectional
}
