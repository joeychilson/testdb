package rootcmd

import (
	"github.com/spf13/cobra"

	gencmd "github.com/joeychilson/testdb/cmd/gen"
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
	return cmd
}
