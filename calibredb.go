package cdb

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"path/filepath"
	"time"
)

// CmdOpt sets options for the calibredb command.
type CmdOpt func(*Command)

// Command represents a calibredb command.
type Command struct {
	CdbCmd     string
	flags      []string
	positional []string
	verbose    bool
	dryRun     bool
}

// Calibredb initializes a command. Path or URL to a library is required.
func Calibredb(path string, opts ...CmdOpt) (*Command, error) {
	cmd := &Command{}

	err := CheckLib(path)
	if err != nil {
		return cmd, err
	}

	cmd.flags = append(cmd.flags, "--with-library", path)

	for _, fn := range opts {
		fn(cmd)
	}

	return cmd, nil
}

// Verbose prints the built command and stdout.
func Verbose() CmdOpt {
	return func(cmd *Command) {
		cmd.verbose = true
	}
}

// DryRun prints the built command.
func DryRun() CmdOpt {
	return func(cmd *Command) {
		cmd.dryRun = true
	}
}

// Authenticate is used to pass credentials to the calibre content server.
func Authenticate(user, pass string) CmdOpt {
	return func(c *Command) {
		c.SetFlags("--username", user)
		c.SetFlags("--password", pass)
	}
}

// SetFlags sets the flags for the calibredb command.
func (c *Command) SetFlags(flags ...string) {
	c.flags = append(c.flags, flags...)
}

// SetArgs sets the positional arguments for the calibredb command.
func (c *Command) SetArgs(args ...string) {
	c.positional = append(c.positional, args...)
}

// Build compiles the calibredb command to be run.
func (c *Command) Build() *exec.Cmd {
	var args []string
	args = append(args, c.CdbCmd)
	args = append(args, c.flags...)
	args = append(args, c.positional...)
	return exec.Command("calibredb", args...)
}

// Run executes the calibredb command.
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
		//case "add":
		//sp := strings.Split(out, ": ")
		//output = sp[1]
		default:
			output = out
		}
	}
	if len(stderr.Bytes()) > 0 {
		output += "\n"
		output += stderr.String()
	}
	return output, nil
}

// SrvIsOnline tests if the calibredb content server is available.
func SrvIsOnline(u *url.URL) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

// CheckLib checks that the calibre database exists or the content server is
// online.
func CheckLib(path string) error {
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
