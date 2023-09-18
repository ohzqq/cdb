package cdb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

var lib *Lib
var books []Book

func TestNewLib(t *testing.T) {
	lib = NewLib(
		"audiobooks",
		viper.GetString("libraries.audiobooks.path"),
		IsAudiobooks(),
		PrintQuery(),
	)
}

func booksByID() ([]Book, error) {
	q := lib.NewQuery().GetByID(9783) //.Limit(1) //.GetByID(1)
	return getBooks(q)
}

func TestQueryByID(t *testing.T) {
	books, err := booksByID()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	if len(books) != 1 {
		t.Error("expected one book")
	}
}

func TestSaveMetadata(t *testing.T) {
	books, err := booksByID()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	for _, book := range books {
		name := filepath.Join("testdata", book.Title)
		for _, ext := range []string{".json", ".toml", ".yaml"} {
			err := book.Save(name, ext, true)
			if err != nil {
				fmt.Printf("error %v\n", err)
			}
		}
	}
}

func getBooks(q *Query) ([]Book, error) {
	r, err := lib.GetBooks(q)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	books, err := r.Books()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	for _, b := range books {
		fmt.Printf("%#v\n", b)
	}
	return r.Books()
}

func TestReadMetadata(t *testing.T) {
	files, err := filepath.Glob("testdata/*")
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	books, err := booksByID()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	book := books[0].editableStringMapString()

	for _, file := range files {
		nb := &Book{}
		err := nb.ReadFile(file)
		if err != nil {
			fmt.Printf("error %v\n", err)
		}

		for k, v := range nb.editableStringMapString() {
			if o := book[k]; o != v {
				t.Errorf("original meta %#v != read meta %#v\n", o, v)
			}
		}
	}

}
