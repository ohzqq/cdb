package cdb

import (
	"log"
)

func Remove(lib string, pos ...string) *Command {
	cmd, err := NewCommand(lib)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Opt(Cmd("remove"), PositionalArgs(pos...))

	return cmd
}

func Permanent() Opt {
	return SetFlags("--permanent")
}
