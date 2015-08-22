package extpoints

import (
	"io"
	"log"
)

type AtexitHandler func()
type CommandHandler func()
type CommandLineHandler func() (AtexitHandler, bool, error)

type SetupLog func(*Log)

type Log struct { // also an io.Writer
	*log.Logger           // wrapping a *log.Logger
	Out         io.Writer // an argument for log.New
	Flag        int       // an argument for log.New
}

// satisfying io.Writer interface
func (l *Log) Write(p []byte) (int, error) {
	l.Logger.Printf("%s", p)
	return len(p), nil
}
