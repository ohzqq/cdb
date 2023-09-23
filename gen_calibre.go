// Code generated
// DO NOT EDIT.

package cdb


// ListCommands lists the available calibredb commands.
func ListCommands() []string {
	return []string{
		"show_metadata",
		"saved_searches list",
		"saved_searches remove",
		"add",
		"add_format",
		"embed_metadata",
		"remove",
		"set_metadata",
		"export",
		"remove_format",
		"saved_searches add",
	}
}


// Export initializes the export command with the id ...string paramaters.
func (c *Command) Export(id ...string) *Export {
	c.SetArgs(id...)
	c.CdbCmd = "export"
	cmd := &Export{
		Command: c,
	}
	return cmd
}

// RemoveFormat initializes the remove_format command with the id string, fmt string paramaters.
func (c *Command) RemoveFormat(id string, fmt string) *RemoveFormat {
	c.SetArgs(id, fmt)
	c.CdbCmd = "remove_format"
	cmd := &RemoveFormat{
		Command: c,
	}
	return cmd
}

// SavedSearchesAdd initializes the saved_searches add command with the name string, expression string paramaters.
func (c *Command) SavedSearchesAdd(name string, expression string) *SavedSearchesAdd {
	c.SetArgs(name, expression)
	c.CdbCmd = "saved_searches add"
	cmd := &SavedSearchesAdd{
		Command: c,
	}
	return cmd
}

// SavedSearchesRemove initializes the saved_searches remove command with the name string paramaters.
func (c *Command) SavedSearchesRemove(name string) *SavedSearchesRemove {
	c.SetArgs(name)
	c.CdbCmd = "saved_searches remove"
	cmd := &SavedSearchesRemove{
		Command: c,
	}
	return cmd
}

// Add initializes the add command with the files ...string paramaters.
func (c *Command) Add(files ...string) *Add {
	c.SetArgs(files...)
	c.CdbCmd = "add"
	cmd := &Add{
		Command: c,
	}
	return cmd
}

// AddFormat initializes the add_format command with the id string, file string paramaters.
func (c *Command) AddFormat(id string, file string) *AddFormat {
	c.SetArgs(id, file)
	c.CdbCmd = "add_format"
	cmd := &AddFormat{
		Command: c,
	}
	return cmd
}

// EmbedMetadata initializes the embed_metadata command.
func (c *Command) EmbedMetadata() *EmbedMetadata {
	c.CdbCmd = "embed_metadata"
	cmd := &EmbedMetadata{
		Command: c,
	}
	return cmd
}

// Remove initializes the remove command.
func (c *Command) Remove() *Remove {
	c.CdbCmd = "remove"
	cmd := &Remove{
		Command: c,
	}
	return cmd
}

// SetMetadata initializes the set_metadata command with the id string paramaters.
func (c *Command) SetMetadata(id string) *SetMetadata {
	c.SetArgs(id)
	c.CdbCmd = "set_metadata"
	cmd := &SetMetadata{
		Command: c,
	}
	return cmd
}

// ShowMetadata initializes the show_metadata command with the id string paramaters.
func (c *Command) ShowMetadata(id string) *ShowMetadata {
	c.SetArgs(id)
	c.CdbCmd = "show_metadata"
	cmd := &ShowMetadata{
		Command: c,
	}
	return cmd
}

// SavedSearchesList initializes the saved_searches list command.
func (c *Command) SavedSearchesList() *SavedSearchesList {
	c.CdbCmd = "saved_searches list"
	cmd := &SavedSearchesList{
		Command: c,
	}
	return cmd
}
