package tablecmd

import (
	"context"
	"fmt"

	"github.com/joeychilson/testdb/db"
	"github.com/spf13/cobra"
)

type TableCmd struct {
	pg  *db.Postgres
	sql *db.MySQL
}

func New(pg *db.Postgres, sql *db.MySQL) *TableCmd {
	return &TableCmd{pg: pg, sql: sql}
}

func (t *TableCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "table [flags]",
		Short: "List all columns in a table",
		RunE: func(cmd *cobra.Command, args []string) error {
			table, err := cmd.Flags().GetString("table")
			if err != nil {
				return err
			}
			return t.handleCommand(cmd.Context(), table)
		},
	}

	cmd.Flags().StringP("table", "t", "", "The table to get columns for")
	cmd.MarkFlagRequired("table")
	return cmd
}

func (t *TableCmd) handleCommand(ctx context.Context, table string) error {
	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = $1
	`

	columns, err := t.pg.Query(ctx, query, table)
	if err != nil {
		return err
	}

	fmt.Printf("Table %q:\n", table)
	fmt.Println("-------------------------------------------")
	for columns.Next() {
		var columnName, dataType string
		if err := columns.Scan(&columnName, &dataType); err != nil {
			return err
		}
		fmt.Printf("%s (%s) \n", columnName, dataType)
	}
	return nil
}
