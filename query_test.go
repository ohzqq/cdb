package cdb

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestQueryByID(t *testing.T) {
	lib := NewLib(
		"audiobooks",
		viper.GetString("libraries.audiobooks.path"),
		IsAudiobooks(),
		PrintQuery(),
	)

	s := lib.NewQuery().GetByID(9783) //.Limit(1) //.GetByID(1)
	r, err := lib.GetBooks(s)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	books, err := r.Books()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	if len(books) != 1 {
		t.Error("expected one book")
	}
}
