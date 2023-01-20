package gencmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"

	"github.com/joeychilson/testdb/db"
)

type Config struct {
	Postgres *db.Postgres
	MySQL    *db.MySQL
}

type GenCmd struct {
	cfg *Config
}

func New(cfg *Config) *GenCmd {
	return &GenCmd{cfg: cfg}
}

func (g *GenCmd) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen [flags]",
		Short: "Generate fake data for a table",
		RunE: func(cmd *cobra.Command, args []string) error {
			table, err := cmd.Flags().GetString("table")
			if err != nil {
				return err
			}

			rows, err := cmd.Flags().GetInt("rows")
			if err != nil {
				return err
			}

			return g.handleCommand(cmd.Context(), table, rows)
		},
	}

	cmd.Flags().StringP("table", "t", "", "The table to generate data for")
	cmd.Flags().IntP("rows", "r", 1, "The number of rows to generate")

	cmd.MarkFlagRequired("table")
	return cmd
}

func (g *GenCmd) handleCommand(ctx context.Context, table string, rows int) error {
	log.Printf("Generating %d rows for table %s", rows, table)

	type Column struct {
		Name  string
		Type  string
		Value interface{}
	}

	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = $1
	`

	results, err := g.cfg.Postgres.Query(ctx, query, table)
	if err != nil {
		log.Fatalf("failed to query database: %v", err)
	}

	columns := make([]Column, 0)
	for results.Next() {
		var c Column
		if err := results.Scan(&c.Name, &c.Type); err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}
		if c.Name == "id" {
			continue
		}
		columns = append(columns, c)
	}

	gofakeit.Seed(0)

	for amount := 0; amount < rows; amount++ {
		for i := range columns {
			switch columns[i].Type {
			case "integer":
				columns[i].Value = gofakeit.Uint8()
			case "numeric":
				columns[i].Value = gofakeit.Float32()
			case "character varying":
				columns[i].Value = gofakeit.BuzzWord()
			case "text":
				columns[i].Value = gofakeit.BuzzWord()
			case "boolean":
				columns[i].Value = gofakeit.Bool()
			case "date":
				columns[i].Value = gofakeit.Date()
			case "timestamp with time zone":
				columns[i].Value = gofakeit.Date()
			case "timestamp without time zone":
				columns[i].Value = gofakeit.Date()
			default:
				log.Printf("unknown data type: %s", columns[i].Type)
			}
		}

		columnNames := make([]string, len(columns))
		placeholders := make([]string, len(columns))
		values := make([]interface{}, len(columns))

		for i := range columns {
			columnNames[i] = columns[i].Name
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			values[i] = columns[i].Value
		}

		query := `
			INSERT INTO ` + table + ` 
				(` + strings.Join(columnNames, ", ") + `) 
			VALUES 
				(` + strings.Join(placeholders, ", ") + `)
		`

		_, err := g.cfg.Postgres.Exec(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("failed to insert row into table %s: %v", table, err)
		}
	}

	log.Printf("Generated %d rows for table %s", rows, table)
	return nil
}
