package calibredb

import (
	"github.com/ohzqq/ur"
	"github.com/ohzqq/ur/util"
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
	if ok := util.FileExist(file); !ok {
		err := util.ErrFileNotExist(file)
		ur.HandleError(file, err)
	}
	return Flags("--cover", file)
}

func Duplicates() Opt {
	return Flags("--duplicates")
}

func Empty() Opt {
	return Flags("--empty")
}
