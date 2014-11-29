package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type atexitMaker interface {
	makeAtexitHandler() AtexitHandler
}

type AtexitHandler func()
type CommandHandler func()
type CommandLineHandler func() (AtexitHandler, bool, error)

type makeSub func(*flag.FlagSet, []string) (CommandHandler, error, []string)
type makeCommandHandler func(*flag.FlagSet, ...SetupLogger) (CommandHandler, io.Writer)

type addedCommands struct {
	makes map[string]makeCommandHandler
	names []string
}

// conforms to sort.Interface
func (ac addedCommands) Len() int {
	return len(ac.names)
}
func (ac addedCommands) Swap(i, j int) {
	ac.names[i], ac.names[j] = ac.names[j], ac.names[i]
}
func (ac addedCommands) Less(i, j int) bool {
	return ac.names[i] < ac.names[j]
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
		added []CommandLineHandler
	}{}
)

func AddCommandLine(hfunc func(*flag.FlagSet) CommandLineHandler) {
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
	commands.added.names = append(commands.added.names, name)
}

func setupFlagset(name string, makes makeCommandHandler, loggerSetups []SetupLogger) (*flag.FlagSet, CommandHandler, io.Writer) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	run, output := makes(fs, loggerSetups...)
	return fs, run, output
}

func setup(name string, arguments []string, makes makeCommandHandler, loggerSetups []SetupLogger) (CommandHandler, error, []string) {
	fs, run, output := setupFlagset(name, makes, loggerSetups)
	err := fs.Parse(arguments)
	if err == nil && output != nil {
		fs.SetOutput(output)
	}
	return run, err, fs.Args()
}

func ParseCommand(handlers []CommandHandler, args []string, loggerOptions ...SetupLogger) ([]CommandHandler, error) {
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

func parseCommands() ([]CommandHandler, error) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	return ParseCommand([]CommandHandler{}, flag.Args() /* no SetupLogger passed */)
}

// true is when to abort
func ArgCommands() (bool, AtexitHandler) {
	handlers, err := parseCommands()
	if err != nil {
		log.Fatal(err)
		return true, func() {} // useless return
	}

	finish := []AtexitHandler{}
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

type SetupLogger func(*Logger)

func NewLogger(prefix string, options ...SetupLogger) *Logger {
	logger := &Logger{ // defaults
		Out:  os.Stderr,
		Flag: log.LstdFlags,
	}
	for _, option := range options {
		option(logger)
	}
	logger.Logger = log.New(
		logger.Out,
		prefix,
		logger.Flag,
	)
	return logger
}

type Logger struct { // also an io.Writer
	*log.Logger           // wrapping a *log.Logger
	Out         io.Writer // an argument for log.New
	Flag        int       // an argument for log.New
}

// satisfying io.Writer interface
func (l *Logger) Write(p []byte) (int, error) {
	l.Logger.Printf("%s", p)
	return len(p), nil
}

func (l *Logger) fatalif(err error) {
	if err != nil {
		l.Fatalln(err)
	}
}
