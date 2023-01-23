package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	rootcmd "github.com/joeychilson/testdb/cmd/root"
	"github.com/joeychilson/testdb/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
		os.Exit(1)
	}

	db, err := db.NewPostgres(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Printf("failed to connect to postgres: %v", err)
		os.Exit(1)
	}

	config := &rootcmd.Config{
		Database: db,
	}

	rootcmd := rootcmd.New(config)
	if err := rootcmd.Execute(); err != nil {
		log.Printf("failed to execute root command: %v", err)
		os.Exit(1)
	}
}
