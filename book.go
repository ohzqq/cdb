package cdb

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Book represents a book record.
type Book struct {
	Title        string    `db:"title" yaml:"title" json:"title"`
	Authors      []string  `db:"authors" yaml:"authors,omitempty" json:"authors,omitempty"`
	Narrators    []string  `db:"#narrators" yaml:"#narrators,omitempty" json:"#narrators,omitempty"`
	Series       string    `db:"series" yaml:"series,omitempty" json:"series,omitempty"`
	SeriesIndex  float64   `db:"series_index" yaml:"series_index,omitempty" json:"series_index,omitempty"`
	Tags         []string  `db:"tags" yaml:"tags,omitempty" json:"tags,omitempty"`
	Pubdate      time.Time `db:"pubdate" yaml:"pubdate,omitempty" json:"pubdate,omitempty"`
	Timestamp    time.Time `db:"timestamp" yaml:"timestamp,omitempty" json:"timestamp,omitempty"`
	Duration     int       `db:"#duration" yaml:"#duration,omitempty" json:"#duration,omitempty"`
	Comments     string    `db:"comments" yaml:"comments,omitempty" json:"description,omitempty"`
	Rating       int       `db:"rating" yaml:"rating,omitempty" json:"rating,omitempty"`
	Publisher    string    `db:"publisher" yaml:"publisher,omitempty" json:"publisher,omitempty"`
	Languages    []string  `db:"languages" yaml:"languages,omitempty" json:"languages,omitempty"`
	Cover        string    `db:"cover" yaml:"cover,omitempty" json:"cover,omitempty"`
	Formats      []string  `db:"formats" yaml:"formats,omitempty" json:"formats,omitempty"`
	Identifiers  []string  `db:"identifiers" yaml:"identifiers,omitempty" json:"identifiers,omitempty"`
	LastModified time.Time `db:"last_modified" yaml:"last_modified,omitempty" json:"last_modified,omitempty"`
	ID           int       `db:"id" yaml:"id,omitempty" json:"id,omitempty"`
	AuthorSort   string    `db:"author_sort" yaml:"author_sort,omitempty" json:"author_sort,omitempty"`
	Sort         string    `db:"sort" yaml:"sort,omitempty" json:"sort,omitempty"`
	Path         string    `db:"path" yaml:"path,omitempty" json:"path,omitempty"`
	UUID         string    `db:"uuid,omitempty" yaml:"uuid,omitempty" json:"uuid,omitempty"`
	Source       string    `json:"source,omitempty"`
}

// MarshalJSON satisfies the json.Marshaler interface, producing a somewhat more
// complicated data model.
//func (b *Book) MarshalJSON() ([]byte, error) {
//  return json.Marshal(b.Map())
//}

// CalibredbFlags is a convenience method for returning a slice of metadata
// fields to use with the 'set_metadata' command.
//func (b *Book) CalibredbFlags() []string {
//  var flags []string
//  book := b.StringMap()
//  for _, f := range AllModels().Editable() {
//    if m, ok := book[f]; ok {
//      flags = append(flags, m)
//    }
//  }
//  return flags
//}

//func (b *Book) ToOPDS2Publication() opds2.Publication {
//  meta := b.Map()
//  for k, v := range meta {
//    switch k {
//    case Pubdate:
//      meta["published"] = v
//    case Timestamp:
//      meta["modified"] = v
//    case Languages:
//      meta["language"] = v
//    case Tags:
//      meta["subject"] = v
//    case Authors:
//      meta["author"] = v
//    case Narrators:
//      meta["narrator"] = v
//    case Publisher, Duration, Title, "source":
//      meta[k] = v
//    }
//  }

//  id := strconv.Itoa(b.ID)
//  meta["identifier"] = filepath.Join("/", b.Source, "books", id)

//  pub := opds2.NewPublication(map[string]any{"metadata": meta})

