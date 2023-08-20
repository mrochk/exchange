package uid

import (
	"testing"
)

func TestNewUID(t *testing.T) {
	set := map[int64]struct{}{}
	g, toInsert := NewUIDGenerator(), 100

	for i := 0; i < toInsert; i++ {
		set[g.NewUID()] = struct{}{}
	}

	if len(set) != toInsert {
		t.Fatalf("Expected %d but got %d.\n", toInsert, len(set))
	}
}
