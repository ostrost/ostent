package commands

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ostrost/ostent"
)

type version struct {
	logger *Logger
}

func (v version) run() {
	v.logger.Println(ostent.VERSION)
}

func newVersion() *version {
	return &version{
		logger: NewLogger(func(l *Logger) {
			l.Out = os.Stdout
			l.Flag = 0
		}),
	}
}

func versionCommand(_ *flag.FlagSet) (CommandHandler, io.Writer) {
	v := newVersion()
	return v.run, v.logger
}

func versionCommandLine(cli *flag.FlagSet) commandLineHandler {
	var fv bool
	cli.BoolVar(&fv, "v", false, "version")
	return func() (atexitHandler, bool, error) {
		if fv {
			newVersion().run()
			return nil, true, nil
		}
		return nil, false, nil
	}
}

func testCommandLine(cli *flag.FlagSet) commandLineHandler {
	var fv bool
	cli.BoolVar(&fv, "z", false, "z test")
	return func() (atexitHandler, bool, error) {
		if fv {
			fmt.Println("Z")
		}
		return nil, false, nil
	}
}

func init() {
	AddCommand("version", versionCommand)
	AddCommandLine(versionCommandLine)
	// AddCommandLine(testCommandLine)
}