//  if v := b.Formats; v != "" {
//    for _, f := range strings.Split(v, ", ") {
//      l := map[string]any{
//        "href": filepath.Join("/", b.Source, f),
//        "rel": []string{
//          "http://opds-spec.org/acquisition",
//          filepath.Join("/", b.Source, b.Path, filepath.Base(f)),
//        },
//      }
//      pub.Links = append(pub.Links, opds2.NewLink(l))
//    }
//  }

//  if v := b.Cover; v != "" {
//    l := map[string]any{
//      "href": filepath.Join("/", b.Source, v),
//      "rel": []string{
//        "cover",
//        filepath.Join("/", b.Source, b.Path, filepath.Base(v)),
//      },
//    }
//    pub.Images = append(pub.Images, opds2.NewLink(l))
//  }

//  if v := b.Series; v != "" {
//    series := map[string]any{
//      "name":     b.Series,
//      "position": b.SeriesIndex,
//    }
//    pub.BelongsToSeries(series)
//  }

//  return pub
//}

const (
	namesSep = " & "
	catSep   = ", "
)

func (b *Book) URL() string {
	id := strconv.Itoa(b.ID)
	return filepath.Join("/", b.Source, "books", id)
}

func joinCat(v []string) string {
	return strings.Join(v, catSep)
}

func joinNames(v []string) string {
	return strings.Join(v, namesSep)
}

//func (b *Book) ToOPDS1Publication() opds1.Entry {
//  var book opds1.Entry

//  book.Identifier = b.URL()
//  book.ID = b.URL()

//  if v := b.Title; v != "" {
//    book.Title = v
//  }

//  if v := b.Series; v != "" {
//    book.Series = []opds1.Serie{
//      opds1.Serie{
//        Name:     b.Series,
//        Position: float32(b.SeriesIndex),
//      },
//    }
//  }

//  if v := b.Pubdate; v != "" {
//    t, err := time.Parse(time.RFC3339, v)
//    if err != nil {
//      t = time.Now()
//    }
//    book.Published = &t
//  }

//  if v := b.Comments; v != "" {
//    book.Summary = opds1.Content{
//      Content:     b.Comments,
//      ContentType: "html",
//    }
//  }

//  if v := b.Publisher; v != "" {
//    book.Publisher = v
//  }

//  if v := b.Cover; v != "" {
//    cover := opds1.Link{
//      Rel:      "http://opds-spec.org/image",
//      Href:     filepath.Join("/", b.Source, v),
//      TypeLink: "image/jpeg",
//    }
//    book.Links = append(book.Links, cover)
//  }

//  if v := b.Formats; v != "" {
//    for _, f := range strings.Split(v, ", ") {
//      l := opds1.Link{
//        Href:     filepath.Join("/", b.Source, f),
//        Rel:      "http://opds-spec.org/acquisition",
//        TypeLink: mime.TypeByExtension(filepath.Ext(f)),
//      }
//      book.Links = append(book.Links, l)
//    }
//  }

//  if v := b.Timestamp; v != "" {
//    t, err := time.Parse(time.RFC3339, v)
//    if err != nil {
//      t = time.Now()
//    }
//    book.Updated = &t
//  }

//  if v := b.Languages; v != "" {
//    book.Language = v
//  }

//  if v := b.Authors; v != "" {
//    for _, a := range strings.Split(v, " & ") {
//      auth := opds1.Author{
//        Name: a,
//      }
//      book.Author = append(book.Author, auth)
//    }
//  }

//  if v := b.Tags; v != "" {
//    for _, t := range strings.Split(v, ", ") {
//      tag := opds1.Category{
//        Term:  t,
//        Label: t,
//      }
//      book.Category = append(book.Category, tag)
//    }
//  }

//  return book
//}

// Map converts a book record to map[string]any.
//func (b *Book) Map() map[string]any {
//  book := make(map[string]any, 22)

//  for l, v := range b.sharedMap() {
//    book[l] = v
//  }

