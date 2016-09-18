package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "ostent",
	Short:        "Ostent is a metrics tool.",
	Long:         "Ostent is a server tool to collect, display and report system metrics.",

	PersistentPostRunE: persistentPostRuns.runE,
	PersistentPreRunE:  persistentPreRuns.runE,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.Flags().StringVar(&cfgFile, "config", "", "config filename (default \"$HOME/.ostent.toml\")")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//- RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initFlags()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ostent") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	viper.SetConfigFile(cfgFile) // NB AFTER .SetConfigName for flag value
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
