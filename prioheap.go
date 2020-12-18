package generationk

// An OhlcHeap is a min-heap of ints.
type OhlcHeap []OHLC

func (h OhlcHeap) Len() int           { return len(h) }
func (h OhlcHeap) Less(i, j int) bool { return h[i].Time.Before(h[j].Time) }
func (h OhlcHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

//Push is used to to put items on the heap
func (h *OhlcHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(OHLC))
}

//Pop is used to remove items from the heap
func (h *OhlcHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
