package main

import (
	"log"

	"github.com/ohzqq/cdb/cmd/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
