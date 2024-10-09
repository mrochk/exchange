package orderbook

import (
	"errors"
	"fmt"

	"github.com/mrochk/exchange/limit"
	"github.com/mrochk/exchange/limits"
	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/uid"
)

/* The order-book is a collection of bid and ask limits containing orders. */
type OrderBook struct {
	base           string  // Base asset.
	quote          string  // Quote asset.
	Price          float64 // Price at which the last market order was executed.
	MidPrice       float64 // (Lowest ask limit price + Highest bid l. p.) / 2
	NumberOfOrders int     // The number of orders placed.
	AskLimitsSize  float64 // The sum of the qty of the ask limit orders.
	BidLimitsSize  float64 // The sum of the qty of the bid limit orders.
	uidGenerator   *uid.UIDGenerator
	AskLimits      limits.Limits
	BidLimits      limits.Limits
	askLimitsMap   map[float64]*limit.Limit
	bidLimitsMap   map[float64]*limit.Limit
}

/* Returns and allocates a new OrderBook. */
func NewOrderBook(base, quote string) *OrderBook {
	return &OrderBook{
		base:           base,
		quote:          quote,
		Price:          -1,
		MidPrice:       -1,
		NumberOfOrders: 0,
		AskLimitsSize:  0,
		BidLimitsSize:  0,
		uidGenerator:   uid.NewUIDGenerator(),
		AskLimits:      limits.NewLimits(),
		BidLimits:      limits.NewLimits(),
		askLimitsMap:   make(map[float64]*limit.Limit),
		bidLimitsMap:   make(map[float64]*limit.Limit),
	}
}

/* Place a limit order. */
func (ob *OrderBook) PlaceOrder(t order.OrderType, price float64, qty float64,
	issuer string) (error, int64) {
	if !ob.canPlaceOrder(price, t) {
		msg := fmt.Sprintf("can not place this type of order (%s) at this"+
			" price (%.2f)", t, price)
		return errors.New(msg), 0
	}
	o := order.NewOrder(ob.uidGenerator.NewUID(), t, qty, issuer)
	id := o.Identifier
	err := ob.placeOrder(price, o)
	if err != nil {
		return err, 0
	}
	if t == order.Buy {
		ob.BidLimitsSize += qty
	} else {
		ob.AskLimitsSize += qty
	}
	ob.updateMidPrice()
	ob.NumberOfOrders++
	return nil, id
}

/* Execute a market order. */
func (ob *OrderBook) ExecuteOrder(t order.OrderType, qty float64,
	issuer string) error {
	if t == order.Buy && qty > ob.BidLimitsSize {
		msg := fmt.Sprintf("order quantity too big (%.2f)", qty)
		return errors.New(msg)
	} else if t == order.Buy && qty > ob.AskLimitsSize {
		msg := fmt.Sprintf("order quantity too big (%.2f)", qty)
		return errors.New(msg)
	}
	o := order.NewOrder(ob.uidGenerator.NewUID(), t, qty, issuer)
	err := ob.executeOrder(o)
	if err != nil {
		return err
	}
	if t == order.Buy {
		ob.BidLimitsSize -= qty
	} else {
		ob.AskLimitsSize -= qty
	}
	return nil
}

func (ob *OrderBook) CancelOrder(t order.OrderType, price float64,
	orderID int64) bool {
	if t == order.Buy {
		l := ob.bidLimitsMap[price]
		return l.DeleteOrder(orderID)
	}
	l := ob.askLimitsMap[price]
	return l.DeleteOrder(orderID)
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
			ob.BidLimits = ob.BidLimits.Insert(l)
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
			ob.AskLimits = ob.AskLimits.Insert(l)
		} else {
			ob.askLimitsMap[price].AddOrder(o)
		}
	}
	return nil
}

