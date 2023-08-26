package cmd

import (
	"fmt"

	"github.com/ohzqq/ur"
	"github.com/ohzqq/ur/cdb"
	"github.com/ohzqq/ur/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm ID...",
	Short: "remove a book from the library",
	Run: func(cmd *cobra.Command, args []string) {
		remove := cdb.Remove(viper.GetString("lib"), args...)
		_, err := remove.Run()
		ur.HandleError("", err)

		err = search.Index(viper.GetString("lib")).DeleteDocuments(args)
		ur.HandleError("index delete", err)

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
