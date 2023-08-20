package orderq

import (
	"testing"

	"github.com/mrochk/exchange/order"
)

func TestOrderQ(t *testing.T) {
	q := NewOrderQ()

	for i := 0.0; i < 100.0; i++ {
		q.Insert(order.NewOrder(int64(i), order.Buy, i, "me"))
	}

	if q.Size != 100 {
		t.Fatal("q.Size != 100")
	}

	var last int64 = -10e9

	for !q.Empty() {
		o := q.PopFirstOrder()
		if last > o.Identifier {
			t.Fatal("elements in queue not in LIFO order")
		}
		last = o.Identifier
	}
}
