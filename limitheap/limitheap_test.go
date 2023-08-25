package limitheap

import (
	"container/heap"
	"testing"

	"github.com/mrochk/exchange/limit"
)

func TestLimitHeap(t *testing.T) {
	hp := NewLimitHeap()

	l := limit.NewLimit(limit.Bid, 1)
	ll := limit.NewLimit(limit.Bid, 2)

	heap.Push(&hp, l)
	heap.Push(&hp, ll)

	if heap.Pop(&hp).(*limit.Limit).Price != 2 {
		t.Fail()
	}

	hp, l, ll = nil, nil, nil
	hp = NewLimitHeap()

	l = limit.NewLimit(limit.Ask, 1)
	ll = limit.NewLimit(limit.Ask, 2)

	heap.Push(&hp, l)
	heap.Push(&hp, ll)

	if heap.Pop(&hp).(*limit.Limit).Price != 1 {
		t.Fail()
	}
}
