package app

import (
	"binance-candlestick-service/config"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"binance-candlestick-service/proto"
	"google.golang.org/grpc"
)

func StreamCandles(config *config.Config) error {
	serverAddress := fmt.Sprintf("%s%s", "localhost", config.Server.Port)

	symbols := config.Symbols
	// Establish a connection to the gRPC server
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewOHLCStreamerClient(conn)

	var wg sync.WaitGroup

	// Start streaming for each symbol
	for _, symbol := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			streamCandlesticks(client, symbol)
		}(symbol)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press Ctrl+C to exit...")
	<-signalChan

	fmt.Println("Received interrupt signal, shutting down...")

	// Wait for all streaming goroutines to complete
	wg.Wait()
	return nil
}

func streamCandlesticks(client proto.OHLCStreamerClient, symbol string) {
	stream, err := client.StreamCandlesticks(context.Background(), &proto.CandlestickRequest{
		Symbol: symbol,
	})
	if err != nil {
		log.Printf("Failed to start streaming for symbol %s: %v", symbol, err)
		return
	}

	fmt.Printf("Started streaming for symbol: %s\n", symbol)

	for {
		candlestick, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive data for symbol %s: %v", symbol, err)
			return
		}
		fmt.Printf("Received candlestick for %s: %+v\n", symbol, candlestick)
	}
}
