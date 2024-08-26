package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func resetDatabaseCMD(ctx context.Context, dbConn *sqlx.DB) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "reset-database",
		Short: "Reset the database schema",
		Run: func(cmd *cobra.Command, args []string) {

			sqlFilePath := filepath.Join("schema", "schema.sql")
			sqlBytes, err := os.ReadFile(sqlFilePath)
			if err != nil {
				log.Fatal(ctx, "failed to read SQL file: ", err)
			}
			sql := string(sqlBytes)

			_, err = dbConn.ExecContext(ctx, sql)
			if err != nil {
				log.Fatalf("failed to execute SQL file: %v", err)
			}

			log.Printf("database schema reset successfully.")
		},
	}

	return cmd
}
