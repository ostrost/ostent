package main

import (
	"compress/gzip"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/share/assets"
)

// extractAssets into `ostent.VERSION' dir:
// - the dir is created, otherwise bails out
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func extractAssets() {
	logger := localLog{log.New(os.Stderr, "[ostent extract-assets] ", log.LstdFlags)}

	dest := ostent.VERSION
	if _, err := os.Stat(dest); err == nil {
		logger.Fatalf("%s: File exists\n", dest)
	}
	logger.fatal(os.Mkdir(dest, os.ModePerm))
	for _, path := range assets.AssetNames() {
		text, err := assets.Asset(path)
		if err != nil {
			logger.Printf("Unexpected: %s: %s", path, err)
			continue
		}
		full := filepath.Join(dest, path)
		dir := filepath.Dir(full)
		if _, err := os.Stat(dir); err != nil {
			logger.fatal(os.MkdirAll(dir, os.ModePerm))
		}

		file, err := os.Create(full)
		logger.fatal(err)

		if _, err := file.Write(text); err != nil {
			logger.fatal(err)
		}
		file.Close()

		gz := len(text) > 1024
		if !gz {
			continue
		}

		gzfile, err := os.Create(full + ".gz")
		logger.fatal(err)

		gzwriter := gzip.NewWriter(gzfile)
		if _, err := gzwriter.Write(text); err != nil {
			logger.fatal(err)
		}
		gzwriter.Close()
		gzfile.Close()
	}
}

type localLog struct {
	*log.Logger
}

func (l *localLog) fatal(err error) {
	if err != nil {
		l.Fatalln(err)
	}
}

type command struct{}

func (c *command) Run() {
	if c == nil {
		log.Fatalf("No such command")
	}
	extractAssets()
}

var commands = map[string]*command{
	"extract-assets": &command{},
}

func argCommand() (*command, bool) {
	if flag.NArg() == 0 {
		return nil, false
	}
	return commands[flag.Arg(0)], true
}
