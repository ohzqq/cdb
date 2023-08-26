package calibredb

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/ohzqq/ur/cfg"
	"github.com/ohzqq/ur/util"
	"github.com/spf13/viper"
)

type Opt func(*Command)

type Command struct {
	CdbCmd     string
	flags      []string
	positional []string
	tmp        *os.File
	verbose    bool
}

func New(lib string, args ...Opt) *Command {
	cmd := &Command{}
	args = append(args, WithLib(lib))
	for _, fn := range args {
		fn(cmd)
	}

	return cmd
}

func Cmd(cmd string) Opt {
	return func(c *Command) {
		c.CdbCmd = cmd
	}
}

func (c *Command) Opt(opts ...Opt) {
	for _, fn := range opts {
		fn(c)
	}
}

func Flags(flag ...string) Opt {
	return func(c *Command) {
		c.flags = append(c.flags, flag...)
	}
}

func PositionalArgs(args ...string) Opt {
	return func(c *Command) {
		c.positional = append(c.positional, args...)
	}
}

func Verbose() Opt {
	return func(cmd *Command) {
		cmd.verbose = true
	}
}

func WithLib(name string) Opt {
	flags := []string{
		"--with-library",
	}
	lib := cfg.Lib(name)

	var uri *url.URL
	var err error
	if cal := viper.GetString("calibre.url"); cal != "" {
		uri, err = url.Parse(cal)
		if err != nil {
			log.Fatal(err)
		}
		uri.Fragment = lib.Name()
	}

	switch util.SrvIsOnline(uri) {
	case false:
		flags = append(flags, lib.Path())
	default:
		flags = append(flags, uri.String())
		if u := viper.GetString("calibre.username"); u != "" {
			flags = append(flags, "--username", u)
		}
		if p := viper.GetString("calibre.password"); p != "" {
			flags = append(flags, "--password", p)
		}
	}
	return Flags(flags...)
}

func (c *Command) Build() *exec.Cmd {
	return exec.Command("calibredb", c.ParseArgs()...)
}

func (c *Command) ParseArgs() []string {
	var args []string
	args = append(args, c.CdbCmd)
	args = append(args, c.flags...)
	args = append(args, c.positional...)
	return args
}

func (c *Command) DryRun() {
	cmd := c.Build()
	fmt.Println(cmd.String())
}

func (c *Command) Run() (string, error) {
	if c.tmp != nil {
		defer os.Remove(c.tmp.Name())
	}

	cmd := c.Build()

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("%v finished with error: %v\n", c.CdbCmd, stderr.String())
	}

	if viper.GetBool("verbose") {
		fmt.Println(cmd.String())
		fmt.Println(stdout.String())
	}

	var output string
	if len(stdout.Bytes()) > 0 {
		out := stdout.String()
		switch c.CdbCmd {
		case "add":
			sp := strings.Split(out, ": ")
			output = sp[1]
		case "search", "saved_searches":
			output = out
		}
	}
	return output, nil
}
