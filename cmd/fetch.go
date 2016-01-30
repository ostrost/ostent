package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/ostent"
)

func init() {
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "json", Short: "Output ALL data in JSON", RunE: FetchRunEFunc("")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "cpu", Short: "Output CPU data", RunE: FetchRunEFunc("cpu")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "df", Short: "Output DF data", RunE: FetchRunEFunc("diskUsage")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "if", Short: "Output IF data", RunE: FetchRunEFunc("ifaddrs")})
	OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
		Use: "mem", Short: "Output memory data", RunE: FetchRunEFunc("memory")})
	/*
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "la", Short: "Output loadavg data", RunE: FetchRunEFunc("loadavg")})
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "proc", Short: "Output processes data", RunE: FetchRunEFunc("procs")})
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "uptime", Short: "Output uptime data", RunE: FetchRunEFunc("uptime")})
		OstentCmd.AddCommand(&cobra.Command{SilenceUsage: true,
			Use: "hostname", Short: "Output hostname data", RunE: FetchRunEFunc("hostname")})
	  // */
}

func FetchRunEFunc(key string) func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		text, err := ostent.Fetch(OstentBind.ClientString(), key)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", text)
		return nil
	}
}
