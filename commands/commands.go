package commands

import (
	"flag"
	"io"
	"log"
	"sync"
)

type atexitMaker interface {
	makeAtexitHandler() atexitHandler
}

type atexitHandler func()
type commandHandler func()
type commandLineHandler func() bool

type makeSub func(*flag.FlagSet, []string) (commandHandler, error, []string)
type makeCommandHandler func(*flag.FlagSet) (commandHandler, io.Writer)

type addedCommands struct {
	submap              map[string]makeSub
	setups              map[string]makeCommandHandler
	commandLineHandlers []commandLineHandler
	names               []string
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
			submap: make(map[string]makeSub),
			setups: make(map[string]makeCommandHandler),
		},
	}

	defaults = struct {
		mutex sync.Mutex
		added map[string]atexitMaker
	}{
		added: make(map[string]atexitMaker),
	}
)

func AddCommandLine(hfunc func(*flag.FlagSet) commandLineHandler) {
	s := hfunc(flag.CommandLine) // NB global flag.CommandLine
	if s == nil {
		return
	}
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	commands.added.commandLineHandlers = append(
		commands.added.commandLineHandlers, s)
}

func AddFlaggedCommand(name string, makes makeCommandHandler) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	commands.added.setups[name] = makes
	commands.added.names = append(commands.added.names, name)
}

func setupFlagset(name string, makes makeCommandHandler) (*flag.FlagSet, commandHandler, io.Writer) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	run, output := makes(fs)
	return fs, run, output
}

func setup(name string, makes makeCommandHandler, arguments []string) (commandHandler, error, []string) {
	fs, run, output := setupFlagset(name, makes)
	err := fs.Parse(arguments)
	if err == nil && output != nil {
		fs.SetOutput(output)
	}
	return run, err, fs.Args()
}

func AddCommand(name string, makes makeSub) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	commands.added.submap[name] = makes
	commands.added.names = append(commands.added.names, name)
}

func AddDefault(name string, def atexitMaker) {
	defaults.mutex.Lock()
	defer defaults.mutex.Unlock()
	defaults.added[name] = def
}

func Defaults() atexitHandler {
	defaults.mutex.Lock()
	defer defaults.mutex.Unlock()
	finish := []atexitHandler{}
	for _, maker := range defaults.added {
		if fin := maker.makeAtexitHandler(); fin != nil {
			finish = append(finish, fin)
		}
	}
	return func() {
		for _, fin := range finish {
			fin()
		}
	}
}

func parseCommand(handlers []commandHandler, args []string) ([]commandHandler, bool) {
	if len(args) == 0 || args[0] == "" {
		return handlers, false
	}
	name := args[0]
	if ctor, ok := commands.added.setups[name]; ok {
		if handler, err, nextargs := setup(name, ctor, args[1:]); err == nil {
			return parseCommand(append(handlers, handler), nextargs)
		}
	} else if ctor, ok := commands.added.submap[name]; ok {
		fs := flag.NewFlagSet(name, flag.ContinueOnError)
		if handler, err, nextargs := ctor(fs, args[1:]); err == nil {
			return parseCommand(append(handlers, handler), nextargs)
		}
		// else { /* log.Printf("%s: %s\n", name, err)
		// printed already by flag package // */ }
	} else {
		log.Fatalf("%s: No such command\n", name)
	}
	return handlers, true
}

func parseCommands() ([]commandHandler, bool) {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	return parseCommand([]commandHandler{}, flag.Args())
}

// true is when to abort
func ArgCommands() bool {
	handlers, errd := parseCommands()
	if errd {
		return true
	}
	if stop := func() bool {
		commands.mutex.Lock()
		defer commands.mutex.Unlock()

		if len(commands.added.commandLineHandlers) > 0 {
			stop := false
			for _, clh := range commands.added.commandLineHandlers {
				if clh() {
					stop = true
				}
			}
			if stop {
				return true
			}
		}
		return false
	}(); stop {
		return true
	}

	if len(handlers) == 0 {
		return false
	}
	for _, handler := range handlers {
		handler()
	}
	return true
}

type loggerWriter struct {
	*log.Logger
}

func (lw *loggerWriter) Write(p []byte) (int, error) {
	lw.Logger.Printf("%s", p)
	return len(p), nil
}
