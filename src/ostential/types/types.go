package types
import (
	"net/url"
)

type SEQ int
func(seq SEQ) Sign(t bool) bool { // used in sortable_*.go
	if seq < 0 {
		return t
	}
	return !t
}

type CPU struct {
	List []Core
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

type DiskData struct {
	DiskName    string
	ShortDiskName string

	Total       string // with units
	Used        string // with units
	Avail       string // with units
	UsePercent  string // as a string, with "%"
	Inodes      string // with units
	Iused       string // with units
	Ifree       string // with units
	IusePercent string // as a string, with "%"
	DirName     string

	 UsePercentClass string
	IusePercentClass string
}

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

type DeltaInterface struct {
	Name     string
	In       string // with units
	Out      string // with units
	DeltaIn  string // with units
	DeltaOut string // with units
}

type Interfaces struct {
	List []DeltaInterface
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
	Name     string

	User     string
	Size     string // with units
	Resident string // with units
}
