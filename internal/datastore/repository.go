package datastore

import (
	"binance-candlestick-service/internal/business/ohlc"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OHLCRepo struct {
	db DatabaseExecutor
}

type DatabaseExecutor interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

func NewOHLCRepo(db DatabaseExecutor) *OHLCRepo {
	return &OHLCRepo{db: db}
}

func (repo *OHLCRepo) SaveOHLC(ohlc ohlc.OHLC) error {
	query := `
        INSERT INTO candlesticks (
            symbol, 
            open, 
            high, 
            low, 
            close, 
            volume, 
            start_timestamp, 
            end_timestamp
        ) VALUES (
            :symbol, :open, :high, :low, :close, :volume, :start_timestamp, :end_timestamp
        )`

	_, err := repo.db.NamedExec(query, map[string]interface{}{
		"symbol":          ohlc.Symbol,
		"open":            ohlc.Open,
		"high":            ohlc.High,
		"low":             ohlc.Low,
		"close":           ohlc.Close,
		"volume":          ohlc.Volume,
		"start_timestamp": ohlc.StartTimeStamp,
		"end_timestamp":   ohlc.EndTimeStamp,
	})
	if err != nil {
		return fmt.Errorf("error saving OHLC data: %v", err)
	}
	return nil
}
