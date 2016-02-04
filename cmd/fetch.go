package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
)

func init() {
	var (
		// fetchServer is the flag value.
		fetchServer = flags.NewBind("127.0.0.1", 8050)
		fetchKeys   []string
		fetchCont   bool
		fetchCmd    = &cobra.Command{
			SilenceUsage: true,
			Use:          "fetch",
			Short:        "Print collected data",
			RunE: func(*cobra.Command, []string) error {
				return ostent.Fetch(fetchServer.String(), fetchKeys, !fetchCont)
			},
		}
	)
	fetchCmd.Flags().VarP(&fetchServer, "server", "s", "Server `address` to fetch from")
	fetchCmd.Flags().StringSliceVarP(&fetchKeys, "key", "k", nil, "Reduce data with `key(s)`")
	fetchCmd.Flags().BoolVarP(&fetchCont, "continue", "c", false, "Continuous printing")
	OstentCmd.AddCommand(fetchCmd)
}
