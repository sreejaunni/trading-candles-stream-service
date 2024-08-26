package grpc

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"binance-candlestick-service/internal/business/ohlc"
	"binance-candlestick-service/proto"
	"github.com/stretchr/testify/assert"
)

type mockStream struct {
	proto.OHLCStreamer_StreamCandlesticksServer
	ch chan *proto.Candlestick
}

func (m *mockStream) Send(candlestick *proto.Candlestick) error {
	m.ch <- candlestick
	return nil
}

func TestStreamCandlesticks(t *testing.T) {
	server := NewServer()

	ch := make(chan *proto.Candlestick, 1)
	stream := &mockStream{ch: ch}

	req := &proto.CandlestickRequest{
		Symbol: "BTCUSDT",
	}

	go func() {
		err := server.StreamCandlesticks(req, stream)
		require.NoError(t, err, "StreamCandlesticks should not return an error")
	}()

	time.Sleep(100 * time.Millisecond)

	ohlcData := ohlc.OHLC{
		Symbol:    "BTCUSDT",
		Open:      100.0,
		High:      110.0,
		Low:       90.0,
		Close:     105.0,
		Volume:    1000,
		Timestamp: time.Now(),
	}
	server.BroadcastOHLC(ohlcData)

	select {
	case received := <-ch:

		t.Logf("Received data: %+v", received)
		assert.Equal(t, "BTCUSDT", received.Symbol, "Symbol should match")
		assert.Equal(t, 100.0, received.Open, "Open price should match")
		assert.Equal(t, 110.0, received.High, "High price should match")
		assert.Equal(t, 90.0, received.Low, "Low price should match")
		assert.Equal(t, 105.0, received.Close, "Close price should match")
		assert.Equal(t, int64(ohlcData.Timestamp.Unix()), received.Timestamp, "Timestamp should match")
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for candlestick")
	}
}

func TestBroadcastOHLC(t *testing.T) {
	server := NewServer()

	ch1 := make(chan ohlc.OHLC, 1)
	ch2 := make(chan ohlc.OHLC, 1)

	server.mu.Lock()
	server.ohlcStreams["BTCUSDT"] = append(server.ohlcStreams["BTCUSDT"], ch1, ch2)
	server.mu.Unlock()

	ohlcData := ohlc.OHLC{
		Symbol:    "BTCUSDT",
		Open:      100.0,
		High:      110.0,
		Low:       90.0,
		Close:     105.0,
		Volume:    1000,
		Timestamp: time.Now(),
	}
	server.BroadcastOHLC(ohlcData)

	select {
	case received := <-ch1:
		assert.Equal(t, ohlcData, received)
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for OHLC on ch1")
	}

	select {
	case received := <-ch2:
		assert.Equal(t, ohlcData, received)
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for OHLC on ch2")
	}
}
