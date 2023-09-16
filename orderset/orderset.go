package orderset

import (
	"errors"
)

type none struct{}

type Set map[int64]none

func NewSet() Set {
	return make(map[int64]none)
}

func (s Set) Exists(id int64) bool {
	_, exists := s[id]
	return exists
}

func (s Set) Insert(id int64) error {
	_, exists := s[id]
	if !exists {
		s[id] = none{}
		return nil
	}
	return errors.New("element already in set")
}

func (s *Set) Delete(id int64) {
	_, exists := (*s)[id]
	if exists {
		delete(*s, id)
	}
}
