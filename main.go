package main

import (
	"fmt"

	"github.com/mrochk/exchange/orderbook"
)

func main() {
	ob := orderbook.NewOrderBook("BTC", "USD")
	fmt.Println(ob)
}
