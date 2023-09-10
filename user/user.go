package user

import (
	"github.com/mrochk/exchange/order"
)

type User struct {
	username string
	userID   int64
	orders   map[int64]*order.Order
}

func NewUser(name string, uid int64) *User {
	return &User{
		username: name,
		userID:   uid,
		orders:   make(map[int64]*order.Order),
	}
}

func (u *User) AddOrder(orderID int64, order *order.Order) {
	u.orders[orderID] = order
}
