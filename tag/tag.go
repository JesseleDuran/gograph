package tag

import "github.com/paulmach/osm"

//Tag represents any metadata that can be of useful information.
type Tag map[string]string

//FromOSMTags transforms a slice of OSM tags to a map of tags where the key and
//the value are the respective ones of the map.
func FromOSMTags(tt osm.Tags) Tag {
	tag := make(Tag)
	for _, t := range tt {
		tag[t.Key] = t.Value
	}
	return tag
}
