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

/*
Returns the first element of the order queue,
or nil if the queue is empty.
*/
func (q *OrderQ) GetFirstOrder() *order.Order {
	if q.Empty() {
		return nil
	}
	return q.first.data
}

/*
Returns and removes the first element of the order queue,
or nil if the queue is empty.
*/
func (q *OrderQ) PopFirstOrder() *order.Order {
	if q.Empty() {
		return nil
	}
	result := q.first.data
	q.first = q.first.next
	q.Size--
	return result
}

func (q OrderQ) Empty() bool {
	return q.Size == 0
}