//  if v := b.Authors; v != "" {
//    book[Authors] = splitNames(v)
//  }
//  if v := b.Narrators; v != "" {
//    book[Narrators] = splitNames(v)
//  }
//  if v := b.Tags; v != "" {
//    book[Tags] = splitCat(v)
//  }
//  if v := b.Languages; v != "" {
//    book[Languages] = splitCat(v)
//  }
//  if v := b.Formats; v != "" {
//    book[Formats] = splitCat(v)
//  }
//  if v := b.Identifiers; v != "" {
//    book[Identifiers] = splitCat(v)
//  }
//  if v := b.Duration; v != "" {
//    i, err := strconv.Atoi(v)
//    if err != nil {
//      i = 0
//    }
//    book[Duration] = i
//  }
//  if v := b.Rating; v != "" {
//    i, err := strconv.Atoi(v)
//    if err != nil {
//      i = 0
//    }
//    book[Rating] = i
//  }
//  if v := b.ID; v != 0 {
//    book[ID] = v
//  }
//  if v := b.SeriesIndex; v >= 0 {
//    book[SeriesIndex] = v
//  }
//  return book
//}

// StringMap converts a book record to map[string]string.
//func (b *Book) StringMap() map[string]string {
//  book := make(map[string]string, 22)

//  for l, v := range b.sharedMap() {
//    book[l] = v
//  }

//  if v := b.Authors; v != "" {
//    book[Authors] = v
//  }
//  if v := b.Narrators; v != "" {
//    book[Narrators] = v
//  }
//  if v := b.Tags; v != "" {
//    book[Tags] = v
//  }
//  if v := b.Languages; v != "" {
//    book[Languages] = v
//  }
//  if v := b.Formats; v != "" {
//    book[Formats] = v
//  }
//  if v := b.Identifiers; v != "" {
//    book[Identifiers] = v
//  }
//  if v := b.Rating; v != "" {
//    book[Rating] = v
//  }
//  if v := b.Duration; v != "" {
//    i, err := strconv.Atoi(v)
//    if err != nil {
//      i = 0
//    }
//    book[Duration] = time.Unix(int64(i), 0).Format(time.TimeOnly)
//  }
//  if v := b.ID; v != 0 {
//    book[ID] = strconv.Itoa(v)
//  }
//  if v := b.SeriesIndex; v >= 0 {
//    book[SeriesIndex] = strconv.FormatFloat(v, 'f', -1, 64)
//  }

//  return book
//}

//func splitNames(v string) []any {
//  var names []any
//  for _, n := range strings.Split(v, " & ") {
//    names = append(names, n)
//  }
//  return names
//}

//func splitCat(v string) []any {
//  var names []any
//  for _, n := range strings.Split(v, ", ") {
//    names = append(names, n)
//  }
//  return names
//}

//func (b *Book) sharedMap() map[string]string {
//  book := make(map[string]string, 13)

//  if v := b.Title; v != "" {
//    book[Title] = v
//  }
//  if v := b.Series; v != "" {
//    book[Series] = v
//  }
//  if v := b.Pubdate; v != "" {
//    book[Pubdate] = v
//  }
//  if v := b.Comments; v != "" {
//    book[Comments] = v
//  }
//  if v := b.Publisher; v != "" {
//    book[Publisher] = v
//  }
//  if v := b.Cover; v != "" {
//    book[Cover] = v
//  }
//  if v := b.LastModified; v != "" {
//    book[LastModified] = v
//  }
//  if v := b.Timestamp; v != "" {
//    book[Timestamp] = v
//  }
//  if v := b.AuthorSort; v != "" {
//    book[AuthorSort] = v
//  }
//  if v := b.Sort; v != "" {
//    book[Sort] = v
//  }
//  if v := b.Path; v != "" {
//    book[Path] = v
//  }
//  if v := b.UUID; v != "" {
//    book[UUID] = v
//  }
//  if v := b.Source; v != "" {
//    book["source"] = v
//  }
//  return book
//}
