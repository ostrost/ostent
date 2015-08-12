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
		switch by.Body {
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
		if by.Head {
			return !r
		}
		return r
	}
}

// LessProcFunc makes a 'less' func for operating.MetricProc comparison.
func LessProcFunc(by *params.Num, uids map[uint]string) func(operating.MetricProc, operating.MetricProc) bool {
	return func(a, b operating.MetricProc) bool {
		r := false
		switch by.Body {
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
			by.Alpha = true
			r = a.Name < b.Name
		case enums.UID:
			r = a.UID > b.UID
		case enums.USER:
			by.Alpha = true
			r = username(uids, a.UID) < username(uids, b.UID)
		}
		if by.Head {
			return !r
		}
		return r
	}
}
