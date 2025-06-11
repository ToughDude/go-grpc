package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ToughDude/go-grpc.git/client"
	"github.com/ToughDude/go-grpc.git/proto"
)

func main() {
	jsonPort := flag.String("json_port", "3002", "the server port")
	grpcPort := flag.String("grpc_port", "3003", "the server port")
	flag.Parse()

	svc := loggingService{priceService{}}

	// Start gRPC server first
	go func() {
		log.Printf("Starting gRPC server on :%s", *grpcPort)
		if err := makeGRPCServerAndRun(context.Background(), ":"+*grpcPort, svc); err != nil {
			log.Fatal(err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(time.Second)

	// Create gRPC client
	grpcClient, err := client.NewGRPCClient(":" + *grpcPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gRPC client created")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			resp, err := 	grpcClient.FetchPrice(context.Background(), &proto.PriceRequest{Ticker: "ETH"})
			if err != nil {
				log.Printf("Error fetching price: %v", err)
			} else {
				log.Printf("Price for %s: %f", resp.Ticker, resp.Price)
			}
		}
	}()



	// Start JSON API server (blocking call at the end)
	log.Printf("Starting JSON API server on :%s", *jsonPort)
	server := NewJSONAPIServer(":"+*jsonPort, svc)
	log.Fatal(server.Run())
}
