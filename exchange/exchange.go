package exchange

import (
	"errors"
	"fmt"

	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/orderbook"
)

type Exchange struct {
	orderbooks map[string]*orderbook.OrderBook
}

func NewExchange() *Exchange {
	return &Exchange{
		orderbooks: make(map[string]*orderbook.OrderBook),
	}
}

func (e Exchange) orderbookExists(obID string) bool {
	_, k := e.orderbooks[obID]
	return k
}

func (e Exchange) NewOrderBook(base string, quote string) error {
	obID := base + "/" + quote
	if e.orderbookExists(obID) {
		msg := fmt.Sprintf("Order book %s already exists.", obID)
		return errors.New(msg)
	}
	e.orderbooks[obID] = orderbook.NewOrderBook(base, quote)
	return nil
}

func (e *Exchange) PlaceOrder(base string, quote string, t order.OrderType,
	price float64, qty float64, issuer string) error {
	obID := base + "/" + quote
	if !e.orderbookExists(obID) {
		msg := fmt.Sprintf("Order book %s, does not exist.", obID)
		return errors.New(msg)
	}
	return e.orderbooks[obID].PlaceOrder(t, price, qty, issuer)
}

func (e *Exchange) ExecuteOrder(base string, quote string, t order.OrderType,
	qty float64, issuer string) error {
	obID := base + "/" + quote
	if !e.orderbookExists(obID) {
		msg := fmt.Sprintf("Order book %s, does not exist.", obID)
		return errors.New(msg)
	}
	return e.orderbooks[obID].ExecuteOrder(t, qty, issuer)
}

func (e Exchange) GetOrderBook(base string, quote string) *orderbook.OrderBook {
	return e.orderbooks[base+"/"+quote]
}
