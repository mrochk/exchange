package set

import "fmt"

type Set[T comparable] struct {
	elements map[T]struct{}
}

func NewSet[T comparable]() Set[T] {
	new := Set[T]{make(map[T]struct{})}
	return new
}

func (s Set[comparable]) Insert(e comparable) {
	s.elements[e] = struct{}{}
}

func (s Set[comparable]) Size() int {
	return len(s.elements)
}

func (s Set[comparable]) Empty() bool {
	return s.Size() == 0
}

func (s Set[comparable]) String() string {
	ret := "{ "
	count := len(s.elements)
	for k := range s.elements {
		count--
		if count > 0 {
			ret += fmt.Sprint(k) + ", "
		} else {
			ret += fmt.Sprint(k)
		}
	}
	return ret + " }"
}
