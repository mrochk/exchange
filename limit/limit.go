package limit

import (
	"errors"
	"fmt"

	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/orderq"
)

type Limit struct {
	LType  LimitType
	Price  float64
	Size   float64
	orders orderq.OrderQ
}

type LimitType int

const (
	Bid LimitType = iota
	Ask
)

func NewLimit(t LimitType, price float64) *Limit {
	return &Limit{
		LType:  t,
		Price:  price,
		Size:   0,
		orders: orderq.NewOrderQ(),
	}
}

func (l *Limit) AddOrder(o *order.Order) error {
	if !l.validOrder(o) {
		msg := fmt.Sprintf("limit type (%s) incompatible with order type (%s)",
			l.LType, o.OType)
		return errors.New(msg)
	}
	l.orders.Insert(o)
	l.Size += o.Quantity
	return nil
}

func (l *Limit) PopFirstOrder() *order.Order {
	l.Size -= l.orders.GetFirstOrder().Quantity
	return l.orders.PopFirstOrder()
}

func (l *Limit) GetFirstOrder() *order.Order {
	return l.orders.GetFirstOrder()
}

func (l Limit) OrdersCount() int {
	return l.orders.Size
}
func (l *Limit) validOrder(o *order.Order) bool {
	A := (l.LType == Bid && o.OType == order.Buy)
	B := (l.LType == Ask && o.OType == order.Sell)
	return A || B
}

func (t LimitType) String() string {
	if t == Bid {
		return "BID"
	}
	return "ASK"
}

func (l Limit) String() string {
	ret := fmt.Sprintf("\nLimit at %.1f,\n    type: %s, ", l.Price, l.LType)
	ret += fmt.Sprintf("\n    size: %.1f, ", l.Size)
	ret += fmt.Sprintf("\n    n_orders: %d\n", l.orders.Size)
	return ret
}
