package types
import (
	"net/url"
	"html/template"
)

type SEQ int
func(seq SEQ) AnyOf(list []SEQ) bool {
	for _, s := range list {
		if s == seq {
			return true
		}
	}
	return false
}
func(seq SEQ) Sign(t bool) bool { // used in sortable_*.go
	if seq < 0 {
		return t
	}
	return !t
}

func NewDataMeta() DataMeta {
	return DataMeta{
		Expandable: new(bool),
		More:       new(int),
	}
}

type DataMeta struct {
	Expandable *bool `json:",omitempty"`
	More       *int  `json:",omitempty"`
}

type CPU struct {
	List []Core
	DataMeta
}

type Core struct {
	N    string

	User uint // percent without "%"
	Sys  uint // percent without "%"
	Idle uint // percent without "%"

	UserClass string
	 SysClass string
	IdleClass string
}

type DiskMeta struct {
	DiskNameHTML template.HTML
	DirNameHTML  template.HTML
	DirNameKey string
}

type DiskBytes struct {
	DiskMeta
	Total       string // with units
	Used        string // with units
	Avail       string // with units
	UsePercent  string // as a string, with "%"
	UsePercentClass string
}

type DiskInodes struct {
	DiskMeta
	Inodes      string // with units
	Iused       string // with units
	Ifree       string // with units
	IusePercent string // as a string, with "%"
	IusePercentClass string
}

type DisksinBytes struct {
	List []DiskBytes
}
type DisksinInodes struct {
	List []DiskInodes
}

// type DiskTable struct {
// 	List  []DiskData
// 	Links *DiskLinkattrs `json:",omitempty"`
// 	HaveCollapsed bool
// }

type Attr struct {
	Href, Class string
}
func(la Linkattrs) Attr(seq SEQ) Attr {
	base := url.Values{}
	for k, v := range la.Base {
		base[k] = v
	}
	ascp := la._attr(base, seq)
	class := "state"
	if ascp != nil {
		class += " "+ map[bool]string{true: "asc", false: "desc"}[*ascp]
	}
	return Attr{
		Href:  "?" + base.Encode(),
		Class: class,
	}
}

func(la Linkattrs) _attr(base url.Values, seq SEQ) *bool {
	unlessreverse := func(t bool) *bool {
		if la.Bimap.SEQ2REVERSE[seq] {
			t = !t
		}
		return &t
	}

	if la.Pname == "" {
		if seq == la.Bimap.Default_seq {
			return unlessreverse(false)
		}
		return nil
	}

	seqstring := la.Bimap.SEQ2STRING[seq]
	values, have_param := base[la.Pname]
	base.Set(la.Pname, seqstring)

	if !have_param { // no parameter in url
		if seq == la.Bimap.Default_seq {
			return unlessreverse(false)
		}
		return nil
	}

	pos, neg := values[0], values[0]
	if neg[0] == '-' {
		pos = neg[1:]
		neg = neg[1:]
	} else {
		neg = "-" + neg
	}

	var ascr *bool
	if pos == seqstring {
		t := neg[0] != '-'
		if seq == la.Bimap.Default_seq {
			t = true
		}
		ascr = unlessreverse(t)
		base.Set(la.Pname, neg)
	}
	if seq == la.Bimap.Default_seq {
		base.Del(la.Pname)
	}
	return ascr
}

type Linkattrs struct {
	Base url.Values
	Pname string
	Bimap Biseqmap
	Seq SEQ
}

type InterfaceMeta struct {
	NameKey  string
	NameHTML template.HTML
}

type Interface struct {
	InterfaceMeta
	In       string // with units
	Out      string // with units
	DeltaIn  string // with units
	DeltaOut string // with units
}

type Interfaces struct { // TODO Networks
	List []Interface
}

type ProcInfo struct {
	PID      uint

	Priority int
	Nice     int

	Time     uint64
	Name     string

	Uid      uint

	Size        uint64
	Resident    uint64
}

type ProcData struct {
	PID      uint

	Priority int
	Nice     int

	Time     string
	NameRaw  string
	NameHTML template.HTML

	UserHTML template.HTML
	Size     string // with units
	Resident string // with units
}
