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
		case client.DFSIZE:
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
func LessProcFunc(num client.Number) func(operating.MetricProc, operating.MetricProc) bool {
	return func(a, b operating.MetricProc) bool {
		r := false
		switch client.UintPS(num.Uint) {
		case client.PID:
			r = a.PID < b.PID
		case client.PRI:
			r = a.Priority < b.Priority
		case client.NICE:
			r = a.Nice < b.Nice
		case client.PSSIZE:
			r = a.Size < b.Size
		case client.RES:
			r = a.Resident < b.Resident
		case client.TIME:
			r = a.Time < b.Time
		case client.NAME:
			r = a.Name < b.Name
		case client.UID:
			r = a.UID < b.UID
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

func (la Links) DFdiskName() client.Attr { return la.EncodeNU("df", client.FS) }
func (la Links) DFdirName() client.Attr  { return la.EncodeNU("df", client.MP) }
func (la Links) DFtotal() client.Attr    { return la.EncodeNU("df", client.DFSIZE) }
func (la Links) DFused() client.Attr     { return la.EncodeNU("df", client.USED) }
func (la Links) DFavail() client.Attr    { return la.EncodeNU("df", client.AVAIL) }

// PS part

func (la Links) PSPID() client.Attr      { return la.EncodeNU("ps", client.PID) }
func (la Links) PSpriority() client.Attr { return la.EncodeNU("ps", client.PRI) }
func (la Links) PSnice() client.Attr     { return la.EncodeNU("ps", client.NICE) }
func (la Links) PStime() client.Attr     { return la.EncodeNU("ps", client.TIME) }
func (la Links) PSname() client.Attr     { return la.EncodeNU("ps", client.NAME) }
func (la Links) PSuser() client.Attr     { return la.EncodeNU("ps", client.UID) }
func (la Links) PSsize() client.Attr     { return la.EncodeNU("ps", client.PSSIZE) }
func (la Links) PSresident() client.Attr { return la.EncodeNU("ps", client.RES) }

// MarshalJSON satisfying json.Marshaler interface.
func (la Links) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]client.Attr{
		// DF
		"DFdiskName": la.DFdiskName(),
		"DFdirName":  la.DFdirName(),
		"DFtotal":    la.DFtotal(),
		"DFused":     la.DFused(),
		"DFavail":    la.DFavail(),

		// PS
		"PSPID":      la.PSPID(),
		"PSpriority": la.PSpriority(),
		"PSnice":     la.PSnice(),
		"PStime":     la.PStime(),
		"PSname":     la.PSname(),
		"PSuser":     la.PSuser(),
		"PSsize":     la.PSsize(),
		"PSresident": la.PSresident(),
	})
}

type PStable struct {
	List []operating.ProcData `json:",omitempty"`
}
