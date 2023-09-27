package cdb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestEncoder(t *testing.T) {
	books, err := booksByID()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	for _, b := range books {
		//enc := NewEncoder(b, EncodeJSON)
		enc := NewEncoder(b, EncodeJSON, EditableOnly())
		err := enc.WriteTo(os.Stdout)
		if err != nil {
			t.Error(err)
		}
		err = enc.Save("")
		if err != nil {
			t.Error(err)
		}

		//enc.FEncode(os.Stdout, EncodeYAML(enc))
		//for _, f := range []string{".json", ".yaml", ".yml", ".toml"} {
		//b.SetEncoder(WithFormat(f))
		//b.Encode(os.Stdout)
		//}
	}
}

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
