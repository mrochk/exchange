package exchange

import (
	"errors"
	"fmt"

	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/orderbook"
	"github.com/mrochk/exchange/uid"
	"github.com/mrochk/exchange/user"
)

type Exchange struct {
	OrderBooks      map[string]*orderbook.OrderBook
	users           map[int64]*user.User
	userIDGenerator *uid.UIDGenerator
}

func NewExchange() *Exchange {
	return &Exchange{
		OrderBooks:      make(map[string]*orderbook.OrderBook),
		users:           make(map[int64]*user.User),
		userIDGenerator: uid.NewUIDGenerator(),
	}
}

func (e Exchange) NewOrderBook(base string, quote string) error {
	obID := base + "/" + quote
	if e.OrderbookExists(obID) {
		msg := fmt.Sprintf("order book %s already exists", obID)
		return errors.New(msg)
	}
	e.OrderBooks[obID] = orderbook.NewOrderBook(base, quote)
	return nil
}

func (e *Exchange) RegisterUser(username string) int64 {
	uid := e.userIDGenerator.NewUID()
	e.users[uid] = user.NewUser(username, uid)
	return uid
}

func (e *Exchange) PlaceOrder(base string, quote string, t order.OrderType,
	price float64, qty float64, issuer string) (error, int64) {
	obID := base + "/" + quote
	if !e.OrderbookExists(obID) {
		msg := fmt.Sprintf("order book %s, does not exist", obID)
		return errors.New(msg), 0
	}
	return e.OrderBooks[obID].PlaceOrder(t, price, qty, issuer)
}

func (e *Exchange) ExecuteOrder(base string, quote string, t order.OrderType,
	qty float64, issuer string) error {
	obID := base + "/" + quote
	if !e.OrderbookExists(obID) {
		msg := fmt.Sprintf("order book %s, does not exist", obID)
		return errors.New(msg)
	}
	return e.OrderBooks[obID].ExecuteOrder(t, qty, issuer)
}

func (e Exchange) GetOrderBook(base string, quote string) *orderbook.OrderBook {
	return e.OrderBooks[base+"/"+quote]
}

func (e Exchange) OrderbookExists(obID string) bool {
	_, k := e.OrderBooks[obID]
	return k
}
