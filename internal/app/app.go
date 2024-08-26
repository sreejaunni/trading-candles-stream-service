package app

import (
	"binance-candlestick-service/config"
	"binance-candlestick-service/internal/business/ohlc"
	db "binance-candlestick-service/internal/datastore"
	grpcServer "binance-candlestick-service/internal/grpc"
	"binance-candlestick-service/internal/pkg/binance"
	"binance-candlestick-service/proto"
	"binance-candlestick-service/utils"
	"github.com/jmoiron/sqlx"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

func Run(cfg *config.Config, dbConn *sqlx.DB) error {
	timeframe := time.Minute

	// Create a new OHLCRepo instance
	ohlcRepo := db.NewOHLCRepo(dbConn)
	ohlcCompleteChan := make(chan ohlc.OHLC)

	srv := grpcServer.NewServer()

	go func() {
		for ohlc := range ohlcCompleteChan {
			utils.PrintOHLC(ohlc)
			srv.BroadcastOHLC(ohlc)
			if err := ohlcRepo.SaveOHLC(ohlc); err != nil {
				log.Printf("Failed to save OHLC data: %v", err)
			}
		}
	}()

	symbols := cfg.Symbols
	for _, symbol := range symbols {
		tickChan := make(chan binance.TradeData)
		client, err := binance.NewBinanceClient(symbol)
		if err != nil {
			return err
		}

		go client.StartTickStreaming(tickChan)

		aggregator := ohlc.NewAggregator(symbol, timeframe, tickChan, ohlcCompleteChan)
		go aggregator.Start()
	}

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterOHLCStreamerServer(grpcServer, srv)
	log.Printf("Server listening on port %v", lis.Addr())
	return grpcServer.Serve(lis)
}
