package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ToughDude/go-grpc.git/proto"
	"github.com/ToughDude/go-grpc.git/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(remoteAddr string) (proto.PriceFetcherClient, error) {
	conn, err := grpc.Dial(remoteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewPriceFetcherClient(conn), nil
}

type Client struct {
	endpoint string
}

func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func (c *Client) FetchPrice(ticker string) (*types.PriceResponse, error) {
	endpoint := fmt.Sprintf("%s?ticker=%s", c.endpoint, ticker)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResp := map[string]any{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("service responded with a non 200 status code: %s", errResp["error"])
	}

	priceResp := new(types.PriceResponse)
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return nil, err
	}
	resp.Body.Close()

	return priceResp, nil
}
