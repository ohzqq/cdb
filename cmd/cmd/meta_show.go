package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// metaShowCmd represents the show command
var metaShowCmd = &cobra.Command{
	Use:   "show [ID]",
	Short: "show a book's metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		lib := cdb.GetLib(viper.GetString("lib"))
		s := lib.NewSearch().GetByID(id)
		r, err := s.Results()
		if err != nil {
			log.Fatal(err)
		}

		for _, b := range r {
			if cmd.Flags().Changed("save") {
				f, err := os.Create(metaFile + ".yaml")
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
	metaShowCmd.Flags().StringVarP(&metaFile, "save", "s", "meta", "save meta to disk")
}
