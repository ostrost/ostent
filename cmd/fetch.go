package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/flags"
	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

func init() {
	var (
		// fetchKeys is the flag value.
		fetchKeys = params.NewFetchKeys(flags.NewBind("127.0.0.1", 8050))
		fetchCmd  = &cobra.Command{
			SilenceUsage: true,
			Use:          "fetch",
			Short:        "Print collected data",
			Long: `Fetch retrieves data from one or many ostent servers for print. Server addresses
and request details may be specified with --keys flag. The value is parsed as a
partial URL: most details can be omitted for defaults. The fragment is treated
as key(s) to reduce data with. Join multiple keys with "#".`,
			Example: `ostent fetch --keys \#uptime
ostent fetch --keys '?times=-1#uptime'
ostent fetch --keys 'http://10.0.0.4:8050?cpun=1#cpu'
ostent fetch --keys 'http://10.0.0.5:8050#hostname,http://10.0.0.6:8050#hostname'`,
			RunE: func(*cobra.Command, []string) error { return ostent.Fetch(fetchKeys) },
		}
	)
	fetchCmd.Flags().VarP(fetchKeys, "keys", "k", "Ostent server(s) `endpoint(s)`")
	OstentCmd.AddCommand(fetchCmd)
}
