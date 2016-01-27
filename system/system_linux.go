// +build linux

package system

import (
	"fmt"
	"os/exec"
	"strings"

	sigar "github.com/ostrost/gosigar"
	metrics "github.com/rcrowley/go-metrics"
)

// Distrib returns system distribution and version string.
func Distrib() (string, error) {
	// https://unix.stackexchange.com/q/6345
	std, err := exec.Command("lsb_release", "--id", "--release").CombinedOutput()
	if err != nil {
		return "", err
	}
	distrib, release := "", ""
	// strings.TrimRight(string(std), "\n\t ")
	for _, s := range strings.Split(string(std), "\n") {
		b := strings.Split(s, "\t")
		if len(b) == 2 {
			if b[0] == "Distributor ID:" {
				distrib = b[1]
				continue
			} else if b[0] == "Release:" {
				release = b[1]
				continue
			}
		}
	}
	if distrib == "" {
		return "", fmt.Errorf("Could not identify Distributor ID from lsb_release output")
	}
	if release == "" {
		return distrib, fmt.Errorf("Could not identify Release from lsb_release output")
	}
	return distrib + " " + release, nil
}

// ProcName returns argv[0].
// Actually, unless it's darwin, the procName itself is returned.
func ProcName(_ int, procName string) string {
	return procName // from /proc/_/stat, may be shortened
}

type ExtraMetricRAM struct {
	Used     metrics.Gauge
	Buffered metrics.Gauge
	Cached   metrics.Gauge
}

func NewMetricRAM(r metrics.Registry) *MetricRAM {
	return ExtraNewMetricRAM(r, &ExtraMetricRAM{
		Used:     metrics.NewRegisteredGauge("memory.memory-used", r),
		Buffered: metrics.NewRegisteredGauge("memory.memory-buffered", r),
		Cached:   metrics.NewRegisteredGauge("memory.memory-cached", r),
	})
}

func (emr *ExtraMetricRAM) UpdateRAM(got sigar.Mem, extra1, extra2 uint64) {
	emr.Used.Update(int64(got.ActualUsed))
	emr.Buffered.Update(int64(extra1))
	emr.Cached.Update(int64(extra2))
}

/* **************************************************************** */

type ExtraMetricCPU struct {
	Wait    *GaugePercent
	Irq     *GaugePercent
	SoftIrq *GaugePercent
	Stolen  *GaugePercent
}

func (emc *ExtraMetricCPU) UpdateCPU(sigarCpu sigar.Cpu, totalDelta int64) {
	emc.Wait.UpdatePercent(totalDelta, sigarCpu.Wait)
	emc.Irq.UpdatePercent(totalDelta, sigarCpu.Irq)
	emc.SoftIrq.UpdatePercent(totalDelta, sigarCpu.SoftIrq)
	emc.Stolen.UpdatePercent(totalDelta, sigarCpu.Stolen)
}

func NewMetricCPU(r metrics.Registry, name string) *MetricCPU {
	return ExtraNewMetricCPU(r, name, &ExtraMetricCPU{
		Wait:    NewGaugePercent(name+".wait", r),
		Irq:     NewGaugePercent(name+".interrupt", r),
		SoftIrq: NewGaugePercent(name+".softirq", r),
		Stolen:  NewGaugePercent(name+".steal", r),
	})
}
