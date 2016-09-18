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

	"github.com/ostrost/ostent/share/assets"
)

// extractassetsCmd represents the extractassets command
var extractassetsCmd = &cobra.Command{
	Use:     "extractassets",
	Short:   "Extract ostent embeded assets & manage symlinks in current directory.",
	PreRunE: extractassetsPreRunE,
	RunE:    extractassetsRunE,
}

func init() {
	RootCmd.AddCommand(extractassetsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// extractassetsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// extractassetsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	extractassetsCmd.Flags().StringVar(&extractDestDir, "destdir",
		OstentVersion /* default is this */, "destrination directory for extraction")
}

var (
	// extractDestDir is a flag value.
	extractDestDir string

	// eaLog is a logger to log this command's messages with.
	eaLog = log.New(os.Stderr, "[ostent extract-assets] ", log.LstdFlags)
)

func extractassetsPreRunE(*cobra.Command, []string) error {
	if extractDestDir == "" {
		return fmt.Errorf("--destdir wasn't provided")
	}
	return nil
}

// extractassetsRunE does the following:
// - creates the dest directory
// - every asset is saved as a file
// - every asset is gzipped saved as a file + .gz if it's size is above threshold
func extractassetsRunE(*cobra.Command, []string) error {
	if _, err := os.Stat(extractDestDir); err == nil {
		return fmt.Errorf("%s: File exists\n", extractDestDir)
	}
	// RestoreAssets (among other things) creates DestDir.
	if err := assets.RestoreAssets(extractDestDir, ""); err != nil {
		return err
	}
	for _, name := range assets.AssetNames() {
		if err := extractGzip(name); err != nil {
			return err
		}
	}
	return nil
}

func extractGzip(name string) error {
	text, err := assets.Asset(name)
	if err != nil {
		eaLog.Printf("assets.Asset: %s: %s", name, err)
		return nil // continue
	}
	full := filepath.Join(extractDestDir, name)
	if name == "favicon.ico" || name == "robots.txt" {
		if err = extractSymlink(name, full); err != nil {
			eaLog.Printf("extractSymlink: %s: %s", name, err)
			return nil // continue
		}
	}

	now := time.Now()
	if err = os.Chtimes(full, now, now); err != nil {
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

func extractSymlink(name, full string) error {
	if dest, err := os.Readlink(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		eaLog.Printf("Removing symlink %q pointing to %q", name, dest)
		if err := os.Remove(name); err != nil {
			return err
		}
	}
	return os.Symlink(full, name)
	// no need to os.Chtimes as os.Symlink will set the times to about now
}

func watchConfig() {}
