package ostential
import (
	"ostential/types"
)

type clientState struct {
	   HideMemory *bool `json:",omitempty"`

	  HideNetwork *bool `json:",omitempty"`
	ExpandNetwork *bool `json:",omitempty"`

	      HideCPU *bool `json:",omitempty"`
	    ExpandCPU *bool `json:",omitempty"`

	    HideDisks *bool `json:",omitempty"`
	  ExpandDisks *bool `json:",omitempty"`

	HideProcesses *bool `json:",omitempty"`

	CurrentNetworkTab *types.SEQ `json:",omitempty"`
	CurrentDisksTab   *types.SEQ `json:",omitempty"`

	NetworkTabs *networkTabs `json:",omitempty"` // immutable
	  DisksTabs   *disksTabs `json:",omitempty"` // immutable

	// UserProcesses string `json:omitempty`
}

type disksTabs struct {
	DisksinBytes  types.SEQ
	DisksinInodes types.SEQ
}
type networkTabs struct {
	NetworkinPackets  types.SEQ
	NetworkinErrors   types.SEQ
	NetworkinBytes    types.SEQ
}

func(nt *networkTabs) merge(src *networkTabs) { if src != nil { *nt = *src } }
func(dt *disksTabs)   merge(src *disksTabs)   { if src != nil { *dt = *src } }

func(_  clientState) merge_bool(dest, src *bool)    { if src != nil { *dest = *src } }
func(_  clientState) mergeSEQ(dest, src *types.SEQ) { if src != nil { *dest = *src } }

func(cs *clientState) Merge(ps clientState) {
	cs.merge_bool(cs.HideMemory,      ps.HideMemory)
	cs.merge_bool(cs.HideNetwork,     ps.HideNetwork)
	cs.merge_bool(cs.ExpandNetwork,   ps.ExpandNetwork)
	cs.merge_bool(cs.HideCPU,         ps.HideCPU)
	cs.merge_bool(cs.ExpandCPU,       ps.ExpandCPU)
	cs.merge_bool(cs.HideDisks,       ps.HideDisks)
	cs.merge_bool(cs.ExpandDisks,     ps.ExpandDisks)
	cs.merge_bool(cs.HideProcesses,   ps.HideProcesses)
	cs.mergeSEQ(cs.CurrentNetworkTab, ps.CurrentNetworkTab)
	cs.mergeSEQ(cs.CurrentDisksTab,   ps.CurrentDisksTab)
	cs.NetworkTabs.merge(ps.NetworkTabs)
	cs.DisksTabs  .merge(ps.DisksTabs)
}

const (
	____NTABID types.SEQ = iota
	NPACKETS_TABID
	 NERRORS_TABID
	  NBYTES_TABID
)
var NETWORK_TABS = []types.SEQ{
	NPACKETS_TABID,
	NERRORS_TABID,
	NBYTES_TABID,
}
const (
	____DTABID types.SEQ = iota
	DINODES_TABID
	 DBYTES_TABID
)
var DISKS_TABS = []types.SEQ{
	DINODES_TABID,
	DBYTES_TABID,
}

func defaultClientState() clientState {
	cs := clientState{}

	cs.HideMemory    = new(bool)
	cs.HideNetwork   = new(bool)
	cs.ExpandNetwork = new(bool)
	cs.HideCPU       = new(bool)
	cs.ExpandCPU     = new(bool)
	cs.HideDisks     = new(bool)
	cs.ExpandDisks   = new(bool)
	cs.HideProcesses = new(bool)

	cs.CurrentNetworkTab = new(types.SEQ)
	cs.CurrentDisksTab   = new(types.SEQ)
	*cs.CurrentNetworkTab = NBYTES_TABID
	*cs.CurrentDisksTab   = DBYTES_TABID

	 cs.DisksTabs = &disksTabs{ // immutable
		 DisksinBytes:  DBYTES_TABID,
		 DisksinInodes: DINODES_TABID,
	 }

	cs.NetworkTabs = &networkTabs{ // immutable
		NetworkinPackets: NPACKETS_TABID,
		NetworkinErrors:  NERRORS_TABID,
		NetworkinBytes:   NBYTES_TABID,
	}
	// cs.UserProcesses = "" // default

	return cs
}
