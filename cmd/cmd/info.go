package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "info about your calibre libs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("All Libraries: %v\n", cdb.ListLibraries())
		fmt.Printf("default Lib: %v\n", cdb.DefaultLibrary())
		lib := cdb.GetLib("audiobooks")
		//fmt.Printf("db models: %v\n", lib.Models)
		s := lib.NewQuery().Limit(5)
		r, err := s.Results()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("search: %+V\n", r)
		fmt.Printf("editable %v\n", cdb.AllModels().Editable())
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
