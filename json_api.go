package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"github.com/ToughDude/go-grpc.git/types"
)

// Define custom type for context keys
type contextKey string

// Constants for context keys
const requestIDKey contextKey = "requestID"

type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func MakeAPIFunc(fn APIFunc) http.HandlerFunc {
	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(ctx, requestIDKey, rand.Intn(1000000))
		if err := fn(ctx, w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
	}
}

// WriteJSON writes the data as JSON to the response writer
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type JSONAPIServer struct {
	listenAddr string
	svc        PriceService
}

func NewJSONAPIServer(listenAddr string, svc PriceService) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() error {
	http.HandleFunc("/", MakeAPIFunc(s.HandleFetchPrice))
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *JSONAPIServer) HandleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	if len(ticker) == 0 {
		return fmt.Errorf("ticker is required")
	}

	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, types.PriceResponse{
		Ticker: ticker,
		Price:  price,
	})
}