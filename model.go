package cdb

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/exp/slices"
)

// Models is a map of Model.
type Models map[string]Model

const (
	Authors       = "authors"
	AuthorSort    = "author_sort"
	Comments      = "comments"
	Cover         = "cover"
	CustomColumns = "custom_columns"
	Formats       = "formats"
	ID            = "id"
	Identifiers   = "identifiers"
	Languages     = "languages"
	LastModified  = "last_modified"
	Path          = "path"
	Pubdate       = "pubdate"
	Publisher     = "publisher"
	Rating        = "rating"
	Series        = "series"
	SeriesIndex   = "series_index"
	Sort          = "sort"
	SortAs        = "sort_as"
	Tags          = "tags"
	Timestamp     = "timestamp"
	Title         = "title"
	UUID          = "uuid"
	Duration      = "#duration"
	Narrators     = "#narrators"
)

// Model represents a book field.
type Model struct {
	CategorySort string `json:"category_sort" db:"category_sort"`
	Column       string `json:"column" db:"column"`
	Colnum       int    `json:"colnum" db:"id"`
	IsCategory   bool   `json:"is_category" db:"is_category"`
	IsCustom     bool   `json:"is_custom" db:"is_custom"`
	IsEditable   bool   `json:"is_editable" db:"is_editable"`
	IsNames      bool   `json:"is_names" db:"is_names"`
	JoinTable    string `json:"join_table" db:"join_table"`
	Label        string `json:"label" db:"label"`
	LinkColumn   string `json:"link_column" db:"link_column"`
	Name         string `json:"name" db:"name"`
	Table        string `json:"table" db:"table"`
}

// DefaultModels returns the default book fields.
func DefaultModels() Models {
	return modelMeta
}

// AudiobookModels returns the default book fields with Duration and Narrators.
func AudiobookModels() Models {
	models := Models{
		Duration:  durationModel,
		Narrators: narratorsModel,
	}
	return models
}

// AllModels returns all possible book fields.
func AllModels() Models {
	models := DefaultModels()
	for l, m := range AudiobookModels() {
		models[l] = m
	}
	return models
}

// ToSql returns a model's sql expression as a string.
func (m Model) ToSql() string {
	switch {
	case slices.Contains(manyToOneModels, m.Label), m.Label == Duration:
		stmt, _ := manyToOne(m)
		return stmt
	case slices.Contains(manyToManyModels, m.Label), m.Label == Narrators:
		stmt, _ := manyToMany(m)
		return stmt
	case m.Label == Cover:
		stmt, _ := coverStmt("")
		return stmt
	default:
		stmt, _ := booksColumn(m)
		return stmt
	}
}

// Editable returns the list of editable book fields.
func (m Models) Editable() []string {
	var edit []string
	for l, mod := range m {
		if mod.IsEditable {
			edit = append(edit, l)
		}
	}
	return edit
}

func booksColumn(m Model) (string, []any) {
	return fmt.Sprintf(`IFNULL(%s, '') %s`, m.Column, m.Label), []any{}
}

func groupConcat(m Model) string {
	sep := defaultSep
	if m.Label == Authors || m.Label == Narrators {
		sep = namesSep
	}
	return fmt.Sprintf("IFNULL(GROUP_CONCAT(%s, '%s'), '')", m.Column, sep)
}

func manyToOne(m Model) (string, []any) {
	sel := sq.Select(groupConcat(m)).
		Prefix("(").
		From(m.Table).
		Where("book=books.id").
		Suffix(fmt.Sprint(") '", m.Label, "'"))

	return toSql(sel)
}

func manyToMany(m Model) (string, []any) {
	col := sq.Select(groupConcat(m)).
		From(m.Table).
		Where(fmt.Sprint(m.Table, ".id"))

	join := sq.Select(m.LinkColumn).
		From(m.JoinTable).
		Where("book IN (books.id)")

	cs, _ := toSql(col)
	js, _ := toSql(join)

	stmt := fmt.Sprintf("(%s IN (%s)) '%s'", cs, js, m.Label)

	return stmt, []any{}
}

func toSql(sel sq.SelectBuilder) (string, []any) {
	stmt, args, err := sel.ToSql()
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	return stmt, args
}

func coverStmt(lib string) (string, []any) {
	sel := sq.Case("has_cover").
		When("1", `IFNULL('books/' || books.id || '/cover.jpg', '')`).
		Else("''")

	stmt, args, _ := sel.ToSql()
	return stmt + " cover", args
}

const (
	namesSep   = " & "
	defaultSep = ", "
)

func calculateOffset(page, limit int) uint64 {
	cur := page - 1
	if cur < 0 {
		cur = 0
	}
	return uint64(cur * limit)
}

var oneToOneModels = []string{
	AuthorSort,
	ID,
	LastModified,
	Path,
	Pubdate,
	SeriesIndex,
	Sort,
	Timestamp,
	Title,
	UUID,
}

var manyToOneModels = []string{
	Comments,
	Identifiers,
	Formats,
}

