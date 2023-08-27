package command

import (
	"log"

	"github.com/ohzqq/cdb"
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

func Cover(file string) Opt {
	if ok := cdb.FileExist(file); !ok {
		err := cdb.ErrFileNotExist(file)
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
