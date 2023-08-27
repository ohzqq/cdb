//go:build exclude

package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/cdb/calibredb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm ID...",
	Short: "remove a book from the library",
	Run: func(cmd *cobra.Command, args []string) {
		remove := calibredb.Remove(viper.GetString("lib"), args...)
		_, err := remove.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s removed: %v\n", viper.GetString("lib"), args)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
