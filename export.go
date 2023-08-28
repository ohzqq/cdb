package cdb

import "log"

// Export saves a book to disk
func Export(lib, pos string, args ...Opt) *Command {
	cmd, err := NewCommand(lib, args...)
	if err != nil {
		log.Fatal(err)
	}

	cmd.Opt(Cmd("export"), PositionalArgs(pos))

	return cmd
}

func All() Opt {
	return SetFlags("--all")
}

func DontAsciize() Opt {
	return SetFlags("--dont-asciize")
}

func DontSaveCover() Opt {
	return SetFlags("--dont-save-cover")
}

func DontUpdateMetadata() Opt {
	return SetFlags("--dont-update-metadata")
}

func DontWriteOpf() Opt {
	return SetFlags("--dont-write-opf")
}

func Progress() Opt {
	return SetFlags("--progress")
}

func ReplaceWhitespace() Opt {
	return SetFlags("--replace-whitespace")
}

func SingleDir() Opt {
	return SetFlags("--single-dir")
}

func ToLowercase() Opt {
	return SetFlags("--to-lowercase")
}

func ExportFormats(formats string) Opt {
	return SetFlags("--formats", formats)
}

func Template(template string) Opt {
	return SetFlags("--template", template)
}

func TimeFmt(format string) Opt {
	return SetFlags("--timefmt", format)
}

func ToDir(dir string) Opt {
	return SetFlags("--to-dir", dir)
}
