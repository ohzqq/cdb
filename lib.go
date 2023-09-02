package cdb

import (
	"fmt"
	"log"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
	"github.com/ohzqq/libopds2-go/opds1"
	"github.com/ohzqq/libopds2-go/opds2"
)

// Lib represents a calibre library.
type Lib struct {
	db           *DB
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
		db:   &DB{},
	}
	for _, opt := range opts {
		opt(lib)
	}
	return lib
}

func (l *Lib) OPDS1Feed() opds1.Feed {
	feed := opds1.Feed{}
	return feed
}

func (l *Lib) OPDS2Feed() opds2.Feed {
	feed := opds2.Feed{}
	return feed
}

func (l *Lib) GetPubs(q sq.Sqlizer) ([]opds2.Publication, error) {
	var pubs []opds2.Publication
	stmt, args, err := q.ToSql()
	if err != nil {
		return pubs, err
	}
	books, err := l.db.getBooks(stmt, args)
	if err != nil {
		return pubs, err
	}
	for _, book := range books {
		pubs = append(pubs, book.ToOPDS2Publication())
	}
	return pubs, nil
}

// GetBooks runs a database query. Takes an squirrel.Sqlizer interface to
// generate the sql.
func (l *Lib) GetBooks(q sq.Sqlizer) ([]*Book, error) {
	stmt, args, err := q.ToSql()
	if err != nil {
		log.Fatal(err)
	}
	return l.db.getBooks(stmt, args)
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

	return l.db.getPreferences(stmt, args...)
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
	err := l.db.Connect(p)
	if err != nil {
		log.Fatal(err)
	}

	models := DefaultModels()

	if l.isAudiobooks {
		am, err := l.db.getAudiobookColumns()
		if err != nil {
			log.Fatal(err)
		}
		for l, m := range am {
			models[l] = m
		}
	}

	cols := []string{
		fmt.Sprintf("'%s' source", l.Name),
	}
	for _, m := range models {
		cols = append(cols, m.ToSql())
	}

	return NewQuery(cols)
}
