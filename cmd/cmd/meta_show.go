package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// metaShowCmd represents the show command
var metaShowCmd = &cobra.Command{
	Use:   "show [IDs]",
	Short: "show a book's metadata",
	Run: func(cmd *cobra.Command, args []string) {
		lib := cdb.GetLib("audiobooks")
		s := lib.NewSearch().GetByID(1)
		r, err := s.Results()
		if err != nil {
			log.Fatal(err)
		}
		for _, b := range r {
			var buf bytes.Buffer
			err := yaml.NewEncoder(&buf).Encode(b)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(buf.String())
		}
	},
}

func init() {
	metaCmd.AddCommand(metaShowCmd)
}
