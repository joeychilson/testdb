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

	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = $1
	`

	results, err := g.cfg.Postgres.Query(ctx, query, table)
	if err != nil {
		log.Fatalf("failed to query database: %v", err)
	}

	columns := make([]*Column, 0)
	for results.Next() {
		c := &Column{}
		if err := results.Scan(&c.Name, &c.Type); err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}
		if c.Name == "id" {
			continue
		}
		columns = append(columns, c)
	}

	for amount := 0; amount < rows; amount++ {
		for i := range columns {
			genValue(columns[i])
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

type Column struct {
	Name  string
	Type  string
	Value interface{}
}

func genValue(column *Column) {
	gofakeit.Seed(0)

	fakers := map[string]func(){
		"uuid":                        func() { column.Value = gofakeit.UUID() },
		"character varying":           func() { column.Value = gofakeit.BuzzWord() },
		"text":                        func() { column.Value = gofakeit.BuzzWord() },
		"integer":                     func() { column.Value = gofakeit.Int8() },
		"bigint":                      func() { column.Value = gofakeit.Int32() },
		"double precision":            func() { column.Value = gofakeit.Float32() },
		"numeric":                     func() { column.Value = gofakeit.Float32() },
		"boolean":                     func() { column.Value = gofakeit.Bool() },
		"inet":                        func() { column.Value = gofakeit.IPv4Address() },
		"macaddr":                     func() { column.Value = gofakeit.MacAddress() },
		"bytea":                       func() { column.Value = gofakeit.Letter() },
		"json":                        func() { column.Value = genJSON() },
		"jsonb":                       func() { column.Value = genJSON() },
		"xml":                         func() { column.Value = genXML() },
		"date":                        func() { column.Value = gofakeit.Date() },
		"time with time zone":         func() { column.Value = gofakeit.Date() },
		"time without time zone":      func() { column.Value = gofakeit.Date() },
		"timestamp with time zone":    func() { column.Value = gofakeit.Date() },
		"timestamp without time zone": func() { column.Value = gofakeit.Date() },
	}

	genFunc, ok := fakers[column.Type]
	if !ok {
		log.Printf("unknown data type: %s", column.Type)
		return
	}
	genFunc()
}

func genJSON() []byte {
	json, err := gofakeit.JSON(&gofakeit.JSONOptions{
		Type:     "array",
		RowCount: 1,
		Indent:   true,
		Fields: []gofakeit.Field{
			{
				Name:     "id",
				Function: "uuid",
			},
			{
				Name:     "name",
				Function: "buzzword",
			},
		},
	})
	if err != nil {
		log.Printf("error generating json: %s", err)
		return nil
	}
	return json
}

func genXML() []byte {
	xml, err := gofakeit.XML(&gofakeit.XMLOptions{
		Type:          "single",
		RootElement:   "xml",
		RecordElement: "record",
		RowCount:      2,
		Indent:        true,
		Fields: []gofakeit.Field{
			{
				Name:     "id",
				Function: "uuid",
			},
			{
				Name:     "name",
				Function: "buzzword",
			},
		},
	})
	if err != nil {
		log.Printf("error generating xml: %s", err)
		return nil
	}
	return xml
}
