package generate

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/brianvoe/gofakeit"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/joeychilson/testdb/cmd/testdb"
)

type Generate struct {
	config *testdb.Config
	table  string
	amount int
}

func New(config *testdb.Config) *ffcli.Command {
	g := Generate{
		config: config,
	}

	fs := flag.NewFlagSet("testdb generate", flag.ExitOnError)
	fs.StringVar(&g.table, "table", "", "the table name to generate data for")
	fs.IntVar(&g.amount, "amount", 1, "the amount of data to generate")

	return &ffcli.Command{
		Name:       "generate",
		ShortUsage: "testdb generate [flags] <key> <value data...>",
		ShortHelp:  "Generate data for a table in a database",
		FlagSet:    fs,
		Exec:       g.Exec,
	}
}

func (g *Generate) Exec(ctx context.Context, args []string) error {
	if g.table == "" {
		return errors.New("table is required")
	}

	if g.amount < 1 {
		return errors.New("amount must be greater than 0")
	}

	log.Printf("Generating %d rows for table %s", g.amount, g.table)

	type Column struct {
		Name  string
		Type  string
		Value interface{}
	}

	query := `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1`
	rows, err := g.config.Postgres.Query(ctx, query, g.table)
	if err != nil {
		log.Fatalf("failed to query database: %v", err)
	}

	columns := make([]Column, 0)
	for rows.Next() {
		var c Column
		if err := rows.Scan(&c.Name, &c.Type); err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}
		columns = append(columns, c)
	}

	gofakeit.Seed(0)

	for amount := 0; amount < g.amount; amount++ {
		for i := range columns {
			switch columns[i].Type {
			case "integer":
				columns[i].Value = gofakeit.Uint32()
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

		query := `INSERT INTO ` + g.table + ` (` + strings.Join(columnNames, ", ") + `) VALUES (` + strings.Join(placeholders, ", ") + `)`

		_, err := g.config.Postgres.Exec(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("failed to insert row into table %s: %v", g.table, err)
		}
	}

	log.Printf("Generated %d rows for table %s", g.amount, g.table)
	return nil
}
