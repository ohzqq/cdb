package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var savedAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a saved search",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	savedCmd.AddCommand(savedAddCmd)
}
