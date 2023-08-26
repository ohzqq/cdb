package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/cdb/calibredb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addFormatCmd represents the format command
var addFormatCmd = &cobra.Command{
	Use:   "formats [id]",
	Short: "add formats to book",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("formats") {
			addFormats(args[0], otherFiles)
		}
	},
}

func addFormats(id string, formats []string) {
	for _, f := range formats {
		add := calibredb.AddFormat(viper.GetString("lib"), id, f)
		_, err := add.Run()
		if err == nil {
			fmt.Printf("added file: %s\n", f)
		} else {
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func init() {
	addCmd.AddCommand(addFormatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// formatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// formatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
