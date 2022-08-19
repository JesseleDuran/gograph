package vector

import (
	"log"
	"testing"
)

func TestVector_Unit(t *testing.T) {

	v := Vector{Components: []float64{
		3, 2,
	}}

	x := Vector{Components: []float64{
		1, 3,
	}}

	log.Println(v.Project(x))
}
