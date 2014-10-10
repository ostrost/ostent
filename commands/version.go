package commands

import (
	"flag"
	"log"
	"os"

	"github.com/ostrost/ostent"
)

type version struct {
	logger *loggerWriter
	Flag   bool
}

func (v version) Run() {
	if v.Flag {
		v.logger.Println(ostent.VERSION)
	}
}

func FlagSetNewVersion(fs *flag.FlagSet) *version {
	v := version{
		logger: &loggerWriter{log.New(os.Stdout, "", 0)},
	}
	fs.BoolVar(&v.Flag, "v", false, "version")
	return &v
}

func versionCommand(fs *flag.FlagSet, arguments []string) (sub, error, []string) {
	v := FlagSetNewVersion(fs)
	v.Flag = true
	fs.SetOutput(v.logger)
	err := fs.Parse(arguments)
	return v.Run, err, fs.Args()
}

func init() {
	AddCommand("version", versionCommand)
}
