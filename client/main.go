package main

import (
	"context"
	"log"

	pb "Homework_60/genproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:1108"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTradeServiceClient(conn)

	stream, err := c.StreamTrades(context.Background())
	if err != nil {
		log.Fatalf("could not open stream: %v", err)
	}

	trades := []pb.TradeRequest{
		{Id: 1, Symbol: "AAPL", Quantity: 10, Price: 150.00},
		{Id: 2, Symbol: "GOOGL", Quantity: 5, Price: 1000.00},
		{Id: 3, Symbol: "MSFT", Quantity: 15, Price: 200.00},
	}

	for i:=0; i<len(trades);i++ {
		if err := stream.Send(&trades[i]); err != nil {
			log.Fatalf("could not send trade: %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}

	log.Printf("Total Trades: %d, Total Amount: %f", response.TotalTrades, response.TotalAmount)
}
