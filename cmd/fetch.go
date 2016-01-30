package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
)

func init() {
	var getKey string
	getCmd := &cobra.Command{SilenceUsage: true,
		Use: "get", Short: "Print collected data", RunE: FetchRunEFunc("", &getKey)}
	getCmd.PersistentFlags().StringVarP(&getKey, "key", "k", "", "Reduce data with `key`")
	OstentCmd.AddCommand(getCmd)
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "cpu", Short: "Print collected cpu data", RunE: FetchRunEFunc("cpu")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "df", Short: "Print collected df data", RunE: FetchRunEFunc("df")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "mem", Short: "Print collected mem data", RunE: FetchRunEFunc("mem")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "netio", Short: "Print collected netio data", RunE: FetchRunEFunc("netio")})
	/*
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "la", Short: "Print collected la (loadavg) data", RunE: FetchRunEFunc("la")})
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "proc", Short: "Print collected proc data", RunE: FetchRunEFunc("proc")})
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "uptime", Short: "Print collected uptime", RunE: FetchRunEFunc("uptime")})
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "hostname", Short: "Print collected hostname", RunE: FetchRunEFunc("hostname")})
	  // */
}

func FetchRunEFunc(key string, pkeys ...*string) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		if len(pkeys) != 0 {
			key = *pkeys[0]
		}
		text, err := ostent.Fetch(OstentBind.ClientString(), key)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", text)
		return nil
	}
}
