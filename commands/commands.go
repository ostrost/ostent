package commands

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

type atexitMaker interface {
	makeAtexitHandler() atexitHandler
}

type atexitHandler func()
type CommandHandler func()
type commandLineHandler func() (atexitHandler, bool, error)

type makeSub func(*flag.FlagSet, []string) (CommandHandler, error, []string)
type makeCommandHandler func(*flag.FlagSet) (CommandHandler, io.Writer)

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
		added []commandLineHandler
	}{}
)

func AddCommandLine(hfunc func(*flag.FlagSet) commandLineHandler) {
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

func setupFlagset(name string, makes makeCommandHandler) (*flag.FlagSet, CommandHandler, io.Writer) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	run, output := makes(fs)
	return fs, run, output
}

func setup(name string, makes makeCommandHandler, arguments []string) (CommandHandler, error, []string) {
	fs, run, output := setupFlagset(name, makes)
	err := fs.Parse(arguments)
	if err == nil && output != nil {
		fs.SetOutput(output)
	}
	return run, err, fs.Args()
}

func parseCommand(handlers []CommandHandler, args []string) ([]CommandHandler, bool) {
	if len(args) == 0 || args[0] == "" {
		return handlers, false
	}
	name := args[0]
	if ctor, ok := commands.added.makes[name]; ok {
		if handler, err, nextargs := setup(name, ctor, args[1:]); err == nil {
			return parseCommand(append(handlers, handler), nextargs)
		}
	} else {
		log.Fatalf("%s: No such command\n", name)
	}
	return handlers, true
}

func parseCommands() ([]CommandHandler, bool) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	return parseCommand([]CommandHandler{}, flag.Args())
}

// true is when to abort
func ArgCommands() (bool, atexitHandler) {
	finish := []atexitHandler{}
	atexit := func() {
		for _, exit := range finish {
			exit()
		}
	}

	handlers, errd := parseCommands()
	if errd {
		return true, atexit
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

func NewLogger(options ...func(*Logger)) *Logger {
	logger := &Logger{ // defaults
		Out:  os.Stderr,
		Flag: log.LstdFlags,
	}
	for _, option := range options {
		option(logger)
	}
	logger.Logger = log.New(
		logger.Out,
		logger.Prefix,
		logger.Flag,
	)
	return logger
}

type Logger struct { // also an io.Writer
	*log.Logger           // wrapping a *log.Logger
	Out         io.Writer // an argument for log.New
	Prefix      string    // an argument for log.New
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
