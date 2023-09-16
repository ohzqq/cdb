package cdb

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Book represents a book record.
type Book struct {
	Title        string     `db:"title" yaml:"title" toml:"title" json:"title"`
	Authors      []string   `db:"authors" yaml:"authors,omitempty" toml:"authors,omitempty" json:"authors,omitempty"`
	Narrators    []string   `db:"narrators" yaml:"narrators,omitempty" toml:"narrators,omitempty" json:"narrators,omitempty"`
	Series       string     `db:"series" yaml:"series,omitempty" toml:"series,omitempty" json:"series,omitempty"`
	SeriesIndex  float64    `db:"series_index" yaml:"series_index,omitempty" toml:"series_index,omitempty" json:"series_index,omitempty"`
	Tags         []string   `db:"tags" yaml:"tags,omitempty" toml:"tags,omitempty" json:"tags,omitempty"`
	Pubdate      *time.Time `db:"pubdate" yaml:"pubdate,omitempty" toml:"pubdate,omitempty" json:"pubdate,omitempty"`
	Timestamp    *time.Time `db:"timestamp" yaml:"timestamp,omitempty" toml:"timestamp,omitempty" json:"timestamp,omitempty"`
	Duration     string     `db:"duration" yaml:"duration,omitempty" toml:"duration,omitempty" json:"duration,omitempty"`
	Comments     string     `db:"comments" yaml:"comments,omitempty" toml:"comments,omitempty" json:"comments,omitempty"`
	Rating       int        `db:"rating" yaml:"rating,omitempty" toml:"rating,omitempty" json:"rating,omitempty"`
	Publisher    string     `db:"publisher" yaml:"publisher,omitempty" toml:"publisher,omitempty" json:"publisher,omitempty"`
	Languages    []string   `db:"languages" yaml:"languages,omitempty" toml:"languages,omitempty" json:"languages,omitempty"`
	Cover        string     `db:"cover" yaml:"cover,omitempty" toml:"cover,omitempty" json:"cover,omitempty"`
	Formats      []string   `db:"formats" yaml:"formats,omitempty" toml:"formats,omitempty" json:"formats,omitempty"`
	Identifiers  []string   `db:"identifiers" yaml:"identifiers,omitempty" toml:"identifiers,omitempty" json:"identifiers,omitempty"`
	LastModified *time.Time `db:"last_modified" yaml:"last_modified,omitempty" toml:"last_modified,omitempty" json:"last_modified,omitempty"`
	ID           int        `db:"id" yaml:"id,omitempty" toml:"id,omitempty" json:"id,omitempty"`
	AuthorSort   string     `db:"author_sort" yaml:"author_sort,omitempty" toml:"author_sort,omitempty" json:"author_sort,omitempty"`
	Sort         string     `db:"sort" yaml:"sort,omitempty" toml:"sort,omitempty" json:"sort,omitempty"`
	Path         string     `db:"path" yaml:"path,omitempty" toml:"path,omitempty" json:"path,omitempty"`
	UUID         string     `db:"uuid,omitempty" yaml:"uuid,omitempty" toml:"uuid,omitempty" json:"uuid,omitempty"`
	Source       string     `json:"source,omitempty" yaml:"-" toml:"-"`
}

type BookEncoder interface {
	Encode(v any) error
}

type BookDecoder interface {
	Decode(v any) error
}

// URL sets the path for a *url.URL and returns a string, by default returns a
// path.
func (b *Book) URL(urlopt ...*url.URL) string {
	u := &url.URL{}
	if len(urlopt) > 0 {
		u = urlopt[0]
	}
	id := strconv.Itoa(b.ID)
	u.Path = filepath.Join("/", b.Source, "books", id)
	return u.String()
}

func (b Book) Save(f string) error {
	if f != ".yaml" || f != ".yml" || f != ".toml" || f != ".json" {
		return fmt.Errorf("only yaml, toml, and json can be written\n")
	}

	file, err := os.Create(b.Title + f)
	if err != nil {
		return err
	}

	return nil
}

func (b Book) Encode(opts ...EncoderOption) {
}

