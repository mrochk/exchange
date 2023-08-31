package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mrochk/exchange/exchange"
	"github.com/mrochk/exchange/order"
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
			writeJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *Server) Run() {
	router := http.NewServeMux()

	router.HandleFunc("/getorderbooks", makeHTTPHandleFunc(s.handleGetOrderBooks))
	router.HandleFunc("/createorderbook", makeHTTPHandleFunc(s.handleCreateOrderBook))
	router.HandleFunc("/placeorder", makeHTTPHandleFunc(s.handlePlaceOrder))

	fmt.Printf("Server running on \"%s\"...\n", s.listenaddr)
	err := http.ListenAndServe(s.listenaddr, router)
	if err != nil {
		fmt.Println(err)
	}
}

type getOrderBooksParams struct {
	OrderBooks []string `json:"order_books"`
}

func (s *Server) handleGetOrderBooks(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	list, i := make([]string, len(s.exchange.OrderBooks)), 0
	for k := range s.exchange.OrderBooks {
		list[i] = k
		i++
	}

	toWrite := getOrderBooksParams{list}
	writeJSON(w, http.StatusOK, toWrite)
	return nil
}

func (s *Server) handleCreateOrderBook(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	base := r.Header.Get("base")
	quote := r.Header.Get("quote")

	if base == "" || quote == "" {
		msg := fmt.Sprintf("empty <base> (%s) or <quote> (%s) key",
			base, quote)
		return errors.New(msg)
	}

	err := s.exchange.NewOrderBook(base, quote)
	if err != nil {
		return err
	}

	obID := base + "/" + quote
	writeJSON(w, http.StatusOK, fmt.Sprintf("orderbook %s created", obID))
	return nil
}

func (s *Server) handlePlaceOrder(w http.ResponseWriter, r *http.Request) error {
	var err error

	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	base := r.Header.Get("base")
	quote := r.Header.Get("quote")

	if base == "" || quote == "" {
		msg := fmt.Sprintf("empty <base> (%s) or <quote> (%s) key",
			base, quote)
		return errors.New(msg)
	}

	obID := base + "/" + quote
	if !s.exchange.OrderbookExists(obID) {
		msg := fmt.Sprintf("orderbook %s does not exist", obID)
		return errors.New(msg)
	}

	var ordertype order.OrderType
	t := r.Header.Get("type")
	if t == "BUY" {
		ordertype = order.Buy
	} else if t == "SELL" {
		ordertype = order.Sell
	} else {
		msg := fmt.Sprintf("invalid type value (%s)", t)
		return errors.New(msg)
	}

	var price float64
	p := r.Header.Get("price")
	if p == "" {
		msg := fmt.Sprintf("invalid type value (%s)", p)
		return errors.New(msg)
	} else {
		price, err = strconv.ParseFloat(p, 64)
		if err != nil {
			return err
		}
	}

	var qty float64
	q := r.Header.Get("quantity")
	if q == "" {
		msg := fmt.Sprintf("invalid quantity value (%s)", q)
		return errors.New(msg)
	} else {
		qty, err = strconv.ParseFloat(q, 64)
		if err != nil {
			return err
		}
	}

	var issuer string
	i := r.Header.Get("issuer")
	if i == "" {
		msg := fmt.Sprintf("invalid issuer value (%s)", i)
		return errors.New(msg)
	}
	issuer = i

	err = s.exchange.PlaceOrder(base, quote, ordertype, price, qty, issuer)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, "order placed")
	return nil
}
