package r2

import (
	"log"
	"testing"
)

func TestSegment_Intersects(t *testing.T) {
	a := Segment{
		A: Point{
			X: 0,
			Y: 0,
		},
		B: Point{
			X: 3,
			Y: 2,
		},
	}

	b := Segment{
		A: Point{
			X: 0,
			Y: 0,
		},
		B: Point{
			X: 10,
			Y: 10,
		},
	}

	c := Segment{
		A: Point{
			X: 0,
			Y: 0,
		},
		B: Point{
			X: 0,
			Y: 10,
		},
	}

	d := Segment{
		A: Point{
			X: 10,
			Y: 0,
		},
		B: Point{
			X: 10,
			Y: 10,
		},
	}

	e := Segment{
		A: Point{
			X: 0,
			Y: 0,
		},
		B: Point{
			X: 10,
			Y: 0,
		},
	}

	f := Segment{
		A: Point{
			X: 11,
			Y: 0,
		},
		B: Point{
			X: 20,
			Y: 0,
		},
	}

	log.Println(a.Intersects(b))
	log.Println(c.Intersects(d))
	log.Println(e.Intersects(f))

	log.Println(a.Project(Point{
		X: 1,
		Y: 3,
	}))
}

func TestSegment_BoundingBox(t *testing.T) {
	A := PointFromCoordinates(-7.84334579506468, -58.37169016477876, 1)
	B := PointFromCoordinates(-7.720063004335502, -56.92109895193889, 2)
	reference := PointFromCoordinates(-8.1988808, -60.7969823, 0)
	segment := SegmentFromPoints(A, B)
	projection := segment.Project(reference)
	log.Println(projection.PointToCoordinates())
}

func TestSegment_Midpoint(t *testing.T) {
	s := Segment{
		A: Point{X: 2, Y: 3},
		B: Point{X: 8, Y: 6},
	}
	p := s.Midpoint()
	if p.X != 5 || p.Y != 4.5 {
		t.Fatalf("got: %v", p)
	}
}
