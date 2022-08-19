package heap

import (
	"errors"
)

var (
	ErrHeapEmpty = errors.New("heap is empty")
)

//Heap represents a priority heap based on the weight of its nodes.
type Heap struct {
	items Nodes
	size  int
}

//Create creates an empty heap of capacity N.
func Create() Heap {
	return Heap{
		items: make(Nodes, 0),
		size:  0,
	}
}

// Insert adds an element to the heap. Assigns the items in the first free
// position, calls heapifyUp to restore heap condition, and increases the
// counter of total of current data.
func (h *Heap) Insert(n Node) {
	h.items = append(h.items, n)
	h.size++
	h.heapifyUp()
}

// Min returns the minimum item of the heap.
func (h *Heap) Min() (Node, error) {
	if !h.IsEmpty() {
		return h.items[0], nil
	}
	return Node{}, ErrHeapEmpty
}

// DeleteMin removes the first element. Extracts the root item and then calls
// heapifyDown to restore heap condition.
func (h *Heap) DeleteMin() error {
	if h.IsEmpty() {
		return ErrHeapEmpty
	}
	h.items[0] = h.items[h.size-1]
	h.size--
	h.items = h.items[:len(h.items)-1]
	h.heapifyDown(0)
	return nil
}

// parentIndex returns the parent index of i.
func parentIndex(i int) int { return (i - 1) / 2 }

// leftChildIndex returns left child index of i.
func leftChildIndex(i int) int { return 2*i + 1 }

// rightChildIndex returns the right child index of i.
func rightChildIndex(i int) int { return 2*i + 2 }

// hasLeftChild returns true if i has a left child.
func (h *Heap) hasLeftChild(i int) bool { return leftChildIndex(i) < h.size }

// hasRightChild returns true if i has a right child.
func (h *Heap) hasRightChild(i int) bool { return rightChildIndex(i) < h.size }

// hasParent returns true if i has a parent.
func (h *Heap) hasParent(i int) bool { return parentIndex(i) >= 0 }

// leftChild returns the left child of i.
func (h *Heap) leftChild(i int) Node { return h.items[leftChildIndex(i)] }

// rightChild returns the right child of i.
func (h *Heap) rightChild(i int) Node { return h.items[rightChildIndex(i)] }

// parent returns true parent of i.
func (h *Heap) parent(i int) Node { return h.items[parentIndex(i)] }

// heapifyUp performs the upward movement. Starts with the index of the last
// item added and, as long as the parent is bigger than the current item, it
// performs a swap and keep moving.
func (h *Heap) heapifyUp() {
	i := h.size - 1
	for h.hasParent(i) && h.parent(i).Cost > h.items[i].Cost {
		temp := h.items[i]
		//swap
		h.items[i] = h.parent(i)
		h.items[parentIndex(i)] = temp
		i = parentIndex(i)
	}
}

// the new root should heapifyDown through the path of minimum values. The function
// compares the root with the min of its children, if the root is greater,
// they are swapped, this ends until the heap condition is not violated, or
// reaches the last level of the tree.
func (h *Heap) heapifyDown(i int) {
	// as long as there's any child, fix the heap.
	for h.hasLeftChild(i) {
		smallerChildIndex := leftChildIndex(i)

		// if results that the right child is even smaller than the left child,
		// then that's the smaller child.
		if h.hasRightChild(i) && h.rightChild(i).Cost < h.leftChild(i).Cost {
			smallerChildIndex = rightChildIndex(i)
		}

		// if the current item is smaller than the smaller of its two children,
		// then the heap condition is done.
		if h.items[i].Cost < h.items[smallerChildIndex].Cost {
			break
		} else {
			//swap
			temp := h.items[i]
			h.items[i] = h.items[smallerChildIndex]
			h.items[smallerChildIndex] = temp
		}
		i = smallerChildIndex
	}
}

func (h *Heap) IsEmpty() bool {
	return h.size == 0
}
