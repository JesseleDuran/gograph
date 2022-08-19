package r2

import (
	"github.com/golang/geo/r2"
	"gograph/nearest_edge/r1"
)

type Rect struct {
	X r1.Interval
	Y r1.Interval
}

type Rects []Rect

//func FromSegments(segments ...Edge) Rect {
//	rectangle := Rect{}
//	for _, s := range segments {
//		for _, p := range s.Points() {
//			rectangle = rectangle.AddPoint(p)
//		}
//	}
//	return rectangle
//}

func RectFromSegments(segments ...Segment) Rect {
	rectangle := r2.EmptyRect()
	for _, s := range segments {
		for _, p := range s.Points() {
			rectangle = rectangle.AddPoint(r2.Point{
				X: p.X,
				Y: p.Y,
			})
		}
	}
	return Rect{
		X: r1.Interval{
			Min: rectangle.X.Lo,
			Max: rectangle.X.Hi,
		},
		Y: r1.Interval{
			Min: rectangle.Y.Lo,
			Max: rectangle.Y.Hi,
		},
	}
}

// AddPoint expands the rectangle to include the given point. The rectangle is
// expanded by the minimum amount possible.
func (r Rect) AddPoint(p Point) Rect {
	return Rect{r.X.AddPoint(p.X), r.Y.AddPoint(p.Y)}
}

func (r Rect) Split() [4]Rect {
	midX := (r.X.Max + r.X.Min) / 2
	midY := (r.Y.Max + r.Y.Min) / 2
	q1 := Rect{
		X: r1.Interval{Min: r.X.Min, Max: midX},
		Y: r1.Interval{Min: midY, Max: r.Y.Max},
	}
	q2 := Rect{
		X: r1.Interval{Min: midX, Max: r.X.Max},
		Y: r1.Interval{Min: midY, Max: r.Y.Max},
	}
	q3 := Rect{
		X: r1.Interval{Min: midX, Max: r.X.Max},
		Y: r1.Interval{Min: r.Y.Min, Max: midY},
	}
	q4 := Rect{
		X: r1.Interval{Min: r.X.Min, Max: midX},
		Y: r1.Interval{Min: r.Y.Min, Max: midY},
	}

	//ul = 0, ur =1, ll = 2, lr = 3
	return [4]Rect{q1, q2, q3, q4}
}

func (r Rect) Centroid() Point {
	midX := (r.X.Max + r.X.Min) / 2
	midY := (r.Y.Max + r.Y.Min) / 2
	return Point{
		X: midX,
		Y: midY,
	}
}

func (r Rect) Middle() Rect {
	padding := 0.00001
	centroid := r.Centroid()
	X := r1.Interval{
		Min: centroid.X - padding,
		Max: r.Centroid().X + padding,
	}
	Y := r1.Interval{
		Min: centroid.Y - padding,
		Max: centroid.Y + padding,
	}
	return Rect{
		X: X,
		Y: Y,
	}
}

func (r Rect) Intercepts(rect2 Rect) bool {
	return r.X.Intercepts(rect2.X) && r.Y.Intercepts(rect2.Y)
}

func (r Rect) Contains(p Point) bool {
	//log.Println(p.PointToCoordinates())
	//log.Println(r.RectToCoordinates())
	//log.Println(p.Y)
	//log.Println(r.Y)
	//log.Println(r.X.Contains(p.X), r.Y.Contains(p.Y))
	return r.X.Contains(p.X) && r.Y.Contains(p.Y)
}

func (r Rect) ClampPoint(p Point) Point {
	return Point{X: r.X.ClampPoint(p.X), Y: r.Y.ClampPoint(p.Y)}
}

func (r Rect) RectToCoordinates() [4][2]float64 {
	return [4][2]float64{
		Point{X: r.X.Min, Y: r.Y.Min}.PointToCoordinates(),
		Point{X: r.X.Min, Y: r.Y.Max}.PointToCoordinates(),
		Point{X: r.X.Max, Y: r.Y.Max}.PointToCoordinates(),
		Point{X: r.X.Max, Y: r.Y.Min}.PointToCoordinates(),
	}
}
