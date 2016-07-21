package ostent

import (
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/system"
)

// DFSort is to facilitate DFSlice sorting.
type DFSort struct {
	Dfk     *params.Num
	DFSlice DFSlice
}

// Len, Swap and Less satisfy sorting interface.
func (ds DFSort) Len() int      { return len(ds.DFSlice) }
func (ds DFSort) Swap(i, j int) { ds.DFSlice[i], ds.DFSlice[j] = ds.DFSlice[j], ds.DFSlice[i] }
func (ds DFSort) Less(i, j int) (r bool) {
	if match, isa, cmpr := dfCmp(ds.Dfk.Absolute, ds.DFSlice[i], ds.DFSlice[j]); match {
		ds.Dfk.Alpha, r = isa, cmpr
	}
	if ds.Dfk.Negative {
		return !r
	}
	return r
}

func dfCmp(k int, a, b *system.MetricDF) (bool, bool, bool) {
	switch k {
	case params.FS:
		return true, true, a.DevName.Snapshot().Value() < b.DevName.Snapshot().Value()
	case params.MP:
		return true, true, a.DirName.Snapshot().Value() < b.DirName.Snapshot().Value()

	case params.TOTAL:
		return true, false, a.Total.Snapshot().Value() > b.Total.Snapshot().Value()
	case params.USED:
		return true, false, a.Used.Snapshot().Value() > b.Used.Snapshot().Value()
	case params.AVAIL:
		return true, false, a.Avail.Snapshot().Value() > b.Avail.Snapshot().Value()
	case params.USEPCT:
		var (
			vau = a.Used.Snapshot().Value()  // float64
			vbu = b.Used.Snapshot().Value()  // float64
			vat = a.Total.Snapshot().Value() // int64
			vbt = b.Total.Snapshot().Value() // int64
		)
		return true, false, (vau / float64(vat)) > (vbu / float64(vbt))
	}
	return false, false, false
}

// PSSort is to facilitate PSSlice sorting.
type PSSort struct {
	Psk     *params.Num
	PSSlice PSSlice
	UIDs    map[uint]string
}

// Len, Swap and Less satisfy sorting interface.
func (ps PSSort) Len() int      { return len(ps.PSSlice) }
func (ps PSSort) Swap(i, j int) { ps.PSSlice[i], ps.PSSlice[j] = ps.PSSlice[j], ps.PSSlice[i] }
func (ps PSSort) Less(i, j int) (r bool) {
	k, a, b := ps.Psk.Absolute, ps.PSSlice[i], ps.PSSlice[j]
	if match, isa, cmpr := psCmp(k, a, b); match {
		ps.Psk.Alpha, r = isa, cmpr
	} else if k == params.USER {
		ps.Psk.Alpha, r = true, username(ps.UIDs, a.UID) < username(ps.UIDs, b.UID)
	}
	if ps.Psk.Negative {
		return !r
	}
	return r
}

func psCmp(k int, a, b *system.PSInfo) (bool, bool, bool) {
	switch k {
	case params.PID:
		return true, false, a.PID > b.PID
	case params.PRI:
		return true, false, a.Priority > b.Priority
	case params.NICE:
		return true, false, a.Nice > b.Nice
	case params.VIRT:
		return true, false, a.Size > b.Size
	case params.RES:
		return true, false, a.Resident > b.Resident
	case params.TIME:
		return true, false, a.Time > b.Time
	case params.UID:
		return true, false, a.UID > b.UID

	case params.NAME: // alpha
		return true, true, a.Name < b.Name
	}
	return false, false, false
}
