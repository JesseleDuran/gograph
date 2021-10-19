package graph

import (
	"log"
	"testing"
)

func TestDeleteRepeated(t *testing.T) {
	b := DeleteRepeated([]int64{1, 3, 2}, []int64{5, 2, 3})
	log.Println(b)
}
