package commands

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/ostrost/ostent/types"
	"github.com/rcrowley/goagain"
)

type webserver struct {
	logger       *loggerWriter
	BindValue    types.BindValue
	ServeFunc    func(net.Listener)
	FirstRunFunc func() bool
	ShutdownFunc func() bool
}

func (_ webserver) GoneAgain() bool {
	return os.Getenv("GOAGAIN_PPID") != ""
}

func (_ webserver) GoAgain() {
	syscall.Kill(os.Getpid(), syscall.SIGUSR2)
}

func (wr webserver) NetListen() net.Listener {
	listen, err := net.Listen("tcp", wr.BindValue.String())
	if err != nil {
		wr.logger.Fatal(err)
	}
	return listen
}

func (wr webserver) Run() {
	listen, err := goagain.Listener()
	if err != nil {
		listen = wr.NetListen()

		if wr.FirstRunFunc != nil && wr.FirstRunFunc() { // had upgrade
			go func() { // delayed kill
				time.Sleep(time.Second) // not before goagain.Wait
				wr.GoAgain()
				// goagain.ForkExec(listen)
			}()
		} else if wr.ServeFunc != nil {
			wr.ServeFunc(listen)
		}
	} else {
		if wr.ServeFunc != nil {
			wr.ServeFunc(listen)
		}

		if err := goagain.Kill(); err != nil {
			wr.logger.Fatalln(err)
		}
	}

	if _, err := goagain.Wait(listen); err != nil { // signals before won't be catched
		wr.logger.Fatalln(err)
	}

	// shutting down

	if wr.ShutdownFunc != nil && wr.ShutdownFunc() {
		time.Sleep(time.Second) // wait for an affect
	}

	if err := listen.Close(); err != nil {
		wr.logger.Fatalln(err)
	}
	time.Sleep(time.Second)
}

// LogInit sets up global log
func InitStdLog() {
	log.SetPrefix(fmt.Sprintf("[%d][ostent] ", os.Getpid()))
	// goagain logging is useless without pid prefix
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func FlagSetNewWebserver(fs *flag.FlagSet) *webserver {
	logger := log.New(os.Stderr,
		fmt.Sprintf("[%d][ostent webserver] ", os.Getpid()),
		log.LstdFlags|log.Lmicroseconds)
	wr := webserver{
		logger:    &loggerWriter{logger},
		BindValue: types.NewBindValue(":8050", "8050"),
	}
	fs.Var(&wr.BindValue, "b", "short for bind")
	fs.Var(&wr.BindValue, "bind", "Bind address")
	return &wr
}

var _ = /* webserverCommand */ func(fs *flag.FlagSet, arguments []string) (sub, error, []string) {
	wr := FlagSetNewWebserver(fs)
	fs.SetOutput(wr.logger)
	err := fs.Parse(arguments)
	return wr.Run, err, fs.Args()
}

func init() {
	/* // "webserver" is not a cli command, at least for now
	AddCommand("webserver", webserverCommand) // */
}
