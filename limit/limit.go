package limit

import (
	"errors"

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
		return errors.New("non compatible order")
	}
	l.orders.Insert(o)
	l.Size += o.Quantity
	return nil
}

func (l *Limit) PopFirstOrder() *order.Order {
	return l.orders.PopFirstOrder()
}

func (l *Limit) validOrder(o *order.Order) bool {
	a := (l.LType == Bid && o.OType == order.Buy)
	b := (l.LType == Ask && o.OType == order.Sell)
	return a || b
}
