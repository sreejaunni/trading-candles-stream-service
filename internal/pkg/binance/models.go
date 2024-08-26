package binance

type TradeData struct {
	EventType        string `json:"e"`
	EventTime        int64  `json:"E"`
	Symbol           string `json:"s"`
	AggregateTradeID int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     int64  `json:"f"`
	LastTradeID      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
	Ignore           bool   `json:"M"`
}
