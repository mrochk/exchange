package limit

import (
	"errors"
	"fmt"

	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/queue"
)

type Limit struct {
	LimitType LimitType
	size      float64
	Price     float64
	orders    queue.Queue[*order.Order]
}

type LimitType int

const (
	Bid LimitType = iota
	Ask
)

func NewLimit(price float64, t LimitType) *Limit {
	return &Limit{
		LimitType: t,
		size:      0,
		Price:     price,
		orders:    queue.NewQueue[*order.Order](),
	}
}

func (l *Limit) AddOrder(o *order.Order) error {
	if l.compatibleOrder(o) {
		l.orders.Insert(o)
		return nil
	}
	return errors.New(
		fmt.Sprintf("uncompatible limit and order type (%s and %s)",
			fmt.Sprint(l.LimitType), fmt.Sprint(o.OrderType)))
}

func (l *Limit) PopFirstOrder() *order.Order {
	if l.orders.Empty() {
		return nil
	}
	ret := l.orders.Head()
	l.orders.Pop()
	return ret
}

func (l *Limit) compatibleOrder(o *order.Order) bool {
	a := (l.LimitType == Ask && o.OrderType == order.Sell)
	b := (l.LimitType == Bid && o.OrderType == order.Buy)
	return a || b
}

func (t LimitType) String() string {
	if t == Bid {
		return "Bid"
	}
	return "Ask"
}

func (l Limit) String() string {
	return fmt.Sprintf("Limit at %.2f\n", l.Price) + fmt.Sprint(l.orders)
}
