package cdb

type Book struct {
	Title        string  `db:"title" yaml:"title"`
	Authors      string  `db:"authors" yaml:"authors,omitempty"`
	Narrators    string  `db:"#narrators" yaml:"#narrators,omitempty"`
	Series       string  `db:"series" yaml:"series,omitempty"`
	SeriesIndex  float32 `db:"series_index" yaml:"series_index,omitempty"`
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
