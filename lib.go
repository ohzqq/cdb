package cdb

import (
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
)

type Lib struct {
	Name       string
	Path       string
	Audiobooks bool
}

func GetLib(name string) *Lib {
	lib := &Lib{}
	err := viper.UnmarshalKey("libraries."+name, lib)
	if err != nil {
		log.Fatal(err)
	}
	lib.Name = name
	return lib
}

func (l *Lib) ConnectDB() (*DB, error) {
	db, err := Configure(l.Name, l.Path)
	if err != nil {
		return db, err
	}

	if l.Audiobooks {
		err := db.IsAudiobooks()
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

func (l *Lib) NewSearch() *Search {
	db, err := l.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	search := &Search{
		db: db,
	}

	var cols []string
	for _, m := range db.Models {
		cols = append(cols, m.ToSql())
	}

	search.query = sq.Select(strings.Join(cols, ",\n")).
		From("books")

	return search
}
