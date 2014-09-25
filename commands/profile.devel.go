// +build !production
// http://blog.golang.org/profiling-go-programs

package commands

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

type memprofile struct {
	logger *loggerWriter
	output string
}

func FlagSetNewMEMProfile(fs *flag.FlagSet) *memprofile { // fs better be flag.CommandLine
	mp := memprofile{logger: &loggerWriter{log.New(os.Stderr, "[ostent memprofile] ", log.LstdFlags)}}
	fs.StringVar(&mp.output, "memprofile", "", "MEM profile output file.")
	return &mp
}

func (mp memprofile) MakeDeferrer() deferred {
	if mp.output == "" {
		return nil
	}

	f, err := os.Create(mp.output)
	if err != nil {
		mp.logger.Fatal(err)
	}

	return func() {
		mp.logger.Print("Writing MEM profile")
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}

/* ******************************************************************************** */

type cpuprofile struct {
	logger *loggerWriter
	output string
}

func FlagSetNewCPUProfile(fs *flag.FlagSet) *cpuprofile { // fs better be flag.CommandLine
	cp := cpuprofile{logger: &loggerWriter{log.New(os.Stderr, "[ostent cpuprofile] ", log.LstdFlags)}}
	fs.StringVar(&cp.output, "cpuprofile", "", "CPU profile output file.")
	return &cp
}

func (cp cpuprofile) MakeDeferrer() deferred {
	if cp.output == "" {
		return nil
	}

	f, err := os.Create(cp.output)
	if err != nil {
		cp.logger.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	return func() {
		cp.logger.Print("Writing CPU profile")
		pprof.StopCPUProfile()
		f.Close()
	}
}

/*
func cpuprofileCommand(fs *flag.FlagSet, arguments []string) (deferredFunc, error, []string) {
	cp := FlagSetNewCPUProfile(fs)
	// no flags
	fs.SetOutput(cp.logger)
	err := fs.Parse(arguments)
	return cp..., err, fs.Args()
} // */

func init() {
	AddDefault("cpuprofile", FlagSetNewCPUProfile(flag.CommandLine))
	AddDefault("memprofile", FlagSetNewMEMProfile(flag.CommandLine))
}
