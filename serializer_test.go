package cdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielgtaylor/casing"
)

func TestEncoder(t *testing.T) {
	books, err := booksByID()
	if err != nil {
		t.Errorf("error %v\n", err)
	}

	for _, b := range books {
		enc := NewEncoder(&b, EditableOnly())
		enc.Encoder(EncodeJSON)
		err := enc.WriteTo(os.Stdout)
		if err != nil {
			t.Error(err)
		}
	}
	//testSave(t, books)
}

const bookData = `testdata/book/`

func testSave(t *testing.T, books []Book) {
	for _, b := range books {
		for _, init := range []EncoderInit{EncodeJSON, EncodeYAML, EncodeTOML} {
			enc := NewEncoder(&b, EditableOnly())
			enc.Encoder(init)
			name := filepath.Join(bookData, casing.Snake(b.Title))

			file, err := os.Create(name + enc.Format)
			if err != nil {
				t.Error(err)
			}
			defer file.Close()

			err = enc.WriteTo(file)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestDecoder(t *testing.T) {
	books, err := booksByID()
	if err != nil {
		t.Errorf("error %v\n", err)
	}

	og := books[0]

	for _, book := range decodeBooks(t) {
		if book.Title != og.Title {
			t.Errorf("decoded title was %s, expected %s\n", book.Title, og.Title)
		}
	}
}

func decodeBooks(t *testing.T) []Book {
	files, err := filepath.Glob("testdata/book/*")
	if err != nil {
		t.Errorf("error %v\n", err)
	}

	var books []Book
	for _, file := range files {
		var init DecoderInit
		ext := filepath.Ext(file)
		switch ext {
		case ".yaml", ".yml":
			init = DecodeYAML
		case ".toml":
			init = DecodeTOML
		case ".json":
			init = DecodeJSON
		}
		var book Book
		s := NewEncoder(&book).Decoder(init)

		err := s.ReadFile(file)
		if err != nil {
			t.Errorf("%v", err)
		}
		books = append(books, book)
	}
	return books
}
