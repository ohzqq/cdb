package cdb

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Opt func(*Command)

type Flag func() []string

type Command struct {
	CdbCmd     string
	flags      []string
	positional []string
	tmp        *os.File
	verbose    bool
	dryRun     bool
}

type CalibredbCmd func() (string, []string)

func Calibredb(path string, opts []Opt, global []Flag, cdb CalibredbCmd, pos ...string) (*Command, error) {
	cmd := &Command{
		positional: pos,
	}

	err := checkLib(path)
	if err != nil {
		return cmd, err
	}

	cmd.flags = append(cmd.flags, "--with-library", path)

	for _, flag := range global {
		cmd.flags = append(cmd.flags, flag()...)
	}

	for _, fn := range opts {
		fn(cmd)
	}

	c, f := cdb()
	cmd.CdbCmd = c
	cmd.flags = append(cmd.flags, f...)

	return cmd, nil
}

func NewCommand(path string, args ...Opt) (*Command, error) {
	cmd := &Command{}

	err := checkLib(path)
	if err != nil {
		return cmd, err
	}

	cmd.flags = append(cmd.flags, "--with-library", path)

	for _, fn := range args {
		fn(cmd)
	}

	return cmd, nil
}

func checkLib(path string) error {
	uri, err := url.Parse(path)
	if err != nil {
		return err
	}

	if uri.Scheme != "" {
		if !SrvIsOnline(uri) {
			return fmt.Errorf("server is offline")
		}
		return nil
	}

	if ok := FileExist(filepath.Join(path, "metadata.db")); !ok {
		return ErrFileNotExist(path)
	}

	return nil
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

func SetFlags(flag ...string) Opt {
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

func DryRun() Opt {
	return func(cmd *Command) {
		cmd.dryRun = true
	}
}

func Username(name string) Flag {
	return func() []string {
		return []string{"--username", name}
	}
}

func Password(pass string) Flag {
	return func() []string {
		return []string{"--password", pass}
	}
}

func WithUsername(name string) Opt {
	return SetFlags("--username", name)
}

func WithPassword(pass string) Opt {
	return SetFlags("--password", pass)
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

	if c.dryRun {
		return cmd.String(), nil
	}

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

	if c.verbose {
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
		default:
			//case "search", "saved_searches":
			output = out
		}
	}
	return output, nil
}

func SrvIsOnline(u *url.URL) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
