package queue

import "fmt"

type Queue[T any] struct {
	first    *Cell[T]
	last     *Cell[T]
	Size     int
	errValue T
}

type Cell[T any] struct {
	next *Cell[T]
	data T
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{
		first: nil,
		last:  nil,
		Size:  0,
	}
}

func newCell[T any](e T) *Cell[T] {
	return &Cell[T]{
		next: nil,
		data: e,
	}
}

func (q *Queue[T]) Insert(e T) {
	if q.Empty() {
		q.first = newCell[T](e)
		q.last = q.first
	} else {
		q.last.next = newCell[T](e)
		q.last = q.last.next
	}
	q.Size++
}

func (q Queue[T]) Head() T {
	return q.first.data
}

func (q *Queue[T]) Pop() {
	if !q.Empty() {
		if q.first.next != nil {
			q.first = q.first.next
		} else {
			q.first = nil
		}
		q.Size--
	}
}
func (q Queue[T]) Empty() bool {
	return q.Size == 0
}

func printCellsRec[T any](c *Cell[T]) string {
	if c.next != nil {
		return fmt.Sprint(c.data) + "\n" + printCellsRec(c.next)
	}
	return fmt.Sprint(c.data) + " ]"
}

func (q Queue[T]) String() string {
	if q.first != nil {
		return "[\n" + printCellsRec[T](q.first)
	}
	return "[]"
}
