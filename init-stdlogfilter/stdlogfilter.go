// Package stdlogfilter patches "log" package standard Logger so that
// any logging of lines containing " handling " via the logger are discarded.
// "github.com/rcrowley/go-tigertonic" is to blame.
package stdlogfilter

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func init() {
	log.SetOutput(newlogger())
}

func newlogger() io.Writer {
	lf := logFiltered{}
	var reader io.Reader
	reader, lf.writer = io.Pipe()
	lf.scanner = bufio.NewScanner(reader)
	go lf.read()
	return &lf
}

type logFiltered struct {
	writer  io.Writer
	scanner *bufio.Scanner
}

func (lf *logFiltered) Write(p []byte) (int, error) {
	return lf.writer.Write(p)
}

func (lf *logFiltered) read() {
	for {
		if !lf.scanner.Scan() {
			if err := lf.scanner.Err(); err != nil {
				log.New(os.Stderr, "", log.LstdFlags).Printf("bufio.Scanner.Scan Err: %s", err)
			}
			continue
		}
		text := lf.scanner.Text()
		if strings.Contains(text, " handling ") {
			continue
		}
		os.Stderr.WriteString(text + "\n")
	}
}
