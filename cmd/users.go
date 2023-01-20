package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"

	"github.com/joeychilson/testdb/db"
	"github.com/joeychilson/testdb/db/pgsql"
)

type User struct {
	FirstName string `fake:"{firstname}"`
	LastName  string `fake:"{lastname}"`
	Email     string `fake:"{email}"`
}

func main() {
	ctx := context.Background()

	_ = godotenv.Load()

	db, err := db.NewPostgres(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	var u User
	gofakeit.Struct(&u)

	params := pgsql.CreateUserParams{
		FirstName: sql.NullString{String: u.FirstName, Valid: u.FirstName != ""},
		LastName:  sql.NullString{String: u.LastName, Valid: u.LastName != ""},
		Email:     sql.NullString{String: u.Email, Valid: u.Email != ""},
	}

	id, err := db.Queries.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Printf("created user with id: %d", id)
}
