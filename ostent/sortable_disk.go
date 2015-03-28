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
