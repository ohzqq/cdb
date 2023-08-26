package cdb

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Lib struct {
	Name       string
	Path       string
	Audiobooks bool
	db         *DB
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
	db, err := Configure(l.Name, l.Path, l.Audiobooks)
}

func ErrFileNotExist(path string) error {
	if !FileExist(path) {
		return fmt.Errorf("%v does not exist or cannot be found, check the path in the config, error: \n", path)
	}
	return nil
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
