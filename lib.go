package cdb

import (
	"log"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
)

type Lib struct {
	db           *DB
	Name         string
	Path         string
	isAudiobooks bool
}

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

func (l *Lib) GetBooks(q sq.Sqlizer) ([]*Book, error) {
	stmt, args, err := q.ToSql()
	if err != nil {
		log.Fatal(err)
	}
	return l.db.getBooks(stmt, args)
}

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

	var cols []string
	for _, m := range models {
		cols = append(cols, m.ToSql())
	}

	return NewQuery(cols)
}
