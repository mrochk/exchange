package orderq

import "github.com/mrochk/exchange/order"

type OrderQ struct {
	Size  int
	first *OrderQCell
	last  *OrderQCell
}

type OrderQCell struct {
	data *order.Order
	next *OrderQCell
}

func NewOrderQ() OrderQ {
	return OrderQ{0, nil, nil}
}

func (q *OrderQ) Insert(o *order.Order) {
	if q.Empty() {
		q.first = &OrderQCell{o, nil}
		q.last = q.first
	} else {
		q.last.next = &OrderQCell{o, nil}
		q.last = q.last.next
	}
	q.Size++
}

func (q *OrderQ) GetFirstOrder() *order.Order {
	if q.Empty() {
		return nil
	}
	ret := q.first.data
	return ret
}

// Must check if return value is not nil (empty queue).
func (q *OrderQ) PopFirstOrder() *order.Order {
	if q.Empty() {
		return nil
	}
	ret := q.first.data
	q.first = q.first.next
	q.Size--
	return ret
}

func (q OrderQ) Empty() bool {
	return q.Size == 0
}
