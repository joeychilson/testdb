package testdb

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/joeychilson/testdb/db"
)

type Config struct {
	Postgres *db.Postgres
	MySQL    *db.MySQL
}

func New() (*ffcli.Command, *Config) {
	var cfg Config

	fs := flag.NewFlagSet("testdb", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "testdb",
		ShortUsage: "testdb [flags] <subcommand> [flags] [<arg>...]",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}, &cfg
}

func (c *Config) Exec(context.Context, []string) error {
	return flag.ErrHelp
}
