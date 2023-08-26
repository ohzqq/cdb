package calibredb

import "github.com/spf13/viper"

// Export saves a book to disk
func Export(lib, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("export"), PositionalArgs(pos))
	if viper.IsSet("calibre.cdb.export") {
		cmd.Opt(Flags(viper.GetStringSlice("calibre.cdb.export")...))
	}
	return cmd
}

func All() Opt {
	return Flags("--all")
}

func DontAsciize() Opt {
	return Flags("--dont-asciize")
}

func DontSaveCover() Opt {
	return Flags("--dont-save-cover")
}

func DontUpdateMetadata() Opt {
	return Flags("--dont-update-metadata")
}

func DontWriteOpf() Opt {
	return Flags("--dont-write-opf")
}

func Progress() Opt {
	return Flags("--progress")
}

func ReplaceWhitespace() Opt {
	return Flags("--replace-whitespace")
}

func SingleDir() Opt {
	return Flags("--single-dir")
}

func ToLowercase() Opt {
	return Flags("--to-lowercase")
}

func Formats(formats string) Opt {
	return Flags("--formats", formats)
}

func Template(template string) Opt {
	return Flags("--template", template)
}

func TimeFmt(format string) Opt {
	return Flags("--timefmt", format)
}

func ToDir(dir string) Opt {
	return Flags("--to-dir", dir)
}
