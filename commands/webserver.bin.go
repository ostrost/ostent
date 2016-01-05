// +build bin

package commands

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/rcrowley/goagain"
)

// InitStdLog sets up global log.
func InitStdLog() {
	log.SetPrefix(fmt.Sprintf("[%d][ostent] ", os.Getpid()))
	// goagain logging is useless without pid prefix
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func (_ webserver) GoneAgain() bool {
	return os.Getenv("GOAGAIN_PPID") != ""
}

func (_ webserver) GoAgain() {
	syscall.Kill(os.Getpid(), syscall.SIGUSR2)
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
			wr.Log.Fatalln(err)
		}
	}

	if _, err := goagain.Wait(listen); err != nil { // signals before won't be catched
		wr.Log.Fatalln(err)
	}

	// shutting down

	if wr.ShutdownFunc != nil && wr.ShutdownFunc() {
		time.Sleep(time.Second) // wait for an affect
	}

	if err := listen.Close(); err != nil {
		wr.Log.Fatalln(err)
	}
	time.Sleep(time.Second)
}
