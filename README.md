# testdb

A repo for creating test databases in PostgreSQL and/or MySQL, with tools for testing and debugging code.

**This is mostly for personal use, but feel free to use whatever you want in the repo.**

## Features

- [x] Docker containers for PostgreSQL and MySQL
- [x] Migrations with dbmate
- [x] Generate Go code from queries with sqlc
- [x] Automatically generate fake data for any PostgreSQL table (MySQL coming soon)

## Requirements

- [dbmate](https://github.com/amacneil/dbmate) (for migrations)
- [sqlc](https://github.com/kyleconroy/sqlc) (for generating Go code from SQL queries)

## Usage

### Clone

```bash
git clone https://github.com/joeychilson/testdb
```

### Setup

```bash
# contains POSTGRES_URL for containers
cp .env.example .env
```

### Migration

```bash
# create new migration file
make dbnew name={migration_name}

# migrates database to latest version
make dbup

# migrates database to previous version
make dbdown

```

### Container

```bash
# start docker container
make up

# stop docker container
make down

# destory docker container
make stop
```

### Generate Go code

```bash
# generate Go code for both PostgreSQL and MySQL queries
make sqlc
```

### Generate fake data

```bash
# generate fake data for any postgres sql tables
# only supports postgres for now
go run cmd/main.go autogen -t table_name -r 100

# generate realistic date for music schema
# this will generate 2 artists, 3 albums per artist, and 8 songs per album
go run cmd/main.go gen music -r 2 -a 3 -s 8
```

### Example using generated Go code

```go
package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"

	"github.com/joeychilson/testdb/db"
	"github.com/joeychilson/testdb/db/sqlc"
	"github.com/joeychilson/testdb/gen"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load()

	db, err := db.NewPostgres(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	gofakeit.Seed(0)

	var genArtist *gen.Artist
	gofakeit.Struct(&genArtist)

	artist := sqlc.CreateArtistParams{
		Name:  genArtist.Name,
		Image: sql.NullString{String: genArtist.Image, Valid: genArtist.Image != ""},
	}

	artistsID, err := g.db.CreateArtist(ctx, artist)
	if err != nil {
		log.Fatalf("failed to create artist: %v", err)
	}

	log.Printf("Created artist %s with ID", genArtist.Name, artistsID)
}
```
