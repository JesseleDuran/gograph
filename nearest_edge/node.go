package nearest_edge

import (
	"github.com/JesseleDuran/gograph/nearest_edge/r2"
	"math"
)

type Position int64

const (
	Left Position = iota
	Right
	Up
	Down
	All
)

func (p Position) Value() []int {
	switch p {
	case Left:
		return []int{0, 3}
	case Right:
		return []int{1, 2}
	case Up:
		return []int{0, 1}
	case Down:
		return []int{2, 3}
	case All:
		return []int{0, 1, 2, 3}
	}
	return []int{0, 1, 2, 3}
}

type Node struct {
	Quadrant r2.Rect //1 quadrant
	Segments r2.Segments
	Depth    int
	Children [4]*Node
}

type NearestResult struct {
	Segment    r2.Segment
	Distance   float64
	Projection r2.Point
}

type GeoNearestResult struct {
	Segment    GeoSegment
	Distance   float64
	Projection GeoPoint
}

type Branch []*Node

func FromGeoSegments(geoSegments ...GeoSegment) Node {
	segments := make(r2.Segments, 0)
	for _, s := range geoSegments {
		segments = append(segments, r2.Segment{
			A: s.A.ToR2(),
			B: s.B.ToR2(),
		})
	}
	return FromSegments(segments...)
}

func FromSegments(segments ...r2.Segment) Node {
	node := Node{
		Quadrant: r2.RectFromSegments(segments...),
	}
	for _, segment := range segments {
		node.Insert(segment)
	}
	return node
}

func (n *Node) Insert(segment r2.Segment) bool {
	if !n.Quadrant.Intercepts(segment.BoundingBox()) || n.Depth == 20 {
		return false
	}

	if n.isLeaf() && n.hasCapacity() && n.Quadrant.Intercepts(segment.BoundingBox()) {
		n.Segments = append(n.Segments, segment)
		return true
	}
	n.rebalance(segment)
	return true
}

// rebalance a node to find space for a given segment.
// the rebalancing process consists in add the node segment + the given
// segment on any of its children.
func (n *Node) rebalance(segment r2.Segment) {
	for _, e := range append(n.Segments, segment) {
		bbSegment := e.BoundingBox()
		middle := n.Quadrant.Middle()
		division := n.Quadrant.Split()
		v := make([]int, 0)
		if bbSegment.Intercepts(n.Quadrant) {
			if bbSegment.X.Max <= middle.X.Min {
				v = Left.Value()
			} else if bbSegment.X.Min >= middle.X.Max {
				v = Right.Value()
			} else if bbSegment.Y.Max <= middle.Y.Min {
				v = Down.Value()
			} else if bbSegment.Y.Min >= middle.Y.Max {
				v = Up.Value()
			} else {
				v = All.Value()
			}
			for _, id := range v {
				n.FindOrCreateChild(division[id], id).Insert(e)
			}
		}
	}
	n.Segments = nil
}

func (n Node) GeoQuery(lat, lng float64, ignore []int32) GeoNearestResult {
	result := n.Query(r2.PointFromCoordinates(lat, lng, 0), ignore)
	return GeoNearestResult{
		Segment: GeoSegment{
			A: GeoPoint{
				Coordinates: result.Segment.A.PointToCoordinates(),
				ID:          result.Segment.A.ID,
			},
			B: GeoPoint{
				Coordinates: result.Segment.B.PointToCoordinates(),
				ID:          result.Segment.B.ID,
			},
		},
		Distance:   Distance([2]float64{lat, lng}, result.Projection.PointToCoordinates()),
		Projection: GeoPoint{Coordinates: result.Projection.PointToCoordinates()},
	}
}

func (n Node) Query(p r2.Point, ignore []int32) NearestResult {
	branch := n.BranchFromPoint(p)
	if len(branch) == 0 {
		return NearestResult{}
	}
	ignored := make(map[int32]bool)
	for _, id := range ignore {
		ignored[id] = true
	}
	return Range(branch, ignored, &r2.Circle{
		Center: p,
		Radius: branch.lastNode().minDistance(p, ignored),
	})
}

func Range(branch Branch, ignore map[int32]bool, circle *r2.Circle) NearestResult {
	result := NearestResult{Segment: r2.Segment{}, Distance: math.MaxFloat64}
	for j := len(branch) - 1; j >= 0; j-- {
		current := branch[j]
		nearest := current.NearestSegment(circle, ignore, branch.nextNode(j))
		if nearest.Distance < result.Distance {
			result = nearest
		}
	}
	return result
}

func (n *Node) FindOrCreateChild(rect r2.Rect, id int) *Node {
	if n.Children[id] == nil {
		n.Children[id] = &Node{Quadrant: rect, Depth: n.Depth + 1}
	}
	return n.Children[id]
}

func (n Node) Child(rect r2.Rect) *Node {
	for _, child := range n.Children {
		if child.Quadrant.Centroid().Equals(rect.Centroid()) {
			return child
		}
	}
	return nil
}

func (n Node) NearestSegment(circle *r2.Circle, ignore map[int32]bool, avoid *Node) NearestResult {
	result := NearestResult{Segment: r2.Segment{}, Distance: math.MaxFloat64}
	if n.isLeaf() {
		for _, e := range n.Segments {
			_, okA := ignore[e.A.ID]
			_, okB := ignore[e.B.ID]
			if okA || okB {
				continue
			}
			projection := e.Project(circle.Center)
			d := projection.Distance(circle.Center)
			if d <= circle.Radius {
				result = NearestResult{Segment: e, Distance: d, Projection: projection}
				circle.Expand(d)
			}
		}
	} else {
		for _, c := range n.Children {
			if c != nil && c != avoid && circle.IntersectsRect(c.Quadrant) {
				nearest := c.NearestSegment(circle, ignore, avoid)
				if nearest.Distance < result.Distance {
					result = nearest
				}
			}
		}
	}
	return result
}

func (n Node) isLeaf() bool {
	return n.Children[0] == nil && n.Children[1] == nil && n.Children[2] == nil && n.Children[3] == nil
}

func (n Node) hasCapacity() bool {
	return n.Segments == nil || len(n.Segments) <= 10
}

// minDistance calculate the min distance between the given and node edges.
func (n Node) minDistance(r r2.Point, ignore map[int32]bool) float64 {
	result := math.MaxFloat64
	for _, e := range n.Segments {
		_, okA := ignore[e.A.ID]
		_, okB := ignore[e.B.ID]
		if okA || okB {
			continue
		}
		p := e.Project(r)
		d := p.Distance(r)
		if d < result {
			result = d
		}
	}
	return result
}

func (n *Node) BranchFromPoint(p r2.Point) Branch {
	if !n.Quadrant.Contains(p) {
		return Branch{}
	}
	nodes := make(Branch, 0)
	nodes = append(nodes, n)
	if n.isLeaf() {
		return nodes
	}
	for _, child := range n.Children {
		if child != nil && child.Quadrant.Contains(p) {
			nodes = append(nodes, child.BranchFromPoint(p)...)
		}
	}
	return nodes
}

func (b Branch) nextNode(i int) *Node {
	lastIndex := len(b) - 1
	if i+1 >= lastIndex {
		return b[lastIndex]
	}
	if i+1 <= 0 {
		return b[0]
	}
	return b[i+1]
}

func (b Branch) lastNode() *Node {
	if len(b) == 1 {
		return b[0]
	}
	return b[len(b)-1]
}
