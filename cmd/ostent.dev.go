// +build !bin

package cmd

import (
	"github.com/ostrost/ostent/cmd/cmdcobra"
)

func init() {
	OstentCmd.PersistentFlags().StringVar(&cmdcobra.ProfileHeapOutput, "profile-heap", "",
		"Profiling heap output `filename`")
	OstentCmd.PersistentFlags().StringVar(&cmdcobra.ProfileCPUOutput, "profile-cpu", "",
		"Profiling CPU output `filename`")
	cmdcobra.PreRuns.Add(cmdcobra.ProfileHeapRun)
	cmdcobra.PreRuns.Add(cmdcobra.ProfileCPURun)
}
