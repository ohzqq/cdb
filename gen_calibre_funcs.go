// Code generated
// DO NOT EDIT.

package cdb

// RemoveFormat represents 'calibredb remove_format'.
type RemoveFormat struct {
	*Command
}

// SetMetadata represents 'calibredb set_metadata'.
type SetMetadata struct {
	*Command
}

// Field sets the --field [string] flag for 'calibredb set_metadata'.
func (c *SetMetadata) Field(v string) *SetMetadata {
	c.SetFlags("--field", v)
	return c
}

// ListFields sets the --list_fields flag for 'calibredb set_metadata'.
func (c *SetMetadata) ListFields() *SetMetadata {
	c.SetFlags("--list_fields")
	return c
}

// SavedSearchesList represents 'calibredb saved_searches list'.
type SavedSearchesList struct {
	*Command
}

// SavedSearchesAdd represents 'calibredb saved_searches add'.
type SavedSearchesAdd struct {
	*Command
}

// SavedSearchesRemove represents 'calibredb saved_searches remove'.
type SavedSearchesRemove struct {
	*Command
}

// Add represents 'calibredb add'.
type Add struct {
	*Command
}

// Add sets the --add [string] flag for 'calibredb add'.
func (c *Add) Add(v string) *Add {
	c.SetFlags("--add", v)
	return c
}

// Automerge sets the --automerge [string] flag for 'calibredb add'.
func (c *Add) Automerge(v string) *Add {
	c.SetFlags("--automerge", v)
	return c
}

// Duplicates sets the --duplicates flag for 'calibredb add'.
func (c *Add) Duplicates() *Add {
	c.SetFlags("--duplicates")
	return c
}

// Empty sets the --empty flag for 'calibredb add'.
func (c *Add) Empty() *Add {
	c.SetFlags("--empty")
	return c
}

// Ignore sets the --ignore [string] flag for 'calibredb add'.
func (c *Add) Ignore(v string) *Add {
	c.SetFlags("--ignore", v)
	return c
}

// OneBookPerDirectory sets the --one_book_per_directory flag for 'calibredb add'.
func (c *Add) OneBookPerDirectory() *Add {
	c.SetFlags("--one_book_per_directory")
	return c
}

// Recurse sets the --recurse flag for 'calibredb add'.
func (c *Add) Recurse() *Add {
	c.SetFlags("--recurse")
	return c
}

// EmbedMetadata represents 'calibredb embed_metadata'.
type EmbedMetadata struct {
	*Command
}

// OnlyFormats sets the --only_formats flag for 'calibredb embed_metadata'.
func (c *EmbedMetadata) OnlyFormats() *EmbedMetadata {
	c.SetFlags("--only_formats")
	return c
}

// Remove represents 'calibredb remove'.
type Remove struct {
	*Command
}

// Permanent sets the --permanent flag for 'calibredb remove'.
func (c *Remove) Permanent() *Remove {
	c.SetFlags("--permanent")
	return c
}

// AddFormat represents 'calibredb add_format'.
type AddFormat struct {
	*Command
}

// DontReplace sets the --dont_replace flag for 'calibredb add_format'.
func (c *AddFormat) DontReplace() *AddFormat {
	c.SetFlags("--dont_replace")
	return c
}

// Export represents 'calibredb export'.
type Export struct {
	*Command
}

// DontAsciiize sets the --dont_asciiize flag for 'calibredb export'.
func (c *Export) DontAsciiize() *Export {
	c.SetFlags("--dont_asciiize")
	return c
}

// DontSaveCover sets the --dont_save_cover flag for 'calibredb export'.
func (c *Export) DontSaveCover() *Export {
	c.SetFlags("--dont_save_cover")
	return c
}

// DontUpdateMetadata sets the --dont_update_metadata flag for 'calibredb export'.
func (c *Export) DontUpdateMetadata() *Export {
	c.SetFlags("--dont_update_metadata")
	return c
}

// DontWriteOpf sets the --dont_write_opf flag for 'calibredb export'.
func (c *Export) DontWriteOpf() *Export {
	c.SetFlags("--dont_write_opf")
	return c
}

// Formats sets the --formats [string] flag for 'calibredb export'.
func (c *Export) Formats(v string) *Export {
	c.SetFlags("--formats", v)
	return c
}

// ReplaceWhitespace sets the --replace_whitespace flag for 'calibredb export'.
func (c *Export) ReplaceWhitespace() *Export {
	c.SetFlags("--replace_whitespace")
	return c
}

// SingleDir sets the --single_dir flag for 'calibredb export'.
func (c *Export) SingleDir() *Export {
	c.SetFlags("--single_dir")
	return c
}

// Template sets the --template [string] flag for 'calibredb export'.
func (c *Export) Template(v string) *Export {
	c.SetFlags("--template", v)
	return c
}

// Timefmt sets the --timefmt [string] flag for 'calibredb export'.
func (c *Export) Timefmt(v string) *Export {
	c.SetFlags("--timefmt", v)
	return c
}

// ToDir sets the --to_dir [string] flag for 'calibredb export'.
func (c *Export) ToDir(v string) *Export {
	c.SetFlags("--to_dir", v)
	return c
}

// ToLowercase sets the --to_lowercase flag for 'calibredb export'.
func (c *Export) ToLowercase() *Export {
	c.SetFlags("--to_lowercase")
	return c
}

// ShowMetadata represents 'calibredb show_metadata'.
type ShowMetadata struct {
	*Command
}

// AsOpf sets the --as_opf flag for 'calibredb show_metadata'.
func (c *ShowMetadata) AsOpf() *ShowMetadata {
	c.SetFlags("--as_opf")
	return c
}	
	