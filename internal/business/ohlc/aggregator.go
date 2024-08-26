package ohlc

import (
	"binance-candlestick-service/internal/pkg/binance"
	"log"
	"strconv"
	"time"
)

type Aggregator struct {
	symbol           string
	timeframe        time.Duration
	tickChan         chan binance.TradeData
	ohlcCompleteChan chan OHLC

	openPrice     float64
	highPrice     float64
	lowPrice      float64
	closePrice    float64
	volume        float64
	intervalStart time.Time
}

func NewAggregator(symbol string, timeframe time.Duration, tickChan chan binance.TradeData, ohlcCompleteChan chan OHLC) *Aggregator {
	return &Aggregator{
		symbol:           symbol,
		timeframe:        timeframe,
		tickChan:         tickChan,
		ohlcCompleteChan: ohlcCompleteChan,
		intervalStart:    time.Now().Truncate(timeframe),
	}
}

func (a *Aggregator) Start() {
	for tick := range a.tickChan {
		now := time.Now()
		price, err := strconv.ParseFloat(tick.Price, 64)
		if err != nil {
			log.Printf("Error parsing price '%s': %v", tick.Price, err)
			continue
		}
		quantity, err := strconv.ParseFloat(tick.Quantity, 64)
		if err != nil {
			log.Printf("Error parsing quantity '%s': %v", tick.Quantity, err)
			continue
		}

		if now.Sub(a.intervalStart) >= a.timeframe {
			// Complete the current OHLC bar
			a.completeBar()
			// Start a new bar
			a.intervalStart = now.Truncate(a.timeframe)
			a.openPrice = price
			a.highPrice = price
			a.lowPrice = price
			a.volume = 0
		}

		// Update OHLC values
		if price > a.highPrice {
			a.highPrice = price
		}
		if price < a.lowPrice || a.lowPrice == 0 {
			a.lowPrice = price
		}
		a.closePrice = price
		a.volume += quantity
	}
}

func (a *Aggregator) completeBar() {
	ohlc := OHLC{
		Symbol:         a.symbol,
		Open:           a.openPrice,
		High:           a.highPrice,
		Low:            a.lowPrice,
		Close:          a.closePrice,
		Volume:         a.volume,
		Timestamp:      a.intervalStart.Add(a.timeframe).UTC(),
		StartTimeStamp: a.intervalStart.UTC(),
		EndTimeStamp:   a.intervalStart.Add(a.timeframe).UTC(),
	}
	a.ohlcCompleteChan <- ohlc
}
