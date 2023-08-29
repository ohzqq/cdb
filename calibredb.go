package cdb

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CmdOpt func(*Command)

type Command struct {
	CdbCmd     string
	flags      []string
	positional []string
	verbose    bool
	dryRun     bool
}

func Calibredb(path string, opts ...CmdOpt) (*Command, error) {
	cmd := &Command{}

	err := checkLib(path)
	if err != nil {
		return cmd, err
	}

	cmd.flags = append(cmd.flags, "--with-library", path)

	for _, fn := range opts {
		fn(cmd)
	}

	return cmd, nil
}

func Verbose() CmdOpt {
	return func(cmd *Command) {
		cmd.verbose = true
	}
}

func DryRun() CmdOpt {
	return func(cmd *Command) {
		cmd.dryRun = true
	}
}

func (c *Command) Authenticate(user, pass string) *Command {
	c.SetFlags("--username", user)
	c.SetFlags("--password", pass)
	return c
}

func (c *Command) SetFlags(flags ...string) {
	c.flags = append(c.flags, flags...)
}

func (c *Command) SetArgs(args ...string) {
	c.positional = append(c.positional, args...)
}

func (c *Command) Build() *exec.Cmd {
	return exec.Command("calibredb", c.parseArgs()...)
}

func (c *Command) parseArgs() []string {
	var args []string
	args = append(args, c.CdbCmd)
	args = append(args, c.flags...)
	args = append(args, c.positional...)
	return args
}

func (c *Command) Run() (string, error) {
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
