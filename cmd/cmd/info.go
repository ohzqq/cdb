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
		fmt.Printf("Lib: %v\n", lib)
		err := lib.ConnectDB()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("db is connected: %v\n", lib.IsConnected())
		//fmt.Printf("db models: %v\n", lib.Models)
		s := lib.NewSearch()
		fmt.Printf("search: %+V\n", s)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
