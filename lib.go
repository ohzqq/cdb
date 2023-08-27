package cdb

import (
	"log"
	"path/filepath"

	sq "github.com/Masterminds/squirrel"
)

type Lib struct {
	db           *DB
	query        *Query
	Name         string
	Path         string
	Audiobooks   bool
	isAudiobooks bool
}

func NewLib(name, path string, opts ...Opt) *Lib {
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

//func (lib *Lib) GetSavedSearches() map[string]string {
//  base := GetDB(lib)
//  saved := base.GetPreference("savedSearches")
//  var searches map[string]string
//  err := json.Unmarshal(saved, &searches)
//  ur.HandleError("saved searches query", err)
//  return searches
//}

//func (lib *Lib) GetHiddenCategories() []string {
//  base := GetDB(lib.Name)
//  hidden := base.GetPreference("hiddenCategories")
//  var cats []string
//  err := json.Unmarshal(hidden, &cats)
//  ur.HandleError("hidden categories query", err)
//  return cats
//}
