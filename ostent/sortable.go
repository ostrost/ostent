package ostent

import (
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/params/enums"
)

// DFSort is to facilitate DFSlice sorting.
type DFSort struct {
	Dfk     *params.Num
	DFSlice DFSlice
}

// Len, Swap and Less satisfy sorting interface.
func (ds DFSort) Len() int      { return len(ds.DFSlice) }
func (ds DFSort) Swap(i, j int) { ds.DFSlice[i], ds.DFSlice[j] = ds.DFSlice[j], ds.DFSlice[i] }
func (ds DFSort) Less(i, j int) bool {
	if a, b := ds.DFSlice[i], ds.DFSlice[j]; true {
		r := false
		switch ds.Dfk.Absolute {
		case enums.FS:
			ds.Dfk.Alpha = true
			r = a.DevName.Snapshot().Value() < b.DevName.Snapshot().Value()
		case enums.TOTAL:
			r = a.Total.Snapshot().Value() > b.Total.Snapshot().Value()
		case enums.USED:
			r = a.Used.Snapshot().Value() > b.Used.Snapshot().Value()
		case enums.AVAIL:
			r = a.Avail.Snapshot().Value() > b.Avail.Snapshot().Value()
		case enums.MP:
			ds.Dfk.Alpha = true
			r = a.DirName.Snapshot().Value() < b.DirName.Snapshot().Value()
		}
		if ds.Dfk.Negative {
			return !r
		}
		return r
	}
	return false
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
func (ps PSSort) Less(i, j int) bool {
	if a, b := ps.PSSlice[i], ps.PSSlice[j]; true {
		r := false
		switch ps.Psk.Absolute {
		case enums.PID:
			r = a.PID > b.PID
		case enums.PRI:
			r = a.Priority > b.Priority
		case enums.NICE:
			r = a.Nice > b.Nice
		case enums.VIRT:
			r = a.Size > b.Size
		case enums.RES:
			r = a.Resident > b.Resident
		case enums.TIME:
			r = a.Time > b.Time
		case enums.NAME:
			ps.Psk.Alpha = true
			r = a.Name < b.Name
		case enums.UID:
			r = a.UID > b.UID
		case enums.USER:
			ps.Psk.Alpha = true
			r = username(ps.UIDs, a.UID) < username(ps.UIDs, b.UID)
		}
		if ps.Psk.Negative {
			return !r
		}
		return r
	}
	return false
}
