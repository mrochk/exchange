package order

import "fmt"

type Order struct {
	Quantity   float64
	Identifier int64
	OrderType  OrderType
}

type OrderType int

const (
	Buy OrderType = iota
	Sell
)

func NewOrder(qty float64, id int64, t OrderType) *Order {
	return &Order{
		Quantity:   qty,
		Identifier: id,
		OrderType:  t,
	}
}

func (o Order) String() string {
	ret := fmt.Sprintf("Order %d\n", o.Identifier)
	ret += fmt.Sprintf("    qty: %.2f\n", o.Quantity)
	ret += fmt.Sprintf("    type: %s\n", fmt.Sprint(o.OrderType))
	return ret
}

func (t OrderType) String() string {
	if t == Buy {
		return "Buy"
	}
	return "Sell"
}
