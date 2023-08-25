package limitheap

import (
	"container/heap"

	"github.com/mrochk/exchange/limit"
)

type LimitHeap []*limit.Limit

func NewLimitHeap() LimitHeap {
	ret := make(LimitHeap, 0)
	return ret
}

func (l LimitHeap) PopLimit() *limit.Limit {
	if len(l) == 0 {
		return nil
	}
	return heap.Pop(&l).(*limit.Limit)
}

/* std/container/heap interface impl. */

func (hp *LimitHeap) Push(x any) {
	*hp = append(*hp, x.(*limit.Limit))
}

func (hp *LimitHeap) Pop() any {
	old := *hp
	length := len(old)
	ret := old[length-1]
	*hp = old[0 : length-1]
	return ret
}

func (hp LimitHeap) Len() int {
	return len(hp)
}

func (hp LimitHeap) Less(i, j int) bool {
	// True if p(h[i]) > p(h[j]).
	if hp[i].LType == limit.Ask {
		return hp[i].Price < hp[j].Price
	}
	return hp[i].Price > hp[j].Price
}

func (hp LimitHeap) Swap(i, j int) {
	hp[i], hp[j] = hp[j], hp[i]
}
