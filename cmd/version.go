package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// VersionFlag is true when displaying version.
var VersionFlag bool

// OstentVersionCmd represents the version command
var OstentVersionCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "version",
	Short:        "Print ostent version.",
	RunE:         OstentVersionRunE,
}

func init() {
	OstentCmd.AddCommand(OstentVersionCmd)
}

func OstentVersionRunE(*cobra.Command, []string) error {
	log.New(os.Stdout, "", 0).Printf("Ostent version %+v", OstentVersion)
	return nil
}

func OstentVersionRun() error {
	if !VersionFlag {
		return nil
	}
	if err := OstentVersionRunE(nil, nil); err != nil {
		return err
	}
	os.Exit(0) // NB
	return nil
}

// OstentVersion of the latest known release. Compared to bin mode.
// MUST BE semver compatible: no two digits ("X.Y") allowed.
const OstentVersion = "0.6.2"
