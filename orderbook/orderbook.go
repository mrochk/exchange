/*
The order-book is a collection of bid and ask limits containing orders.
*/
package orderbook

import (
	"fmt"

	"github.com/mrochk/exchange/limit"
	"github.com/mrochk/exchange/limitheap"
	"github.com/mrochk/exchange/order"
)

type OrderBook struct {
	bidLimitsMap  map[float64]*limit.Limit
	bidLimitsHeap limitheap.LimitHeap
	askLimitsMap  map[float64]*limit.Limit
	askLimitsHeap limitheap.LimitHeap
	base          string
	quote         string
	size          int
}

func NewOrderBook(base, quote string) *OrderBook {
	return &OrderBook{
		bidLimitsMap:  make(map[float64]*limit.Limit),
		bidLimitsHeap: *limitheap.NewLimitHeap(),
		askLimitsMap:  make(map[float64]*limit.Limit),
		askLimitsHeap: *limitheap.NewLimitHeap(),
		base:          base,
		quote:         quote,
		size:          0,
	}
}

// Place a limit order.
func PlaceLimitOrder(price float64, o *order.Order) error {
	return nil
}

// Place a market order.
func ExecuteMarketOrder(o *order.Order) error {
	return nil
}

func (ob *OrderBook) String() string {
	ret := "Orderbook " + ob.base + " to " + ob.quote + "\n"
	ret += "    size: " + fmt.Sprint(ob.size)
	ret += "\n______________"
	return ret
}
