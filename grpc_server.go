package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/ToughDude/go-grpc.git/proto"
)

func makeGRPCServerAndRun(ctx context.Context, listenAddr string, svc PriceService) error {
	options := []grpc.ServerOption{}
	server := grpc.NewServer(options...)
	proto.RegisterPriceFetcherServer(server, NewGRPCPriceService(svc))
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	return server.Serve(listener)
}

type GRPCPriceFetcherServer struct {
	proto.UnimplementedPriceFetcherServer
	svc PriceService
}

func NewGRPCPriceService(svc PriceService) *GRPCPriceFetcherServer {
	return &GRPCPriceFetcherServer{
		svc: svc,
	}
}

func (s *GRPCPriceFetcherServer) FetchPrice(ctx context.Context, req *proto.PriceRequest) (*proto.PriceResponse, error) {
	price, err := s.svc.FetchPrice(ctx, req.Ticker)
	if err != nil {
		return nil, err
	}
	return &proto.PriceResponse{
		Ticker: req.Ticker,
		Price:  price,
	}, nil
}