func (ob *OrderBook) executeOrder(o *order.Order) error {
	if o.OType == order.Sell {
		if len(ob.BidLimits) == 0 {
			msg := fmt.Sprintf("Can not execute order %d, "+
				"no bid limits.", o.Identifier)
			return errors.New(msg)
		}

		for len(ob.BidLimits) > 0 && o.Quantity >= ob.BidLimits[0].Size {
			o.Quantity -= ob.BidLimits[0].Size
			ob.NumberOfOrders -= ob.BidLimits[0].OrdersCount()
			ob.Price = ob.BidLimits[0].Price
			ob.BidLimits.DeleteFirst()
		}

		for len(ob.BidLimits) > 0 && ob.BidLimits[0].Size > 0 &&
			o.Quantity >= ob.BidLimits[0].GetFirstOrder().Quantity {
			ob.Price = ob.BidLimits[0].Price
			ord := ob.BidLimits[0].PopFirstOrder()
			o.Quantity -= ord.Quantity
			ob.NumberOfOrders--
		}

		if o.Quantity > 0 {
			ob.BidLimits[0].GetFirstOrder().Quantity -= o.Quantity
			ob.BidLimits[0].Size -= o.Quantity
			o.Quantity = 0
			ob.Price = ob.BidLimits[0].Price
		}
	} else /* Buy Order */ {
		if len(ob.AskLimits) == 0 {
			msg := fmt.Sprintf("Can not execute order %d, "+
				"no ask limits.", o.Identifier)
			return errors.New(msg)
		}

		for len(ob.AskLimits) > 0 && o.Quantity >= ob.AskLimits[0].Size {
			o.Quantity -= ob.AskLimits[0].Size
			ob.NumberOfOrders -= ob.AskLimits[0].OrdersCount()
			ob.Price = ob.AskLimits[0].Price
			ob.AskLimits.DeleteFirst()
		}

		for len(ob.AskLimits) > 0 && ob.AskLimits[0].Size > 0 &&
			o.Quantity >= ob.AskLimits[0].GetFirstOrder().Quantity {
			ord := ob.AskLimits[0].PopFirstOrder()
			ob.Price = ob.AskLimits[0].Price
			o.Quantity -= ord.Quantity
			ob.NumberOfOrders--
		}

		if o.Quantity > 0 {
			ob.AskLimits[0].GetFirstOrder().Quantity -= o.Quantity
			ob.AskLimits[0].Size -= o.Quantity
			o.Quantity = 0
			ob.Price = ob.AskLimits[0].Price
		}
	}
	return nil
}

func (ob *OrderBook) canPlaceOrder(price float64, t order.OrderType) bool {
	if t == order.Buy {
		// Does not make sense to place a buy limit order
		// higher or equal than the smallest ask limit.
		empty := len(ob.AskLimits) == 0
		if !empty && !(price < ob.AskLimits[0].Price) {
			return false
		}
	} else /* Sell Order */ {
		// Same for a sell limit order lower of equal than
		// the highest ask limit.
		empty := len(ob.BidLimits) == 0
		if !empty && !(price > ob.BidLimits[0].Price) {
			return false
		}
	}
	return true
}

func (ob *OrderBook) updateMidPrice() {
	if len(ob.AskLimits) > 0 && len(ob.BidLimits) > 0 {
		askPrice := ob.AskLimits[0].Price
		bidPrice := ob.BidLimits[0].Price
		ob.MidPrice = (askPrice + bidPrice) / 2
	}
}

func (ob OrderBook) String() string {
	spread := ob.AskLimits[0].Price - ob.BidLimits[0].Price

	ret := fmt.Sprintf("___________________\nOrder book %s/%s:\n",
		ob.base, ob.quote)

	var lim int
	if len(ob.AskLimits) <= 10 {
		lim = len(ob.AskLimits)
	} else {
		lim = 10
	}

	for i := lim - 1; i >= 0; i-- {
		ret += fmt.Sprintf("\033[31m[%.3f] orders: %d, size: %.1f"+
			"\033[0m\n", ob.AskLimits[i].Price,
			ob.AskLimits[i].OrdersCount(), ob.AskLimits[i].Size)
	}

	ret += fmt.Sprintf("\n\033[0;34mMidprice: %.2f, Price: %.2f, Spread: %.2f\nNumber of"+
		" orders: %d\nAsk limits size: %.1f, Bid limits size: %.1f\033[0m\n\n",
		ob.MidPrice, ob.Price, spread, ob.NumberOfOrders, ob.AskLimitsSize,
		ob.BidLimitsSize)

	if len(ob.BidLimits) <= 10 {
		lim = len(ob.BidLimits)
	} else {
		lim = 10
	}

	for i := 0; i < lim; i++ {
		ret += fmt.Sprintf("\033[32m[%.3f] orders: %d, size: %.1f"+
			"\033[0m\n", ob.BidLimits[i].Price,
			ob.BidLimits[i].OrdersCount(), ob.BidLimits[i].Size)
	}

	return ret
}
