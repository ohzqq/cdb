package cmd

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ohzqq/ur"
	"github.com/ohzqq/ur/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	set := cdb.SetMetadata(viper.GetString("lib"), id)
	f, err := os.Open(file)
	ur.HandleError("", err)

	fields := make(map[string]string)
	_, err = toml.NewDecoder(f).Decode(&fields)
	ur.HandleError("", err)

	set.Opt(cdb.Fields(fields))
	out, err := set.Run()
	ur.HandleError("", err)
	if out != "" {
		fmt.Println(out)
	}
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
