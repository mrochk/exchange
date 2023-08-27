package orderbook

import (
	"errors"
	"fmt"
	"os"

	"github.com/mrochk/exchange/limit"
	"github.com/mrochk/exchange/limits"
	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/uid"
)

// An order-book is a collection of bid and ask limits containing orders.
type OrderBook struct {
	base         string  // Base currency.
	quote        string  // Quote currency.
	Price        float64 // Price at which the last market order was executed.
	MidPrice     float64 // (Lowest ask limit price + Highest bid l. p.) / 2
	Size         int     // The number of orders placed.
	uidGenerator *uid.UIDGenerator
	askLimits    limits.Limits
	bidLimits    limits.Limits
	askLimitsMap map[float64]*limit.Limit
	bidLimitsMap map[float64]*limit.Limit
}

func NewOrderBook(base, quote string) *OrderBook {
	return &OrderBook{
		base:         base,
		quote:        quote,
		Price:        -1,
		MidPrice:     -1,
		Size:         0,
		uidGenerator: uid.NewUIDGenerator(),
		askLimits:    limits.NewLimits(),
		bidLimits:    limits.NewLimits(),
		askLimitsMap: make(map[float64]*limit.Limit),
		bidLimitsMap: make(map[float64]*limit.Limit),
	}
}

/* Place limit order. */
func (ob *OrderBook) PlaceOrder(t order.OrderType, price float64, qty float64,
	issuer string) error {
	if !ob.canPlaceOrder(price, t) {
		msg := fmt.Sprintf("can not place this type of order (%s) at this"+
			"price (%.2f)", t, price)
		return errors.New(msg)
	}
	o := order.NewOrder(ob.uidGenerator.NewUID(), t, qty, issuer)
	err := ob.placeOrder(price, o)
	if err != nil {
		return err
	}
	ob.updateMidPrice()
	ob.Size++
	return nil
}

func (ob *OrderBook) ExecuteOrder(t order.OrderType, qty float64,
	issuer string) error {
	o := order.NewOrder(ob.uidGenerator.NewUID(), t, qty, issuer)
	err := ob.executeOrder(o)
	if err != nil {
		return err
	}
	if o.Quantity > 0 {
		fmt.Fprintf(os.Stderr, "Order %d filled partially.\n", o.Identifier)
	}
	return nil
}

/* Execute market order. */
func (ob *OrderBook) executeOrder(o *order.Order) error {
	if o.OType == order.Sell {
		if len(ob.bidLimits) == 0 {
			msg := fmt.Sprintf("Can not execute order %d, "+
				"no bid limits.", o.Timestamp)
			return errors.New(msg)
		}

		for len(ob.bidLimits) > 0 && o.Quantity >= ob.bidLimits[0].Size {
			o.Quantity -= ob.bidLimits[0].Size
			ob.Size -= ob.bidLimits[0].OrdersCount()
			ob.Price = ob.bidLimits[0].Price
			ob.bidLimits.DeleteFirst()
		}

		for len(ob.bidLimits) > 0 && ob.bidLimits[0].Size > 0 &&
			o.Quantity >= ob.bidLimits[0].GetFirstOrder().Quantity {
			ob.Price = ob.bidLimits[0].Price
			ord := ob.bidLimits[0].PopFirstOrder()
			o.Quantity -= ord.Quantity
			ob.Size--
		}

		if o.Quantity > 0 {
			ob.bidLimits[0].GetFirstOrder().Quantity -= o.Quantity
			ob.bidLimits[0].Size -= o.Quantity
			o.Quantity = 0
		}
	} else /* Buy Order */ {
		if len(ob.askLimits) == 0 {
			msg := fmt.Sprintf("Can not execute order %d, "+
				"no ask limits.", o.Timestamp)
			return errors.New(msg)
		}

		for len(ob.askLimits) > 0 && o.Quantity >= ob.askLimits[0].Size {
			o.Quantity -= ob.askLimits[0].Size
			ob.Size -= ob.askLimits[0].OrdersCount()
			ob.Price = ob.askLimits[0].Price
			ob.askLimits.DeleteFirst()
		}

		for len(ob.askLimits) > 0 && ob.askLimits[0].Size > 0 &&
			o.Quantity >= ob.askLimits[0].GetFirstOrder().Quantity {
			ord := ob.askLimits[0].PopFirstOrder()
			ob.Price = ob.askLimits[0].Price
			o.Quantity -= ord.Quantity
			ob.Size--
		}

		if o.Quantity > 0 {
			fmt.Println(o.Quantity)
			ob.askLimits[0].GetFirstOrder().Quantity -= o.Quantity
			ob.askLimits[0].Size -= o.Quantity
			o.Quantity = 0
		}
	}
	return nil
}

func (ob *OrderBook) placeOrder(price float64, o *order.Order) error {
	if o.OType == order.Buy {
		if _, exists := ob.bidLimitsMap[price]; !exists {
			l := limit.NewLimit(limit.Bid, price)
			err := l.AddOrder(o)
			if err != nil {
				return err
			}
			ob.bidLimitsMap[price] = l
			ob.bidLimits = ob.bidLimits.Insert(l)
		} else {
			ob.bidLimitsMap[price].AddOrder(o)
		}
	} else /* Sell Order */ {
		if _, exists := ob.askLimitsMap[price]; !exists {
			l := limit.NewLimit(limit.Ask, price)
			err := l.AddOrder(o)
			if err != nil {
				return err
			}
			ob.askLimitsMap[price] = l
			ob.askLimits = ob.askLimits.Insert(l)
		} else {
			ob.askLimitsMap[price].AddOrder(o)
		}
	}
	return nil
}

func (ob *OrderBook) canPlaceOrder(price float64, t order.OrderType) bool {
	if t == order.Buy {
		// Does not make sense to place a buy limit order
		// higher or equal than the smallest ask limit.
		empty := len(ob.askLimits) == 0
		if !empty && !(price < ob.askLimits[0].Price) {
			return false
		}
	} else /* Sell Order */ {
		// Same for a sell limit order lower of equal than
		// the highest ask limit.
		empty := len(ob.bidLimits) == 0
		if !empty && !(price > ob.bidLimits[0].Price) {
			return false
		}
	}
	return true
}

func (ob *OrderBook) updateMidPrice() {
	if len(ob.askLimits) > 0 && len(ob.bidLimits) > 0 {
		askPrice := ob.askLimits[0].Price
		bidPrice := ob.bidLimits[0].Price
		ob.MidPrice = (askPrice + bidPrice) / 2
	}
}

func (ob OrderBook) String() string {
	ret := fmt.Sprintf("OrderBook %s / %s:\n\n", ob.base, ob.quote)
	var lim int
	if len(ob.askLimits) <= 10 {
		lim = len(ob.askLimits)
	} else {
		lim = 10
	}
	for i := lim - 1; i >= 0; i-- {
		ret += fmt.Sprintf("[%.2f] orders: %d, size: %.1f\n", ob.askLimits[i].Price, ob.askLimits[i].OrdersCount(), ob.askLimits[i].Size)
	}
	ret += fmt.Sprintf("\nMidprice: %.2f, Price: %.2f, Number of orders: %d\n\n", ob.MidPrice, ob.Price, ob.Size)

	if len(ob.bidLimits) <= 10 {
		lim = len(ob.bidLimits)
	} else {
		lim = 10
	}

	for i := 0; i < lim; i++ {
		ret += fmt.Sprintf("[%.2f] orders: %d, size: %.1f\n", ob.bidLimits[i].Price, ob.bidLimits[i].OrdersCount(), ob.bidLimits[i].Size)
	}

	return ret
}
