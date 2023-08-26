package calibredb

import "github.com/spf13/viper"

// AddFormat adds a format to an existing book
func AddFormat(lib, id, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("add_format"), PositionalArgs(id, pos))
	if viper.IsSet("calibre.cdb.add_format") {
		cmd.Opt(Flags(viper.GetStringSlice("calibre.cdb.add_format")...))
	}
	return cmd
}

func DontReplace() Opt {
	return Flags("--dont-replace")
}
