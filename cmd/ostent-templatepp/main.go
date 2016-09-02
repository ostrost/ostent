package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/templateutil/templatepipe"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "ostent-templatepp",
	Short:        "Template preprocessor for ostent templates",
	Long: `Template preprocessor of ostent templates

ostent-templatepp is an utility program to deal with ostent templates:
it takes a template to read defines from and a template to base the output on.`,
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
	//- cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//- RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ostent-templatepp.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//- RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initFlags()
}

func initFlags() {
	RootCmd.RunE = runE
	RootCmd.PreRunE = preRunE
	RootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "An output file")
	RootCmd.Flags().StringVarP(&definesFromFile, "definesfrom", "d", "", "The html template file with defines")
	RootCmd.Flags().StringVarP(&inputTemplateFile, "template", "t", "", "The text template file to apply")
}

var (
	outputFile        string
	definesFromFile   string
	inputTemplateFile string
)

func preRunE(*cobra.Command, []string) error {
	if definesFromFile == "" {
		if inputTemplateFile == "" {
			return fmt.Errorf("--definesfrom and --template were not provided")
		}
		return fmt.Errorf("--definesfrom was not provided")
	}
	if inputTemplateFile == "" {
		return fmt.Errorf("--template was not provided")
	}
	return nil
}

func runE(*cobra.Command, []string) error {
	return templatepipe.Convert(
		inputTemplateFile,
		definesFromFile,
		nil,
		outputFile,
	)
}

func main() { Execute() }
