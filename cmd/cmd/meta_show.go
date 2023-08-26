package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metaShowCmd represents the show command
var metaShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show a book's metadata",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("show called")
	},
}

func init() {
	metaCmd.AddCommand(metaShowCmd)
}
