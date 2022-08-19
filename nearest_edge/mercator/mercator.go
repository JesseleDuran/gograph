package mercator

import "github.com/golang/geo/s2"

var Projector = NewMercatorProjection()

func NewMercatorProjection() s2.Projection {
	return s2.NewMercatorProjection(4775228.75015334)
}
