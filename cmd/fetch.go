package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
)

func init() {
	var fetchCmdFlags FetchFlags
	fetchCmd := &cobra.Command{SilenceUsage: true,
		Use: "fetch", Short: "Print collected/partial data"}
	fetchCmd.PersistentFlags().StringVarP(&fetchCmdFlags.K, "key", "k", "",
		"Reduce data with `key`")
	AddFetchCommand(&fetchCmdFlags, fetchCmd)

	AddFetchCommand(&FetchFlags{K: "cpu"}, &cobra.Command{SilenceUsage: true,
		Use: "cpu", Short: "Print collected cpu data"})
	AddFetchCommand(&FetchFlags{K: "df"}, &cobra.Command{SilenceUsage: true,
		Use: "df", Short: "Print collected df data"})
	AddFetchCommand(&FetchFlags{K: "mem"}, &cobra.Command{SilenceUsage: true,
		Use: "mem", Short: "Print collected mem data"})
	AddFetchCommand(&FetchFlags{K: "netio"}, &cobra.Command{SilenceUsage: true,
		Use: "netio", Short: "Print collected netio data"})
	/*
		AddFetchCommand(&FetchFlags{K: "la"}, &cobra.Command{SilenceUsage: true,
			Use: "la", Short: "Print collected la (loadavg) data"})
		AddFetchCommand(&FetchFlags{K: "proc"}, &cobra.Command{SilenceUsage: true,
			Use: "proc", Short: "Print collected proc data"})
		AddFetchCommand(&FetchFlags{K: "uptime"}, &cobra.Command{SilenceUsage: true,
			Use: "uptime", Short: "Print collected uptime"})
		AddFetchCommand(&FetchFlags{K: "hostname"}, &cobra.Command{SilenceUsage: true,
			Use: "hostname", Short: "Print collected hostname"})
		// */
}

// AddFetchCommand sets up cmd and adds it to OstentCmd.
func AddFetchCommand(flags *FetchFlags, cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&flags.Cont, "continue", "c", false,
		"Continuous printing")
	cmd.RunE = flags.Run
	OstentCmd.AddCommand(cmd)
}

// FetchFlags is the flags for any fetch command.
type FetchFlags struct {
	K    string
	Cont bool
}

// Run is to be run with the fetch command with the flags.
func (flags *FetchFlags) Run(*cobra.Command, []string) error {
	return ostent.Fetch(OstentBind.ClientString(), flags.K, flags.Cont)
}
