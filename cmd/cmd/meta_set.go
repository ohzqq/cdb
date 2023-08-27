package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ohzqq/cdb"
	"github.com/ohzqq/cdb/calibredb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// metaSetCmd represents the set command
var metaSetCmd = &cobra.Command{
	Use:    "set",
	Short:  "set metadata for book",
	PreRun: debug,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("list-fields") {
			models := cdb.DefaultModels()
			if cdb.GetLib(viper.GetString("lib")).Audiobooks {
				models = cdb.AudiobookModels()
			}
			fmt.Println(models.Editable())
			os.Exit(0)
		}

		if len(args) != 2 {
			log.Fatal("this command requires two arguments")
		}
		id := args[0]
		file := args[1]

		setMeta(id, file)
	},
}

func setMeta(id, file string) {
	set := calibredb.SetMetadata(viper.GetString("lib"), id)
	f, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	b := cdb.FromYAML(f)

	set.Opt(calibredb.Fields(b.StringMap()))
	out, err := set.Run()
	if err != nil {
		log.Fatal(err)
	}
	if out != "" {
		fmt.Println(out)
	}
}

func decodeMeta(r io.Reader) map[string]string {
	fields := make(map[string]string)
	err := yaml.NewDecoder(r).Decode(&fields)
	if err != nil {
		log.Fatal(err)
	}
	return fields
}

func init() {
	metaCmd.AddCommand(metaSetCmd)
	metaSetCmd.Flags().StringVarP(&metaFile, "meta", "m", "", "metadata for book")
	metaSetCmd.Flags().Bool("list-fields", false, "show editable metadata fields")
}
