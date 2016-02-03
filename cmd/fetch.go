package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
)

func init() {
	var fetchCmdFlags FetchFlags
	fetchCmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "fetch",
		Short:        "Print collected data",
		Run:          fetchCmdFlags.Run,
	}
	fetchCmd.PersistentFlags().StringVarP(&fetchCmdFlags.K, "key", "k", "",
		"Reduce data with `key`")
	fetchCmd.PersistentFlags().BoolVarP(&fetchCmdFlags.Cont, "continue", "c", false,
		"Continuous printing")
	OstentCmd.AddCommand(fetchCmd)
}

// FetchFlags is the flags for fetch command.
type FetchFlags struct {
	K    string
	Cont bool
}

// Run is to be run with the fetch command with the flags.
func (flags *FetchFlags) Run(*cobra.Command, []string) error {
	return ostent.Fetch(OstentBind.ClientString(), flags.K, flags.Cont)
}
