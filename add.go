package cdb

import (
	"log"
)

func Add(lib, pos string, args ...Opt) *Command {
	cmd, err := NewCommand(lib, args...)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Opt(Cmd("add"), PositionalArgs(pos))

	return cmd
}

func Automerge() Opt {
	return Flags("--automerge")
}

func AddCover(file string) Opt {
	if ok := FileExist(file); !ok {
		err := ErrFileNotExist(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	return Flags("--cover", file)
}

func Duplicates() Opt {
	return Flags("--duplicates")
}

func Empty() Opt {
	return Flags("--empty")
}
