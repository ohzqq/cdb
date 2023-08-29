package cdb

import (
	"strconv"
	"strings"
)

// Book represents a book record.
type Book struct {
	Title        string  `db:"title" yaml:"title" json:"title"`
	Authors      string  `db:"authors" yaml:"authors,omitempty" json:"authors,omitempty"`
	Narrators    string  `db:"#narrators" yaml:"#narrators,omitempty" json:"#narrators,omitempty"`
	Series       string  `db:"series" yaml:"series,omitempty" json:"series,omitempty"`
	SeriesIndex  float64 `db:"series_index" yaml:"series_index,omitempty" json:"series_index,omitempty"`
	Tags         string  `db:"tags" yaml:"tags,omitempty" json:"tags,omitempty"`
	Pubdate      string  `db:"pubdate" yaml:"pubdate,omitempty" json:"pubdate,omitempty"`
	Timestamp    string  `db:"timestamp" yaml:"timestamp,omitempty" json:"timestamp,omitempty"`
	Duration     string  `db:"#duration" yaml:"#duration,omitempty" json:"#duration,omitempty"`
	Comments     string  `db:"comments" yaml:"comments,omitempty" json:"comments,omitempty"`
	Rating       string  `db:"rating" yaml:"rating,omitempty" json:"rating,omitempty"`
	Publisher    string  `db:"publisher" yaml:"publisher,omitempty" json:"publisher,omitempty"`
	Languages    string  `db:"languages" yaml:"languages,omitempty" json:"languages,omitempty"`
	Cover        string  `db:"cover" yaml:"cover,omitempty" json:"cover,omitempty"`
	Formats      string  `db:"formats" yaml:"formats,omitempty" json:"formats,omitempty"`
	Identifiers  string  `db:"identifiers" yaml:"identifiers,omitempty" json:"identifiers,omitempty"`
	LastModified string  `db:"last_modified" yaml:"last_modified,omitempty" json:"last_modified,omitempty"`
	ID           int     `db:"id" yaml:"id,omitempty" yaml:"id,omitempty"`
	AuthorSort   string  `db:"author_sort" yaml:"author_sort,omitempty" json:"author_sort,omitempty"`
	Sort         string  `db:"sort" yaml:"sort,omitempty" json:"sort,omitempty"`
	Path         string  `db:"path" yaml:"path,omitempty" json:"path,omitempty"`
	UUID         string  `db:"uuid,omitempty" yaml:"uuid,omitempty" json:"uuid,omitempty"`
}

// Map converts a book record to map[string]any.
func (b *Book) Map() map[string]any {
	book := make(map[string]any, 22)

	for l, v := range b.sharedMap() {
		book[l] = v
	}

	if v := b.Authors; v != "" {
		book[Authors] = splitNames(v)
	}
	if v := b.Narrators; v != "" {
		book[Narrators] = splitNames(v)
	}
	if v := b.Tags; v != "" {
		book[Tags] = splitCat(v)
	}
	if v := b.Languages; v != "" {
		book[Languages] = splitCat(v)
	}
	if v := b.Formats; v != "" {
		book[Formats] = splitCat(v)
	}
	if v := b.Identifiers; v != "" {
		book[Identifiers] = splitCat(v)
	}
	if v := b.Rating; v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			i = 0
		}
		book[Rating] = i
	}
	if v := b.ID; v != 0 {
		book[ID] = v
	}
	if v := b.SeriesIndex; v >= 0 {
		book[SeriesIndex] = v
	}
	return book
}

// StringMap converts a book record to map[string]string.
func (b *Book) StringMap() map[string]string {
	book := make(map[string]string, 22)

	for l, v := range b.sharedMap() {
		book[l] = v
	}

	if v := b.Authors; v != "" {
		book[Authors] = v
	}
	if v := b.Narrators; v != "" {
		book[Narrators] = v
	}
	if v := b.Tags; v != "" {
		book[Tags] = v
	}
	if v := b.Languages; v != "" {
		book[Languages] = v
	}
	if v := b.Formats; v != "" {
		book[Formats] = v
	}
	if v := b.Identifiers; v != "" {
		book[Identifiers] = v
	}
	if v := b.Rating; v != "" {
		book[Rating] = v
	}
	if v := b.ID; v != 0 {
		book[ID] = strconv.Itoa(v)
	}
	if v := b.SeriesIndex; v >= 0 {
		book[SeriesIndex] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return book
}

func splitNames(v string) []string {
	return strings.Split(v, " & ")
}

func splitCat(v string) []string {
	return strings.Split(v, ", ")
}

func (b *Book) sharedMap() map[string]string {
	book := make(map[string]string, 13)

	if v := b.Title; v != "" {
		book[Title] = v
	}
	if v := b.Series; v != "" {
		book[Series] = v
	}
	if v := b.Pubdate; v != "" {
		book[Pubdate] = v
	}
	if v := b.Duration; v != "" {
		book[Duration] = v
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
	if v := b.LastModified; v != "" {
		book[LastModified] = v
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
	return book
}
