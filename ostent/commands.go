package main

import (
	"compress/gzip"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

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
	logger.fatalif(os.Mkdir(dest, os.ModePerm))
	for _, path := range assets.AssetNames() {
		text, err := assets.Asset(path)
		if err != nil {
			logger.Printf("Unexpected: %s: %s", path, err)
			continue
		}
		full := filepath.Join(dest, path)
		dir := filepath.Dir(full)
		if _, err := os.Stat(dir); err != nil {
			logger.fatalif(os.MkdirAll(dir, os.ModePerm))
		}

		file, err := os.Create(full)
		logger.fatalif(err)

		_, err = file.Write(text)
		logger.fatalif(err)
		file.Close()

		gz := len(text) > 1024
		if !gz {
			continue
		}

		gzfile, err := os.Create(full + ".gz")
		logger.fatalif(err)

		gzwriter := gzip.NewWriter(gzfile)
		_, err = gzwriter.Write(text)
		logger.fatalif(err)

		gzwriter.Close()
		gzfile.Close()
	}
}

type localLog struct {
	*log.Logger
}

func (l *localLog) fatalif(err error) {
	if err != nil {
		l.Fatalln(err)
	}
}

func init() {
	AddCommand("extract-assets", extractAssets)
}

type command func()

var (
	comutex  sync.Mutex
	commands = make(map[string]command)
)

func AddCommand(name string, fun func()) {
	comutex.Lock()
	defer comutex.Unlock()
	commands[name] = fun
}

func ArgCommand() command {
	if flag.NArg() == 0 {
		return nil
	}
	name := flag.Arg(0)
	enoent := func() {
		log.Fatalln("%s: No such command", name)
	}
	if name == "" {
		return enoent
	}
	comutex.Lock()
	defer comutex.Unlock()
	if fun, ok := commands[name]; ok {
		return fun
	}
	return enoent
}
