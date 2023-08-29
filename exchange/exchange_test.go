package exchange

import (
	"testing"
)

func TestOBExists(t *testing.T) {
	e := NewExchange()
	e.NewOrderBook("ETH", "BTC")

	err := e.NewOrderBook("ETH", "BTC")
	if err == nil {
		t.Fatal()
	}
}
