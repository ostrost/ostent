package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
)

func init() {
	var (
		fetchKeys []string
		fetchCont bool
		fetchCmd  = &cobra.Command{
			SilenceUsage: true,
			Use:          "fetch",
			Short:        "Print collected data",
			RunE: func(*cobra.Command, []string) error {
				return ostent.Fetch(OstentBind.ClientString(), fetchKeys, !fetchCont)
			},
		}
	)
	fetchCmd.PersistentFlags().StringSliceVarP(&fetchKeys, "key", "k", nil,
		"Reduce data with `key`")
	fetchCmd.PersistentFlags().BoolVarP(&fetchCont, "continue", "c", false,
		"Continuous printing")
	OstentCmd.AddCommand(fetchCmd)
}
