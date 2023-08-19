package uid

import "time"

type UIDGenerator struct {
	last      int64
	increment int64
}

func NewUIDGenerator() UIDGenerator {
	return UIDGenerator{0, 0}
}

func (u *UIDGenerator) NewUID() int64 {
	now := time.Now().Unix()
	if now == u.last {
		now += u.increment + 1
		u.increment++
	} else {
		u.last = now
		u.increment = 0
	}
	return now
}
