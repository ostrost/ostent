package stdlogfilter

import (
	"bytes"
	"log"
	"testing"
)

func TestLogFiltering(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(NewLogFiltered(buf), "", log.LstdFlags)
	logger.Printf("one\n")
	logger.Printf("x handling y\n")
	logger.Printf("two\n")
	if bytes.Contains(buf.Bytes(), []byte(" handling ")) {
		t.Errorf("LogFilterted did not filter")
	}
}
