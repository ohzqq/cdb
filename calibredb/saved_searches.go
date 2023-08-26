package calibredb

func SavedSearchesList(lib string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("saved_searches"), PositionalArgs("list"))
	return cmd
}

func SavedSearchesAdd(lib, name, exp string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("saved_searches"), PositionalArgs("add", name, exp))
	return cmd
}

func SavedSearchesRemove(lib, pos string, args ...Opt) *Command {
	cmd := New(lib, args...)
	cmd.Opt(Cmd("saved_searches"), PositionalArgs("remove", pos))
	return cmd
}
