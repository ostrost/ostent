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
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

// gendocCmd represents the gendoc command
var gendocCmd = &cobra.Command{
	Use:   "gendoc",
	Short: "Generate ostent commands docs.",
	RunE:  gendocRunE,
}

func init() {
	RootCmd.AddCommand(gendocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVar(&profileHeapOutput, "profile-heap", "",
		"Profiling heap output `filename`")
	RootCmd.PersistentFlags().StringVar(&profileCPUOutput, "profile-cpu", "",
		"Profiling CPU output `filename`")
	persistentPreRuns.add(profileHeapRun)
	persistentPreRuns.add(profileCPURun)

	pkg, err := build.Import(pkgPath, "", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	gendocDir = filepath.Join(pkg.Dir, "doc")
	gendocCmd.Flags().StringVar(&gendocDir, "directory", gendocDir,
		"Output `directory` for saving docs")
}

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

// pkgPath defined for looking up the package directory.
const pkgPath = "github.com/ostrost/ostent"

var gendocDir string // flag value

func gendocRunE(*cobra.Command, []string) error {
	RootCmd.DisableAutoGenTag = true
	if cmd, _, err := RootCmd.Find([]string{"gendoc"}); err == nil {
		// err is gone
		RootCmd.RemoveCommand(cmd)
	}
	if err := doc.GenMarkdownTree(RootCmd, gendocDir); err != nil {
		return err
	}
	mdfile := filepath.Join(gendocDir,
		strings.Replace(RootCmd.CommandPath(), " ", "_", -1)+".md")
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

func watchConfig() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("%s config file changed\n", e.Name)
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	})
	viper.WatchConfig()
}
