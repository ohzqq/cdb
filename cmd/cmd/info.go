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
		//fmt.Printf("default Lib: %v\n", cdb.DefaultLibrary())
		//stmt, arg := cdb.GetPreferences()
		//fmt.Printf("default Lib: %v\n %v\n", stmt, arg)
		lib := cdb.NewLib(
			"audiobooks",
			viper.GetString("libraries.audiobooks.path"),
			cdb.IsAudiobooks(),
		)
		fmt.Printf("All Libraries: %v\n", lib)

		s := lib.NewQuery().GetByID(1)
		r, err := lib.GetBooks(s)
		if err != nil {
			log.Fatal(err)
		}

		for _, b := range r {
			fmt.Printf("%v\n", b.ID)
			setMeta(lib.Path, "1", b)
		}

		//c, err := calibredb.NewCommand(testURL, calibredb.DryRun())
		//if err != nil {
		//log.Fatal(err)
		//}
		//c.Run()
	},
}

func setMeta(path, id string, b *cdb.Book) {
	set := cdb.SetMetadata(path, id, b.StringMap())

	set.Opt(cdb.DryRun())
	out, err := set.Run()
	if err != nil {
		log.Fatal(err)
	}
	if out != "" {
		fmt.Println(out)
	}
}

var testURL = "http://localhost:8080/"

func init() {
	rootCmd.AddCommand(infoCmd)
}
