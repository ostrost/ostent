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

type Links client.Linkattrs // TODO MAYBE become a struct with inline map[string]client.Attr and prefill for MarshalJSON

func (la Links) DiskName() client.Attr { return client.Linkattrs(la).Attr(client.DFFS) }
func (la Links) Total() client.Attr    { return client.Linkattrs(la).Attr(client.DFSIZE) }
func (la Links) Used() client.Attr     { return client.Linkattrs(la).Attr(client.DFUSED) }
func (la Links) Avail() client.Attr    { return client.Linkattrs(la).Attr(client.DFAVAIL) }
func (la Links) DirName() client.Attr  { return client.Linkattrs(la).Attr(client.DFMP) }

// MarshalJSON satisfying json.Marshaler interface.
func (la Links) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]client.Attr{
		"DiskName": la.DiskName(),
		"Total":    la.Total(),
		"Used":     la.Used(),
		"Avail":    la.Avail(),
		"DirName":  la.DirName(),
	})
}

type PSlinks client.Linkattrs

func (la PSlinks) PID() client.Attr      { return client.Linkattrs(la).Attr(client.PSPID) }
func (la PSlinks) Priority() client.Attr { return client.Linkattrs(la).Attr(client.PSPRI) }
func (la PSlinks) Nice() client.Attr     { return client.Linkattrs(la).Attr(client.PSNICE) }
func (la PSlinks) Time() client.Attr     { return client.Linkattrs(la).Attr(client.PSTIME) }
func (la PSlinks) Name() client.Attr     { return client.Linkattrs(la).Attr(client.PSNAME) }
func (la PSlinks) User() client.Attr     { return client.Linkattrs(la).Attr(client.PSUID) }
func (la PSlinks) Size() client.Attr     { return client.Linkattrs(la).Attr(client.PSSIZE) }
func (la PSlinks) Resident() client.Attr { return client.Linkattrs(la).Attr(client.PSRES) }

func (la PSlinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]client.Attr{
		"PID":      la.PID(),
		"Priority": la.Priority(),
		"Nice":     la.Nice(),
		"Time":     la.Time(),
		"Name":     la.Name(),
		"User":     la.User(),
		"Size":     la.Size(),
		"Resident": la.Resident(),
	})
}

type PStable struct {
	List []operating.ProcData `json:",omitempty"`
}
