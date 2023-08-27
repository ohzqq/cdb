package cdb

import (
	"log"
)

// AddFormat adds a format to an existing book
func AddFormat(lib, id, pos string, args ...Opt) *Command {
	cmd, err := NewCommand(lib, args...)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Opt(Cmd("add_format"), PositionalArgs(id, pos))
	return cmd
}

func DontReplace() Opt {
	return Flags("--dont-replace")
}
