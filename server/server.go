package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mrochk/exchange/exchange"
	"github.com/mrochk/exchange/limits"
)

type Server struct {
	listenaddr string
	exchange   *exchange.Exchange
}

type serverFunc func(http.ResponseWriter, *http.Request) error

func NewServer(addr string, port string, exchange *exchange.Exchange) *Server {
	return &Server{
		listenaddr: addr + ":" + port,
		exchange:   exchange,
	}
}

func makeHTTPHandleFunc(f serverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *Server) Run() {
	router := http.NewServeMux()

	router.HandleFunc("/price", makeHTTPHandleFunc(s.handlePrice))
	router.HandleFunc("/test", makeHTTPHandleFunc(s.handleTest))

	fmt.Printf("Server running on port %s...\n", s.listenaddr)
	err := http.ListenAndServe(s.listenaddr, router)
	if err != nil {
		fmt.Println(err)
	}
}

type priceParams struct {
	Price float64 `json:"price"`
}

func (s *Server) handlePrice(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	price := s.exchange.GetOrderBook("ETH", "BTC").Price
	WriteJSON(w, http.StatusOK, priceParams{Price: price})
	return nil
}

type testParams struct {
	AskLimits limits.Limits `json:"ask_limits"`
}

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	a := make(map[int]int)
	a[1] = 2
	ret := testParams{AskLimits: s.exchange.GetOrderBook("ETH", "BTC").AskLimits}
	WriteJSON(w, http.StatusOK, ret)
	return nil
}
