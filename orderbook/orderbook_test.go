package orderbook

import (
	"testing"

	"github.com/mrochk/exchange/order"
)

func TestPlaceOrder(t *testing.T) {
	var err error

	ob := NewOrderBook("EUR", "USD")

	for i := 1; i <= 2_000; i++ {
		err = ob.PlaceOrder(order.Buy, 10_000-float64(i), float64(i), "me")
		if err != nil {
			t.Fatal(err)
		}
		err = ob.PlaceOrder(order.Sell, 10_000+float64(i), float64(i), "me")
		if err != nil {
			t.Fatal(err)
		}
	}

	if ob.Size != 4_000 {
		t.Fatalf("expected 4_000, got %d", ob.Size)
	}

	fst_ask_limit := ob.askLimitsHeap.PopLimit()
	fst_bid_limit := ob.bidLimitsHeap.PopLimit()

	if fst_ask_limit.Price != 10_001 {
		t.Fatalf("expected 10_001, got %.2f", fst_ask_limit.Price)
	}
	if fst_bid_limit.Price != 9_999 {
		t.Fatalf("expected 9_999, got %.2f", fst_bid_limit.Price)
	}

	snd_ask_limit := ob.askLimitsHeap.PopLimit()
	snd_bid_limit := ob.bidLimitsHeap.PopLimit()

	if snd_ask_limit.Price != 10_002 {
		t.Fatalf("expected 10_001, got %.2f", snd_ask_limit.Price)
	}
	if snd_bid_limit.Price != 9_998 {
		t.Fatalf("expected 9_999, got %.2f", snd_bid_limit.Price)
	}

	err = ob.PlaceOrder(order.Buy, 9_997, 97, "me")
	if err != nil {
		t.Fatal(err)
	}

	if ob.Size != 4_001 {
		t.Fatalf("expected 4_000, got %d", ob.Size)
	}

	if ob.BidLimitsMap[9_997].Size != 100 {
		t.Fatalf("expected 100, got %.2f", ob.BidLimitsMap[9_997].Size)
	}

}
