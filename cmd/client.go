package main

import (
	"binance-candlestick-service/config"
	"binance-candlestick-service/internal/app"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"log"
)

func clientCMD(ctx context.Context, cfg *config.Config, db *sqlx.DB) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start-client",
		Short: "run client",
		Run: func(cmd *cobra.Command, args []string) {

			log.Print("...printing candlesticks...")

			if err := app.StreamCandles(cfg); err != nil {
				log.Fatalf("Application error: %v", err)
			}
		},
	}

	return cmd
}
