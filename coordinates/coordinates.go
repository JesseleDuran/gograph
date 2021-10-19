package coordinates

import (
	"fmt"
	"math"
	"strconv"

	"github.com/golang/geo/s2"
	"github.com/umahmood/haversine"
)

type Coordinates struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

func Distance(a, b Coordinates) float64 {
	_, km := haversine.Distance(
		haversine.Coord{Lat: a.Lat, Lon: a.Lng},
		haversine.Coord{Lat: b.Lat, Lon: b.Lng},
	)
	return km * 1000
}

func FromStrings(latS, lngS string) (Coordinates, error) {
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil || math.IsNaN(lat) {
		return Coordinates{}, fmt.Errorf("invalid lat")
	}
	lng, err := strconv.ParseFloat(lngS, 64)
	if err != nil || math.IsNaN(lng) {
		return Coordinates{}, fmt.Errorf("invalid lng")
	}
	return Coordinates{Lat: lat, Lng: lng}, nil
}

func FromS2LatLng(ll s2.LatLng) Coordinates {
	return Coordinates{
		Lat: ll.Lat.Degrees(),
		Lng: ll.Lng.Degrees(),
	}
}

func (c Coordinates) ToCellID() s2.CellID {
	return s2.CellFromPoint(s2.PointFromLatLng(
		s2.LatLngFromDegrees(c.Lat, c.Lng))).ID().Parent(20) //19
}
