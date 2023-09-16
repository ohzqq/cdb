package cdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestConfigInit(t *testing.T) {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	cfgDir := filepath.Join(home, ".config/cdb/")
	viper.AddConfigPath(cfgDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		t.Errorf("error %v\n", err)
	}
}
