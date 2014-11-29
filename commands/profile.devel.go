// +build !production
// http://blog.golang.org/profiling-go-programs

package commands

import (
	"flag"
	"os"
	"runtime/pprof"
)

type memprofile struct {
	logger *Logger
	output string
	_f     *os.File
}

func memProfileCommandLine(cli *flag.FlagSet) CommandLineHandler {
	mp := &memprofile{
		logger: NewLogger("[ostent memprofile] "),
	}
	cli.StringVar(&mp.output, "memprofile", "", "MEM profile output file")
	return func() (AtexitHandler, bool, error) {
		if mp.output == "" {
			return nil, false, nil
		}
		return mp.atexit, false, mp.run()
	}
}

func (mp *memprofile) atexit() {
	mp.logger.Print("Writing MEM profile")
	if err := pprof.WriteHeapProfile(mp._f); err != nil {
		mp.logger.Print(err) // just print
	}
	if err := mp._f.Close(); err != nil {
		mp.logger.Print(err) // just print
	}
}

func (mp *memprofile) run() (err error) {
	mp._f, err = os.Create(mp.output)
	if err != nil {
		mp.logger.Fatal(err) // log with the mp logger
	}
	return err
}

/* ******************************************************************************** */

type cpuprofile struct {
	logger *Logger
	output string
	_f     *os.File
}

func cpuProfileCommandLine(cli *flag.FlagSet) CommandLineHandler {
	cp := &cpuprofile{
		logger: NewLogger("[ostent cpuprofile] "),
	}
	cli.StringVar(&cp.output, "cpuprofile", "", "CPU profile output file")
	return func() (AtexitHandler, bool, error) {
		if cp.output == "" {
			return nil, false, nil
		}
		return cp.atexit, false, cp.run()
	}
}

func (cp *cpuprofile) atexit() {
	cp.logger.Print("Writing CPU profile")
	pprof.StopCPUProfile()
	if err := cp._f.Close(); err != nil {
		cp.logger.Print(err) // just print
	}
}

func (cp *cpuprofile) run() (err error) {
	cp._f, err = os.Create(cp.output)
	if err == nil {
		err = pprof.StartCPUProfile(cp._f)
	}
	if err != nil {
		cp.logger.Fatal(err) // log with the cp logger
	}
	return err
}

func init() {
	AddCommandLine(cpuProfileCommandLine)
	AddCommandLine(memProfileCommandLine)
}
