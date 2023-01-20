package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/joeychilson/testdb/db/pgsql"
)

type Postgres struct {
	*pgsql.Queries
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
		Queries: pgsql.New(conn),
		pool:    conn,
	}, nil
}

func (p *Postgres) Tx(ctx context.Context, fn func(*pgsql.Queries) error) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	qtx := p.Queries.WithTx(tx)

	err = fn(qtx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
