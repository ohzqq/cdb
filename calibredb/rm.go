package calibredb

import (
	"github.com/spf13/viper"
)

func Remove(lib string, pos ...string) *Command {
	cmd := New(lib)
	cmd.Opt(Cmd("remove"), PositionalArgs(pos...))
	if viper.IsSet("calibre.cdb.remove") {
		cmd.Opt(Flags(viper.GetStringSlice("calibre.cdb.remove")...))
	}
	return cmd
}

func Permanent() Opt {
	return Flags("--permanent")
}
