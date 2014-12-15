package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/types"
)

type diskOrder struct {
	disks   []MetricDF
	seq     types.SEQ
	reverse bool
}

func (do diskOrder) Len() int {
	return len(do.disks)
}

func (do diskOrder) Swap(i, j int) {
	do.disks[i], do.disks[j] = do.disks[j], do.disks[i]
}

func (do diskOrder) Less(i, j int) bool {
	t := false
	switch do.seq {
	case client.DFFS, -client.DFFS:
		t = do.seq.Sign(do.disks[i].DevName.Snapshot().Value() < do.disks[j].DevName.Snapshot().Value())
	case client.DFSIZE, -client.DFSIZE:
		t = do.seq.Sign(do.disks[i].Total.Snapshot().Value() < do.disks[j].Total.Snapshot().Value())
	case client.DFUSED, -client.DFUSED:
		t = do.seq.Sign(do.disks[i].Used.Snapshot().Value() < do.disks[j].Used.Snapshot().Value())
	case client.DFAVAIL, -client.DFAVAIL:
		t = do.seq.Sign(do.disks[i].Avail.Snapshot().Value() < do.disks[j].Avail.Snapshot().Value())
	case client.DFMP, -client.DFMP:
		t = do.seq.Sign(do.disks[i].DirName.Snapshot().Value() < do.disks[j].DirName.Snapshot().Value())
	}
	if do.reverse {
		return !t
	}
	return t
}

type DFlinks types.Linkattrs

func (la DFlinks) DiskName() types.Attr { return types.Linkattrs(la).Attr(client.DFFS) }
func (la DFlinks) Total() types.Attr    { return types.Linkattrs(la).Attr(client.DFSIZE) }
func (la DFlinks) Used() types.Attr     { return types.Linkattrs(la).Attr(client.DFUSED) }
func (la DFlinks) Avail() types.Attr    { return types.Linkattrs(la).Attr(client.DFAVAIL) }
func (la DFlinks) DirName() types.Attr  { return types.Linkattrs(la).Attr(client.DFMP) }

func (la DFlinks) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]types.Attr{
		"DiskName": types.Linkattrs(la).Attr(client.DFFS),
		"Total":    types.Linkattrs(la).Attr(client.DFSIZE),
		"Used":     types.Linkattrs(la).Attr(client.DFUSED),
		"Avail":    types.Linkattrs(la).Attr(client.DFAVAIL),
		"DirName":  types.Linkattrs(la).Attr(client.DFMP),
	})
}
