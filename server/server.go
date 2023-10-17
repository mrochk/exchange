package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mrochk/exchange/exchange"
	"github.com/mrochk/exchange/limits"
	"github.com/mrochk/exchange/order"
)

type Server struct {
	ListenAddr string
	exchange   *exchange.Exchange
}

func NewServer(addr string, port int, exchange *exchange.Exchange) *Server {
	return &Server{
		ListenAddr: addr + ":" + strconv.FormatInt(int64(port), 10),
		exchange:   exchange,
	}
}

type serverFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f serverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return encoder.Encode(v)
}

func (s *Server) Run() {
	router := http.NewServeMux()

	router.HandleFunc("/orderbooks", makeHTTPHandleFunc(s.handleOrderBooks))
	router.HandleFunc("/registeruser", makeHTTPHandleFunc(s.handleRegisterUser))
	router.HandleFunc("/placeorder", makeHTTPHandleFunc(s.handlePlaceOrder))
	router.HandleFunc("/executeorder", makeHTTPHandleFunc(s.handleExecuteOrder))
	router.HandleFunc("/orderbookdata", makeHTTPHandleFunc(s.handleOrderBookData))
	router.HandleFunc("/cancelorder", makeHTTPHandleFunc(s.handleCancelOrder))

	if err := http.ListenAndServe(s.ListenAddr, router); err != nil {
		log.Fatal(err)
	}
}

type orderBooksRes struct {
	OrderBooks []string `json:"order_books"`
}

func (s *Server) handleOrderBooks(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	list, i := make([]string, len(s.exchange.OrderBooks)), 0
	for k := range s.exchange.OrderBooks {
		list[i] = k
		i++
	}

	writeJSON(w, http.StatusOK, orderBooksRes{list})
	return nil
}

type placeOrderReq struct {
	executeOrderReq
	Price float64 `json:"price"`
}

type placeOrderRes struct {
	OrderID int64 `json:"order_id"`
}

func (s *Server) checkPlaceOrderReq(req placeOrderReq) error {
	if req.Base == "" || req.Quote == "" {
		msg := fmt.Sprintf("empty <base> (%s) or <quote> (%s) key",
			req.Base, req.Base)
		return errors.New(msg)
	}

	obID := req.Base + "/" + req.Quote
	if !s.exchange.OrderbookExists(obID) {
		msg := fmt.Sprintf("orderbook %s does not exist", obID)
		return errors.New(msg)
	}

	if req.Type == "" {
		msg := fmt.Sprintf("invalid type value (%s)", req.Type)
		return errors.New(msg)
	}

	if req.Price <= 0.0 {
		msg := fmt.Sprintf("invalid price value (%f)", req.Price)
		return errors.New(msg)
	}

	if req.Qty <= 0.0 {
		msg := fmt.Sprintf("invalid quantity value (%f)", req.Qty)
		return errors.New(msg)
	}

	if req.Issuer == "" {
		msg := fmt.Sprintf("empty issuer value (%s)", req.Issuer)
		return errors.New(msg)
	}

	return nil
}

func (s *Server) handlePlaceOrder(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("Method not allowed: %s", r.Method)
		return errors.New(msg)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	writeJSON(w, 200, len(data))

	req := placeOrderReq{}
	json.Unmarshal(data, &req)

	err = s.checkPlaceOrderReq(req)
	if err != nil {
		return err
	}

	var otype order.OrderType
	if req.Type == "BUY" {
		otype = order.Buy
	} else {
		otype = order.Sell
	}

	err, id := s.exchange.PlaceOrder(req.Base, req.Quote, otype, req.Price, req.Qty, req.Issuer)

	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, &placeOrderRes{OrderID: id})
	return nil
}

type executeOrderReq struct {
	Base   string  `json:"base"`
	Quote  string  `json:"quote"`
	Type   string  `json:"type"`
	Qty    float64 `json:"quantity"`
	Issuer string  `json:"issuer"`
}

