package tradeserver

import (
	"database/sql"
	"io"

	pb "Homework_60/genproto"
)

type TradeServer struct {
	pb.UnimplementedTradeServiceServer
	db *sql.DB
}

func NewTradeServer(db *sql.DB) *TradeServer {
	return &TradeServer{db: db}
}

func (s *TradeServer) StreamTrades(stream pb.TradeService_StreamTradesServer) error {
	var totalTrades int32
	var totalAmount float64

	for {
		trade, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		totalTrades++
		totalAmount += float64(trade.Quantity) * trade.Price

		_, err = s.db.Exec(
			"INSERT INTO trades (symbol, quantity, price) VALUES ($1, $2, $3)",
			trade.Symbol, trade.Quantity, trade.Price)
		if err != nil {
			return err
		}
	}

	return stream.SendAndClose(&pb.TradeResponse{
		TotalTrades: totalTrades,
		TotalAmount: totalAmount,
	})
}
