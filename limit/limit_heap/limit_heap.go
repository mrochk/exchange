/*
   Implementing the std/container/heap interface.
*/

package limit

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
	/* When dealing with ask limits, we prioritize the limit having
	the lowest price. */
	if h[i].LimitType == limit.Ask {
		return h[i].Price < h[j].Price
	}
	/* When dealing with bid limits, we do the contrary. */
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
