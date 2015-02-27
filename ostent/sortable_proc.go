// Package ostent is the library part of ostent cmd.
package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/types"
)

// SortCritProc is a distinct types.SeqNReverse type.
type SortCritProc types.SeqNReverse

// LessProc is a 'less' func for types.MetricProc comparison.
func (crit SortCritProc) LessProc(a, b types.MetricProc) bool {
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

type PSlinks types.Linkattrs

func (la PSlinks) PID() types.Attr      { return types.Linkattrs(la).Attr(client.PSPID) }
func (la PSlinks) Priority() types.Attr { return types.Linkattrs(la).Attr(client.PSPRI) }
func (la PSlinks) Nice() types.Attr     { return types.Linkattrs(la).Attr(client.PSNICE) }
func (la PSlinks) Time() types.Attr     { return types.Linkattrs(la).Attr(client.PSTIME) }
func (la PSlinks) Name() types.Attr     { return types.Linkattrs(la).Attr(client.PSNAME) }
func (la PSlinks) User() types.Attr     { return types.Linkattrs(la).Attr(client.PSUID) }
func (la PSlinks) Size() types.Attr     { return types.Linkattrs(la).Attr(client.PSSIZE) }
func (la PSlinks) Resident() types.Attr { return types.Linkattrs(la).Attr(client.PSRES) }

func (la PSlinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]types.Attr{
		"PID":      types.Linkattrs(la).Attr(client.PSPID),
		"Priority": types.Linkattrs(la).Attr(client.PSPRI),
		"Nice":     types.Linkattrs(la).Attr(client.PSNICE),
		"Time":     types.Linkattrs(la).Attr(client.PSTIME),
		"Name":     types.Linkattrs(la).Attr(client.PSNAME),
		"User":     types.Linkattrs(la).Attr(client.PSUID),
		"Size":     types.Linkattrs(la).Attr(client.PSSIZE),
		"Resident": types.Linkattrs(la).Attr(client.PSRES),
	})
}

type PStable struct {
	List []types.ProcData `json:",omitempty"`
}
