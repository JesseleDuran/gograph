package resources

import "github.com/paulmach/osm"

var Ways = make(map[int64]osm.Way)
var Relations = make(map[int64]osm.Relation)
var Nodes = make(map[int64]osm.Node)
