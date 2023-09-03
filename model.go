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
	CategorySort string `yaml:"category_sort" db:"category_sort"`
	Column       string `yaml:"column" db:"column"`
	Colnum       int    `yaml:"colnum" db:"id"`
	IsCategory   bool   `yaml:"is_category" db:"is_category"`
	IsCustom     bool   `yaml:"is_custom" db:"is_custom"`
	IsEditable   bool   `yaml:"is_editable" db:"is_editable"`
	IsNames      bool   `yaml:"is_names" db:"is_names"`
	JoinTable    string `yaml:"join_table" db:"join_table"`
	Label        string `yaml:"label" db:"label"`
	LinkColumn   string `yaml:"link_column" db:"link_column"`
	Name         string `yaml:"name" db:"name"`
	Table        string `yaml:"table" db:"table"`
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

func (m Model) ToSqlJSON() string {
	switch {
	case slices.Contains(manyToOneModels, m.Label), m.Label == Duration:
		stmt, _ := manyToOneJSON(m)
		return stmt
	case slices.Contains(manyToManyModels, m.Label), m.Label == Narrators:
		stmt, _ := manyToManyJSON(m)
		return stmt
	case m.Label == Cover:
		stmt, _ := coverStmtJSON("")
		return stmt
	default:
		stmt := m.colJSON()
		return stmt
	}
}

func (m Model) colJSON() string {
	switch m.Label {
	case Authors, Narrators, Tags, Formats, Identifiers, Languages:
		return fmt.Sprintf("IFNULL(JSON_GROUP_ARRAY(%s), '[]')", m.Column)
	case Duration, Rating:
		return fmt.Sprintf("JSON(IFNULL(GROUP_CONCAT(%s), 0))", m.Column)
	case Comments, Publisher, Series:
		return fmt.Sprintf("JSON_QUOTE(IFNULL(GROUP_CONCAT(%s), ''))", m.Column)
	case SeriesIndex:
		return fmt.Sprintf(`'%s', JSON(IFNULL(%s, 0))`, m.Label, m.Column)
	default:
		return fmt.Sprintf(`'%s', JSON_QUOTE(IFNULL(%s, ''))`, m.Label, m.Column)
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

func booksColJSON(m Model) (string, []any) {
	return fmt.Sprintf(`'%s', JSON_QUOTE(IFNULL(%s, ''))`, m.Label, m.Column), []any{}
}

func groupConcat(m Model) string {
	sep := catSep
	if m.Label == Authors || m.Label == Narrators {
		sep = namesSep
	}
	return fmt.Sprintf("IFNULL(GROUP_CONCAT(%s, '%s'), '')", m.Column, sep)
}

func groupArray(m Model) string {
	return fmt.Sprintf("IFNULL(JSON_GROUP_ARRAY(%s), '[]')", m.Column)
}

func manyToOneJSON(m Model) (string, []any) {
	sel := sq.Select(m.colJSON()).
		Prefix(fmt.Sprint("'", m.Label, "', (")).
		From(m.Table).
		Where("book=books.id").
		Suffix(")")

	return toSql(sel)
}

func manyToOne(m Model) (string, []any) {
	sel := sq.Select(groupConcat(m)).
		Prefix("(").
		From(m.Table).
		Where("book=books.id").
		Suffix(fmt.Sprint(") '", m.Label, "'"))

	return toSql(sel)
}

func manyToManyJSON(m Model) (string, []any) {
	col := sq.Select(m.colJSON()).
		From(m.Table).
		Where(fmt.Sprint(m.Table, ".id"))

	join := sq.Select(m.LinkColumn).
		From(m.JoinTable).
		Where("book IN (books.id)")

	cs, _ := toSql(col)
	js, _ := toSql(join)

	stmt := fmt.Sprintf("'%s', (%s IN (%s))", m.Label, cs, js)

	return stmt, []any{}
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

func coverStmtJSON(lib string) (string, []any) {
	sel := sq.Case("has_cover").
		When("1", `IFNULL('books/' || books.id || '/cover.jpg', '')`).
		Else("''")

	stmt, args, _ := sel.ToSql()
	return "'cover', " + stmt, args
}

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
	Column:       "strftime('%s', value)",
	CategorySort: "value",
	IsCustom:     true,
	IsEditable:   true,
	Label:        "#duration",
	Name:         "duration",
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
	Name:         "narrator",
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
		Name:         "author",
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
		Name:       "description",
		Table:      "comments",
	},

	Cover: Model{
		Column:     "cover",
		IsEditable: false,
		IsNames:    false,
		Label:      "cover",
		Name:       "images",
	},

	Formats: Model{
		Column:       `'books/' || books.id || '/' || data.name || '.' || lower(format)`,
		CategorySort: "format",
		IsCategory:   true,
		IsEditable:   false,
		Label:        "formats",
		Name:         "links",
		Table:        "data",
	},

	ID: Model{
		CategorySort: "",
		Column:       "id",
		IsEditable:   false,
		Label:        "id",
		Name:         "identifier",
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
		Name:         "language",
		Table:        "languages",
	},

	LastModified: Model{
		CategorySort: "last_modified",
		Column:       "strftime('%Y-%m-%dT%H:%M:%S', last_modified) || 'Z'",
		IsEditable:   false,
		Label:        "last_modified",
		Name:         "lastModified",
		Table:        "books",
	},

	Path: Model{
		CategorySort: "path",
		Column:       "path",
		IsEditable:   false,
		Label:        "path",
		Name:         "path",
	},

	Pubdate: Model{
		CategorySort: "pubdate",
		Column:       "strftime('%Y-%m-%dT%H:%M:%S', pubdate) || 'Z'",
		IsEditable:   true,
		Label:        "pubdate",
		Name:         "published",
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
		Name:         "publisher",
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
		Name:         "rating",
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
		Name:         "series",
		Table:        "series",
	},

	SeriesIndex: Model{
		CategorySort: "series_index",
		Column:       "books.series_index",
		IsEditable:   true,
		Label:        "series_index",
		Name:         "position",
	},

	Sort: Model{
		CategorySort: "sort",
		Column:       "sort",
		IsEditable:   true,
		Label:        "sort",
		Name:         "sortAs",
	},

	Tags: Model{
		CategorySort: "name",
		Column:       "name",
		IsCategory:   true,
		IsEditable:   true,
		JoinTable:    "books_tags_link",
		Label:        "tags",
		LinkColumn:   "tag",
		Name:         "subject",
		Table:        "tags",
	},

	Timestamp: Model{
		CategorySort: "timestamp",
		Column:       "strftime('%Y-%m-%dT%H:%M:%S', timestamp) || 'Z'",
		IsEditable:   false,
		Label:        "timestamp",
		Name:         "modified",
		Table:        "books",
	},

	Title: Model{
		CategorySort: "sort",
		Column:       "title",
		IsEditable:   true,
		Label:        "title",
		Name:         "title",
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
