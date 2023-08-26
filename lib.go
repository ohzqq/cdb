package cdb

import (
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
)

type Lib struct {
	*DB
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

func (l *Lib) ConnectDB() error {
	db, err := Configure(l.Name, l.Path)
	if err != nil {
		return err
	}

	if l.Audiobooks {
		err := db.IsAudiobooks()
		if err != nil {
			return err
		}
	}

	l.DB = db
	return nil
}

func (l *Lib) NewSearch() *Search {
	search := &Search{}

	var cols []string
	for _, m := range l.Models {
		cols = append(cols, m.ToSql())
	}

	search.query = sq.Select(strings.Join(cols, ",\n")).
		From("books")

	return search
}
