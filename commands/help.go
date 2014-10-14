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
	logger    *loggerWriter
	isCommand bool
	listing   string
}

func (h help) usage(k string, sfunc setupFunc) {
	fs, _, _ := setupFlagset(k, sfunc)
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
		for _, x := range commands.mapsub.keys {
			if x == h.listing {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("%s: No such command\n", h.listing)
		} else {
			h.logger.Println("Usage of command:")
			h.logger.Printf("   %s\n", h.listing)
			if sfunc, ok := commands.mapsub.setups[h.listing]; ok {
				h.usage(h.listing, sfunc)
			}
		}
		return
	}
	sort.Stable(commands.mapsub)
	fstline := "Commands available:"
	if !h.isCommand {
		fstline = fmt.Sprintf("Commands of %s:", os.Args[0]) // as in usage
	}
	h.logger.Println(fstline)
	for _, k := range commands.mapsub.keys {
		h.logger.Printf("   %s\n", k)
		if sfunc, ok := commands.mapsub.setups[k]; ok {
			h.usage(k, sfunc)
		}
	}
}

func newHelp(logout io.Writer) *help {
	return &help{
		logger: &loggerWriter{log.New(logout, "", 0)},
	}
}

func setupCommands(fs *flag.FlagSet) (sub, io.Writer) {
	h := newHelp(os.Stdout)
	h.isCommand = true
	fs.StringVar(&h.listing, "h", "", "A command")
	return h.Run, h.logger
}

func init() {
	AddFlaggedCommand("commands", setupCommands)
	// AddCommand("help", helpCommand)
}
