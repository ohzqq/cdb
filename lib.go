package cdb

import (
	"fmt"
	"log"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
)

// Lib represents a calibre library.
type Lib struct {
	*DB
	path         string
	Name         string
	Path         string
	IsAudiobooks bool
}

// Option is a set a library option.
type Option func(*Lib)

// IsAudiobooks marks a library as containing audiobooks.
func IsAudiobooks() Option {
	return func(l *Lib) {
		l.IsAudiobooks = true
	}
}

// PrintQuery prints the database query to stdout.
func PrintQuery() Option {
	return func(l *Lib) {
		l.printQuery = true
	}
}

// NewLib initializes a library.
func NewLib(path string, opts ...Option) *Lib {
	lib := &Lib{
		path: path,
		DB:   &DB{},
	}
	for _, opt := range opts {
		opt(lib)
	}
	return lib
}

// GetBooks runs a database query. Takes an squirrel.Sqlizer interface to
// generate the sql.
func (l *Lib) GetBooks(q sq.Sqlizer) (Records, error) {
	stmt, args, err := q.ToSql()
	if err != nil {
		log.Fatal(err)
	}
	return l.bookQuery(stmt, args)
}

// NewQuery initializes a database query.
func (l *Lib) NewQuery() *Query {
	p := filepath.Join(l.path, metaDB)
	err := l.Connect(p)
	if err != nil {
		log.Fatal(err)
	}

	name := filepath.Base(l.path)
	cols := []string{
		fmt.Sprintf("'source', JSON_QUOTE('%s')", name),
	}
	for _, m := range l.models() {
		cols = append(cols, m.ToSqlJSON())
	}

	return NewQuery(cols)
}

func (l *Lib) models() Models {
	models := DefaultModels()
	if l.IsAudiobooks {
		am, err := l.getAudiobookColumns()
		if err != nil {
			log.Fatal(err)
		}
		for l, m := range am {
			models[l] = m
		}
	}
	return models
}

// GetPreference gets a calibre preference.
func (l *Lib) GetPreference(c string) ([]byte, error) {
	col := "JSON(val) " + c

	sel := sq.Select(col).
		From("preferences").
		Where(sq.Eq{"key": c})

	stmt, args, err := sel.ToSql()
	if err != nil {
		return []byte{}, err
	}

	return l.getPreferences(stmt, args...)
}

// ListPreferences lists the retrievable calibre preferences.
func ListPreferences() []string {
	return []string{
		SavedSearches,
		DisplayFields,
		HiddenCategories,
		FieldMetadata,
	}
}
