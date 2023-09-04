package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

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

	e.PlaceOrder("ETH", "BTC", order.Sell, 1000, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 999, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 1000, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 1000, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 999, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 999, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 800, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 700, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Sell, 600, 10, "maxime")

	e.PlaceOrder("ETH", "BTC", order.Buy, 500, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 580, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 599, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 400, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 450, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 400, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 500, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 589, 10, "maxime")
	e.PlaceOrder("ETH", "BTC", order.Buy, 100, 10, "maxime")

	addr := "127.0.0.1"
	port := 8080
	s := server.NewServer(addr, port, e)
	go s.Run()
	fmt.Printf("Server running on \"%s:8080\"...\n", addr)
	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Println(e.GetOrderBook("ETH", "BTC"))
		time.Sleep(time.Second)
	}
}