var manyToManyModels = []string{
	Authors,
	Languages,
	Publisher,
	Rating,
	Series,
	Tags,
}

var durationModel = Model{
	Column:       "value",
	CategorySort: "value",
	IsCustom:     true,
	IsEditable:   true,
	Label:        "#duration",
	Name:         "Duration",
}

var narratorsModel = Model{
	CategorySort: "value",
	Column:       "value",
	IsCategory:   true,
	IsCustom:     true,
	IsEditable:   true,
	IsNames:      true,
	Label:        "#narrators",
	LinkColumn:   "value",
	Name:         "Narrators",
}

var modelMeta = Models{
	Authors: Model{
		CategorySort: "sort",
		Column:       "name",
		IsCategory:   true,
		IsEditable:   true,
		IsNames:      true,
		JoinTable:    "books_authors_link",
		Label:        "authors",
		LinkColumn:   "author",
		Name:         "Authors",
		Table:        "authors",
	},

	AuthorSort: Model{
		Column:     "author_sort",
		IsEditable: true,
		Label:      "author_sort",
		Name:       "AuthorSort",
	},

	Comments: Model{
		Column:     "text",
		IsEditable: true,
		Label:      "comments",
		Name:       "Description",
		Table:      "comments",
	},

	Cover: Model{
		Column:     "cover",
		IsEditable: false,
		IsNames:    false,
		Label:      "cover",
		Name:       "Cover",
	},

	Formats: Model{
		Column:       `'books/' || books.id || '/' || data.name || '.' || lower(format)`,
		CategorySort: "format",
		IsCategory:   true,
		IsEditable:   false,
		Label:        "formats",
		Name:         "Formats",
		Table:        "data",
	},

	ID: Model{
		CategorySort: "",
		Column:       "id",
		IsEditable:   false,
		Label:        "id",
		Name:         "ID",
		Table:        "books",
	},

	Identifiers: Model{
		Column:       "type || ':' || val",
		CategorySort: "val",
		IsCategory:   true,
		IsEditable:   true,
		Label:        "identifiers",
		Name:         "Identifiers",
		Table:        "identifiers",
	},

	Languages: Model{
		CategorySort: "lang_code",
		Column:       "lang_code",
		IsCategory:   true,
		IsEditable:   true,
		JoinTable:    "books_languages_link",
		Label:        "languages",
		LinkColumn:   "lang_code",
		Name:         "Languages",
		Table:        "languages",
	},

	LastModified: Model{
		CategorySort: "last_modified",
		Column:       "date(last_modified)",
		IsEditable:   false,
		Label:        "last_modified",
		Name:         "Modified",
		Table:        "books",
	},

	Path: Model{
		CategorySort: "path",
		Column:       "path",
		IsEditable:   false,
		Label:        "path",
		Name:         "Path",
	},

	Pubdate: Model{
		CategorySort: "pubdate",
		Column:       "date(pubdate)",
		IsEditable:   true,
		Label:        "pubdate",
		Name:         "Published",
		Table:        "books",
	},

	Publisher: Model{
		CategorySort: "name",
		Column:       "name",
		IsCategory:   true,
		IsEditable:   true,
		JoinTable:    "books_publishers_link",
		Label:        "publisher",
		LinkColumn:   "publisher",
		Name:         "Publisher",
		Table:        "publishers",
	},

	Rating: Model{
		CategorySort: "rating",
		Column:       "rating",
		IsCategory:   true,
		IsEditable:   true,
		JoinTable:    "books_ratings_link",
		Label:        "rating",
		LinkColumn:   "rating",
		Name:         "Rating",
		Table:        "ratings",
	},

	Series: Model{
		CategorySort: "name",
		Column:       "name",
		IsCategory:   true,
		IsEditable:   true,
		JoinTable:    "books_series_link",
		Label:        "series",
		LinkColumn:   "series",
		Name:         "Series",
		Table:        "series",
	},

	SeriesIndex: Model{
		CategorySort: "series_index",
		Column:       "books.series_index",
		IsEditable:   true,
		Label:        "series_index",
		Name:         "Position",
	},

	Sort: Model{
		CategorySort: "sort",
		Column:       "sort",
		IsEditable:   true,
		Label:        "sort",
		Name:         "SortAs",
	},

	Tags: Model{
		CategorySort: "name",
		Column:       "name",
		IsCategory:   true,
		IsEditable:   true,
		JoinTable:    "books_tags_link",
		Label:        "tags",
		LinkColumn:   "tag",
		Name:         "Tags",
		Table:        "tags",
	},

	Timestamp: Model{
		CategorySort: "timestamp",
		Column:       "date(timestamp)",
		IsEditable:   false,
		Label:        "timestamp",
		Name:         "Added",
		Table:        "books",
	},

	Title: Model{
		CategorySort: "sort",
		Column:       "title",
		IsEditable:   true,
		Label:        "title",
		Name:         "Title",
		Table:        "books",
	},

	UUID: Model{
		CategorySort: "uuid",
		Column:       "uuid",
		IsEditable:   false,
		Label:        "uuid",
		Name:         "UUID",
	},
}
