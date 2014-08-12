package types

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
