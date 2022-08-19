package r2

import (
	"math"
)

type Circle struct {
	Center Point
	Radius float64
}

func (c Circle) Contains(p Point) bool {
	return math.Pow(p.X-c.Center.X, 2)+math.Pow(p.Y-c.Center.Y, 2) < c.Radius
}

func (c Circle) IntersectsRect(rect Rect) bool {
	closestPoint := rect.ClampPoint(c.Center)
	if closestPoint.Distance(c.Center) <= c.Radius {
		return true
	}
	return false
}

func (c *Circle) Expand(r float64) {
	c.Radius = r
}
