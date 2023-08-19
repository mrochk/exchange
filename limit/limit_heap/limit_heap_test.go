package limit

import (
	"container/heap"
	"testing"

	"github.com/mrochk/exchange/limit"
)

func TestLimitHeap(t *testing.T) {
	h := NewLimitHeap()

	for i := 0; i < 100; i++ {
		heap.Push(&h, limit.NewLimit(float64(i)*10, limit.Ask))
	}

	last := -10e9

	// Ask limits should always be ranked in ascending order.
	for len(h) > 0 {
		new := heap.Pop(&h).(*limit.Limit).Price
		if last > new {
			t.Error("ask limits should always be ranked in ascending order")
		}
		last = new
	}

	last = 10e9

	for i := 0; i < 100; i++ {
		heap.Push(&h, limit.NewLimit(float64(i)*10, limit.Bid))
	}
	heap.Init(&h)

	// Bid limits should always be ranked in non-ascending order.
	for len(h) > 0 {
		new := heap.Pop(&h).(*limit.Limit).Price
		if last < new {
			t.Error("bid limits should always be ranked in non-ascending order")
		}
		last = new
	}
}
