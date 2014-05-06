package ostential
import (
	"ostential/types"

	"encoding/json"
)

type procOrder struct {
	procs []types.ProcInfo
	seq types.SEQ
	reverse bool
}
func(po procOrder) Len() int {
	return len(po.procs)
}
func(po procOrder) Swap(i, j int) {
	po.procs[i], po.procs[j] = po.procs[j], po.procs[i]
}
func(po procOrder) Less(i, j int) bool {
	t := false
	switch po.seq {
	case PSPID,  -PSPID:  t = po.seq.Sign(po.procs[i].PID      < po.procs[j].PID)
	case PSPRI,  -PSPRI:  t = po.seq.Sign(po.procs[i].Priority < po.procs[j].Priority)
	case PSNICE, -PSNICE: t = po.seq.Sign(po.procs[i].Nice     < po.procs[j].Nice)
	case PSSIZE, -PSSIZE: t = po.seq.Sign(po.procs[i].Size     < po.procs[j].Size)
	case PSRES,  -PSRES:  t = po.seq.Sign(po.procs[i].Resident < po.procs[j].Resident)
	case PSTIME, -PSTIME: t = po.seq.Sign(po.procs[i].Time     < po.procs[j].Time)
	case PSNAME, -PSNAME: t = po.seq.Sign(po.procs[i].Name     < po.procs[j].Name)
	case PSUID,  -PSUID:  t = po.seq.Sign(po.procs[i].Uid      < po.procs[j].Uid)
	}
	if po.reverse {
		return !t
	}
	return t
}
const (
____PSIOTA		types.SEQ = iota
	PSPID
    PSPRI
    PSNICE
    PSSIZE
    PSRES
	PSTIME
	PSNAME
	PSUID
)

type ProcLinkattrs types.Linkattrs
func(la ProcLinkattrs) PID()      types.Attr { return types.Linkattrs(la).Attr(PSPID ); }
func(la ProcLinkattrs) Priority() types.Attr { return types.Linkattrs(la).Attr(PSPRI ); }
func(la ProcLinkattrs) Nice()     types.Attr { return types.Linkattrs(la).Attr(PSNICE); }
func(la ProcLinkattrs) Time()     types.Attr { return types.Linkattrs(la).Attr(PSTIME); }
func(la ProcLinkattrs) Name()     types.Attr { return types.Linkattrs(la).Attr(PSNAME); }
func(la ProcLinkattrs) User()     types.Attr { return types.Linkattrs(la).Attr(PSUID ); }
func(la ProcLinkattrs) Size()     types.Attr { return types.Linkattrs(la).Attr(PSSIZE); }
func(la ProcLinkattrs) Resident() types.Attr { return types.Linkattrs(la).Attr(PSRES ); }

func(la ProcLinkattrs) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]types.Attr{
		"PID":      types.Linkattrs(la).Attr(PSPID),
		"Priority": types.Linkattrs(la).Attr(PSPRI),
		"Nice":     types.Linkattrs(la).Attr(PSNICE),
		"Time":     types.Linkattrs(la).Attr(PSTIME),
		"Name":     types.Linkattrs(la).Attr(PSNAME),
		"User":     types.Linkattrs(la).Attr(PSUID),
		"Size":     types.Linkattrs(la).Attr(PSSIZE),
		"Resident": types.Linkattrs(la).Attr(PSRES),
	})
}

type ProcTable struct {
	List  []types.ProcData
	Links *ProcLinkattrs `json:",omitempty"`
	MoreText      string `json:",omitempty"` // should never be empty, sanity check
	NotExpandable  *bool `json:",omitempty"`
}
