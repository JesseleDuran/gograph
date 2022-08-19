package nearest_edge

import (
	"github.com/JesseleDuran/gograph/nearest_edge/r2"
	"github.com/umahmood/haversine"
)

type GeoSegment struct {
	A, B GeoPoint
}

type GeoPoint struct {
	Coordinates [2]float64 //lat, lng
	ID          int32
}

type GeoSegments []GeoSegment

func GeoPointFromCoords(lat, lng float64, id int32) GeoPoint {
	return GeoPoint{
		Coordinates: [2]float64{lat, lng},
		ID:          id,
	}
}

func (c GeoPoint) ToR2() r2.Point {
	return r2.PointFromCoordinates(c.Coordinates[0], c.Coordinates[1], c.ID)
}

func Distance(a, b [2]float64) float64 {
	_, km := haversine.Distance(
		haversine.Coord{Lat: a[0], Lon: a[1]},
		haversine.Coord{Lat: b[0], Lon: b[1]},
	)
	return km * 1000
}
