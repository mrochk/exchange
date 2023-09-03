package uid

import (
	"time"
)

type UIDGenerator struct {
	last      int64
	increment int64
}

func NewUIDGenerator() *UIDGenerator {
	return &UIDGenerator{
		last:      0,
		increment: 0,
	}
}

func (u *UIDGenerator) NewUID() int64 {
	new := time.Now().Unix()
	if new <= u.last+u.increment {
		u.increment++
		new += u.increment
	} else {
		u.last, u.increment = new, 0
	}
	return new
}
