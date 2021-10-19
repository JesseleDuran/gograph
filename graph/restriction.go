package graph

import (
	"context"
	"github.com/JesseleDuran/osm-graph-parser/tag"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmxml"
	"log"
	"os"
)

var ways = make(map[int64]osm.Way)
var relations = make(map[int64]osm.Relation)
var nodes = make(map[int64]osm.Node)

type Restriction struct {
	From       []int64
	FromBridge Bridge
	Via        []int64
	ToBridge   Bridge
	To         []int64
	ID         int64
	Type       string
}

type Bridge struct {
	From int64
	To   int64
}

//key is first last node of from and then the first nodo of via. Value the restriction.
type Restrictions []Restriction

func (restriction Restriction) IsNegate() bool {
	if restriction.Type == "no_left_turn" || restriction.Type == "no_right_turn" || restriction.Type == "no_straight_on" {
		return true
	}
	return false
}

func (restriction Restriction) IsExclusive() bool {
	if restriction.Type == "only_left_turn" || restriction.Type == "only_right_turn" ||
		restriction.Type == "only_straight_on" || restriction.Type == "only_u_turn" {
		return true
	}
	return false
}

func (restriction Restriction) Clean() Restriction {
	return Restriction{
		From: DeleteRepeated(restriction.From, restriction.Via),
		Via:  restriction.Via,
		To:   DeleteRepeated(restriction.To, restriction.Via),
		ID:   restriction.ID,
		Type: restriction.Type,
	}
}

// Delete elements in B that are in A and in B.
func DeleteRepeated(B, A []int64) []int64 {
	inter := make(map[int64]bool)
	hash := make(map[int64]bool)
	for _, e := range B {
		hash[e] = true
	}
	for _, e := range A {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter[e] = true
		}
	}
	result := make([]int64, 0)
	for _, value := range B {
		if !inter[value] {
			result = append(result, value)
		}
	}
	return result
}

func RestrictionsFromFile(path string, graph GraphV2) (Restrictions, error) {
	f, err := os.Open(path)
	if err != nil {
		return Restrictions{}, err
	}
	defer f.Close()
	scanner := osmxml.New(context.Background(), f)
	defer scanner.Close()

	for scanner.Scan() {
		o := scanner.Object()
		switch o.ObjectID().Type() {
		case "node":
			n := *o.(*osm.Node)
			nodes[n.ID.FeatureID().Ref()] = n

		case "way":
			w := *o.(*osm.Way)
			auxTags := tag.FromOSMTags(w.Tags)
			v, ok := auxTags["highway"]
			v1, _ := auxTags["access"]
			//v1, ok1 := auxTags["motor_vehicle"];
			//v2, ok2 := auxTags["motorcar"];
			if ok && (v == "motorway" ||
				v == "motorway_link" ||
				v == "trunk" ||
				v == "trunk_link" ||
				v == "primary" ||
				v == "primary_link" ||
				v == "secondary" ||
				v == "secondary_link" ||
				v == "tertiary" ||
				v == "tertiary_link" ||
				v == "residential" ||
				v == "service" ||
				v == "unclassified" ||
				v == "living_street") && (v1 != "no") {
				ways[w.ID.FeatureID().Ref()] = w
			}

		case "relation":
			r := *o.(*osm.Relation)
			auxTags := tag.FromOSMTags(r.Tags)
			v, ok := auxTags["type"]
			if ok && v == "restriction" {
				relations[r.ID.FeatureID().Ref()] = r
			}

		default:
			continue
		}
	}

	err = scanner.Err()
	if err != nil {
		return Restrictions{}, err
	}

	rr := make(Restrictions, 0)
	for _, v := range relations {
		r := RestrictionsFromRelations(v, graph)
		if r.ID == 0 || len(r.From) == 0 || len(r.To) == 0 || len(r.Via) == 0 {
			continue
		}
		rr = append(rr, r)
	}

	return rr, nil
}

func RestrictionsFromRelations(r osm.Relation, graph GraphV2) Restriction {
	auxTags := tag.FromOSMTags(r.Tags)
	aux := Restriction{
		From: make([]int64, 0),
		To:   make([]int64, 0),
		Via:  make([]int64, 0),
		ID:   r.ID.FeatureID().Ref(),
		Type: auxTags["restriction"],
	}

	for _, m := range r.Members {
		switch m.Role {
		case "from":
			if m.Type == "way" {
				way := ways[m.Ref]
				for _, n := range way.Nodes {
					aux.From = append(aux.From, n.ID.FeatureID().Ref())
				}

			} else {
				log.Println("error", m.Type)
			}

		case "via":
			if m.Type == "node" {
				aux.Via = append(aux.Via, m.Ref)
			} else if m.Type == "way" {
				way := ways[m.Ref]
				for _, n := range way.Nodes {
					aux.Via = append(aux.Via, n.ID.FeatureID().Ref())
				}
			}

		case "to":
			if m.Type == "way" {
				way := ways[m.Ref]
				for _, n := range way.Nodes {
					aux.To = append(aux.To, n.ID.FeatureID().Ref())
				}
			} else {
				log.Println("error", m.Type)
			}
		}
	}
	if len(aux.From) == 0 || len(aux.To) == 0 || len(aux.Via) == 0 {
		return Restriction{}
	}
	aux = aux.Clean()
	aux.FromBridge = graph.Bridge(aux.From, aux.Via)
	aux.ToBridge = graph.Bridge(aux.Via, aux.To)
	return aux

}
