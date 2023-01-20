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

func (m *MySQL) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *MySQL) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *MySQL) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return m.db.QueryRowContext(ctx, query, args...)
}

func (m *MySQL) Tx(ctx context.Context, fn func(*mysql.Queries) error) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	qtx := m.WithTx(tx)

	err = fn(qtx)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
