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
	Name         string
	Path         string
	isAudiobooks bool
}

// Option is a set a library option.
type Option func(*Lib)

// IsAudiobooks marks a library as containing audiobooks.
func IsAudiobooks() Option {
	return func(l *Lib) {
		l.isAudiobooks = true
	}
}

// NewLib initializes a library.
func NewLib(name, path string, opts ...Option) *Lib {
	lib := &Lib{
		Name: name,
		Path: path,
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

// NewQuery initializes a database query.
func (l *Lib) NewQuery() *Query {
	p := filepath.Join(l.Path, l.Name, metaDB)
	err := l.Connect(p)
	if err != nil {
		log.Fatal(err)
	}

	models := DefaultModels()

	if l.isAudiobooks {
		am, err := l.getAudiobookColumns()
		if err != nil {
			log.Fatal(err)
		}
		for l, m := range am {
			models[l] = m
		}
	}

	cols := []string{
		//fmt.Sprintf("'%s' source", l.Name),
		fmt.Sprintf("'source', JSON_QUOTE('%s')", l.Name),
	}
	for _, m := range models {
		//fmt.Println(m.ToSqlJSON())
		//cols = append(cols, m.ToSql())
		cols = append(cols, m.ToSqlJSON())
	}

	return NewQueryJSON(cols)
}
