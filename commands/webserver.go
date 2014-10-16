package commands

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ostrost/ostent/types"
)

type webserver struct {
	logger       *loggerWriter
	BindValue    types.BindValue
	ServeFunc    func(net.Listener)
	FirstRunFunc func() bool
	ShutdownFunc func() bool
}

func (wr webserver) NetListen() net.Listener {
	listen, err := net.Listen("tcp", wr.BindValue.String())
	if err != nil {
		wr.logger.Fatal(err)
	}
	return listen
}

// LogInit sets up global log
func InitStdLog() {
	log.SetPrefix(fmt.Sprintf("[%d][ostent] ", os.Getpid()))
	// goagain logging is useless without pid prefix
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func NewWebserver() *webserver {
	return &webserver{
		logger: &loggerWriter{
			log.New(os.Stderr,
				fmt.Sprintf("[%d][ostent webserver] ", os.Getpid()),
				log.LstdFlags|log.Lmicroseconds),
		},
		BindValue: types.NewBindValue(":8050", "8050"),
	}
}

func (ws *webserver) AddCommandLine() *webserver {
	AddCommandLine(func(cli *flag.FlagSet) commandLineHandler {
		cli.Var(&ws.BindValue, "b", "short for bind")
		cli.Var(&ws.BindValue, "bind", "Bind address")
		return nil
	})
	return ws
}
