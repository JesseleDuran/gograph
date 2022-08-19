package vector

import (
	"math"
)

type Vector struct {
	Components []float64
}

func (v Vector) Dot(x Vector) float64 {
	var result float64
	for i := 0; i < len(v.Components); i++ {
		result += v.Components[i] * x.Components[i]
	}
	return result
}

func (v Vector) Magnitude() float64 {
	magnitude := 0.0
	for _, c := range v.Components {
		magnitude += c * c
	}
	return math.Sqrt(magnitude)
}

// Unit retrieves a vector with magnitude always equals to 1.
func (v Vector) Unit() Vector {
	result := Vector{Components: make([]float64, 0, len(v.Components))}
	magnitude := v.Magnitude()
	for _, c := range v.Components {
		result.Components = append(result.Components, c/magnitude)
	}
	return result
}

func (v *Vector) Translate(x Vector) Vector {
	result := Vector{Components: make([]float64, 0, len(x.Components))}

	// T: R^n -> R^n
	for i := 0; i < len(v.Components); i++ {
		result.Components = append(result.Components, v.Components[i]+x.Components[i])
	}

	return result
}

func (v Vector) Scale(factor float64) Vector {
	result := Vector{
		Components: make([]float64, 0, len(v.Components)),
	}
	for _, c := range v.Components {
		result.Components = append(result.Components, c*factor)
	}
	return result
}

func (v Vector) Project(x Vector) Vector {
	unit := v.Unit()
	return unit.Scale(x.Dot(unit))

}

func (v Vector) Sum(x Vector) Vector {
	return Vector{Components: []float64{}}
}

func (v Vector) Sub(x Vector) Vector {
	return Vector{Components: []float64{}}
}

func (v Vector) R2Determinant(w Vector) float64 {
	return v.Components[0]*w.Components[1] - w.Components[0]*v.Components[1]
}

func (v Vector) Between(x Vector) bool {
	return v.Dot(x) > 0 && v.Dot(x) < v.Dot(v)
}
