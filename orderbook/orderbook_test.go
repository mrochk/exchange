package orderbook

import (
	"fmt"
	"testing"

	"github.com/mrochk/exchange/order"
)

/*
func TestPlaceOrder(t *testing.T) {
	ob := NewOrderBook("A", "B")

	err := ob.PlaceOrder(order.Buy, 100, 10, "me")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ob.bidLimits)

	err = ob.PlaceOrder(order.Sell, 150, 10, "me")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ob.askLimits)

	err = ob.PlaceOrder(order.Buy, 100, 10, "me")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ob.bidLimits)
}
*/

func TestExecuteOrder(t *testing.T) {
	ob := NewOrderBook("BTC", "ETH")

	err := ob.PlaceOrder(order.Buy, 9_000, 100, "maxime")
	if err != nil {
		fmt.Println(err)
	}

	err = ob.PlaceOrder(order.Buy, 9_000, 50, "maxime")
	if err != nil {
		fmt.Println(err)
	}

	err = ob.PlaceOrder(order.Sell, 11_000, 10, "maxime")
	if err != nil {
		fmt.Println(err)
	}

	err = ob.PlaceOrder(order.Sell, 11_000, 10, "maxime")
	if err != nil {
		fmt.Println(err)
	}

	err = ob.PlaceOrder(order.Sell, 11_000, 80, "maxime")
	if err != nil {
		fmt.Println(err)
	}

	err = ob.PlaceOrder(order.Sell, 12_000, 100, "maxime")
	if err != nil {
		fmt.Println(err)
	}

	err = ob.ExecuteOrder(order.Sell, 50, "maxime")
	if err != nil {
		fmt.Println(err)
	}
}
