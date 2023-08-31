package cdb

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ohzqq/libopds2-go/opds2"
)

// Book represents a book record.
type Book struct {
	Title     string `db:"title" yaml:"title" json:"title"`
	Authors   string `db:"authors" yaml:"authors,omitempty" json:"author,omitempty"`
	Narrators string `db:"#narrators" yaml:"#narrators,omitempty" json:"narrator,omitempty"`
	Series    string `db:"series" yaml:"series,omitempty" json:"belongs_to.series,omitempty"`
	//SeriesIndex  float64 `db:"series_index" yaml:"series_index,omitempty" json:"series_index,omitempty"`
	SeriesIndex float64 `db:"series_index" yaml:"series_index,omitempty" json:"-"`
	Tags        string  `db:"tags" yaml:"tags,omitempty" json:"subject,omitempty"`
	Pubdate     string  `db:"pubdate" yaml:"pubdate,omitempty" json:"published,omitempty"`
	Timestamp   string  `db:"timestamp" yaml:"timestamp,omitempty" json:"modified,omitempty"`
	Duration    string  `db:"#duration" yaml:"#duration,omitempty" json:"duration,omitempty"`
	Comments    string  `db:"comments" yaml:"comments,omitempty" json:"description,omitempty"`
	//Rating       string  `db:"rating" yaml:"rating,omitempty" json:"rating,omitempty"`
	Rating       string `db:"rating" yaml:"rating,omitempty" json:"rating,omitempty"`
	Publisher    string `db:"publisher" yaml:"publisher,omitempty" json:"publisher,omitempty"`
	Languages    string `db:"languages" yaml:"languages,omitempty" json:"language,omitempty"`
	Cover        string `db:"cover" yaml:"cover,omitempty" json:"cover,omitempty"`
	Formats      string `db:"formats" yaml:"formats,omitempty" json:"formats,omitempty"`
	Identifiers  string `db:"identifiers" yaml:"identifiers,omitempty" json:"identifiers,omitempty"`
	LastModified string `db:"last_modified" yaml:"last_modified,omitempty" json:"last_modified,omitempty"`
	ID           int    `db:"id" yaml:"id,omitempty" json:"id,omitempty"`
	AuthorSort   string `db:"author_sort" yaml:"author_sort,omitempty" json:"author_sort,omitempty"`
	Sort         string `db:"sort" yaml:"sort,omitempty" json:"sort,omitempty"`
	Path         string `db:"path" yaml:"path,omitempty" json:"path,omitempty"`
	UUID         string `db:"uuid,omitempty" yaml:"uuid,omitempty" json:"uuid,omitempty"`
}

// MarshalJSON satisfies the json.Marshaler interface, producing a somewhat more
// complicated data model.
//func (b *Book) MarshalJSON() ([]byte, error) {
//  return json.Marshal(b.Map())
//}

// CalibredbFlags is a convenience method for returning a slice of metadata
// fields to use with the 'set_metadata' command.
func (b *Book) CalibredbFlags() []string {
	var flags []string
	book := b.StringMap()
	for _, f := range AllModels().Editable() {
		if m, ok := book[f]; ok {
			flags = append(flags, m)
		}
	}
	return flags
}

func (b *Book) toOPDS(lib string) opds2.Publication {
	meta := make(map[string]any)
	var img []any
	var links []any

	id := strconv.Itoa(b.ID)
	meta["identifier"] = filepath.Join("/", lib, "books", id)
	meta["source"] = lib

	if v := b.Authors; v != "" {
		meta["author"] = splitNames(v)
	}
	if v := b.Narrators; v != "" {
		meta["narrator"] = splitNames(v)
	}
	if v := b.Tags; v != "" {
		meta["subject"] = splitCat(v)
	}
	if v := b.Languages; v != "" {
		meta["language"] = splitCat(v)
	}
	if v := b.Formats; v != "" {
		for _, f := range splitCat(v) {
			l := make(map[string]any)
			l["href"] = filepath.Join("/", lib, f)
			l["rel"] = []string{
				"http://opds-spec.org/acquisition",
				filepath.Join("/", lib, b.Path, filepath.Base(f)),
			}
			links = append(links, l)
		}
	}
	if v := b.Series; v != "" {
		meta["belongs_to"] = map[string]any{
			"series": []any{
				map[string]any{
					"name":     b.Series,
					"position": b.SeriesIndex,
				},
			},
		}
	}
	if v := b.Cover; v != "" {
		l := make(map[string]any)
		l["href"] = filepath.Join("/", lib, v)
		l["rel"] = []string{
			"cover",
			filepath.Join("/", lib, b.Path, filepath.Base(v)),
		}
		img = append(img, l)
	}
	for k, v := range b.Map() {
		switch k {
		case Pubdate:
			meta["published"] = v
		case Timestamp:
			meta["modified"] = v
		case Publisher, Duration, Title:
			meta[k] = v
		}
	}
	pub := map[string]any{
		"metadata": meta,
		"images":   img,
		"links":    links,
	}
	return opds2.ParsePublication(pub)
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
	if v := b.Duration; v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			i = 0
		}
		book[Duration] = i
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
	if v := b.Duration; v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			i = 0
		}
		book[Duration] = time.Unix(int64(i), 0).Format(time.TimeOnly)
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
	if v := b.Timestamp; v != "" {
		book[Timestamp] = v
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
