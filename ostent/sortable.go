// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/system/operating"
)

// SortCritDisk is a distinct client.SeqNReverse type.
type SortCritDisk client.SeqNReverse

// LessDisk is a 'less' func for operating.MetricDF comparison.
func (crit SortCritDisk) LessDisk(a, b operating.MetricDF) bool {
	t := false
	switch crit.SEQ {
	case client.DFFS, -client.DFFS:
		t = crit.SEQ.Sign(a.DevName.Snapshot().Value() < b.DevName.Snapshot().Value())
	case client.DFSIZE, -client.DFSIZE:
		t = crit.SEQ.Sign(a.Total.Snapshot().Value() < b.Total.Snapshot().Value())
	case client.DFUSED, -client.DFUSED:
		t = crit.SEQ.Sign(a.Used.Snapshot().Value() < b.Used.Snapshot().Value())
	case client.DFAVAIL, -client.DFAVAIL:
		t = crit.SEQ.Sign(a.Avail.Snapshot().Value() < b.Avail.Snapshot().Value())
	case client.DFMP, -client.DFMP:
		t = crit.SEQ.Sign(a.DirName.Snapshot().Value() < b.DirName.Snapshot().Value())
	}
	if crit.Reverse {
		return !t
	}
	return t
}

// SortCritProc is a distinct client.SeqNReverse type.
type SortCritProc client.SeqNReverse

// LessProc is a 'less' func for operating.MetricProc comparison.
func (crit SortCritProc) LessProc(a, b operating.MetricProc) bool {
	t := false
	switch crit.SEQ {
	case client.PSPID, -client.PSPID:
		t = crit.SEQ.Sign(a.PID < b.PID)
	case client.PSPRI, -client.PSPRI:
		t = crit.SEQ.Sign(a.Priority < b.Priority)
	case client.PSNICE, -client.PSNICE:
		t = crit.SEQ.Sign(a.Nice < b.Nice)
	case client.PSSIZE, -client.PSSIZE:
		t = crit.SEQ.Sign(a.Size < b.Size)
	case client.PSRES, -client.PSRES:
		t = crit.SEQ.Sign(a.Resident < b.Resident)
	case client.PSTIME, -client.PSTIME:
		t = crit.SEQ.Sign(a.Time < b.Time)
	case client.PSNAME, -client.PSNAME:
		t = crit.SEQ.Sign(a.Name < b.Name)
	case client.PSUID, -client.PSUID:
		t = crit.SEQ.Sign(a.UID < b.UID)
	}
	if crit.Reverse {
		return !t
	}
	return t
}

type Links struct {
	client.Linkattrs
	// TODO MAYBE add inline map[string]client.Attr and prefill for MarshalJSON
}

// DF part

func (la Links) DFdiskName() client.Attr { return la.Attr("df", client.DFFS) }
func (la Links) DFdirName() client.Attr  { return la.Attr("df", client.DFMP) }
func (la Links) DFtotal() client.Attr    { return la.Attr("df", client.DFSIZE) }
func (la Links) DFused() client.Attr     { return la.Attr("df", client.DFUSED) }
func (la Links) DFavail() client.Attr    { return la.Attr("df", client.DFAVAIL) }

// PS part

func (la Links) PSPID() client.Attr      { return la.Attr("ps", client.PSPID) }
func (la Links) PSpriority() client.Attr { return la.Attr("ps", client.PSPRI) }
func (la Links) PSnice() client.Attr     { return la.Attr("ps", client.PSNICE) }
func (la Links) PStime() client.Attr     { return la.Attr("ps", client.PSTIME) }
func (la Links) PSname() client.Attr     { return la.Attr("ps", client.PSNAME) }
func (la Links) PSuser() client.Attr     { return la.Attr("ps", client.PSUID) }
func (la Links) PSsize() client.Attr     { return la.Attr("ps", client.PSSIZE) }
func (la Links) PSresident() client.Attr { return la.Attr("ps", client.PSRES) }

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
