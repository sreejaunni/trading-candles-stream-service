package binance

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type BinanceClient struct {
	Conn *websocket.Conn
}

func NewBinanceClient(symbol string) (*BinanceClient, error) {
	u := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws/" + symbol + "@aggTrade"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &BinanceClient{Conn: conn}, nil
}

func (c *BinanceClient) StartTickStreaming(tickChan chan<- TradeData) {
	defer c.Conn.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var trade TradeData
		err = json.Unmarshal(message, &trade)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		tickChan <- trade
	}
}
