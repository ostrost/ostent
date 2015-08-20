package ostent

import (
	"github.com/ostrost/ostent/params"
	"github.com/ostrost/ostent/params/enums"
	"github.com/ostrost/ostent/system/operating"
)

// LessDiskFunc makes a 'less' func for operating.MetricDF comparison.
func LessDiskFunc(by *params.Num) func(operating.MetricDF, operating.MetricDF) bool {
	return func(a, b operating.MetricDF) bool {
		r := false
		switch by.Absolute {
		case enums.FS:
			by.Alpha = true
			r = a.DevName.Snapshot().Value() < b.DevName.Snapshot().Value()
		case enums.TOTAL:
			r = a.Total.Snapshot().Value() > b.Total.Snapshot().Value()
		case enums.USED:
			r = a.Used.Snapshot().Value() > b.Used.Snapshot().Value()
		case enums.AVAIL:
			r = a.Avail.Snapshot().Value() > b.Avail.Snapshot().Value()
		case enums.MP:
			by.Alpha = true
			r = a.DirName.Snapshot().Value() < b.DirName.Snapshot().Value()
		}
		if by.Negative {
			return !r
		}
		return r
	}
}

// ProcSort is to facilitate ProcSlice sorting.
type ProcSort struct {
	Psn       *params.Num
	ProcSlice ProcSlice
	UIDs      map[uint]string
}

func (ps ProcSort) Len() int { return len(ps.ProcSlice) }

func (ps ProcSort) Swap(i, j int) { ps.ProcSlice[i], ps.ProcSlice[j] = ps.ProcSlice[j], ps.ProcSlice[i] }

// Less is sorting interface.
func (ps ProcSort) Less(i, j int) bool {
	if a, b := ps.ProcSlice[i], ps.ProcSlice[j]; true {
		r := false
		switch ps.Psn.Absolute {
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
			ps.Psn.Alpha = true
			r = a.Name < b.Name
		case enums.UID:
			r = a.UID > b.UID
		case enums.USER:
			ps.Psn.Alpha = true
			r = username(ps.UIDs, a.UID) < username(ps.UIDs, b.UID)
		}
		if ps.Psn.Negative {
			return !r
		}
		return r
	}
	return false
}
