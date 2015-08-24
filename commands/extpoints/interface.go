//go:generate go-extpoints .

package extpoints

import (
	"flag"
	"io"
)

type CommandLine interface {
	SetupFlagSet(*flag.FlagSet) CommandLineHandler
}

type Command interface {
	SetupCommand(*flag.FlagSet, ...SetupLog) (CommandHandler, io.Writer)
}
