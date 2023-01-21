package tablescmd

import (
	"context"
	"fmt"

	"github.com/joeychilson/testdb/db"
	"github.com/spf13/cobra"
)

type TablesCmd struct {
	pg  *db.Postgres
	sql *db.MySQL
}

func New(pg *db.Postgres, sql *db.MySQL) *TablesCmd {
	return &TablesCmd{pg: pg, sql: sql}
}

func (t *TablesCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tables [flags]",
		Short: "List all tables in the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.handleCommand(cmd.Context())
		},
	}
	return cmd
}

func (t *TablesCmd) handleCommand(ctx context.Context) error {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`

	tables, err := t.pg.Query(ctx, query)
	if err != nil {
		return err
	}

	fmt.Println("Tables:")
	fmt.Println("-------")
	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			return err
		}
		fmt.Println(tableName)
	}
	return nil
}
