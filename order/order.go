package order

import (
	"time"
)

type Order struct {
	Identifier int64
	OType      OrderType
	Quantity   float64
	Timestamp  int64
	Issuer     string
	Canceled   bool
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
		Canceled:   false,
	}
}

func (t OrderType) String() string {
	if t == Buy {
		return "BUY"
	}
	return "SELL"
}
