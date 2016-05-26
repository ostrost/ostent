// +build !bin

package cmd

import (
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// http://blog.golang.org/profiling-go-programs

var (
	profileHeapOutput string // flag value
	profileCPUOutput  string // flag value
)

// profileHeapRun is a prerun func for starting heap profile.
func profileHeapRun() error {
	if profileHeapOutput == "" {
		return nil
	}
	file, err := os.Create(profileHeapOutput)
	if err != nil {
		return err
	}
	persistentPostRuns.add(func() error {
		logger := log.New(os.Stderr, "[ostent profile-heap] ", log.LstdFlags)
		if err := pprof.Lookup("heap").WriteTo(file, 1); err != nil {
			logger.Print(err) // just print
		}
		if err := file.Close(); err != nil {
			logger.Print(err) // just print
		}
		return nil
	})
	return nil
}

// profileCPURun is a prerun func for starting CPU profile.
func profileCPURun() error {
	if profileCPUOutput == "" {
		return nil
	}
	file, err := os.Create(profileCPUOutput)
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(file); err != nil {
		return err
	}
	persistentPostRuns.add(func() error {
		logger := log.New(os.Stderr, "[ostent profile-cpu] ", log.LstdFlags)
		logger.Print("Writing CPU profile")
		pprof.StopCPUProfile()
		if err := file.Close(); err != nil {
			logger.Print(err) // just print
		}
		return nil
	})
	return nil
}

func init() {
	OstentCmd.PersistentFlags().StringVar(&profileHeapOutput, "profile-heap", "",
		"Profiling heap output `filename`")
	OstentCmd.PersistentFlags().StringVar(&profileCPUOutput, "profile-cpu", "",
		"Profiling CPU output `filename`")
	persistentPreRuns.add(profileHeapRun)
	persistentPreRuns.add(profileCPURun)

	pkg, err := build.Import(pkgPath, "", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	genDocDir = filepath.Join(pkg.Dir, "doc")
	genDocCmd.Flags().StringVar(&genDocDir, "directory", genDocDir,
		"Output `directory` for saving docs")
	OstentCmd.AddCommand(genDocCmd)
}

// pkgPath defined for looking up the package directory.
const pkgPath = "github.com/ostrost/ostent"

var (
	genDocCmd = &cobra.Command{ // gendoc subcommand
		Use:   "gendoc",
		Short: "Generate ostent commands docs.",
		RunE:  genDocRunE,
	}
	genDocDir string // flag value
)

func genDocRunE(*cobra.Command, []string) error {
	OstentCmd.DisableAutoGenTag = true
	if cmd, _, err := OstentCmd.Find([]string{"gendoc"}); err == nil {
		// err is gone
		OstentCmd.RemoveCommand(cmd)
	}
	if err := doc.GenMarkdownTree(OstentCmd, genDocDir); err != nil {
		return err
	}
	mdfile := filepath.Join(genDocDir,
		strings.Replace(OstentCmd.CommandPath(), " ", "_", -1)+".md")
	text, err := ioutil.ReadFile(mdfile)
	if err != nil {
		return err
	}
	var lines []string
	for _, line := range strings.Split(string(text), "\n") {
		if strings.HasSuffix(line, "# SEE ALSO") {
			break
		}
		if !strings.Contains(line, "--profile-") { // skip dev-only flags
			lines = append(lines, line)
		}
	}
	return ioutil.WriteFile(mdfile, []byte(strings.Join(lines, "\n")), 0600)
}
