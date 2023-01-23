package autogencmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"

	"github.com/joeychilson/testdb/db"
)

type AutoGenCmd struct {
	db *db.Postgres
}

func New(db *db.Postgres) *AutoGenCmd {
	return &AutoGenCmd{db: db}
}

func (g *AutoGenCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autogen [flags]",
		Short: "Auto generate fake data for the provided table",
		RunE: func(cmd *cobra.Command, args []string) error {
			table, err := cmd.Flags().GetString("table")
			if err != nil {
				return err
			}

			rows, err := cmd.Flags().GetInt("rows")
			if err != nil {
				return err
			}

			return g.handleCmd(cmd.Context(), table, rows)
		},
	}

	cmd.Flags().StringP("table", "t", "", "The table to generate data for")
	cmd.Flags().IntP("rows", "r", 1, "The number of rows to generate")

	cmd.MarkFlagRequired("table")
	return cmd
}

func (g *AutoGenCmd) handleCmd(ctx context.Context, table string, rows int) error {
	log.Printf("Generating %d rows for table %s", rows, table)

	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = $1`

	results, err := g.db.Query(ctx, query, table)
	if err != nil {
		log.Fatalf("failed to query database: %v", err)
	}
	defer results.Close()

	columns := make([]*column, 0)
	for results.Next() {
		c := &column{}
		if err := results.Scan(&c.name, &c.dataType); err != nil {
			continue
		}
		if c.name == "id" {
			continue
		}
		columns = append(columns, c)
	}

	for amount := 0; amount < rows; amount++ {
		columnNames := make([]string, len(columns))
		placeholders := make([]string, len(columns))
		values := make([]any, len(columns))

		for i := range columns {
			columnNames[i] = columns[i].name
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			values[i] = genValue(columns[i].dataType)
		}
		query := `
			INSERT INTO ` + table + ` 
				(` + strings.Join(columnNames, ", ") + `) 
			VALUES 
				(` + strings.Join(placeholders, ", ") + `)`

		_, err := g.db.Exec(ctx, query, values...)
		if err != nil {
			return fmt.Errorf("failed to insert row into table %s: %v", table, err)
		}
	}

	log.Printf("Generated %d rows for table %s", rows, table)
	return nil
}

type column struct {
	name     string
	dataType string
	value    interface{}
}

var fakers = map[string]func() any{
	"uuid":                        func() any { return gofakeit.UUID() },
	"character varying":           func() any { return gofakeit.BuzzWord() },
	"text":                        func() any { return gofakeit.BuzzWord() },
	"integer":                     func() any { return gofakeit.Int8() },
	"bigint":                      func() any { return gofakeit.Int32() },
	"double precision":            func() any { return gofakeit.Float32() },
	"numeric":                     func() any { return gofakeit.Float32() },
	"boolean":                     func() any { return gofakeit.Bool() },
	"inet":                        func() any { return gofakeit.IPv4Address() },
	"macaddr":                     func() any { return gofakeit.MacAddress() },
	"bytea":                       func() any { return gofakeit.Letter() },
	"json":                        func() any { return genJSON() },
	"jsonb":                       func() any { return genJSON() },
	"xml":                         func() any { return genXML() },
	"date":                        func() any { return gofakeit.Date() },
	"time with time zone":         func() any { return gofakeit.Date() },
	"time without time zone":      func() any { return gofakeit.Date() },
	"timestamp with time zone":    func() any { return gofakeit.Date() },
	"timestamp without time zone": func() any { return gofakeit.Date() },
}

func genValue(dataType string) any {
	gofakeit.Seed(0)

	genFunc, ok := fakers[dataType]
	if !ok {
		log.Printf("unknown data type: %s", dataType)
		return nil
	}
	return genFunc()
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
