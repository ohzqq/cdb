package gen

import (
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestGenFuncs(t *testing.T) {
	commandFuncs()
	genCommand()
}

func genCommand() {
	cmds := ReadCommands()

	f, err := os.Create("gen_calibre.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(`// Code generated
// DO NOT EDIT.

package cdb

`)
	if err != nil {
		log.Fatal(err)
	}

	ls := cmds.CommandList()
	_, err = f.Write(gofmt(ls))

	b := cmds.CommandFuncs()
	_, err = f.Write(gofmt(b))
}

func commandFuncs() {
	cmds := ReadCommands()

	f, err := os.Create("gen_calibre_funcs.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(`// Code generated
// DO NOT EDIT.

package cdb

`)
	if err != nil {
		log.Fatal(err)
	}

	b := cmds.CommandBuilder()
	_, err = f.Write(gofmt(b))
}

func gofmt(in string) []byte {
	c := exec.Command("gofmt")
	stdin, err := c.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, in)
	}()

	out, err := c.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return out
}
