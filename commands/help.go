package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/ostrost/ostent/commands/extpoints"
)

func UsageFunc(fs *flag.FlagSet) func() {
	return func() {
		// default usage
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults() // flag.PrintDefaults()
		// continued usage:
		newHelp(os.Stderr).Run()
	}
}

type help struct {
	Log       *extpoints.Log
	isCommand bool
	listing   string
}

func (h help) usage(k string, makes makeCommandHandler) {
	fs, _, _ := setupFlagset(k, makes, []extpoints.SetupLog{func(l *extpoints.Log) {
		l.Out = h.Log.Out // although `makes' must not use l.Out outside the Run
	}})
	// fs.Usage is ignored
	fs.VisitAll(func(f *flag.Flag) { // mimics fs.PrintDefaults
		format := "  -%s=%s: %s\n"
		if _, ok := f.Value.(flag.Getter).Get().(string); ok {
			format = "  -%s=%q: %s\n" // put quotes on the value
		}
		format = "   " + format
		h.Log.Printf(format, f.Name, f.DefValue, f.Usage)
	})
}

func (h *help) Run() {
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	if h.listing != "" {
		found := false
		for _, name := range commands.added.Names {
			if name == h.listing {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("%s: No such command\n", h.listing)
		} else {
			h.Log.Println("Usage of command:")
			h.Log.Printf("   %s\n", h.listing)
			if makes, ok := commands.added.makes[h.listing]; ok {
				h.usage(h.listing, makes)
			}
		}
		return
	}
	sort.Strings(commands.added.Names)
	fstline := "Commands available:"
	if !h.isCommand {
		fstline = fmt.Sprintf("Commands of %s:", os.Args[0]) // as in usage
	}
	h.Log.Println(fstline)
	for _, name := range commands.added.Names {
		h.Log.Printf("   %s\n", name)
		if makes, ok := commands.added.makes[name]; ok {
			h.usage(name, makes)
		}
	}
}

func newHelp(logout io.Writer, loggerOptions ...extpoints.SetupLog) *help {
	return &help{
		Log: NewLog("", append([]extpoints.SetupLog{
			func(l *extpoints.Log) {
				l.Out = logout
				l.Flag = 0
			},
		}, loggerOptions...)...),
	}
}

func setupCommands(fs *flag.FlagSet, loggerOptions ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer) {
	h := newHelp(os.Stdout, loggerOptions...)
	h.isCommand = true
	fs.StringVar(&h.listing, "h", "", "A command")
	return h.Run, h.Log
}

func init() {
	AddCommand("commands", setupCommands)
}
