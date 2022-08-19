package r1

import (
	"math"
)

type Interval struct {
	Min, Max float64
}

func (i Interval) Intercepts(i2 Interval) bool {
	if i.Min <= i2.Min {
		return i2.Min <= i.Max && i2.Min <= i2.Max
	}
	return i.Min <= i2.Max && i.Min <= i.Max
}

func (i Interval) Contains(value float64) bool {
	return (i.Min <= value) && (value <= i.Max)
}

// AddPoint returns the interval expanded so that it contains the given point.
func (i Interval) AddPoint(p float64) Interval {
	if i.IsEmpty() {
		return Interval{p, p}
	}
	if p < i.Min {
		return Interval{p, i.Max}
	}
	if p > i.Max {
		return Interval{i.Min, p}
	}
	return i
}

func (i Interval) IsEmpty() bool { return i.Min >= i.Max }

func (i Interval) ClampPoint(p float64) float64 {
	return math.Max(i.Min, math.Min(i.Max, p))
}
