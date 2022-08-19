package r2

import (
	"github.com/JesseleDuran/gograph/nearest_edge/mercator"
	"github.com/golang/geo/r2"
	"github.com/golang/geo/s2"
	"math"
)

type Point struct {
	X, Y float64
	ID   int32 //used for index
}

func (p Point) Move(units Point) Point {
	return Point{
		X: p.X + units.X,
		Y: p.Y + units.Y,
	}
}

func (p Point) Equals(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

func (p Point) Invert() Point {
	return Point{
		X:  p.X * -1,
		Y:  p.Y * -1,
		ID: p.ID,
	}
}

func (p Point) Distance(p2 Point) float64 {
	return math.Pow(p2.X-p.X, 2) + math.Pow(p2.Y-p.Y, 2)
}

func (p Point) PointToCoordinates() [2]float64 {
	ll := mercator.Projector.ToLatLng(r2.Point{
		X: p.X,
		Y: p.Y,
	})
	return [2]float64{ll.Lat.Degrees(), ll.Lng.Degrees()}
}

func PointFromCoordinates(lat, lng float64, id int32) Point {
	pG := mercator.Projector.FromLatLng(s2.LatLngFromDegrees(lat, lng))
	return Point{
		X:  pG.X,
		Y:  pG.Y,
		ID: id,
	}
}
