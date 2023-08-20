package uid

import "time"

type UIDGenerator struct {
	last, increment int64
}

func NewUIDGenerator() *UIDGenerator { return &UIDGenerator{0, 0} }

func (u *UIDGenerator) NewUID() int64 {
	new := time.Now().Unix()
	if new == u.last {
		u.increment++
		new += u.increment
	} else {
		u.last, u.increment = new, 0
	}
	return new
}
