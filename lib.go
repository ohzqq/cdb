package cdb

import (
	"log"

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
