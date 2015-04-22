package client

import "net/url"

// Uint is a positive or 0 number.
type Uint uint

type Number struct {
	Uint
	Negative bool
}

// Uinter defines required (read-only) methods
// for all Uint-derived types interface.
type Uinter interface {
	Touint() Uint
	Marshal() (string, error)
	// MarshalJSON() ([]byte, error)
}

// Upointer defines required (incl. pointer-) methods
// for all Uint-derived types interface.
type Upointer interface {
	Uinter
	Unmarshal(string, *bool) error
	// UnmarshalJSON([]byte) error
}

// Attr type keeps link attributes.
type Attr struct {
	Href, Class, CaretClass string
}

// EncodeNU returns uinter applied Attr.
func (links Links) EncodeNU(pname string, uinter Uinter) Attr {
	base := url.Values{}
	for k, v := range links.Values {
		base[k] = v
	}
	attr := Attr{Class: "state"}
	if cur := links.SetBase(base, pname, uinter); cur != nil {
		attr.CaretClass = "caret"
		attr.Class += " current"
		if *cur {
			attr.Class += " dropup"
		}
	}
	attr.Href = "?" + base.Encode() // sorted by key
	return attr
}

// SetBase modifies the base.
func (links Links) SetBase(base url.Values, pname string, uinter Uinter) *bool {
	this := uinter.Touint()

	// TODO Better name for Decodes/DecodedMap/Decoded
	// because they might be not decoded, but default.
	decoded, haveok := links.Decodes.DecodedMap[pname]
	if !haveok {
		// That's unexpected: SetBase (and Encode) is not supposed to be called without pname decoded prior.
		return nil
	}

	ddef := decoded.Decoder.Default.Uint
	dnum := decoded.Number

	// Default ordering is desc (values are numeric most of the time).
	// Alpha values ordering: asc.
	desc := !decoded.IsAlpha(this)
	if dnum.Negative {
		desc = !desc
	}
	var ret *bool
	if this == dnum.Uint {
		ret = new(bool)
		*ret = !desc
	}
	// for default, opposite of having a parameter is it's absence.
	if this == ddef && decoded.Specified {
		base.Del(pname)
		return ret
	}
	low, err := uinter.Marshal()
	if err != nil { // ignoring the error
		return nil
	}
	if this == dnum.Uint && !dnum.Negative {
		low = "-" + low
	}
	base.Set(pname, low)
	return ret
}

// Links type for link making.
type Links struct {
	url.Values // provides Set(string, string) for Linker interface.
	Decodes    // provides SetDecoded(string, Decoded), SetError(error) for Linker interface.
}

type Decoder struct {
	Default Number
	Alphas  []Uint
}

func (d Decoder) IsAlpha(p Uint) bool {
	for _, u := range d.Alphas {
		if u == p {
			return true
		}
	}
	return false
}

var PS = Decoder{
	Default: Number{Uint: Uint(PID)},
	Alphas:  []Uint{Uint(NAME), Uint(UID)},
}

var DF = Decoder{
	Default: Number{Uint: Uint(FS)},
	Alphas:  []Uint{Uint(FS), Uint(MP)},
}

func NewLinks() *Links {
	return &Links{
		Values:  make(url.Values),
		Decodes: Decodes{DecodedMap: make(DecodedMap)},
	}
}

type Linker interface {
	Set(string, string)
	SetDecoded(string, Decoded)
	SetError(error)
}

func (d Decoder) Decode(form url.Values, pname string, linker Linker, setn *Number, uptr Upointer) error {
	n, spec, err := d.Find(form[pname], pname, linker, uptr)
	if err != nil {
		return err
	}
	*setn = n
	linker.SetDecoded(pname, Decoded{Number: n, Decoder: d, Specified: spec})
	return nil
}

// Find side effects: ui.Unmarshal (ui.Methods[2]) and linker.Set (eg url.Values{}.Set())
func (d Decoder) Find(values []string, pname string, linker Linker, uptr Upointer) (Number, bool, error) {
	if len(values) == 0 || values[0] == "" {
		return d.Default, false, nil
	}
	var negate bool
	in := values[0]
	if in[0] == '-' {
		in = in[1:]
		negate = true
	}
	err := uptr.Unmarshal(in, &negate) // .UnmarshalJSON([]byte(fmt.Sprintf("%q", strings.ToUpper(in))))
	n := Number{}
	if err != nil {
		if rerr, ok := err.(RenamedConstError); ok {
			// The case when err (of type RenamedConstError) is set
			// AND uptr actually holds corresponding ("renamed") value.
			if l, err := uptr.Marshal(); err == nil {
				if negate {
					l = "-" + l
				}
				linker.Set(pname, l)
			}
			linker.SetError(rerr)
		}
		return n, true, err
	}
	n.Uint = uptr.Touint()
	if negate {
		n.Negative = true
	}
	linker.Set(pname, values[0])
	return n, true, err
}

type Decoded struct {
	Number
	Decoder
	Specified bool
}

type DecodedMap map[string]Decoded
type Decodes struct {
	DecodedMap
	RCError error
}

// SetDecoded required by Linker interface.
func (ds *Decodes) SetDecoded(name string, d Decoded) { ds.DecodedMap[name] = d }

func (ds *Decodes) SetError(err error) { ds.RCError = err }
