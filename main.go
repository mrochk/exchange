package main

import (
	"fmt"

	"github.com/mrochk/exchange/exchange"
	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/server"
)

func main() {
	e := exchange.NewExchange()

	err := e.NewOrderBook("ETH", "BTC")
	if err != nil {
		fmt.Println(err)
	}

	err = e.NewOrderBook("EUR", "GBP")
	if err != nil {
		fmt.Println(err)
	}

	e.PlaceOrder("ETH", "BTC", order.Sell, 1000, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 999, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 1000, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 1000, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 999, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 999, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 800, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 700, 10, "me")
	e.PlaceOrder("ETH", "BTC", order.Sell, 600, 10, "me")

	s := server.NewServer("127.0.0.1", "8080", e)
	s.Run()
}
