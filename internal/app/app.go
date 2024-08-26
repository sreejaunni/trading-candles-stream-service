package app

import (
	"binance-candlestick-service/config"
	"binance-candlestick-service/internal/business/ohlc"
	db "binance-candlestick-service/internal/datastore"
	grpcServer "binance-candlestick-service/internal/grpc"
	"binance-candlestick-service/internal/pkg/binance"
	"binance-candlestick-service/proto"
	"binance-candlestick-service/utils"
	"log"
	"net"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc"
)

const TIMEFRAME = time.Minute

func Run(cfg *config.Config, dbConn *sqlx.DB) error {
	logger := log.New(os.Stdout, "", 0)

	ohlcRepo := db.NewOHLCRepo(dbConn)
	ohlcCompleteChan := make(chan ohlc.OHLC)

	srv := grpcServer.NewServer()

	symbols := cfg.Symbols

	// Counter for the number of OHLC bars processed
	barCounter := 0
	barsPerLog := len(symbols)

	go func() {
		for ohlc := range ohlcCompleteChan {

			utils.PrintOHLC(ohlc)

			srv.BroadcastOHLC(ohlc)
			if err := ohlcRepo.SaveOHLC(ohlc); err != nil {
				log.Printf("Failed to save OHLC data: %v", err)
			}

			// Log the number of OHLC bars processed
			barCounter++
			if barCounter >= barsPerLog {

				logger.Printf("#### Processed %d OHLC Bars ####", barCounter)
				logger.Printf("#### Preparing for the Next Batch... ####")
				logger.Printf("==============================================================")

				barCounter = 0
			}
		}
	}()

	for _, symbol := range symbols {
		tickChan := make(chan binance.TradeData)
		client, err := binance.NewBinanceClient(symbol)
		if err != nil {
			return err
		}

		go client.StartTickStreaming(tickChan)

		aggregator := ohlc.NewAggregator(symbol, TIMEFRAME, tickChan, ohlcCompleteChan)
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
