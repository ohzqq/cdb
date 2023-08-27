package cdb

import (
	"log"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Book struct {
	Title        string  `db:"title" yaml:"title"`
	Authors      string  `db:"authors" yaml:"authors,omitempty"`
	Narrators    string  `db:"#narrators" yaml:"#narrators,omitempty"`
	Series       string  `db:"series" yaml:"series,omitempty"`
	SeriesIndex  float64 `db:"series_index" yaml:"series_index,omitempty"`
	Tags         string  `db:"tags" yaml:"tags,omitempty"`
	Pubdate      string  `db:"pubdate" yaml:"pubdate,omitempty"`
	Timestamp    string  `db:"timestamp" yaml:"timestamp,omitempty"`
	Duration     string  `db:"#duration" yaml:"#duration,omitempty"`
	Comments     string  `db:"comments" yaml:"comments,omitempty"`
	Rating       string  `db:"rating" yaml:"rating,omitempty"`
	Publisher    string  `db:"publisher" yaml:"publisher,omitempty"`
	Languages    string  `db:"languages" yaml:"languages,omitempty"`
	Cover        string  `db:"cover" yaml:"-"`
	Formats      string  `db:"formats" yaml:"-"`
	Identifiers  string  `db:"identifiers" yaml:"identifiers,omitempty"`
	LastModified string  `db:"last_modified" yaml:"last_modified,omitempty"`
	ID           int     `db:"id" yaml:"-"`
	AuthorSort   string  `db:"author_sort" yaml:"author_sort,omitempty"`
	Sort         string  `db:"sort" yaml:"sort,omitempty"`
	Path         string  `db:"path" yaml:"-"`
	UUID         string  `db:"uuid,omitempty" yaml:"-"`
}

func (b *Book) ToYAML() []byte {
	d, err := yaml.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func (b *Book) Map() map[string]any {
	book := make(map[string]any, 22)

	if v := b.Title; v != "" {
		book[Title] = v
	}
	if v := b.Authors; v != "" {
		book[Authors] = v
	}
	if v := b.Narrators; v != "" {
		book[Narrators] = v
	}
	if v := b.Series; v != "" {
		book[Series] = v
	}
	if v := b.Tags; v != "" {
		book[Tags] = v
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
	if v := b.Rating; v != "" {
		book[Rating] = v
	}
	if v := b.Publisher; v != "" {
		book[Publisher] = v
	}
	if v := b.Languages; v != "" {
		book[Languages] = v
	}
	if v := b.Cover; v != "" {
		book[Cover] = v
	}
	if v := b.Formats; v != "" {
		book[Formats] = v
	}
	if v := b.Identifiers; v != "" {
		book[Identifiers] = v
	}
	if v := b.LastModified; v != "" {
		book[LastModified] = v
	}
	if v := b.ID; v != 0 {
		book[ID] = v
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
	if v := b.SeriesIndex; v >= 0 {
		book[SeriesIndex] = v
	}

	return book
}

func (b *Book) StringMap() map[string]string {
	book := make(map[string]string, 22)
	for k, v := range b.Map() {
		switch k {
		case ID:
			book[k] = strconv.Itoa(v.(int))
		case SeriesIndex:
			book[k] = strconv.FormatFloat(v.(float64), 'f', -1, 32)
		default:
			book[k] = v.(string)
		}
	}
	return book
}

func FromYAML(d []byte) *Book {
	b := &Book{}
	err := yaml.Unmarshal(d, b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
