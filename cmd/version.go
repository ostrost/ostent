package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "version",
	Short:        "Print ostent version.",
	RunE:         versionRunE,
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// versionFlag is true when displaying version.
var versionFlag bool

func versionRunE(*cobra.Command, []string) error {
	log.New(os.Stdout, "", 0).Printf("Ostent version %+v", OstentVersion)
	return nil
}

func versionRun() error {
	if !versionFlag {
		return nil
	}
	if err := versionRunE(nil, nil); err != nil {
		return err
	}
	os.Exit(0) // NB
	return nil
}

// OstentVersion of the latest known release. Compared to bin mode.
// MUST BE semver compatible: no two digits ("X.Y") allowed.
const OstentVersion = "0.6.2"
