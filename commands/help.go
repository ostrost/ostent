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

func (h *help) UsageFunc(fs *flag.FlagSet) func() {
	return func() {
		// default usage
		fmt.Fprintf(h.Output, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults() // flag.PrintDefaults()
		// continued usage:
		h.Run()
	}
}

type help struct {
	Output    io.Writer
	Log       *extpoints.Log
	isCommand bool
	listing   string
}

func (h help) usage(k string, cmd extpoints.Command) {
	fs, _, _ := setupFlagset(k, cmd, []extpoints.SetupLog{func(l *extpoints.Log) {
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
	if h.listing != "" {
		cmd := extpoints.Commands.Lookup(h.listing)
		if cmd == nil {
			log.Fatalf("%s: No such command\n", h.listing)
			return
		}
		h.Log.Println("Usage of command:")
		h.Log.Printf("   %s\n", h.listing)
		h.usage(h.listing, cmd)
		return
	}
	fstline := "Commands available:"
	if !h.isCommand {
		fstline = fmt.Sprintf("Commands of %s:", os.Args[0]) // as in usage
	}
	h.Log.Println(fstline)
	names := extpoints.Commands.Names()
	sort.Strings(names)
	for _, name := range names {
		h.Log.Printf("   %s\n", name)
		h.usage(name, extpoints.Commands.Lookup(name))
	}
}

func NewHelp(logout io.Writer, loggerOptions ...extpoints.SetupLog) *help {
	return &help{
		Output: logout,
		Log: NewLog("", append([]extpoints.SetupLog{
			func(l *extpoints.Log) {
				l.Out = logout
				l.Flag = 0
			},
		}, loggerOptions...)...),
	}
}

func (_ Helps) SetupCommand(fs *flag.FlagSet, loggerOptions ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer) {
	h := NewHelp(os.Stdout, loggerOptions...)
	h.isCommand = true
	fs.StringVar(&h.listing, "h", "", "A command")
	return h.Run, h.Log
}

type Helps struct{}

func init() {
	extpoints.Commands.Register(Helps{}, "help")
}
