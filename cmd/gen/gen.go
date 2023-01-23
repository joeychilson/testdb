package gen

import (
	"github.com/joeychilson/testdb/db"
	"github.com/spf13/cobra"
)

type GenCmd struct {
	db *db.Postgres
}

func New(db *db.Postgres) *GenCmd {
	return &GenCmd{db: db}
}

func (g *GenCmd) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen [flags]",
		Short: "Generate realistic fake data for the provided schema",
	}

	musicGenCmd := NewMusicGen(g.db)
	cmd.AddCommand(musicGenCmd.MusicGenCommand())
	return cmd
}
