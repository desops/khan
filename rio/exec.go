package rio

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/keegancsmith/shell"
)

type Cmd struct {
	Path string
	Args []string
	Env  [][2]string
	Dir  string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Host    string
	Context context.Context

	config *Config
}

func (config *Config) Command(ctx context.Context, path string, args ...string) *Cmd {
	return &Cmd{
		Path: path,
		Args: args,

		Context: ctx,

		config: config,
	}
}

func (cmd *Cmd) Run() error {
	if cmd.config.Pool == nil {
		c := exec.CommandContext(cmd.Context, cmd.Path, cmd.Args...)
		c.Stdin = cmd.Stdin
		c.Stdout = cmd.Stdout
		c.Stderr = cmd.Stderr
		if len(cmd.Env) > 0 {
			// For now, don't copy everything, just PATH
			c.Env = append(c.Env, "PATH="+os.Getenv("PATH"))
			for _, e := range cmd.Env {
				c.Env = append(c.Env, e[0]+"="+e[1])
			}
		}
		return c.Run()
	}

	session, err := cmd.config.Pool.Get(cmd.config.Host)
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdin = cmd.Stdin
	session.Stdout = cmd.Stdout
	session.Stderr = cmd.Stderr
	for _, e := range cmd.Env {
		if err := session.Setenv(e[0], e[1]); err != nil {
			return err
		}
	}

	cmdline := cmd.Path
	for _, a := range cmd.Args {
		cmdline += " " + shell.ReadableEscapeArg(a)
	}

	if cmd.config.Verbose {
		fmt.Println("ssh", cmd.config.Host, cmdline, cmd.Env)
	}

	err = session.Run(cmdline)

	if cmd.config.Verbose {
		fmt.Println("ssh", cmd.config.Host, cmdline, err)
	}

	return err
}
