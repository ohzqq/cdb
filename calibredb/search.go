package calibredb

import "log"

func Search(lib, pos string, args ...Opt) *Command {
	cmd, err := NewCommand(lib, args...)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Opt(Cmd("search"), PositionalArgs(pos))
	return cmd
}

func Limit(l string) Opt {
	return Flags("--limit", l)
}
