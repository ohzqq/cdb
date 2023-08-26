package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// savedRmCmd represents the rm command
var savedRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "delete a saved search",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rm called")
	},
}

func init() {
	savedCmd.AddCommand(savedRmCmd)
}
