package limit_heap

import (
	"github.com/mrochk/exchange/limit"
)

type LimitHeap []*limit.Limit

func NewLimitHeap() LimitHeap {
	new := make(LimitHeap, 0)
	return new
}

func (h LimitHeap) Len() int {
	return len(h)
}

func (h LimitHeap) Less(i, j int) bool {
	// If we are dealing with ask limits, we wanna
	// rank first the limit having the lowest price.
	if h[i].LimitType == limit.Ask {
		return h[i].Price < h[j].Price
	}
	// If we are dealing with bid limits, we wanna
	// rank first the limit having the highest price.
	return h[i].Price > h[j].Price
}

func (h LimitHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *LimitHeap) Push(x any) {
	*h = append(*h, x.(*limit.Limit))
}

func (h *LimitHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
