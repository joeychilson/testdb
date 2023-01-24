package rootcmd

import (
	"github.com/spf13/cobra"

	autogencmd "github.com/joeychilson/testdb/cmd/autogen"
	columnscmd "github.com/joeychilson/testdb/cmd/columns"
	gencmd "github.com/joeychilson/testdb/cmd/gen"
	tablescmd "github.com/joeychilson/testdb/cmd/tables"
	"github.com/joeychilson/testdb/db"
)

type Config struct {
	Database *db.Postgres
}

func New(config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testdb",
		Short: "The TestDB CLI",
	}

	autogencmd := autogencmd.New(config.Database)
	cmd.AddCommand(autogencmd.Cmd())

	columnscmd := columnscmd.New(config.Database)
	cmd.AddCommand(columnscmd.Cmd())

	genCmd := gencmd.New(config.Database)
	cmd.AddCommand(genCmd.Cmd())

	tablescmd := tablescmd.New(config.Database)
	cmd.AddCommand(tablescmd.Cmd())

	return cmd
}
