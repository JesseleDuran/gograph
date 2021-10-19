package dijkstra

// func TestDijkstra_FromCoordinatesv2(t *testing.T) {
//	g := graph.BuildFromJsonFilev2("testdata/bogota-v1.4-rc.json")
//	g.Nodes.ToGeoJSON()
//	log.Println(g.Nodes[-316230431].Edges)
//	d := DijkstraV2{Graph: g}
//	r := d.FromCoordinates(coordinates.Coordinates{
//		Lat: 4.7108,
//		Lng:  -74.0719,
//
//
//	}, coordinates.Coordinates{
//		Lat:4.7111,
//		Lng: -74.0709,
//
//	})
//	log.Println("termino", r.TotalWeight, len(r.Polyline))
//	for _, c := range r.Polyline {
//		fmt.Println("[", c[0], ",", c[1], "],")
//	}
//	log.Println(r.Polyline)
//	log.Println(len(g.Nodes))
//}
//
//func TestDijkstra_FromCoordinates2(t *testing.T) {
//	g := graph.BuildFromJsonFile1("testdata/osm-graph-deprimido-small-2.json", nil)
//	g.Nodes.ToGeoJSON()
//	n := g.Nodes[10250081240385323008]
//	log.Println(uint64(coordinates.Coordinates{
//		Lat: 4.6821607,
//		Lng: -74.0532025,
//	}.ToCellID()))
//	for layer, node := range n {
//		fmt.Println("nodo layer", layer)
//		for k, vMap := range node.Edges {
//			for layerE, _ := range vMap {
//				fmt.Println(k.ToToken(), layerE)
//			}
//		}
//		fmt.Println(layer, node.Edges)
//	}
//d := Dijkstra{Graph: g}
//r := d.FromCoordinates(coordinates.Coordinates{
//	Lat: 4.684710506427826,
//	Lng: -74.05229330062866,
//
//}, coordinates.Coordinates{
//	Lat:4.680286276849436,
//	Lng: -74.05593587947281,
//})
//log.Println(r.TotalWeight)
//for _, c := range r.Polyline {
//	fmt.Println("[", c[0], ",", c[1], "],")
//}
//log.Println(r.Polyline)
//}

//func TestDijkstra(t *testing.T) {
//  start := time.Now()
//  g := graph.BuildFromJsonFile("testdata/osm-graph-sp-16.json", nil)
//  end := time.Since(start)
//  log.Println("done graph", end.Milliseconds(), len(g.Nodes))
//  //for n, _ := range g.Nodes {
//  //log.Println("node", n.ToToken())
//  //}
//
//  s := s2.CellIDFromToken("94ce4fe4d")
//  e := s2.CellIDFromToken("94ce5bb29")
//  d := Dijkstra{Graph: g}
//  _, prev := d.FromCellIDs(s, e)
//  log.Println(len(prev))
//  path := path(s, e, prev)
//  for _, v := range path {
//    fmt.Println(
//      s2.CellFromCellID(v).ID().ToToken())
//  }
//}
//
//func TestPath(t *testing.T) {
//	previous := map[s2.CellID]s2.CellID{
//		1: 0,
//		2: 1,
//		3: 2,
//		4: 3,
//		5: 4,
//	}
//	path := path(1, 5, previous)
//	log.Println(path)
//}
