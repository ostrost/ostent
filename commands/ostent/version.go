package ostent

import (
	"flag"
	"io"
	"os"

	"github.com/ostrost/ostent/commands"
	"github.com/ostrost/ostent/ostent"
)

type version struct {
	logger *commands.Logger
}

func (v version) run() {
	v.logger.Println(ostent.VERSION)
}

func newVersion(loggerOptions ...commands.SetupLogger) *version {
	return &version{
		logger: commands.NewLogger("", append([]commands.SetupLogger{
			func(l *commands.Logger) {
				l.Out = os.Stdout
				l.Flag = 0
			},
		}, loggerOptions...)...),
	}
}

func versionCommand(_ *flag.FlagSet, loggerOptions ...commands.SetupLogger) (commands.CommandHandler, io.Writer) {
	v := newVersion(loggerOptions...)
	return v.run, v.logger
}

func versionCommandLine(cli *flag.FlagSet) commands.CommandLineHandler {
	var fv bool
	cli.BoolVar(&fv, "v", false, "version")
	return func() (commands.AtexitHandler, bool, error) {
		if fv {
			newVersion().run()
			return nil, true, nil
		}
		return nil, false, nil
	}
}

func init() {
	commands.AddCommand("version", versionCommand)
	commands.AddCommandLine(versionCommandLine)
}
