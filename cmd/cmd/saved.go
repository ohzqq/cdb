package cmd

import (
	"github.com/spf13/cobra"
)

// savedCmd represents the saved command
var savedCmd = &cobra.Command{
	Use:   "saved",
	Short: "access saved searches",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	searchCmd.AddCommand(savedCmd)
}
