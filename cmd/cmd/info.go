package cmd

import (
	"fmt"

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
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
