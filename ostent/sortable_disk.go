package ostent

import (
	"encoding/json"

	"github.com/ostrost/ostent/client"
	"github.com/ostrost/ostent/types"
)

// SortCritDisk is a distinct types.SeqNReverse type.
type SortCritDisk types.SeqNReverse

// LessDisk is a 'less' func for types.MetricDF comparison.
func (crit SortCritDisk) LessDisk(a, b types.MetricDF) bool {
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
