package orderbook

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/mrochk/exchange/limit"
	"github.com/mrochk/exchange/limitheap"
	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/uid"
)

// An order-book is a collection of bid and ask limits containing orders.
type OrderBook struct {
	base          string  // Base currency.
	quote         string  // Quote currency.
	Price         float64 // Price at which the last market order was executed.
	MidPrice      float64 // (Lowest ask limit price + Highest bid l. p.) / 2
	Size          int     // The number of orders placed.
	uidGenerator  *uid.UIDGenerator
	BidLimitsMap  map[float64]*limit.Limit
	bidLimitsHeap limitheap.LimitHeap
	AskLimitsMap  map[float64]*limit.Limit
	askLimitsHeap limitheap.LimitHeap
}

func NewOrderBook(base, quote string) *OrderBook {
	return &OrderBook{
		base:          base,
		quote:         quote,
		Price:         -1,
		MidPrice:      -1,
		Size:          0,
		uidGenerator:  uid.NewUIDGenerator(),
		BidLimitsMap:  make(map[float64]*limit.Limit),
		bidLimitsHeap: limitheap.NewLimitHeap(),
		AskLimitsMap:  make(map[float64]*limit.Limit),
		askLimitsHeap: limitheap.NewLimitHeap(),
	}
}

// Place limit order.
func (ob *OrderBook) PlaceOrder(t order.OrderType, price float64, qty float64,
	issuer string) error {
	if !ob.canPlaceOrder(price, t) {
		msg := fmt.Sprintf("can not place this type of order (%s) at this"+
			"price (%.2f)", t, price)
		return errors.New(msg)
	}
	o := order.NewOrder(ob.uidGenerator.NewUID(), t, qty, issuer)
	ob.placeOrder(price, o)
	ob.updateMidPrice()
	ob.Size++
	return nil
}

// Execute market order.
func ExecuteOrder(o *order.Order) error {
	return nil
}

func (ob *OrderBook) updateMidPrice() {
	if len(ob.askLimitsHeap) > 0 && len(ob.bidLimitsHeap) > 0 {
		askPrice := ob.askLimitsHeap[0].Price
		bidPrice := ob.bidLimitsHeap[0].Price
		ob.MidPrice = (askPrice + bidPrice) / 2
	}
}

func (ob *OrderBook) placeOrder(price float64, o *order.Order) {
	if o.OType == order.Buy {
		if ob.BidLimitsMap[price] == nil {
			l := limit.NewLimit(limit.Bid, price)
			heap.Push(&ob.bidLimitsHeap, l)
			l.AddOrder(o)
			ob.BidLimitsMap[price] = l
		} else {
			ob.BidLimitsMap[price].AddOrder(o)
		}
		return
	}
	if ob.AskLimitsMap[price] == nil {
		l := limit.NewLimit(limit.Ask, price)
		heap.Push(&ob.askLimitsHeap, l)
		l.AddOrder(o)
		ob.AskLimitsMap[price] = l
	} else {
		ob.AskLimitsMap[price].AddOrder(o)
	}
}

func (ob *OrderBook) canPlaceOrder(price float64, t order.OrderType) bool {
	if t == order.Buy {
		// Does not make sense to place a buy limit order
		// higher or equal than the smallest ask limit.
		empty := len(ob.askLimitsHeap) == 0
		if !empty && !(price < ob.askLimitsHeap[0].Price) {
			return false
		}
	} else {
		// Same for a sell limit order lower of equal than
		// the highest ask limit.
		empty := len(ob.bidLimitsHeap) == 0
		if !empty && !(price > ob.bidLimitsHeap[0].Price) {
			return false
		}
	}
	return true
}
