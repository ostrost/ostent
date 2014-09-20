package commands

import (
	"flag"
	"log"
	"sync"
)

type runnable interface {
	Run()
}

type commandMaker func(*flag.FlagSet, []string) (runnable, error, []string)

var (
	comutex  sync.Mutex
	COMMANDS = make(map[string]commandMaker)
)

func AddCommand(name string, fun commandMaker) {
	comutex.Lock()
	defer comutex.Unlock()
	COMMANDS[name] = fun
}

func parseCommand(runs []runnable, args []string) ([]runnable, bool) {
	if len(args) == 0 || args[0] == "" {
		return runs, false
	}
	name := args[0]
	if ctor, ok := COMMANDS[name]; ok {
		fs := flag.NewFlagSet(name, flag.ContinueOnError)
		if run, err, nextargs := ctor(fs, args[1:]); err == nil {
			return parseCommand(append(runs, run), nextargs)
		} // else { /* log.Printf("%s: %s\n", name, err) // printed already by flag package // */ }
	} else {
		log.Fatalf("%s: No such command\n", name)
	}
	return runs, true
}

func parseCommands() ([]runnable, bool) {
	comutex.Lock()
	defer comutex.Unlock()
	return parseCommand([]runnable{}, flag.Args())
}

// true is when to abort
func ArgCommands() bool {
	runs, errd := parseCommands()
	if errd {
		return true
	}
	if len(runs) == 0 {
		return false
	}
	for _, cmd := range runs {
		cmd.Run()
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
