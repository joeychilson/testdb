package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"

	rootcmd "github.com/joeychilson/testdb/cmd/root"
	"github.com/joeychilson/testdb/db"
)

func main() {
	_ = godotenv.Load()

	postgres, _ := db.NewPostgres(context.Background(), os.Getenv("POSTGRES_URL"))
	mysql, _ := db.NewMySQL(context.Background(), os.Getenv("MYSQL_URL"))

	config := &rootcmd.Config{
		Postgres: postgres,
		MySQL:    mysql,
	}

	rootcmd := rootcmd.New(config)
	if err := rootcmd.Execute(); err != nil {
		os.Exit(1)
	}
}
