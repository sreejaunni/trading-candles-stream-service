package grpc

import (
	"sync"

	"binance-candlestick-service/internal/business/ohlc"
	"binance-candlestick-service/proto"
)

type Server struct {
	proto.UnimplementedOHLCStreamerServer
	ohlcStreams map[string][]chan ohlc.OHLC
	mu          sync.Mutex
}

func NewServer() *Server {
	return &Server{
		ohlcStreams: make(map[string][]chan ohlc.OHLC),
	}
}

func (s *Server) StreamCandlesticks(req *proto.CandlestickRequest, stream proto.OHLCStreamer_StreamCandlesticksServer) error {
	s.mu.Lock()
	ch := make(chan ohlc.OHLC)
	s.ohlcStreams[req.Symbol] = append(s.ohlcStreams[req.Symbol], ch)
	s.mu.Unlock()

	for ohlc := range ch {
		err := stream.Send(&proto.Candlestick{
			Symbol:    ohlc.Symbol,
			Open:      ohlc.Open,
			High:      ohlc.High,
			Low:       ohlc.Low,
			Close:     ohlc.Close,
			Volume:    ohlc.Volume,
			Timestamp: ohlc.Timestamp.Unix(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) BroadcastOHLC(ohlc ohlc.OHLC) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channels, ok := s.ohlcStreams[ohlc.Symbol]; ok {
		for _, ch := range channels {
			ch <- ohlc
		}
	}
}
