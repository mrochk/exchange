package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/mrochk/exchange/exchange"
	"github.com/mrochk/exchange/server"
)

func main() {
	ex, addr, port := exchange.NewExchange(), "192.168.1.62", 8080
	ex.NewOrderBook("EUR", "GBP")
	s := server.NewServer(addr, port, ex)

	go s.Run()

	for {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		output := fmt.Sprintf("Server ready to receive requests on port %d...\n\n", port)
		for _, v := range ex.OrderBooks {
			output += fmt.Sprintln(v)
		}
		fmt.Println(output)
		time.Sleep(time.Second / 2)
	}
}