func (s *Server) checkExecuteOrderReq(req executeOrderReq) error {
	if req.Base == "" || req.Quote == "" {
		msg := fmt.Sprintf("empty <base> (%s) or <quote> (%s) key",
			req.Base, req.Base)
		return errors.New(msg)
	}

	obID := req.Base + "/" + req.Quote
	if !s.exchange.OrderbookExists(obID) {
		msg := fmt.Sprintf("orderbook %s does not exist", obID)
		return errors.New(msg)
	}

	if req.Type == "" {
		msg := fmt.Sprintf("invalid type value (%s)", req.Type)
		return errors.New(msg)
	}

	if req.Qty <= 0.0 {
		msg := fmt.Sprintf("invalid quantity value (%f)", req.Qty)
		return errors.New(msg)
	}

	if req.Issuer == "" {
		msg := fmt.Sprintf("empty issuer value (%s)", req.Issuer)
		return errors.New(msg)
	}

	return nil
}

func (s *Server) handleExecuteOrder(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("method not allowed: (%s)", r.Method)
		return errors.New(msg)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	req := executeOrderReq{}
	json.Unmarshal(data, &req)

	err = s.checkExecuteOrderReq(req)
	if err != nil {
		return err
	}

	var otype order.OrderType
	if req.Type == "BUY" {
		otype = order.Buy
	} else {
		otype = order.Sell
	}

	err = s.exchange.ExecuteOrder(req.Base, req.Quote, otype, req.Qty,
		req.Issuer)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, "order executed")
	return nil
}

type orderBookDataReq struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type orderBookDataRes struct {
	Base           string        `json:"base"`
	Quote          string        `json:"quote"`
	Price          float64       `json:"price"`
	MidPrice       float64       `json:"mid_price"`
	NumberOfOrders int           `json:"n_orders"`
	AskLimitsSize  float64       `json:"ask_limits_size"`
	BidLimitsSize  float64       `json:"bid_limits_size"`
	AskLimits      limits.Limits `json:"ask_limits"`
	BidLimits      limits.Limits `json:"bid_limits"`
}

func (s *Server) handleOrderBookData(w http.ResponseWriter,
	r *http.Request) error {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf("HTTP method not allowed: (%s)", r.Method)
		return errors.New(msg)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	req := orderBookDataReq{}
	json.Unmarshal(data, &req)

	if req.Base == "" {
		return errors.New("base key is empty")
	}

	if req.Quote == "" {
		return errors.New("quote key is empty")
	}

	obID := req.Base + "/" + req.Quote
	if !s.exchange.OrderbookExists(obID) {
		msg := fmt.Sprintf("orderbook does not exist for %s",
			obID)
		return errors.New(msg)
	}

	ob := s.exchange.GetOrderBook(req.Base, req.Quote)

	result := orderBookDataRes{
		Base:           req.Base,
		Quote:          req.Quote,
		Price:          ob.Price,
		MidPrice:       ob.MidPrice,
		NumberOfOrders: ob.NumberOfOrders,
		AskLimitsSize:  ob.AskLimitsSize,
		BidLimitsSize:  ob.BidLimitsSize,
		AskLimits:      ob.AskLimits,
		BidLimits:      ob.BidLimits,
	}

	writeJSON(w, http.StatusOK, result)
	return nil
}

type registerUserReq struct {
	Username string `json:"username"`
}

type registerUserRes struct {
	UID int64 `json:"uid"`
}

func (s *Server) handleRegisterUser(w http.ResponseWriter,
	r *http.Request) error {

	if r.Method != http.MethodPost {
		msg := fmt.Sprintf("HTTP method not allowed: (%s)", r.Method)
		return errors.New(msg)
	}

	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	data := registerUserReq{}
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return err
	}
	uid := s.exchange.RegisterUser(data.Username)
	writeJSON(w, http.StatusOK, &registerUserRes{uid})
	return nil
}

type cancelOrderReq struct {
	Base    string  `json:"base"`
	Quote   string  `json:"quote"`
	Price   float64 `json:"price"`
	Type    string  `json:"type"`
	OrderID int64   `json:"order_id"`
}

func (s *Server) handleCancelOrder(w http.ResponseWriter,
	r *http.Request) error {
	buffer, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	data := cancelOrderReq{}
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return err
	}

	var t order.OrderType
	if data.Type == "BUY" {
		t = order.Buy
	} else if data.Type == "SELL" {
		t = order.Buy
	} else {
		return errors.New("invalid type value")
	}

	ob := s.exchange.GetOrderBook(data.Base, data.Quote)
	success := ob.CancelOrder(t, data.Price, data.OrderID)

	if success {
		writeJSON(w, http.StatusOK, "order canceled")
		return nil
	}

	return errors.New("impossible to cancel order")
}
