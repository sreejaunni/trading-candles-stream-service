package ohlc

import (
	"binance-candlestick-service/internal/pkg/binance"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAggregator_PriceParsingError(t *testing.T) {
	tickChan := make(chan binance.TradeData, 1)
	ohlcCompleteChan := make(chan OHLC, 1)
	aggregator := NewAggregator("BTCUSDT", time.Minute, tickChan, ohlcCompleteChan)

	go aggregator.Start()

	tickChan <- binance.TradeData{Price: "invalid_price", Quantity: "1"}

	close(tickChan)

	time.Sleep(time.Second)

	close(ohlcCompleteChan)

	_, ok := <-ohlcCompleteChan
	if ok {
		t.Error("Expected no OHLC bar to be completed due to price parsing error")
	}
}

func TestAggregator_NewBarInterval(t *testing.T) {
	tickChan := make(chan binance.TradeData, 10)
	ohlcCompleteChan := make(chan OHLC, 10)
	aggregator := NewAggregator("BTCUSDT", time.Millisecond*500, tickChan, ohlcCompleteChan)

	go aggregator.Start()

	tickChan <- binance.TradeData{Price: "50000", Quantity: "1", Symbol: "BTCUSDT"}
	time.Sleep(time.Millisecond * 600)
	tickChan <- binance.TradeData{Price: "51000", Quantity: "2", Symbol: "BTCUSDT"}
	time.Sleep(time.Millisecond * 600)

	close(tickChan)
	close(ohlcCompleteChan)

	var ohlc OHLC
	select {
	case ohlc = <-ohlcCompleteChan:
	case <-time.After(time.Second * 2):
		t.Fatal("Timed out waiting for the first OHLC bar")
	}

	assert.Equal(t, "BTCUSDT", ohlc.Symbol)

	assert.Equal(t, 50000.0, ohlc.High)
	assert.Equal(t, 50000.0, ohlc.Low)
	assert.Equal(t, 50000.0, ohlc.Close)
	assert.Equal(t, 1.0, ohlc.Volume)

}

func TestAggregator_InvalidPriceQuantity(t *testing.T) {
	tickChan := make(chan binance.TradeData, 10)
	ohlcCompleteChan := make(chan OHLC, 10)
	aggregator := NewAggregator("BTCUSDT", time.Millisecond*500, tickChan, ohlcCompleteChan)

	go aggregator.Start()

	tickChan <- binance.TradeData{Price: "50000", Quantity: "0", Symbol: "BTCUSDT"}
	time.Sleep(time.Millisecond * 600)
	tickChan <- binance.TradeData{Price: "51000", Quantity: "0", Symbol: "BTCUSDT"}
	time.Sleep(time.Millisecond * 600)

	close(tickChan)
	close(ohlcCompleteChan)

	var ohlc OHLC
	select {
	case ohlc = <-ohlcCompleteChan:
	case <-time.After(time.Second * 2):
		t.Fatal("Timed out waiting for the first OHLC bar")
	}

	assert.Equal(t, "BTCUSDT", ohlc.Symbol)

	assert.Equal(t, 50000.0, ohlc.High)
	assert.Equal(t, 50000.0, ohlc.Low)
	assert.Equal(t, 50000.0, ohlc.Close)
	assert.Equal(t, 0.0, ohlc.Volume)
}

func TestAggregator_GracefulShutdown(t *testing.T) {
	tickChan := make(chan binance.TradeData, 10)
	ohlcCompleteChan := make(chan OHLC, 10)
	aggregator := NewAggregator("BTCUSDT", time.Minute, tickChan, ohlcCompleteChan)

	go aggregator.Start()

	// Send a tick and then close the channel to simulate a shutdown
	tickChan <- binance.TradeData{Price: "50000", Quantity: "1"}
	close(tickChan)
	time.Sleep(time.Second) // Allow time for processing

	// Close OHLC channel and check if shutdown is graceful
	close(ohlcCompleteChan)
	_, ok := <-ohlcCompleteChan
	if ok {
		t.Error("Expected OHLC channel to be closed after aggregator shutdown")
	}
}
