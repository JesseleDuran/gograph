package gograph

import (
	"github.com/JesseleDuran/gograph/nearest_edge"
	"github.com/JesseleDuran/gograph/nearest_edge/r2"
	"github.com/golang/geo/s2"
	"log"
)

// Compress based in COMA: Road Network Compression For Map-Matching
// C is conflict factor threshold, which controls the behavior of
// the technique and trades the compression ratio for the map-matching quality.
func (g *Graph) Compress(C float64) {
	compressedNodes := 0
	originalNodes := len(g.Nodes)
	for _, n := range g.Nodes {
		if g.isVictim(n) && g.CheckConflict(n, C) {
			g.DeleteAndMerge(n)
			compressedNodes++
		}
	}
	//Compression Ratio
	cr := 1.0 - float64(compressedNodes/originalNodes)
	log.Println("cr", compressedNodes, originalNodes, cr)
	g.edgeIndex = g.buildEdgeIndex()
}

func (g Graph) isVictim(n Node) bool {
	return g.isIntermediate(n)
}

// Is an intermediate node if it is connected to only two different nodes.
// Exists 2 cases:
// Case 1: Intermediate node of a one-directional path: ni -> n -> nj.
// Case 2: Intermediate node of a bi-directional path: ni <-> n <-> nj.
func (g Graph) isIntermediate(n Node) bool {
	//Case 1.
	eIn := g.IncomingEdges[n.ID]
	eOut := g.OutgoingEdges[n.ID]
	if len(eIn) == 1 && len(eOut) == 1 {
		neighborInID := eIn[0].ID
		neighborOutID := eOut[0].ID
		// Bidirectional to same node.
		if neighborInID == neighborOutID {
			return false
		}
		return true
	}
	//Case 2.
	//if len(eIn) == 2 && len(eOut) == 2 {
	//	i := 0
	//	for _, edgeIn := range eIn {
	//		for _, edgeOut := range eOut {
	//			if edgeIn.ID == edgeOut.ID {
	//				i++
	//			}
	//		}
	//	}
	//	if i == 2 {
	//		return true
	//	}
	//}
	return false
}

// isFanInOut A node that is connected to more than two other nodes
// with one-directional edges, and there is only one input edge
// and all the remaining edges are output edges, or viceversa.
func (g Graph) isFanInOut(n Node) bool {
	eIn := g.IncomingEdges[n.ID]
	eOut := g.OutgoingEdges[n.ID]
	if len(eOut) == 1 && len(eIn) > 1 || len(eIn) == 1 && len(eOut) > 1 {
		for _, eId := range eIn {
			for _, eOutId := range eOut {
				if eId.ID == eOutId.ID {
					return false
				}
			}
		}
		return true
	}
	return false
}

// CheckConflict first finds out the closest edge to the node n.
// After that, we create a new edge by linking the start node of the input edge
// and the end node of the output edge of the under processing pair of edges
// Next, we get the ratio between the distance from n to the bridge edge, and the distance
// from n to the conflict edge. If this ratio is less than the parameter C,
// the bridge is far from nearby conflicting edges so it is a victim.
func (g Graph) CheckConflict(n Node, C float64) bool {
	lat, lng := s2.CellID(n.Location).LatLng().Lat.Degrees(), s2.CellID(n.Location).LatLng().Lng.Degrees()
	eConflict := g.edgeIndex.GeoQuery(lat, lng, []int32{n.ID})
	for _, eIn := range g.IncomingEdges[n.ID] {
		for _, eOut := range g.OutgoingEdges[n.ID] {
			a, b := s2.CellID(g.Nodes[eIn.ID].Location).LatLng(), s2.CellID(g.Nodes[eOut.ID].Location).LatLng()
			bridge := r2.MakeSegmentFromCoordinates([2]float64{a.Lat.Degrees(), a.Lng.Degrees()}, [2]float64{b.Lat.Degrees(), b.Lng.Degrees()})
			//projection of victim in bridge.
			pBridge := bridge.Project(r2.PointFromCoordinates(lat, lng, 0))
			//distance of victim to bridge.
			dBridge := nearest_edge.Distance([2]float64{lat, lng}, pBridge.PointToCoordinates())
			if ConflictFactor(dBridge, eConflict.Distance) < C {
				//midpoint of bridge.
				midpoint := bridge.Midpoint().PointToCoordinates()
				//closest segment from midpoint of bridge.
				newConflict := g.edgeIndex.GeoQuery(midpoint[0], midpoint[1], []int32{n.ID}).Segment
				eNewConflict := r2.MakeSegmentFromCoordinates(newConflict.A.Coordinates, newConflict.B.Coordinates)
				//project victim node in new conflict.
				pNewConflict := eNewConflict.Project(r2.PointFromCoordinates(lat, lng, 0))
				//distance from victim to the projection of victim in new conflict.
				distanceNewConflict := nearest_edge.Distance([2]float64{lat, lng}, pNewConflict.PointToCoordinates())
				if (newConflict.A.ID != eConflict.Segment.A.ID || newConflict.B.ID != eConflict.Segment.B.ID) || ConflictFactor(dBridge, distanceNewConflict) < C {
					return true
				}
			}
		}
	}
	return false
}

// DeleteAndMerge Perform two things, (1) deleting the victim node and its connected
// edges, and (2) adds the new bridge edge to the graph.
func (g *Graph) DeleteAndMerge(n Node) {
	for _, eId := range g.IncomingEdges[n.ID] {
		for _, eOut := range g.OutgoingEdges[n.ID] {
			w := eId.Weight + eOut.Weight
			g.RelateNodes(g.Nodes[eId.ID], g.Nodes[eOut.ID], w, LeftToRight)
		}
	}
	g.NodeAsCompressed(n.ID)
	g.DeleteRelations(n.ID)
}

func ConflictFactor(dBridge, dCompare float64) float64 {
	if dCompare == 0 {
		return 0
	}
	return dBridge / dCompare
}
