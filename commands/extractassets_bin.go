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

func (_ Assets) SetupCommand(fs *flag.FlagSet, loggerOptions ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer) {
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
	// RestoreAssets (among other things) creates DestDir.
	ae.Check(assets.RestoreAssets(ae.destdir, ""))
	for _, name := range assets.AssetNames() {
		ae.Check(ae.Gzip(name))
	}
}

func (ae *assetsRestore) Gzip(name string) error {
	if true { // indent
		text, err := assets.Asset(name)
		if err != nil {
			ae.Log.Printf("assets.Asset: %s: %s", name, err)
			return nil // continue
		}
		full := filepath.Join(ae.destdir, name)
		if name == "favicon.ico" || name == "robots.txt" {
			if err := ae.Symlink(name, full); err != nil {
				ae.Log.Printf("Symlink: %s: %s", name, err)
				return nil // continue
			}
		}

		now := time.Now()
		if err := os.Chtimes(full, now, now); err != nil {
			return err
		}

		if len(text) <= 1024 {
			return nil // continue
		}

		gzfile, err := os.Create(full + ".gz")
		if err != nil {
			return err
		}

		gzwriter := gzip.NewWriter(gzfile)
		_, err = gzwriter.Write(text)
		if err != nil {
			return err
		}

		gzwriter.Close()
		gzfile.Close()

		return os.Chtimes(full+".gz", now, now)
	}
}

func (ae *assetsRestore) Symlink(name, full string) error {
	if dest, err := os.Readlink(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		ae.Log.Printf("Removing symlink %q pointing to %q", name, dest)
		if err := os.Remove(name); err != nil {
			return err
		}
	}
	return os.Symlink(full, name)
	// no need to os.Chtimes as os.Symlink will set the times to about now
}

func (ae assetsRestore) Check(err error) {
	if err != nil {
		ae.Log.Fatalln(err)
	}
}

type Assets struct{}

func init() {
	extpoints.Commands.Register(Assets{}, "restore-assets")
}
