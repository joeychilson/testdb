package tablecmd

import (
	"context"
	"fmt"

	"github.com/joeychilson/testdb/db"
	"github.com/spf13/cobra"
)

type ColumnsCmd struct {
	db *db.Postgres
}

func New(db *db.Postgres) *ColumnsCmd {
	return &ColumnsCmd{db: db}
}

func (c *ColumnsCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "columns [flags]",
		Short: "List all columns in a table",
		RunE: func(cmd *cobra.Command, args []string) error {
			table, err := cmd.Flags().GetString("table")
			if err != nil {
				return err
			}
			return c.handleCmd(cmd.Context(), table)
		},
	}

	cmd.Flags().StringP("table", "t", "", "The table to get columns for")
	cmd.MarkFlagRequired("table")
	return cmd
}

func (c *ColumnsCmd) handleCmd(ctx context.Context, table string) error {
	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_name = $1
	`

	columns, err := c.db.Query(ctx, query, table)
	if err != nil {
		return err
	}

	fmt.Printf("Columns for %q:\n", table)
	fmt.Println("-------------------------------")
	for columns.Next() {
		var columnName, dataType string
		if err := columns.Scan(&columnName, &dataType); err != nil {
			return err
		}
		fmt.Printf("%s (%s) \n", columnName, dataType)
	}
	return nil
}