// StringMapString converts a book record to map[string]string.
func (b Book) StringMapString() map[string]string {
	m := make(map[string]string)
	for k, v := range b.sharedMap() {
		m[k] = v
	}
	if v := b.Pubdate; v != nil {
		m[Pubdate] = v.Format(time.DateOnly)
	}
	if v := b.LastModified; v != nil {
		m[LastModified] = v.Format(time.DateOnly)
	}
	if v := b.Timestamp; v != nil {
		m[Timestamp] = v.Format(time.DateOnly)
	}
	if v := b.Authors; len(v) != 0 {
		m[Authors] = GetModel(Authors).Join(v)
	}
	if v := b.Narrators; len(v) != 0 {
		m[Narrators] = GetModel(Narrators).Join(v)
	}
	if v := b.Tags; len(v) != 0 {
		m[Tags] = GetModel(Tags).Join(v)
	}
	if v := b.Languages; len(v) != 0 {
		m[Languages] = GetModel(Languages).Join(v)
	}
	if v := b.Formats; len(v) != 0 {
		m[Formats] = GetModel(Formats).Join(v)
	}
	if v := b.Identifiers; len(v) != 0 {
		m[Identifiers] = GetModel(Identifiers).Join(v)
	}
	if v := b.Duration; v != "" {
		m[Duration] = v
	}
	if v := b.Rating; v != 0 {
		m[Rating] = strconv.Itoa(v)
	}
	if v := b.ID; v != 0 {
		m[ID] = strconv.Itoa(v)
	}
	if v := b.SeriesIndex; v >= 0 {
		m[SeriesIndex] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return m
}

// StringMap converts a book record to map[string]any.
func (b Book) StringMap() map[string]any {
	m := make(map[string]any)
	for k, v := range b.sharedMap() {
		m[k] = v
	}
	if v := b.Pubdate; v != nil {
		m[Pubdate] = v
	}
	if v := b.LastModified; v != nil {
		m[LastModified] = v
	}
	if v := b.Timestamp; v != nil {
		m[Timestamp] = v
	}
	if v := b.Authors; len(v) != 0 {
		m[Authors] = v
	}
	if v := b.Narrators; len(v) != 0 {
		m[Narrators] = v
	}
	if v := b.Tags; len(v) != 0 {
		m[Tags] = v
	}
	if v := b.Languages; len(v) != 0 {
		m[Languages] = v
	}
	if v := b.Formats; len(v) != 0 {
		m[Formats] = v
	}
	if v := b.Identifiers; len(v) != 0 {
		m[Identifiers] = v
	}
	if v := b.Duration; v != "" {
		m[Duration] = v
	}
	if v := b.Rating; v != 0 {
		m[Rating] = v
	}
	if v := b.ID; v != 0 {
		m[ID] = v
	}
	if v := b.SeriesIndex; v >= 0 {
		m[SeriesIndex] = v
	}
	return m
}

// CalibredbFlags is a convenience method for returning a slice of metadata
// fields to use with the 'set_metadata' command.
func (b *Book) CalibredbFlags() []string {
	var flags []string
	book := b.StringMapString()
	for l, model := range AllModels().Editable() {
		if m, ok := book[l]; ok {
			if model.IsCustom {
				l = "#" + l
			}
			flags = append(flags, l+":"+m)
		}
	}
	return flags
}

func (b *Book) sharedMap() map[string]string {
	book := make(map[string]string, 13)

	if v := b.Title; v != "" {
		book[Title] = v
	}
	if v := b.Series; v != "" {
		book[Series] = v
	}
	if v := b.Comments; v != "" {
		book[Comments] = v
	}
	if v := b.Publisher; v != "" {
		book[Publisher] = v
	}
	if v := b.Cover; v != "" {
		book[Cover] = v
	}
	if v := b.AuthorSort; v != "" {
		book[AuthorSort] = v
	}
	if v := b.Sort; v != "" {
		book[Sort] = v
	}
	if v := b.Path; v != "" {
		book[Path] = v
	}
	if v := b.UUID; v != "" {
		book[UUID] = v
	}
	if v := b.Source; v != "" {
		book["source"] = v
	}
	return book
}
