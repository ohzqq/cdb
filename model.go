package cdb

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
	SortAs        = "sort"
	Tags          = "tags"
	Timestamp     = "timestamp"
	Title         = "title"
	UUID          = "uuid"
	Duration      = "#duration"
	Narrators     = "#narrators"
)

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

func DefaultModels() Models {
	return modelMeta
}

func AudiobookModels() Models {
	models := Models{
		Duration:  durationModel,
		Narrators: narratorsModel,
	}
	return models
}

func AllModels() Models {
	models := DefaultModels()
	for l, m := range AudiobookModels() {
		models[l] = m
	}
	return models
}

func (m Model) ToSql() string {
	switch {
	case m.isManyToOne():
		stmt, _ := manyToOne(m)
		return stmt
	case m.isManyToMany():
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

func (m Model) isOneToOne() bool {
	for _, l := range oneToOneModels {
		if m.Label == l {
			return true
		}
	}
	return false
}

func (m Model) isManyToOne() bool {
	if m.Label == Duration {
		return true
	}
	for _, l := range manyToOneModels {
		if m.Label == l {
			return true
		}
	}
	return false
}

func (m Model) isManyToMany() bool {
	if m.Label == Narrators {
		return true
	}
	for _, l := range manyToManyModels {
		if m.Label == l {
			return true
		}
	}
	return false
}

func (m Models) Editable() []string {
	var edit []string
	for l, mod := range m {
		if mod.IsEditable {
			edit = append(edit, l)
		}
	}
	return edit
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
