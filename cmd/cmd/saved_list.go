//go:build exclude

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// savedListCmd represents the list command
var savedListCmd = &cobra.Command{
	Use:   "list",
	Short: "list saved searches",
	Run: func(cmd *cobra.Command, args []string) {
		//saved := db.GetSavedSearches(viper.GetString("lib"))
		var titles []string
		var searches []string
		//for t, s := range saved {
		//titles = append(titles, t)
		//searches = append(searches, s)
		//}
		//m := list.New(titles)
		//sel := m.Run()

		//var save string
		//for _, i := range sel {
		//fmt.Printf("saved %+V\n", searches[i])
		//save = searches[i]
		//break
		//}

		//lib := viper.GetString("lib")
		//c := cdb.Search(lib, save)
		//o, err := c.Run()
		//ur.HandleError("", err)

		//q := ur.NewQuery("books", lib).
		//QueryId(o)

		//books := search.Books(q)

		//pubs := tui.ListBooks(books)
		//for _, b := range pubs {
		//fmt.Printf("pubs %+V\n", b.Metadata.Title.String())
		//}
		fmt.Printf("%v %v\n", titles, searches)
	},
}

func init() {
	savedCmd.AddCommand(savedListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
