package commands

import (
	"flag"
	"log"
	"sync"
)

type command func()

var (
	comutex  sync.Mutex
	commands = make(map[string]command)
)

func AddCommand(name string, fun func()) {
	comutex.Lock()
	defer comutex.Unlock()
	commands[name] = fun
}

func ArgCommand() command {
	if flag.NArg() == 0 {
		return nil
	}
	name := flag.Arg(0)
	enoent := func() {
		log.Fatalln("%s: No such command", name)
	}
	if name == "" {
		return enoent
	}
	comutex.Lock()
	defer comutex.Unlock()
	if fun, ok := commands[name]; ok {
		return fun
	}
	return enoent
}
