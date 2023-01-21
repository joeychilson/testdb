package rootcmd

import (
	"github.com/spf13/cobra"

	gencmd "github.com/joeychilson/testdb/cmd/gen"
	tablecmd "github.com/joeychilson/testdb/cmd/table"
	tablescmd "github.com/joeychilson/testdb/cmd/tables"
	"github.com/joeychilson/testdb/db"
)

type Config struct {
	Postgres *db.Postgres
	MySQL    *db.MySQL
}

func New(config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testdb",
		Short: "The TestDB CLI",
	}

	gencmd := gencmd.New(config.Postgres, config.MySQL)
	cmd.AddCommand(gencmd.Command())

	tablescmd := tablescmd.New(config.Postgres, config.MySQL)
	cmd.AddCommand(tablescmd.Command())

	tablecmd := tablecmd.New(config.Postgres, config.MySQL)
	cmd.AddCommand(tablecmd.Command())
	return cmd
}
