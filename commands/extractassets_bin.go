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

type ExtractAssets struct {
	Log     *extpoints.Log
	DestDir string
}

func (_ Assets) SetupCommand(fs *flag.FlagSet, loggerOptions ...extpoints.SetupLog) (extpoints.CommandHandler, io.Writer) {
	ea := &ExtractAssets{
		DestDir: ostent.VERSION,
		Log:     NewLog("[ostent extract-assets] ", loggerOptions...),
	}
	fs.StringVar(&ea.DestDir, "d", ea.DestDir, "Destination directory")
	return ea.Run, ea.Log
}

// Run does the following:
// - creates the dest directory
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func (ea *ExtractAssets) Run() {
	if _, err := os.Stat(ea.DestDir); err == nil {
		ea.Log.Fatalf("%s: File exists\n", ea.DestDir)
	}
	// RestoreAssets (among other things) creates DestDir.
	ea.Check(assets.RestoreAssets(ea.DestDir, ""))
	for _, name := range assets.AssetNames() {
		ea.Check(ea.Gzip(name))
	}
}

func (ea *ExtractAssets) Gzip(name string) error {
	text, err := assets.Asset(name)
	if err != nil {
		ea.Log.Printf("assets.Asset: %s: %s", name, err)
		return nil // continue
	}
	full := filepath.Join(ea.DestDir, name)
	if name == "favicon.ico" || name == "robots.txt" {
		if err := ea.Symlink(name, full); err != nil {
			ea.Log.Printf("Symlink: %s: %s", name, err)
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

func (ea *ExtractAssets) Symlink(name, full string) error {
	if dest, err := os.Readlink(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		ea.Log.Printf("Removing symlink %q pointing to %q", name, dest)
		if err := os.Remove(name); err != nil {
			return err
		}
	}
	return os.Symlink(full, name)
	// no need to os.Chtimes as os.Symlink will set the times to about now
}

func (ea *ExtractAssets) Check(err error) {
	if err != nil {
		ea.Log.Fatalln(err)
	}
}

type Assets struct{}

func init() {
	extpoints.Commands.Register(Assets{}, "extractassets")
}
