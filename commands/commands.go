package commands

import (
	"flag"
	"log"
	"sync"
)

type sub func()

type makeSub func(*flag.FlagSet, []string) (sub, error, []string)

var (
	comutex  sync.Mutex
	COMMANDS = make(map[string]makeSub)
)

func AddCommand(name string, makes makeSub) {
	comutex.Lock()
	defer comutex.Unlock()
	COMMANDS[name] = makes
}

func parseCommand(subs []sub, args []string) ([]sub, bool) {
	if len(args) == 0 || args[0] == "" {
		return subs, false
	}
	name := args[0]
	if ctor, ok := COMMANDS[name]; ok {
		fs := flag.NewFlagSet(name, flag.ContinueOnError)
		if sub, err, nextargs := ctor(fs, args[1:]); err == nil {
			return parseCommand(append(subs, sub), nextargs)
		} // else { /* log.Printf("%s: %s\n", name, err) // printed already by flag package // */ }
	} else {
		log.Fatalf("%s: No such command\n", name)
	}
	return subs, true
}

func parseCommands() ([]sub, bool) {
	comutex.Lock()
	defer comutex.Unlock()
	return parseCommand([]sub{}, flag.Args())
}

// true is when to abort
func ArgCommands() bool {
	subs, errd := parseCommands()
	if errd {
		return true
	}
	if len(subs) == 0 {
		return false
	}
	for _, sub := range subs {
		sub()
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
