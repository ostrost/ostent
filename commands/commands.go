package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ostrost/ostent/commands/extpoints"
)

func setupFlagset(name string, cmd extpoints.Command, loggerSetups []extpoints.SetupLog) (*flag.FlagSet, extpoints.CommandHandler, io.Writer) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	run, output := cmd.SetupCommand(fs, loggerSetups...)
	return fs, run, output
}

func setup(name string, arguments []string, cmd extpoints.Command, loggerSetups []extpoints.SetupLog) (extpoints.CommandHandler, error, []string) {
	fs, run, output := setupFlagset(name, cmd, loggerSetups)
	err := fs.Parse(arguments)
	if err == nil && output != nil {
		fs.SetOutput(output)
	}
	return run, err, fs.Args()
}

func ParseCommand(handlers []extpoints.CommandHandler, args []string, loggerOptions ...extpoints.SetupLog) ([]extpoints.CommandHandler, error) {
	if len(args) == 0 || args[0] == "" {
		return handlers, nil
	}
	name := args[0]
	cmd := extpoints.Commands.Lookup(name)
	if cmd == nil {
		return handlers, fmt.Errorf("%s: No such command\n", name)
	}
	handler, err, nextargs := setup(name, args[1:], cmd, loggerOptions)
	if err != nil {
		return handlers, err
	}
	return ParseCommand(append(handlers, handler), nextargs)
}

func parseCommands() ([]extpoints.CommandHandler, error) {
	return ParseCommand([]extpoints.CommandHandler{}, flag.Args() /* no SetupLog passed */)
}

// true is when to abort
func ArgCommands() (bool, extpoints.AtexitHandler) {
	handlers, err := parseCommands()
	if err != nil {
		log.Fatal(err)
		return true, func() {} // useless return
	}

	finish := []extpoints.AtexitHandler{}
	atexit := func() {
		for _, exit := range finish {
			exit()
		}
	}

	if stop := func() bool {
		if len(CLIHandlers) > 0 {
			stop := false
			for _, clh := range CLIHandlers {
				if clh == nil {
					continue
				}
				atexit, term, err := clh()
				if err != nil {
					// the err must have been logged by clh
					stop = true
				} else if term {
					stop = true
				} else if atexit != nil {
					finish = append(finish, atexit)
				}
			}
			if stop {
				return true
			}
		}
		return false
	}(); stop {
		return true, atexit
	}

	if len(handlers) == 0 {
		return false, atexit
	}
	for _, handler := range handlers {
		handler()
	}
	return true, atexit
}

var CLIHandlers []extpoints.CommandLineHandler

func Parse(fs *flag.FlagSet, arguments []string) {
	for _, cli := range extpoints.CommandLines.All() {
		CLIHandlers = append(CLIHandlers, cli.SetupFlagSet(fs))
	}
	fs.Usage = NewHelp(os.Stderr).UsageFunc(fs)
	fs.Parse(arguments)
}

func NewLog(prefix string, options ...extpoints.SetupLog) *extpoints.Log {
	l := &extpoints.Log{ // defaults
		Out:  os.Stderr,
		Flag: log.LstdFlags,
	}
	for _, option := range options {
		option(l)
	}
	l.Logger = log.New(
		l.Out,
		prefix,
		l.Flag,
	)
	return l
}
