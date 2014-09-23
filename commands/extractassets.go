// !build production

package commands

import (
	"compress/gzip"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/share/assets"
)

func newAssetsExtract(fs *flag.FlagSet, arguments []string) (sub, error, []string) {
	ae := assetsExtract{logger: &loggerWriter{log.New(os.Stderr,
		"[ostent extract-assets] ", log.LstdFlags)}}
	fs.SetOutput(ae.logger)
	err := fs.Parse(arguments)
	return ae.run, err, fs.Args()
}

type assetsExtract struct {
	logger *loggerWriter
}

func (ae assetsExtract) run() {
	extractAssets(ae.logger.Logger)
}

// extractAssets into `ostent.VERSION' dir:
// - the dir is created, otherwise bails out
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func extractAssets(loglogger *log.Logger) {
	logger := localLog{loglogger}

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
	AddCommand("extract-assets", newAssetsExtract)
}
