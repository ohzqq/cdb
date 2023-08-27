//go:build exclude

package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/cdb/calibredb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	metaFile   string
	otherFiles []string
	coverFile  string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "add book to calibre",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		in := args[0]
		id := addBook(in)

		if cmd.Flags().Changed("meta") {
			setMeta(id, metaFile)
		}

		if cmd.Flags().Changed("formats") {
			addFormats(id, otherFiles)
		}

		// q := ur.NewQuery("books", viper.GetString("lib")).PathId(id)
		// err := search.Index(viper.GetString("lib")).Add(q)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
}

func addBook(in string) string {
	add := calibredb.Add(viper.GetString("lib"), in)
	if coverFile != "" {
		add.Opt(calibredb.Cover(coverFile))
	}
	id, err := add.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s imported: %s\n", viper.GetString("lib"), id)

	return id
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&metaFile, "meta", "m", "", "metadata for book")
	addCmd.Flags().StringVarP(&coverFile, "cover", "c", "", "book cover")
	addCmd.PersistentFlags().StringSliceVarP(&otherFiles, "formats", "f", []string{}, "other formats")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
