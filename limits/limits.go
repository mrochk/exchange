package limits

import "github.com/mrochk/exchange/limit"

type Limits []*limit.Limit

func NewLimits() Limits {
	return []*limit.Limit{}
}

func (l *Limits) DeleteFirst() {
	*l = (*l)[1:]
}

func (l Limits) Insert(lim *limit.Limit) Limits {
	if len(l) == 0 {
		return []*limit.Limit{lim}
	}

	var (
		new   = make([]*limit.Limit, 0, len(l)+1)
		index = 0
	)

	if lim.LType == limit.Ask {
		for index < len(l) && l[index].Price < lim.Price {
			index++
		}
		a, b := l[:index], l[index:]
		for i := range a {
			new = append(new, a[i])
		}
		new = append(new, lim)
		for i := range b {
			new = append(new, b[i])
		}
	} else {
		for index < len(l) && l[index].Price > lim.Price {
			index++
		}
		a, b := l[:index], l[index:]
		for i := range a {
			new = append(new, a[i])
		}
		new = append(new, lim)
		for i := range b {
			new = append(new, b[i])
		}
	}
	return new
}
