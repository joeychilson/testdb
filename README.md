# testdb

## Overview

A repo for creating test databases in PostgreSQL and/or MySQL, with Go clients for testing and debugging code.

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
# contains POSTGRES_URL, and MYSQL_URL for containers
cp .env.example .env
```

### Migration

```bash
# create new migration file
make pgdbnew name={migration_name} or make mydbnew name={migration_name}

# migrates database to latest version
make pgdbup or make mydbup

# migrates database to previous version
make pgdbdown or make mydbdown

```

### Container

```bash
# start docker container
make pgup or make myup

# stop docker container
make pgdown or make mydown

# destory docker container
make pgstop or make mystop
```

### Generate Go code

```bash
# generate Go code for both PostgreSQL and MySQL queries
make sqlc
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
```
