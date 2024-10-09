package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/mrochk/exchange/exchange"
	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/server"
)

func main() {
	/*
		Example code where we open an order book with base currency
		EUR and quote currency GBP and then launch the server.
	*/
	exchange, address, port := exchange.NewExchange(), "127.0.0.1", 8080
	exchange.NewOrderBook("EUR", "GBP")
	server := server.NewServer(address, port, exchange)

	buy, sell := order.Buy, order.Sell

	for i, f := 0, 0.1; i < 1000; i++ {
		for y := 0; y < i%100; y++ {
			exchange.PlaceOrder("EUR", "GBP", buy, float64(i)+f, f*10.0, "user_"+string(i%128))
		}
		f += 0.21
	}

	for i, f := 2000, 0.1; i > 1000; i-- {
		for y := 0; y < i%100; y++ {
			exchange.PlaceOrder("EUR", "GBP", sell, float64(i)+f, f*10.0, "user_"+string(i%128))
		}
		f += 0.21
	}

	exchange.ExecuteOrder("EUR", "GBP", buy, 1, "maxime")

	go server.Run()

	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		output := fmt.Sprintf("Server listening on port %d...\n\n", port)
		for _, v := range exchange.OrderBooks {
			output += fmt.Sprintln(v)
		}
		fmt.Println(output)

		n := rand.Intn(2)
		var orderType order.OrderType
		if n == 0 {
			orderType = buy
		} else {
			orderType = sell
		}

		exchange.ExecuteOrder("EUR", "GBP", orderType, float64(rand.Intn(200000)), "maxime")

		time.Sleep(time.Second / 2)

	}
}
