package calibredb

import (
	"log"

	"github.com/ohzqq/cdb"
	"github.com/spf13/viper"
)

func Add(lib, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("add"), PositionalArgs(pos))
	if viper.IsSet("calibre.cdb.add") {
		cmd.Opt(Flags(viper.GetStringSlice("calibre.cdb.add")...))
	}
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
