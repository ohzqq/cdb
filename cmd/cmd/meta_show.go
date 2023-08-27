package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// metaShowCmd represents the show command
var metaShowCmd = &cobra.Command{
	Use:   "show [ID]",
	Short: "show a book's metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		lib := cdb.GetLib("audiobooks")
		s := lib.NewSearch().GetByID(id)
		r, err := s.Results()
		if err != nil {
			log.Fatal(err)
		}

		for _, b := range r {
			if cmd.Flags().Changed("save") {
				f, err := os.Create("test")
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
