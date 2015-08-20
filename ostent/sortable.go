package ostent

import (
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/params/enums"
)

// DiskSort is to facilitate DiskSlice sorting.
type DiskSort struct {
	Dfk       *params.Num
	DiskSlice DiskSlice
}

// Len, Swap and Less satisfy sorting interface.
func (ds DiskSort) Len() int      { return len(ds.DiskSlice) }
func (ds DiskSort) Swap(i, j int) { ds.DiskSlice[i], ds.DiskSlice[j] = ds.DiskSlice[j], ds.DiskSlice[i] }
func (ds DiskSort) Less(i, j int) bool {
	if a, b := ds.DiskSlice[i], ds.DiskSlice[j]; true {
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

// ProcSort is to facilitate ProcSlice sorting.
type ProcSort struct {
	Psk       *params.Num
	ProcSlice ProcSlice
	UIDs      map[uint]string
}

// Len, Swap and Less satisfy sorting interface.
func (ps ProcSort) Len() int      { return len(ps.ProcSlice) }
func (ps ProcSort) Swap(i, j int) { ps.ProcSlice[i], ps.ProcSlice[j] = ps.ProcSlice[j], ps.ProcSlice[i] }
func (ps ProcSort) Less(i, j int) bool {
	if a, b := ps.ProcSlice[i], ps.ProcSlice[j]; true {
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
