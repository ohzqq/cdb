package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/danielgtaylor/casing"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// metaShowCmd represents the show command
var metaShowCmd = &cobra.Command{
	Use:    "show [ID]",
	Short:  "show a book's metadata",
	Args:   cobra.ExactArgs(1),
	PreRun: debug,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		lib := cdb.GetLib(viper.GetString("lib"))
		s := lib.NewSearch().GetByID(id)
		r, err := s.Results()
		if err != nil {
			log.Fatal(err)
		}

		if len(r) < 1 {
			fmt.Printf("no book with id %s\n", id)
			os.Exit(1)
		}

		for _, b := range r {
			if cmd.Flags().Changed("save") {
				f, err := os.Create(casing.Snake(b.Title) + ".yaml")
				if err != nil {
					log.Fatal(err)
				}
				err = yaml.NewEncoder(f).Encode(b)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				var buf bytes.Buffer
				err := yaml.NewEncoder(&buf).Encode(b)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(buf.String())
			}
		}
	},
}

func init() {
	metaCmd.AddCommand(metaShowCmd)
	metaShowCmd.Flags().BoolP("save", "s", false, "save meta to disk")
}
