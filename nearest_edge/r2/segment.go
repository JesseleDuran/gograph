package r2

import (
	geojson "github.com/paulmach/go.geojson"
	"gograph/nearest_edge/r1"
	"gograph/nearest_edge/vector"
	"math"
)

type Segment struct {
	// A is the origin of the segment and B the destiny.
	A, B Point
}

type Segments []Segment

func SegmentFromPoints(A, B Point) Segment {
	return Segment{
		A: A,
		B: B,
	}
}

func MakeSegmentFromCoordinates(a, b [2]float64) Segment {
	return Segment{
		A: PointFromCoordinates(a[0], a[1], 0),
		B: PointFromCoordinates(b[0], b[1], 0),
	}
}

func (s Segment) Points() [2]Point {
	return [2]Point{s.A, s.B}
}

func (s Segment) SegmentToLatLng() [2][2]float64 {
	return [2][2]float64{s.A.PointToCoordinates(), s.B.PointToCoordinates()}
}

func (s Segment) Midpoint() Point {
	return Point{
		X: (s.A.X + s.B.X) / 2,
		Y: (s.A.Y + s.B.Y) / 2,
	}
}

func (s Segment) Project(p Point) Point {
	v := vector.Vector{Components: []float64{
		s.B.X - s.A.X, s.B.Y - s.A.Y,
	}}
	w := vector.Vector{Components: []float64{
		p.X - s.A.X, p.Y - s.A.Y,
	}}
	z := v.Project(w)
	projected := Point{X: z.Components[0] + s.A.X, Y: z.Components[1] + s.A.Y}
	if !v.Between(z) {
		if projected.Distance(s.A) < projected.Distance(s.B) {
			return s.A
		}
		return s.B
	}
	return projected
}

func (s Segment) BoundingBox() Rect {
	eps := 0.0
	minX := math.Min(s.A.X, s.B.X)
	minY := math.Min(s.A.Y, s.B.Y)
	maxX := math.Max(s.A.X, s.B.X)
	maxY := math.Max(s.A.Y, s.B.Y)
	return Rect{
		X: r1.Interval{Min: minX + eps, Max: maxX + eps},
		Y: r1.Interval{Min: minY + eps, Max: maxY + eps},
	}
}

func (s Segment) Intersects(x Segment) bool {
	p1, q1 := s.A, s.B
	p2, q2 := x.A, x.B
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special Cases
	// p1, q1 and p2 are collinear and p2 lies on p1q1
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	// p1, q1 and q2 are collinear and q2 lies on nearest_edge p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	// p2, q2 and p1 are collinear and p1 lies on nearest_edge p2q2
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true

	}

	// p2, q2 and q1 are collinear and q1 lies on nearest_edge p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false

}

func onSegment(p, q, r Point) bool {
	return math.Max(p.X, r.X) >= q.X && q.X >= math.Min(p.X, r.X) && math.Max(p.Y, r.Y) >= q.Y && q.Y >= math.Min(p.Y, r.Y)
}

func (s Segment) onSegment(r Point) bool {
	return math.Max(s.A.X, r.X) >= s.B.X && s.B.X >= math.Min(s.A.X, r.X) && math.Max(s.A.Y, r.Y) >= s.B.Y && s.B.Y >= math.Min(s.A.Y, r.Y)
}

func orientation(p, q, r Point) int {
	det := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	if det == 0 {
		return 0 // collinear
	}
	if det > 0 {
		return 1 // clockwise
	}
	return 2 // counterclockwise
}

func (segments Segments) ToGeoJSON() []*geojson.Feature {
	fc := geojson.NewFeatureCollection()
	ff := []*geojson.Feature{}
	for _, n := range segments {
		l := n.SegmentToLatLng()
		p := geojson.NewLineStringFeature([][]float64{{l[0][1], l[0][0]}, {l[1][1], l[1][0]}})
		ff = append(ff, p)
		p.Properties = map[string]interface{}{
			"k": 1,
		}
		fc.AddFeature(p)
	}
	//bytes, _ := fc.MarshalJSON()
	//fmt.Println(string(bytes))
	return ff
}
