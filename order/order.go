package order

import (
	"fmt"
	"time"
)

type Order struct {
	Identifier int64
	OType      OrderType
	Quantity   float64
	Timestamp  int64
	Issuer     string
}

type OrderType int

const (
	Buy OrderType = iota
	Sell
)

func NewOrder(id int64, t OrderType, qty float64, issuer string) *Order {
	return &Order{
		Identifier: id,
		OType:      t,
		Quantity:   qty,
		Timestamp:  time.Now().Unix(),
		Issuer:     issuer,
	}
}

func (t OrderType) String() string {
	if t == Buy {
		return "BUY"
	}
	return "SELL"
}

func (o Order) String() string {
	ret := fmt.Sprintf("{Order %d, ", o.Identifier)
	ret += fmt.Sprintf("type:  %s, ", o.OType)
	ret += fmt.Sprintf("qty:  %.1f, ", o.Quantity)
	ret += fmt.Sprintf("issuer:  %s}", o.Issuer)
	return ret
}
