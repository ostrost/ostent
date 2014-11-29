// +build production

package commands

import (
	"compress/gzip"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/ostrost/ostent"
	"github.com/ostrost/ostent/share/assets"
)

type assetsExtract struct {
	logger  *Logger
	destdir string
}

func assetsExtractCommand(fs *flag.FlagSet) (CommandHandler, io.Writer) {
	ae := &assetsExtract{
		destdir: ostent.VERSION,
		logger: NewLogger(func(l *Logger) {
			l.Prefix = "[ostent extract-assets] "
		}),
	}
	fs.StringVar(&ae.destdir, "d", ae.destdir, "Destination directory")
	return ae.run, ae.logger
}

// run does the following:
// - creates the dest directory
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func (ae *assetsExtract) run() {
	if _, err := os.Stat(ae.destdir); err == nil {
		ae.logger.Fatalf("%s: File exists\n", ae.destdir)
	}
	ae.logger.fatalif(os.Mkdir(ae.destdir, os.ModePerm))
	for _, path := range assets.AssetNames() {
		text, err := assets.Asset(path)
		if err != nil {
			ae.logger.Printf("Unexpected: %s: %s", path, err)
			continue
		}
		full := filepath.Join(ae.destdir, path)
		dir := filepath.Dir(full)
		if _, err := os.Stat(dir); err != nil {
			ae.logger.fatalif(os.MkdirAll(dir, os.ModePerm))
		}

		file, err := os.Create(full)
		ae.logger.fatalif(err)

		_, err = file.Write(text)
		ae.logger.fatalif(err)
		file.Close()

		gz := len(text) > 1024
		if !gz {
			continue
		}

		gzfile, err := os.Create(full + ".gz")
		ae.logger.fatalif(err)

		gzwriter := gzip.NewWriter(gzfile)
		_, err = gzwriter.Write(text)
		ae.logger.fatalif(err)

		gzwriter.Close()
		gzfile.Close()
	}
}

func init() {
	AddCommand("extract-assets", assetsExtractCommand)
}
