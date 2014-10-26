package types

import (
	"html/template"
	"net/url"

	sigar "github.com/rzab/gosigar"
)

// SEQ is a distinct int type for consts and other uses.
type SEQ int

// AnyOf returns true if the seq is present in the list.
func (seq SEQ) AnyOf(list []SEQ) bool {
	for _, s := range list {
		if s == seq {
			return true
		}
	}
	return false
}

// Sign is a logical operator, useful for sorting.
func (seq SEQ) Sign(t bool) bool { // used in sortable_*.go
	if seq < 0 {
		return t
	}
	return !t
}

// Memory type is a struct of memory metrics.
type Memory struct {
	Kind           string
	Total          string
	Used           string
	Free           string
	UsePercentHTML template.HTML
}

// MEM type has a list of Memory.
type MEM struct {
	List []Memory
}

type CpuList struct {
	sigarList sigar.CpuList
	deltaList *sigar.CpuList
}

func NewCpuList() CpuList {
	cl := sigar.CpuList{}
	cl.Get()
	return CpuList{sigarList: cl}
}

func (cl CpuList) List() sigar.CpuList {
	if cl.deltaList != nil {
		return *cl.deltaList
	}
	return cl.sigarList
}

func (cl CpuList) SigarList() *sigar.CpuList {
	return &cl.sigarList
}

func (cl *CpuList) CalculateDelta(other sigar.CpuList) {
	if len(other.List) == 0 {
		return
	}
	if len(other.List) != len(cl.sigarList.List) {
		return
	}
	cl.deltaList = &sigar.CpuList{
		List: make([]sigar.Cpu, len(cl.sigarList.List)),
	}
	for i, sigarCpu := range cl.sigarList.List {
		cl.deltaList.List[i] = sigarCpu.Delta(other.List[i])
	}
}

// CPU type has a list of Core.
type CPU struct {
	List []Core
}

// Core type is a struct of core metrics.
type Core struct {
	N         string
	User      uint // percent without "%"
	Sys       uint // percent without "%"
	Idle      uint // percent without "%"
	UserClass string
	SysClass  string
	IdleClass string
	// UserSpark string
	// SysSpark  string
	// IdleSpark string
}

// DiskMeta type has common for DiskBytes and DiskInodes fields.
type DiskMeta struct {
	DiskNameHTML template.HTML
	DirNameHTML  template.HTML
	DirNameKey   string
}

// DiskBytes type is a struct of disk bytes metrics.
type DiskBytes struct {
	DiskMeta
	Total           string // with units
	Used            string // with units
	Avail           string // with units
	UsePercent      string // as a string, with "%"
	UsePercentClass string
}

// DiskInodes type is a struct of disk inodes metrics.
type DiskInodes struct {
	DiskMeta
	Inodes           string // with units
	Iused            string // with units
	Ifree            string // with units
	IusePercent      string // as a string, with "%"
	IusePercentClass string
}

// DFbytes type has a list of DiskBytes.
type DFbytes struct {
	List []DiskBytes
}

// DFinodes type has a list of DiskInodes.
type DFinodes struct {
	List []DiskInodes
}

// type DiskTable struct {
// 	List  []DiskData
// 	Links *DiskLinkattrs `json:",omitempty"`
// 	HaveCollapsed bool
// }

// Attr type keeps link attributes.
type Attr struct {
	Href, Class, CaretClass string
}

// Attr returns a seq applied Attr taking the la link and updating/setting the parameter.
func (la Linkattrs) Attr(seq SEQ) Attr {
	base := url.Values{}
	for k, v := range la.Base {
		base[k] = v
	}
	attr := Attr{Class: "state"}
	if ascp := la._attr(base, seq); ascp != nil {
		attr.CaretClass = "caret"
		attr.Class += " current"
		if *ascp {
			attr.Class += " dropup"
		}
	}
	attr.Href = "?" + base.Encode() // la._attr modifies base, DO NOT use prior to the call
	return attr
}

// _attr side effect: modifies the base
func (la Linkattrs) _attr(base url.Values, seq SEQ) *bool {
	unlessreverse := func(t bool) *bool {
		if la.Bimap.SEQ2REVERSE[seq] {
			t = !t
		}
		return &t
	}

	if la.Pname == "" {
		if seq == la.Bimap.DefaultSeq {
			return unlessreverse(false)
		}
		return nil
	}

	seqstring := la.Bimap.SEQ2STRING[seq]
	values, haveParam := base[la.Pname]
	base.Set(la.Pname, seqstring)

	if !haveParam { // no parameter in url
		if seq == la.Bimap.DefaultSeq {
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
		if seq == la.Bimap.DefaultSeq {
			t = true
		}
		ascr = unlessreverse(t)
		base.Set(la.Pname, neg)
	}
	if seq == la.Bimap.DefaultSeq {
		base.Del(la.Pname)
	}
	return ascr
}

// Linkattrs type for link making.
type Linkattrs struct {
	Base  url.Values
	Pname string
	Bimap Biseqmap
}

// InterfaceMeta type has common Interface fields.
type InterfaceMeta struct {
	NameKey  string
	NameHTML template.HTML
}

// Interface type is a struct of interface metrics.
type Interface struct {
	InterfaceMeta
	In       string // with units
	Out      string // with units
	DeltaIn  string // with units
	DeltaOut string // with units
}

// Interfaces type has a list of Interface.
type Interfaces struct {
	List []Interface
}

// ProcInfo type is an internal account of a process.
type ProcInfo struct {
	PID      uint
	Priority int
	Nice     int
	Time     uint64
	Name     string
	UID      uint
	Size     uint64
	Resident uint64
}

// ProcData type is a public (for index context, json marshaling) account of a process.
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
