// +build production

package commands

import (
	"os"
	"syscall"
	"time"

	"github.com/rcrowley/goagain"
)

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
