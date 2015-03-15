// Package stdlogfilter patches "log" package standard Logger so that
// any logging of lines containing " handling " via the logger are discarded.
// "github.com/rcrowley/go-tigertonic" is to blame.
package stdlogfilter

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func init() {
	log.SetOutput(NewLogFiltered(os.Stderr))
}

func NewLogFiltered(out io.Writer) io.Writer {
	reader, writer := io.Pipe()
	lf := LogFiltered{
		Out:     out,
		Writer:  writer,
		Scanner: bufio.NewScanner(reader),
	}
	go lf.read()
	return &lf
}

type LogFiltered struct {
	Out     io.Writer // original out
	Writer  io.Writer
	Scanner *bufio.Scanner
}

func (lf *LogFiltered) Write(p []byte) (int, error) {
	return lf.Writer.Write(p)
}

func (lf *LogFiltered) read() {
	for {
		if !lf.Scanner.Scan() {
			if err := lf.Scanner.Err(); err != nil {
				log.New(os.Stderr, "", log.LstdFlags).Printf("bufio.Scanner.Scan Err: %s", err)
			}
			continue
		}
		text := lf.Scanner.Bytes()
		if bytes.Contains(text, []byte(" handling ")) {
			continue
		}
		lf.Out.Write(append(text, []byte("\n")...))
	}
}
