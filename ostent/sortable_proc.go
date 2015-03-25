// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/system/operating"
)

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
