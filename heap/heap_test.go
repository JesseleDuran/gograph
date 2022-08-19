package heap

import "testing"

func TestHeap_Insert_Extract(t *testing.T) {
	heap := Heap{items: make(Nodes, 0)}
	values := make(Nodes, 0)
	for i := 100; i > 0; i-- {
		values = append(values, Node{
			Value: int32(i),
			Cost:  float32(i),
		})
	}
	for i := 0; i < len(values); i++ {
		heap.Insert(values[i])
	}
	for i := 1; i <= len(values); i++ {
		got, _ := heap.Min()
		heap.DeleteMin()
		expected := float32(i)
		if expected != got.Cost {
			t.Fatalf("Expected %f & got %f", expected, got.Cost)
		}
	}
	if !heap.IsEmpty() {
		t.Fatal("Expected an empty heap", heap.size)
	}
}
