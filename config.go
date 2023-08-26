package cdb

import (
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

func ListLibraries() []string {
	libs := viper.GetStringMap("libraries")
	return maps.Keys(libs)
}

func DefaultLibrary() string {
	return viper.GetString("options.default")
}
