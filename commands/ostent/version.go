package ostent

import (
	"flag"
	"io"
	"os"

	"github.com/ostrost/ostent/commands"
	"github.com/ostrost/ostent/ostent"
)

type Version struct {
	Log *commands.Logger
}

func (v Version) Run() {
	v.Log.Println(ostent.VERSION)
}

func NewVersion(logOptions ...commands.SetupLogger) *Version {
	return &Version{
		Log: commands.NewLogger("", append([]commands.SetupLogger{
			func(l *commands.Logger) {
				l.Out = os.Stdout
				l.Flag = 0
			},
		}, logOptions...)...),
	}
}

func VersionCommand(_ *flag.FlagSet, logOptions ...commands.SetupLogger) (commands.CommandHandler, io.Writer) {
	v := NewVersion(logOptions...)
	return v.Run, v.Log
}

func VersionCLI(cli *flag.FlagSet) commands.CommandLineHandler {
	var fv bool
	cli.BoolVar(&fv, "v", false, "Short for version")
	cli.BoolVar(&fv, "version", false, "Print version and exit")
	return func() (commands.AtexitHandler, bool, error) {
		if fv {
			NewVersion().Run()
			return nil, true, nil
		}
		return nil, false, nil
	}
}

func init() {
	commands.AddCommand("version", VersionCommand)
	commands.AddCommandLine(VersionCLI)
}
