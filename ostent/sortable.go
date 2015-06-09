package ostent

import (
	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/client/enums"
	"github.com/ostrost/ostent/system/operating"
)

// LessDiskFunc makes a 'less' func for operating.MetricDF comparison.
func LessDiskFunc(param client.EnumParam) func(operating.MetricDF, operating.MetricDF) bool {
	return func(a, b operating.MetricDF) bool {
		r := false
		switch enums.UintDF(param.Number.Uint) {
		case enums.FS:
			r = a.DevName.Snapshot().Value() < b.DevName.Snapshot().Value()
		case enums.TOTAL:
			r = a.Total.Snapshot().Value() < b.Total.Snapshot().Value()
		case enums.USED:
			r = a.Used.Snapshot().Value() < b.Used.Snapshot().Value()
		case enums.AVAIL:
			r = a.Avail.Snapshot().Value() < b.Avail.Snapshot().Value()
		case enums.MP:
			r = a.DirName.Snapshot().Value() < b.DirName.Snapshot().Value()
		}
		return param.LessorMore(r)
	}
}

// LessProcFunc makes a 'less' func for operating.MetricProc comparison.
func LessProcFunc(uids map[uint]string, param client.EnumParam) func(operating.MetricProc, operating.MetricProc) bool {
	return func(a, b operating.MetricProc) bool {
		r := false
		switch enums.UintPS(param.Number.Uint) {
		case enums.PID:
			r = a.PID < b.PID
		case enums.PRI:
			r = a.Priority < b.Priority
		case enums.NICE:
			r = a.Nice < b.Nice
		case enums.VIRT:
			r = a.Size < b.Size
		case enums.RES:
			r = a.Resident < b.Resident
		case enums.TIME:
			r = a.Time < b.Time
		case enums.NAME:
			r = a.Name < b.Name
		case enums.UID:
			r = a.UID < b.UID
		case enums.USER:
			r = username(uids, a.UID) < username(uids, b.UID)
		}
		return param.LessorMore(r)
	}
}
