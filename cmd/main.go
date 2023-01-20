package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/joeychilson/testdb/cmd/generate"
	"github.com/joeychilson/testdb/cmd/testdb"
	"github.com/joeychilson/testdb/db"
)

func main() {
	var cmd, config = testdb.New()

	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during load .env file: %v\n", err)
		os.Exit(1)
	}

	postgres, _ := db.NewPostgres(context.Background(), os.Getenv("POSTGRES_URL"))
	mysql, _ := db.NewMySQL(context.Background(), os.Getenv("MYSQL_URL"))

	config.Postgres = postgres
	config.MySQL = mysql

	cmd.Subcommands = []*ffcli.Command{
		generate.New(config),
	}

	if err := cmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error during parse command: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "error during run command: %v\n", err)
		os.Exit(1)
	}
}
