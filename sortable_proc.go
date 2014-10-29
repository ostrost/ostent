package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/types"
)

type procOrder struct {
	procs   []types.ProcInfo
	seq     types.SEQ
	reverse bool
}

func (po procOrder) Len() int {
	return len(po.procs)
}

func (po procOrder) Swap(i, j int) {
	po.procs[i], po.procs[j] = po.procs[j], po.procs[i]
}

func (po procOrder) Less(i, j int) bool {
	t := false
	switch po.seq {
	case client.PSPID, -client.PSPID:
		t = po.seq.Sign(po.procs[i].PID < po.procs[j].PID)
	case client.PSPRI, -client.PSPRI:
		t = po.seq.Sign(po.procs[i].Priority < po.procs[j].Priority)
	case client.PSNICE, -client.PSNICE:
		t = po.seq.Sign(po.procs[i].Nice < po.procs[j].Nice)
	case client.PSSIZE, -client.PSSIZE:
		t = po.seq.Sign(po.procs[i].Size < po.procs[j].Size)
	case client.PSRES, -client.PSRES:
		t = po.seq.Sign(po.procs[i].Resident < po.procs[j].Resident)
	case client.PSTIME, -client.PSTIME:
		t = po.seq.Sign(po.procs[i].Time < po.procs[j].Time)
	case client.PSNAME, -client.PSNAME:
		t = po.seq.Sign(po.procs[i].Name < po.procs[j].Name)
	case client.PSUID, -client.PSUID:
		t = po.seq.Sign(po.procs[i].UID < po.procs[j].UID)
	}
	if po.reverse {
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
