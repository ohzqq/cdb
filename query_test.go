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
	path := viper.GetString("libraries.audiobooks.path")
	p := filepath.Join(path, "audiobooks")
	lib = NewLib(
		p,
		IsAudiobooks(),
		//PrintQuery(),
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

//func TestSaveMetadata(t *testing.T) {
//  books, err := booksByID()
//  if err != nil {
//    fmt.Printf("error %v\n", err)
//  }

//  for _, book := range books {
//    name := filepath.Join("testdata", "book", casing.Snake(book.Title))
//    for _, ext := range []string{".json", ".toml", ".yaml"} {
//      //var opt EncoderOpt
//      switch ext {
//      case ".json":
//        //opt = EncodeJSON
//      case ".toml":
//        //opt = EncodeTOML
//      case ".yaml":
//        err := book.Encode(EncodeYAML).WriteTo(os.Stdout)
//        //err := NewEnc(book, EncodeYAML).WriteTo(os.Stdout)
//        if err != nil {
//          t.Error(err)
//        }
//        //opt = EncodeYAML
//      }
//      println(name)
//      //err := book.Save(name + ext)
//      //if err != nil {
//      //  fmt.Printf("error %v\n", err)
//      //}
//    }
//  }
//}

//func TestCalibreFlags(t *testing.T) {
//  books, err := booksByID()
//  if err != nil {
//    fmt.Printf("error %v\n", err)
//  }
//  for _, book := range books {
//    flags := book.CalibredbFlags()
//  }
//}

func getBooks(q *Query) ([]Book, error) {
	r, err := lib.GetBooks(q)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	//books, err := r.Books()
	//if err != nil {
	//  fmt.Printf("error %v\n", err)
	//}
	//for _, b := range books {
	//  fmt.Printf("%#v\n", b)
	//}
	return r.Books()
}

var flags = []string{
	"authors:Charles Dean & R.A. Mejia",
	"publisher:Single Player Books",
	"#narrators:Daniel Wisniewski & Elizabeth Plant",
	"series_index:2",
	"#duration:12:23:00",
	"languages:eng",
	"series:Upgrade Apocalypse",
	"tags:litrpg",
	"title:The Upgrade Apocalypse: Book 2",
	"author_sort:Charles Dean & R.A. Mejia",
	"pubdate:0101-01-01",
	"sort:Upgrade Apocalypse: Book 2, The",
	"comments:<p>Welcome to the end of the world, where the lines between reality and a fantasy game blur. The survivors must not only face the relentless hordes of the Hell-Cursed, but also the cruel fate of death itself. </p> <p>Join Chedderfield and his companions as they navigate a post-apocalyptic wasteland filled with dark humor, horror, and heart-wrenching deaths. Get ready for twists, turns, and gut-wrenching losses as you follow the journey to uncover the truth behind Archimedes' death and the true meaning of survival in a world where anything can happen.</p> <p>Welcome to hell on Earth. Welcome to <i>The Upgrade Apocalypse!</i></p> <p>Warning: This book contains apocalyptic adventures, creative cursing, jokes about the apocalypse, an interesting card skill system, horror monsters to fight or flee from, and character deaths when they fail their save throws.</p>",
}
