package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joeychilson/testdb/db/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	*mysql.Queries
	db *sql.DB
}

func NewMySQL(ctx context.Context, connStr string) (*MySQL, error) {
	if connStr == "" {
		return nil, fmt.Errorf("connection string not set")
	}
	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &MySQL{
		Queries: mysql.New(conn),
		db:      conn,
	}, nil
}
