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
	_f     *os.File
}

func memProfileCommandLine(cli *flag.FlagSet) commandLineHandler {
	mp := &memprofile{
		logger: &loggerWriter{
			log.New(os.Stderr, "[ostent memprofile] ", log.LstdFlags),
		},
	}
	cli.StringVar(&mp.output, "memprofile", "", "MEM profile output file")
	return func() (atexitHandler, bool, error) {
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
	logger *loggerWriter
	output string
	_f     *os.File
}

func cpuProfileCommandLine(cli *flag.FlagSet) commandLineHandler {
	cp := &cpuprofile{
		logger: &loggerWriter{
			log.New(os.Stderr, "[ostent cpuprofile] ", log.LstdFlags),
		},
	}
	cli.StringVar(&cp.output, "cpuprofile", "", "CPU profile output file")
	return func() (atexitHandler, bool, error) {
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
