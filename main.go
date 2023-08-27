package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/mrochk/exchange/order"
	"github.com/mrochk/exchange/orderbook"
)

func main() {
	ob := orderbook.NewOrderBook("EUR", "USD")

	ob.PlaceOrder(order.Buy, 9_000, 10, "maxime")
	ob.PlaceOrder(order.Buy, 8_000, 10, "maxime")
	ob.PlaceOrder(order.Buy, 7_000, 10, "maxime")

	ob.PlaceOrder(order.Sell, 11_000, 10, "maxime")
	ob.PlaceOrder(order.Sell, 12_000, 10, "maxime")
	ob.PlaceOrder(order.Sell, 13_000, 10, "maxime")

	for {
		clear := exec.Command("clear")
		clear.Stdout = os.Stdout
		err := clear.Run()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(ob)

		priceSell := 10_001 + rand.Float64()*(20_001-10_001)
		priceBuy := 1_000 + rand.Float64()*(10_000-1_000)

		qtyA := 1 + rand.Float64()*(100-1)
		qtyB := 1 + rand.Float64()*(100-1)
		n_orders := int(10 + rand.Float64()*(50-10))

		var buyorsell order.OrderType
		if rand.Int()%2 == 0 {
			buyorsell = order.Buy
		} else {
			buyorsell = order.Sell
		}

		for i := 0; i < n_orders; i++ {
			err = ob.PlaceOrder(order.Buy, priceBuy, qtyA, "maxime")
			if err != nil {
				fmt.Println(err)
			}
		}

		for i := 0; i < n_orders; i++ {
			err = ob.PlaceOrder(order.Sell, priceSell, qtyB, "maxime")
			if err != nil {
				fmt.Println(err)
			}
		}

		t := time.Millisecond * 50
		time.Sleep(t)

		qtyC := 1 + rand.Float64()*(10-1)

		err = ob.ExecuteOrder(buyorsell, qtyC, "maxime")
		if err != nil {
			fmt.Println(err)
		}

		qtyD := 1 + rand.Float64()*(10-1)

		var buyorsellB order.OrderType
		if buyorsell == order.Buy {
			buyorsellB = order.Sell
		} else {
			buyorsellB = order.Buy
		}

		err = ob.ExecuteOrder(buyorsellB, qtyD, "maxime")
		if err != nil {
			fmt.Println(err)
		}

		clear = exec.Command("clear")
		clear.Stdout = os.Stdout
		err = clear.Run()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(ob)

		time.Sleep(t)
	}

}
