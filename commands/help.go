package commands

import (
	"flag"
	"log"
	"os"
)

type help struct {
	logger *loggerWriter
	Flag   bool
}

func (h help) Run() {
	if !h.Flag {
		return
	}
	commands.mutex.Lock()
	defer commands.mutex.Unlock()
	h.logger.Println("Commands available:")
	for k := range commands.mapsub {
		h.logger.Printf("\t%s\n", k)
	}
}

func FlagSetNewHelp(fs *flag.FlagSet) *help {
	h := help{
		logger: &loggerWriter{log.New(os.Stdout, "", 0)},
	}
	// fs.BoolVar(&v.Flag, "h", false, "help")
	return &h
}

func helpCommand(fs *flag.FlagSet, arguments []string) (sub, error, []string) {
	h := FlagSetNewHelp(fs)
	h.Flag = true
	fs.SetOutput(h.logger)
	err := fs.Parse(arguments)
	return h.Run, err, fs.Args()
}

func init() {
	AddCommand("help", helpCommand)
}
