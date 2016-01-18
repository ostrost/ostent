package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ostrost/ostent/templateutil/templatefunc"
	"github.com/ostrost/ostent/templateutil/templatepipe"
)

// TemplateppCmd sets the main command.
var TemplateppCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "ostent-templatepp",
	Short:        "Template preprocessor for ostent templates",
	Long: `Template preprocessor of ostent templates

ostent-templatepp is an utility program to deal with ostent templates:
it takes a template to read defines from and a template to base the output on.`,
}

var (
	outputFile        string
	definesFromFile   string
	inputTemplateFile string
)

// TemplateppPreRunE is to become TemplateppCmd.PreRunE.
func TemplateppPreRunE(*cobra.Command, []string) error {
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

// TemplateppRunE is to become TemplateppCmd.RunE.
// Calls templatepipe.Convert.
func TemplateppRunE(*cobra.Command, []string) error {
	return templatepipe.Convert(
		inputTemplateFile,
		definesFromFile,
		templatefunc.FuncMapHTML(),
		templatefunc.FuncMapJSXL(),
		outputFile,
	)
}

func main() {
	if err := TemplateppCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	TemplateppCmd.RunE = TemplateppRunE
	TemplateppCmd.PreRunE = TemplateppPreRunE
	TemplateppCmd.Flags().StringVarP(&outputFile, "output", "o", "", "An output file")
	TemplateppCmd.Flags().StringVarP(&definesFromFile, "definesfrom", "d", "", "The html template file with defines")
	TemplateppCmd.Flags().StringVarP(&inputTemplateFile, "template", "t", "", "The text template file to apply")
}
