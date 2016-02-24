// +build bin

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rcrowley/goagain"
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/cmd"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/templates"
)

var (
	// BootTime is the boot time.
	BootTime = time.Now()
	// AssetAltModTimeFunc returns BootTime to be asset ModTime.
	AssetAltModTimeFunc = func() time.Time { return BootTime }
	// AgainLog is for goagain logging.
	AgainLog = log.New(os.Stderr, fmt.Sprintf("[%d][ostent webserver] ", os.Getpid()),
		log.LstdFlags|log.Lmicroseconds)
)

func init() {
	log.SetPrefix(fmt.Sprintf("[%d][ostent] ", os.Getpid()))
	// goagain logging is useless without pid prefix
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// OstentRunE is the run of the ostent command.
func OstentRunE(*cobra.Command, []string) error {
	ostent.RunBackground()

	listen, err := goagain.Listener()
	goneagain := err == nil
	if !goneagain {
		listen, err = net.Listen("tcp", cmd.OstentBind.String())
		if err != nil {
			return err
		}
	}

	go func() {
		templates.InitTemplates(nil) // preventive
		// sequential: Serve must wait for InitTemplates
		Serve(listen, true, nil) // true stands for taggedbin
	}()

	if goneagain {
		if err := goagain.Kill(); err != nil {
			AgainLog.Fatalln(err)
		}
	}

	if _, err := goagain.Wait(listen); err != nil { // signals before won't be catched
		AgainLog.Fatalln(err)
	}

	// shutting down

	if ostent.Connections.Reload() {
		time.Sleep(time.Second) // wait for an affect
	}

	if err := listen.Close(); err != nil {
		AgainLog.Fatalln(err)
	}
	time.Sleep(time.Second)
	return nil
}

func main() {
	cmd.OstentCmd.RunE = OstentRunE
	cmd.Execute()
}
