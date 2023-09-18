package cmd

import (
	"fmt"
	"log"

	"github.com/ohzqq/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:    "info",
	Short:  "info about your calibre libs",
	PreRun: debug,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("default Lib: %v\n", cdb.DefaultLibrary())
		//stmt, arg := cdb.GetPreferences()
		//fmt.Printf("default Lib: %v\n %v\n", stmt, arg)
		lib := cdb.NewLib(
			"audiobooks",
			viper.GetString("libraries.audiobooks.path"),
			cdb.IsAudiobooks(),
			cdb.PrintQuery(),
		)
		fmt.Printf("All Libraries: %v\n", lib)

		s := lib.NewQuery().GetByID(9783) //.Limit(1) //.GetByID(1)
		r, err := lib.GetBooks(s)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(string(r))

		books, err := r.Books()
		if err != nil {
			log.Fatal(err)
		}

		for _, b := range books {
			fmt.Printf("%s\n", b)
			fmt.Printf("%s\n", b.CalibredbFlags())
			//ext := ".json"
			//fmt.Printf("%s\n", nb)

			//db
		}
		//println(len(books))
		//  m := b.(string)
		//  d, err := json.Marshal(json.RawMessage(m))
		//  if err != nil {
		//    log.Fatal(err)
		//  }
		//  //bm := make(map[string]any)
		//  //err = json.Unmarshal([]byte(m), &bm)
		//  //if err != nil {
		//  //log.Fatal(err)
		//  //}
		//  fmt.Printf("%+V\n", string(d))
		//  //  t, err := time.Parse(time.RFC3339, bm["modified"].(string))
		//  //  //t, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

		//  //  if err != nil {
		//  //    log.Fatal(err)
		//  //  }
		//  //  delete(bm, "duration")
		//  //p := []any{map[string]any{"metadata": m}}
		//  //f := map[string]any{"publications": p}
		//  //fd, err := json.Marshal(f)
		//  //if err != nil {
		//  //log.Fatal(err)
		//  //}
		//  //feed, err := opds2.ParseBuffer(fd)
		//  //if err != nil {
		//  //log.Fatal(err)
		//  //}
		//  //fmt.Printf("%s\n", feed)
		//  //for _, f := range feed.Publications {
		//  //fmt.Printf("%v\n", f.Metadata)
		//  //}
		//}

		//showMeta()

	},
}

func showMeta() {
	caldb, err := cdb.Calibredb(
		viper.GetString("calibre.url"),
		cdb.Authenticate(viper.GetString("calibre.username"), viper.GetString("calibre.password")),
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := caldb.ShowMetadata("1")
	out, err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	if out != "" {
		println(out)
	}

}

var testURL = "http://localhost:8080/"

func init() {
	rootCmd.AddCommand(infoCmd)
}
