package main

import (
	"binance-candlestick-service/config"
	db "binance-candlestick-service/internal/datastore"
	"context"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func main() {

	ctx := context.Background()

	cfg := config.NewConfig()

	dbConn, err := db.InitDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initilaise database: %v", err)
	}
	defer dbConn.Close()

	time.Local = time.UTC

	rootCmd := &cobra.Command{
		Use:   "root-cmd",
		Short: "Root CMD",
	}

	rootCmd.AddCommand(apiCMD(ctx, cfg, dbConn))
	rootCmd.AddCommand(resetDatabaseCMD(ctx, dbConn))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute root command: %v", err)
	}

	log.Printf("main() exit.")
}
