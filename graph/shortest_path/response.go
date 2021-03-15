package shortest_path

import "github.com/JesseleDuran/osm-graph/coordinates"

type Response struct {
  Leg         Legs
  TotalWeight float64
}

type Leg struct {
  Points [2]Point
}

type Legs []Leg

type Point struct {
  Point coordinates.Coordinates
  Name  string
}