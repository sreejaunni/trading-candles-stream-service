package utils

import (
	"binance-candlestick-service/internal/business/ohlc"
	"fmt"
	"time"
)

func PrintOHLC(ohlc ohlc.OHLC) {
	fmt.Printf("%s\t%.10f\t%.10f\t%.10f\t%.10f\t%.10f\t%s\t%s\n",
		ohlc.Symbol,
		ohlc.Open,
		ohlc.High,
		ohlc.Low,
		ohlc.Close,
		ohlc.Volume,
		ohlc.StartTimeStamp.Format(time.RFC3339),
		ohlc.EndTimeStamp.Format(time.RFC3339),
	)
}
