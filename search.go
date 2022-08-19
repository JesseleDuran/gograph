package gograph

type ShortestPathCriteria struct {
	From        int32
	To          int32
	MaxCost     float32
	InitialCost float32
}

type Distances map[int32]float32

// Cost returns the current weight of a given node ID, if not exists, then returns Infinity.
func (distances Distances) Cost(id int32) float32 {
	if v, ok := distances[id]; ok {
		return v
	}
	return INFINITE
}
