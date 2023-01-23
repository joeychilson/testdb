package tablescmd

import (
	"context"
	"fmt"

	"github.com/joeychilson/testdb/db"
	"github.com/spf13/cobra"
)

type TablesCmd struct {
	db *db.Postgres
}

func New(db *db.Postgres) *TablesCmd {
	return &TablesCmd{db: db}
}

func (t *TablesCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tables [flags]",
		Short: "List all tables in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.handleCmd(cmd.Context())
		},
	}
	return cmd
}

func (t *TablesCmd) handleCmd(ctx context.Context) error {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`

	tables, err := t.db.Query(ctx, query)
	if err != nil {
		return err
	}

	fmt.Println("Tables:")
	fmt.Println("-------------------------------")
	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			return err
		}
		fmt.Println(tableName)
	}
	return nil
}
