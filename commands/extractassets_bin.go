// +build bin

package commands

import (
	"compress/gzip"
	"flag"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ostrost/ostent/commands/extpoints"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
)

type assetsRestore struct {
	Log     *extpoints.Log
	destdir string
}

func assetsRestoreCommand(fs *flag.FlagSet, loggerOptions ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer) {
	ae := &assetsRestore{
		destdir: ostent.VERSION,
		Log:     NewLog("[ostent restore-assets] ", loggerOptions...),
	}
	fs.StringVar(&ae.destdir, "d", ae.destdir, "Destination directory")
	return ae.Run, ae.Log
}

// Run does the following:
// - creates the dest directory
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func (ae *assetsRestore) Run() {
	if _, err := os.Stat(ae.destdir); err == nil {
		ae.Log.Fatalf("%s: File exists\n", ae.destdir)
	}
	ae.Check(assets.RestoreAssets(ae.destdir, ""))
	// ae.Check(os.Mkdir(ae.destdir, os.ModePerm))
	for _, name := range assets.AssetNames() {
		text, err := assets.Asset(name)
		if err != nil {
			ae.Log.Printf("assets.Asset: %s: %s", name, err)
			continue
		}
		full := filepath.Join(ae.destdir, name)
		/* dir := filepath.Dir(full)
		if _, err := os.Stat(dir); err != nil {
			ae.Check(os.MkdirAll(dir, os.ModePerm))
		}

		file, err := os.Create(full)
		ae.Check(err)

		_, err = file.Write(text)
		ae.Check(err)
		file.Close() // */

		now := time.Now()
		ae.Check(os.Chtimes(full, now, now))

		if len(text) <= 1024 {
			continue
		}

		gzfile, err := os.Create(full + ".gz")
		ae.Check(err)

		gzwriter := gzip.NewWriter(gzfile)
		_, err = gzwriter.Write(text)
		ae.Check(err)

		gzwriter.Close()
		gzfile.Close()

		ae.Check(os.Chtimes(full+".gz", now, now))
	}
}

func (ae assetsRestore) Check(err error) {
	if err != nil {
		ae.Log.Fatalln(err)
	}
}

func init() {
	AddCommand("restore-assets", assetsRestoreCommand)
}
