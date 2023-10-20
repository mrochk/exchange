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
	/*
		Example code where we open an order book with base currency
		EUR and quote currency GBP and then launch the server.
	*/
	ex, addr, port := exchange.NewExchange(), "127.0.0.1", 8080
	ex.NewOrderBook("EUR", "GBP")
	s := server.NewServer(addr, port, ex)

	for i, f := 0, 0.1; i < 1000; i++ {
		for y := 0; y < i%100; y++ {
			ex.PlaceOrder("EUR", "GBP", order.Buy, float64(i)+f, f*10.0, "user_"+string(i%128))
		}
		f += 0.21
	}

	for i, f := 2000, 0.1; i > 1000; i-- {
		for y := 0; y < i%100; y++ {
			ex.PlaceOrder("EUR", "GBP", order.Sell, float64(i)+f, f*10.0, "user_"+string(i%128))
		}
		f += 0.21
	}

	go s.Run()

	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		output := fmt.Sprintf("Server listening on port %d...\n\n", port)
		for _, v := range ex.OrderBooks {
			output += fmt.Sprintln(v)
		}
		fmt.Println(output)
		time.Sleep(time.Second / 2)
	}
}
