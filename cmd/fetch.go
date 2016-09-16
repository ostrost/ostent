package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
	"github.com/ostrost/ostent/params"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "fetch",
	Short:        "Print collected data.",
	Long: `Fetch retrieves data from one or many ostent servers for print. Server addresses
and request details may be specified with --keys flag. The value is parsed as a
partial URL: most details can be omitted for defaults. The fragment is treated
as key(s) to reduce data with. Join multiple keys with "#".`,
	Example: `ostent fetch --keys \#uptime
ostent fetch --keys '?times=-1#uptime'
ostent fetch --keys 'http://10.0.0.4:8050?cpun=1#cpu'
ostent fetch --keys 'http://10.0.0.5:8050#hostname,http://10.0.0.6:8050#hostname'`,
	RunE: fetchRun,
}

func init() {
	RootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fetchCmd.PersistentFlags().VarP(fetchKeys, "keys", "k", "Ostent server(s) `endpoint(s)`")
}

var fetchKeys = params.NewFetchKeys(8050)

func fetchRun(*cobra.Command, []string) error {
	if len(fetchKeys.Values) == 0 {
		if err := fetchKeys.Set(""); err != nil {
			return err
		}
	}
	return ostent.Fetch(fetchKeys)
}
