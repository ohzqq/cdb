package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:    "info",
	Short:  "info about your calibre libs",
	PreRun: debug,
	Run: func(cmd *cobra.Command, args []string) {
		lib := cdb.NewLib(
			"audiobooks",
			viper.GetString("libraries.audiobooks.path"),
			cdb.IsAudiobooks(),
		)
		fmt.Printf("All Libraries: %v\n", lib)
	},
}

func showMeta() {
	caldb, err := cdb.Calibredb(
		viper.GetString("calibre.url"),
		cdb.Authenticate(viper.GetString("calibre.username"), viper.GetString("calibre.password")),
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := caldb.ShowMetadata("1")
	out, err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	if out != "" {
		println(out)
	}

}

var testURL = "http://localhost:8080/"

func init() {
	rootCmd.AddCommand(infoCmd)
}
