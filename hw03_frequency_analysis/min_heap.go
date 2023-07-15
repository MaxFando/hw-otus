package hw03frequencyanalysis

import (
	"container/heap"
)

type Pair struct {
	word      string
	frequency int
}

type MinHeap []Pair

func newMinHeap() *MinHeap {
	min := &MinHeap{}
	heap.Init(min)
	return min
}

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Empty() bool {
	return len(h) == 0
}

func (h MinHeap) Less(i, j int) bool {
	if h[i].frequency != h[j].frequency {
		return h[i].frequency < h[j].frequency
	}

	return h[i].word > h[j].word
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h MinHeap) Top() Pair {
	return h[0]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Pair))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
