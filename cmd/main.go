package main

import (
	"binance-candlestick-service/config"
	db "binance-candlestick-service/internal/datastore"
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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
