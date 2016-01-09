// +build bin

package cmd

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/share/assets"
)

var (
	// ExtractDestDir is a flag value.
	ExtractDestDir string

	// EaLog is a logger to log this subcommand's messages with.
	EaLog = log.New(os.Stderr, "[ostent extract-assets] ", log.LstdFlags)
)

// ExtractAssetsCmd represents the extractassets subcommand
var ExtractAssetsCmd = &cobra.Command{
	Use:   "extractassets",
	Short: "extractassets subcommand extracts embeded assets",
	// Long: ``,
	PreRunE: ExtractAssetsPreRunE,
	RunE:    ExtractAssetsRunE,
}

func init() {
	OstentCmd.AddCommand(ExtractAssetsCmd)
	ExtractAssetsCmd.Flags().StringVarP(&ExtractDestDir, "destdir", "d",
		ostent.VERSION /* default is this */, "Destrination directory for extraction")
}

func ExtractAssetsPreRunE(*cobra.Command, []string) error {
	if ExtractDestDir == "" {
		return fmt.Errorf("--destdir wasn't provided")
	}
	return nil
}

// ExtractAssetsRunE does the following:
// - creates the dest directory
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func ExtractAssetsRunE(*cobra.Command, []string) error {
	if _, err := os.Stat(ExtractDestDir); err == nil {
		return fmt.Errorf("%s: File exists\n", ExtractDestDir)
	}
	// RestoreAssets (among other things) creates DestDir.
	if err := assets.RestoreAssets(ExtractDestDir, ""); err != nil {
		return err
	}
	for _, name := range assets.AssetNames() {
		if err := ExtractGzip(name); err != nil {
			return err
		}
	}
	return nil
}

func ExtractGzip(name string) error {
	text, err := assets.Asset(name)
	if err != nil {
		EaLog.Printf("assets.Asset: %s: %s", name, err)
		return nil // continue
	}
	full := filepath.Join(ExtractDestDir, name)
	if name == "favicon.ico" || name == "robots.txt" {
		if err := ExtractSymlink(name, full); err != nil {
			EaLog.Printf("ExtractSymlink: %s: %s", name, err)
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

func ExtractSymlink(name, full string) error {
	if dest, err := os.Readlink(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		EaLog.Printf("Removing symlink %q pointing to %q", name, dest)
		if err := os.Remove(name); err != nil {
			return err
		}
	}
	return os.Symlink(full, name)
	// no need to os.Chtimes as os.Symlink will set the times to about now
}
