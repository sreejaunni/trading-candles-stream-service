package datastore

import (
	"binance-candlestick-service/config"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDatabase(cfg config.DBConfig) (*sqlx.DB, error) {

	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name)

	db, err = sqlx.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
