package term

import "strconv"

// single side range term or double side range term
type RangeTerm struct {
	SRangeTerm *SRangeTerm `parser:"  @@" json:"s_range_term"`
	DRangeTerm *DRangeTerm `parser:"| @@" json:"d_range_term"`
}

func (t *RangeTerm) String() string {
	if t == nil {
		return ""
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.String()
	} else if t.DRangeTerm != nil {
		return t.DRangeTerm.String()
	} else {
		return ""
	}
}

func (t *RangeTerm) ToBound() *Bound {
	if t == nil {
		return nil
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.ToBound()
	} else if t.DRangeTerm != nil {
		return t.DRangeTerm.ToBound()
	} else {
		return nil
	}
}

// bool term: a term is behind of symbol ("+" / "-" / "!")
type PrefixTerm struct {
	BoolSymbol string      `parser:"@(MINUS|PLUS|NOT)?" json:"prefix"`
	SimpleTerm *SingleTerm `parser:"( @@ " json:"simple_term"`
	PhraseTerm *PhraseTerm `parser:"| @@ " json:"phrase_term"`
	RangeTerm  *RangeTerm  `parser:"| @@)" json:"range_term"`
}

// func (t *PrefixTerm) String() string {
// 	if t == nil {
// 		return ""
// 	} else if t.SimpleTerm
// }

type GroupTerm struct {
	PreBoolTerms []PrefixTerm `parser:"LPAREN @@+ RPAREN" json:"pre-bool_terms"`
}

// a term with boost symbol like this foo^2 / "foo bar"^2 / [1 TO 2]^2 or nothing (default 1.0)
type BoostTerm struct {
	SingleTerm  *SingleTerm `parser:"( @@  " json:"simple_term"`
	PhraseTerm  *PhraseTerm `parser:"| @@  " json:"phrase_term"`
	RangeTerm   *RangeTerm  `parser:"| @@  " json:"range_term"`
	GroupTerm   *GroupTerm  `parser:"| @@ )" json:"group_term"`
	BoostSymbol string      `parser:"@BOOST?" json:"boost_symbol"`
}

type FuzzyTerm struct {
	SingleTerm  *SingleTerm `parser:"( @@  " json:"simple_term"`
	PhraseTerm  *PhraseTerm `parser:"| @@ )" json:"phrase_term"`
	FuzzySymbol string      `parser:"( @FUZZY  " json:"fuzzy_symbol"`
	BoostSymbol string      `parser:"| @BOOST)?" json:"boost_symbol"`
}

func (t *FuzzyTerm) Fuzziness() int {
	if t == nil || len(t.FuzzySymbol) == 0 {
		return 0
	} else if t.FuzzySymbol == "~" {
		return 1
	} else {
		var v, _ = strconv.Atoi(t.FuzzySymbol[1:])
		return v
	}
}

func (t *FuzzyTerm) Boost() float64 {
	if t == nil {
		return 0.0
	} else if len(t.BoostSymbol) == 0 {
		return 1.0
	} else {
		var res, _ = strconv.ParseFloat(t.BoostSymbol[1:], 64)
		return res
	}
}

func (t *FuzzyTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.SingleTerm.String() + t.FuzzySymbol + t.BoostSymbol
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.String() + t.FuzzySymbol + t.BoostSymbol
	} else {
		return ""
	}
}
