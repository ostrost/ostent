package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/system/operating"
)

// LessDiskFunc makes a 'less' func for operating.MetricDF comparison.
func LessDiskFunc(num client.Number) func(operating.MetricDF, operating.MetricDF) bool {
	return func(a, b operating.MetricDF) bool {
		r := false
		switch client.UintDF(num.Uint) {
		case client.FS:
			r = a.DevName.Snapshot().Value() < b.DevName.Snapshot().Value()
		case client.TOTAL:
			r = a.Total.Snapshot().Value() < b.Total.Snapshot().Value()
		case client.USED:
			r = a.Used.Snapshot().Value() < b.Used.Snapshot().Value()
		case client.AVAIL:
			r = a.Avail.Snapshot().Value() < b.Avail.Snapshot().Value()
		case client.MP:
			r = a.DirName.Snapshot().Value() < b.DirName.Snapshot().Value()
		}
		// numeric values: reverse "less"
		if !client.DF.IsAlpha(num.Uint) {
			r = !r
		}
		if num.Negative {
			r = !r
		}
		return r
	}
}

// LessProcFunc makes a 'less' func for operating.MetricProc comparison.
func LessProcFunc(uids map[uint]string, num client.Number) func(operating.MetricProc, operating.MetricProc) bool {
	return func(a, b operating.MetricProc) bool {
		r := false
		switch client.UintPS(num.Uint) {
		case client.PID:
			r = a.PID < b.PID
		case client.PRI:
			r = a.Priority < b.Priority
		case client.NICE:
			r = a.Nice < b.Nice
		case client.VIRT:
			r = a.Size < b.Size
		case client.RES:
			r = a.Resident < b.Resident
		case client.TIME:
			r = a.Time < b.Time
		case client.NAME:
			r = a.Name < b.Name
		case client.UID:
			r = a.UID < b.UID
		case client.USER:
			r = username(uids, a.UID) < username(uids, b.UID)
		}
		// numeric values: reverse "less"
		if !client.PS.IsAlpha(num.Uint) {
			r = !r
		}
		if num.Negative {
			r = !r
		}
		return r
	}
}

type Links struct {
	client.Links
	// TODO MAYBE add inline map[string]client.Attr and prefill for MarshalJSON
}

// DF part

func (la Links) DFFS() client.Attr    { return la.EncodeNU("df", client.FS) }
func (la Links) DFMP() client.Attr    { return la.EncodeNU("df", client.MP) }
func (la Links) DFTOTAL() client.Attr { return la.EncodeNU("df", client.TOTAL) }
func (la Links) DFUSED() client.Attr  { return la.EncodeNU("df", client.USED) }
func (la Links) DFAVAIL() client.Attr { return la.EncodeNU("df", client.AVAIL) }

// PS part

func (la Links) PSPID() client.Attr  { return la.EncodeNU("ps", client.PID) }
func (la Links) PSPRI() client.Attr  { return la.EncodeNU("ps", client.PRI) }
func (la Links) PSNICE() client.Attr { return la.EncodeNU("ps", client.NICE) }
func (la Links) PSTIME() client.Attr { return la.EncodeNU("ps", client.TIME) }
func (la Links) PSNAME() client.Attr { return la.EncodeNU("ps", client.NAME) }
func (la Links) PSUID() client.Attr  { return la.EncodeNU("ps", client.UID) }
func (la Links) PSUSER() client.Attr { return la.EncodeNU("ps", client.USER) }
func (la Links) PSVIRT() client.Attr { return la.EncodeNU("ps", client.VIRT) }
func (la Links) PSRES() client.Attr  { return la.EncodeNU("ps", client.RES) }

// MarshalJSON satisfying json.Marshaler interface.
func (la Links) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]client.Attr{
		// DF
		"DFFS":    la.DFFS(),
		"DFMP":    la.DFMP(),
		"DFTOTAL": la.DFTOTAL(),
		"DFUSED":  la.DFUSED(),
		"DFAVAIL": la.DFAVAIL(),

		// PS
		"PSPID":  la.PSPID(),
		"PSPRI":  la.PSPRI(),
		"PSNICE": la.PSNICE(),
		"PSTIME": la.PSTIME(),
		"PSNAME": la.PSNAME(),
		"PSUID":  la.PSUID(),
		"PSUSER": la.PSUSER(),
		"PSVIRT": la.PSVIRT(),
		"PSRES":  la.PSRES(),
	})
}

type PStable struct {
	List []operating.ProcData `json:",omitempty"`
}
