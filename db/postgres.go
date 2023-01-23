package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/joeychilson/testdb/db/sqlc"
)

type Postgres struct {
	*sqlc.Queries
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connStr string) (*Postgres, error) {
	if connStr == "" {
		return nil, fmt.Errorf("connection string not set")
	}
	conn, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		Queries: sqlc.New(conn),
		pool:    conn,
	}, nil
}

func (p *Postgres) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, query, args...)
}

func (p *Postgres) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, query, args...)
}

func (p *Postgres) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, query, args...)
}

func (p *Postgres) Tx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
