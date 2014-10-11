package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func UsageFunc(fs *flag.FlagSet) func() {
	return func() {
		// default usage
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults() // flag.PrintDefaults()
		// continued usage:
		flagSetNewHelp(fs, os.Stderr).Run()
	}
}

type help struct {
	logger    *loggerWriter
	isCommand bool
}

func (h help) Run() {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	sort.Stable(commands.mapsub)
	fstline := "Commands available:"
	if !h.isCommand {
		fstline = fmt.Sprintf("Commands of %s:", os.Args[0]) // as in usage
	}
	h.logger.Println(fstline)
	for _, k := range commands.mapsub.keys {
		h.logger.Printf("   %s\n", k)
	}
}

func flagSetNewHelp(fs *flag.FlagSet, logout io.Writer) *help {
	return &help{
		logger: &loggerWriter{log.New(logout, "", 0)},
	}
}

func helpCommand(fs *flag.FlagSet, arguments []string) (sub, error, []string) {
	h := flagSetNewHelp(fs, os.Stdout)
	h.isCommand = true
	fs.SetOutput(h.logger)
	err := fs.Parse(arguments)
	return h.Run, err, fs.Args()
}

func init() {
	AddCommand("help", helpCommand)
}
