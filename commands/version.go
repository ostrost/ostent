package commands

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ostrost/ostent"
)

type version struct {
	logger *loggerWriter
}

func (v version) Run() {
	v.logger.Println(ostent.VERSION)
}

/*
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
// */

func newVersion() *version {
	return &version{
		logger: &loggerWriter{log.New(os.Stdout, "", 0)},
	}
}

func versionCommand(_ *flag.FlagSet) (sub, io.Writer) {
	v := newVersion()
	return v.Run, v.logger
}

func versionCommandLine(cli *flag.FlagSet) commandLineHandler {
	var fv bool
	cli.BoolVar(&fv, "v", false, "version")
	return func() bool {
		if fv {
			newVersion().Run()
			return true
		}
		return false
	}
}

func testCommandLine(cli *flag.FlagSet) commandLineHandler {
	var fv bool
	cli.BoolVar(&fv, "z", false, "z test")
	return func() bool {
		if fv {
			fmt.Println("Z")
		}
		return false
	}
}

func init() {
	AddFlaggedCommand("version", versionCommand)
	AddCommandLine(versionCommandLine)
	// AddCommandLine(testCommandLine)
}
