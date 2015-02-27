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
		newHelp(os.Stderr).Run()
	}
}

type help struct {
	logger    *Logger
	isCommand bool
	listing   string
}

func (h help) usage(k string, makes makeCommandHandler) {
	fs, _, _ := setupFlagset(k, makes, []SetupLogger{func(l *Logger) {
		l.Out = h.logger.Out // although `makes' must not use l.Out outside the Run
	}})
	// fs.Usage is ignored
	fs.VisitAll(func(f *flag.Flag) { // mimics fs.PrintDefaults
		format := "  -%s=%s: %s\n"
		if _, ok := f.Value.(flag.Getter).Get().(string); ok {
			format = "  -%s=%q: %s\n" // put quotes on the value
		}
		format = "   " + format
		h.logger.Printf(format, f.Name, f.DefValue, f.Usage)
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
			h.logger.Println("Usage of command:")
			h.logger.Printf("   %s\n", h.listing)
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
	h.logger.Println(fstline)
	for _, name := range commands.added.Names {
		h.logger.Printf("   %s\n", name)
		if makes, ok := commands.added.makes[name]; ok {
			h.usage(name, makes)
		}
	}
}

func newHelp(logout io.Writer, loggerOptions ...SetupLogger) *help {
	return &help{
		logger: NewLogger("", append([]SetupLogger{
			func(l *Logger) {
				l.Out = logout
				l.Flag = 0
			},
		}, loggerOptions...)...),
	}
}

func setupCommands(fs *flag.FlagSet, loggerOptions ...SetupLogger) (CommandHandler, io.Writer) {
	h := newHelp(os.Stdout, loggerOptions...)
	h.isCommand = true
	fs.StringVar(&h.listing, "h", "", "A command")
	return h.Run, h.logger
}

func init() {
	AddCommand("commands", setupCommands)
}
