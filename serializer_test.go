package cdb

import (
	"fmt"
	"path/filepath"
	"testing"
)

//func TestEncoder(t *testing.T) {
//  books, err := booksByID()
//  if err != nil {
//    fmt.Printf("error %v\n", err)
//  }

//  for _, b := range books {
//    for _, f := range []string{".json", ".yaml", ".yml", ".toml"} {
//      b.Print(f, true)
//    }
//  }
//}

func TestDecoder(t *testing.T) {
	files, err := filepath.Glob("testdata/book/*")
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	books, err := booksByID()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	for _, file := range files {
		d, err := ReadMetadataFile(file)
		if err != nil {
			t.Errorf("%v", err)
		}

		var book Book
		d.Decode(&book)
		if books[0].Title != book.Title {
			t.Errorf("decoded title was %s, expected %s", book.Title, books[0].Title)
		}
	}
}
