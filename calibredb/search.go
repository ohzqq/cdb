package calibredb

func Search(lib, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("search"), PositionalArgs(pos))
	return cmd
}

func Limit(l string) Opt {
	return Flags("--limit", l)
}
