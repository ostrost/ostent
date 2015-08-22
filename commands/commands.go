package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/ostrost/ostent/commands/extpoints"
)

type atexitMaker interface {
	makeAtexitHandler() extpoints.AtexitHandler
}

type makeSub func(*flag.FlagSet, []string) (extpoints.CommandHandler, error, []string)
type makeCommandHandler func(*flag.FlagSet, ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer)

type addedCommands struct {
	makes map[string]makeCommandHandler
	Names []string
}

var (
	commands = struct {
		mutex sync.Mutex
		added addedCommands
	}{
		added: addedCommands{
			makes: make(map[string]makeCommandHandler),
		},
	}
	commandline = struct {
		mutex sync.Mutex
		added []extpoints.CommandLineHandler
	}{}
)

func AddCommandLine(hfunc func(*flag.FlagSet) extpoints.CommandLineHandler) {
	s := hfunc(flag.CommandLine) // NB global flag.CommandLine
	if s == nil {
		return
	}
	commandline.mutex.Lock()
	defer commandline.mutex.Unlock()
	commandline.added = append(commandline.added, s)
}

func AddCommand(name string, makes makeCommandHandler) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	commands.added.makes[name] = makes
	commands.added.Names = append(commands.added.Names, name)
}

func setupFlagset(name string, makes makeCommandHandler, loggerSetups []extpoints.SetupLog) (*flag.FlagSet, extpoints.CommandHandler, io.Writer) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	run, output := makes(fs, loggerSetups...)
	return fs, run, output
}

func setup(name string, arguments []string, makes makeCommandHandler, loggerSetups []extpoints.SetupLog) (extpoints.CommandHandler, error, []string) {
	fs, run, output := setupFlagset(name, makes, loggerSetups)
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
	ctor, ok := commands.added.makes[name]
	if !ok {
		return handlers, fmt.Errorf("%s: No such command\n", name)
	}
	handler, err, nextargs := setup(name, args[1:], ctor, loggerOptions)
	if err != nil {
		return handlers, err
	}
	return ParseCommand(append(handlers, handler), nextargs)
}

func parseCommands() ([]extpoints.CommandHandler, error) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
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
		commandline.mutex.Lock()
		defer commandline.mutex.Unlock()

		if len(commandline.added) > 0 {
			stop := false
			for _, clh := range commandline.added {
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
