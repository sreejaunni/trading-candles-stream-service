package main

import (
	"binance-candlestick-service/config"
	"binance-candlestick-service/internal/app"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func apiCMD(ctx context.Context, cfg *config.Config, db *sqlx.DB) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "server",
		Short: "run application server",
		Run: func(cmd *cobra.Command, args []string) {

			log.Print("...Aggregating candlesticks...")

			if err := app.Run(cfg, db); err != nil {
				log.Fatalf("Application error: %v", err)
			}
		},
	}

	return cmd
}
