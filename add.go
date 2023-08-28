package cdb

import (
	"log"
)

func AddCmd(lib, pos string, args ...Opt) *Command {
	cmd, err := NewCommand(lib, args...)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Opt(Cmd("add"), PositionalArgs(pos))

	return cmd
}

func Automerge() Opt {
	return SetFlags("--automerge")
}

func AddCover(file string) Opt {
	if ok := FileExist(file); !ok {
		err := ErrFileNotExist(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	return SetFlags("--cover", file)
}

func Duplicates() Opt {
	return SetFlags("--duplicates")
}

func Empty() Opt {
	return SetFlags("--empty")
}
