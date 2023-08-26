package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ohzqq/cdb/calibredb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// metaSetCmd represents the set command
var metaSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set metadata for book",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		file := args[1]
		setMeta(id, file)
	},
}

func setMeta(id, file string) {
	set := calibredb.SetMetadata(viper.GetString("lib"), id)
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	fields := decodeMeta(f)

	set.Opt(calibredb.Fields(fields))
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
