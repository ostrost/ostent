// +build production

package commands

import (
	"compress/gzip"
	"flag"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
)

type assetsRestore struct {
	logger  *Logger
	destdir string
}

func assetsRestoreCommand(fs *flag.FlagSet, loggerOptions ...SetupLogger) (CommandHandler, io.Writer) {
	ae := &assetsRestore{
		destdir: ostent.VERSION,
		logger:  NewLogger("[ostent restore-assets] ", loggerOptions...),
	}
	fs.StringVar(&ae.destdir, "d", ae.destdir, "Destination directory")
	return ae.run, ae.logger
}

// run does the following:
// - creates the dest directory
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func (ae *assetsRestore) run() {
	if _, err := os.Stat(ae.destdir); err == nil {
		ae.logger.Fatalf("%s: File exists\n", ae.destdir)
	}
	ae.logger.fatalif(assets.RestoreAssets(ae.destdir, ""))
	// ae.logger.fatalif(os.Mkdir(ae.destdir, os.ModePerm))
	for _, name := range assets.AssetNames() {
		text, err := assets.Asset(name)
		if err != nil {
			ae.logger.Printf("assets.Asset: %s: %s", name, err)
			continue
		}
		full := filepath.Join(ae.destdir, name)
		/* dir := filepath.Dir(full)
		if _, err := os.Stat(dir); err != nil {
			ae.logger.fatalif(os.MkdirAll(dir, os.ModePerm))
		}

		file, err := os.Create(full)
		ae.logger.fatalif(err)

		_, err = file.Write(text)
		ae.logger.fatalif(err)
		file.Close() // */

		now := time.Now()
		ae.logger.fatalif(os.Chtimes(full, now, now))

		if len(text) <= 1024 {
			continue
		}

		gzfile, err := os.Create(full + ".gz")
		ae.logger.fatalif(err)

		gzwriter := gzip.NewWriter(gzfile)
		_, err = gzwriter.Write(text)
		ae.logger.fatalif(err)

		gzwriter.Close()
		gzfile.Close()

		ae.logger.fatalif(os.Chtimes(full+".gz", now, now))
	}
}

func init() {
	AddCommand("restore-assets", assetsRestoreCommand)
}
