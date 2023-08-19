package uid

import (
	"testing"

	"github.com/mrochk/exchange/set"
)

func TestNewUID(t *testing.T) {
	g, s, toInsert := NewUIDGenerator(), set.NewSet[int64](), 100
	for i := 0; i < toInsert; i++ {
		new := g.NewUID()
		s.Insert(new)
	}
	if s.Size() != toInsert {
		t.Fatalf("Expected %d but got %d.\n", toInsert, s.Size())
	}
}
