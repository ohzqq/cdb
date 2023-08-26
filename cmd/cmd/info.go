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
		fmt.Printf("%v\n", cdb.ListLibraries())
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
