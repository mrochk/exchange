package exchange

import (
	"fmt"
	"testing"
)

func TestOBExists(t *testing.T) {
	e := NewExchange()
	fmt.Println(e.orderbookExists("ETH/BTC"))
	e.NewOrderBook("ETH", "BTC")
	fmt.Println(e.orderbookExists("ETH/BTC"))

	err := e.NewOrderBook("ETH", "BTC")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(e.orderbookExists("ETH/BTC"))
}
