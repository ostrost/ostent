package client

import (
	"net/http"
	"net/url"
)

// SEQ is a distinct int type for consts and other uses.
type SEQ int

// SeqNReverse holds a SEQ and a Reverse bool.
type SeqNReverse struct {
	SEQ     SEQ
	Reverse bool
}

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

/* * Attr & Linkattrs: ******************************************** */

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

func valuesSet(req *http.Request, base url.Values, pname string, bimap Biseqmap) SEQ {
	if params, ok := req.Form[pname]; ok && len(params) > 0 {
		if seq, ok := bimap.STRING2SEQ[params[0]]; ok {
			base.Set(pname, params[0])
			return seq
		}
	}
	return bimap.DefaultSeq
}

func NewLinkAttrs(req *http.Request, base url.Values, pname string, bimap Biseqmap, seq *SEQ) *Linkattrs {
	*seq = valuesSet(req, base, pname, bimap)
	return &Linkattrs{
		Base:  base,
		Pname: pname,
		Bimap: bimap,
	}
}

/* * bimap.go: **************************************************** */

// Seq2string type is a map of string by SEQ
type Seq2string map[SEQ]string

// Biseqmap type holds bi-directional relations between SEQ and string and a DefaultSeq
type Biseqmap struct {
	SEQ2STRING  Seq2string
	STRING2SEQ  map[string]SEQ
	SEQ2REVERSE map[SEQ]bool
	DefaultSeq  SEQ
}

func contains(thiss SEQ, lists []SEQ) bool {
	for _, s := range lists {
		if s == thiss {
			return true
		}
	}
	return false
}

// Seq2bimap makes a Biseqmap with default defSeq. reverse holds a list of SEQ to be reversed.
func Seq2bimap(defSeq SEQ, s2s Seq2string, reverse []SEQ) Biseqmap {
	bi := Biseqmap{
		SEQ2STRING:  Seq2string{},
		STRING2SEQ:  map[string]SEQ{},
		SEQ2REVERSE: map[SEQ]bool{},
	}
	bi.DefaultSeq = defSeq

	for seq, str := range s2s {
		isreverse := contains(seq, reverse)
		bi.SEQ2REVERSE[seq] = isreverse
		bi.SEQ2REVERSE[-seq] = isreverse

		bi.SEQ2STRING[seq] = str
		bi.SEQ2STRING[-seq] = "-" + str

		nseq := seq
		if seq == defSeq {
			nseq = -nseq
		}
		bi.STRING2SEQ[str] = nseq
		bi.STRING2SEQ["-"+str] = -nseq
	}
	return bi
}
