package term

import (
	"strconv"

	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
)

// a simple term
type SimpleTerm struct {
	SingleTerm *SimpleTerm `parser:"  @@" json:"single_term"`
	PhraseTerm *PhraseTerm `parser:"| @@" json:"phrase_term"`
}

func (t *SimpleTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.SingleTerm.String()
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.String()
	} else {
		return ""
	}
}

func (t *SimpleTerm) ValueS() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.SingleTerm.ValueS()
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.ValueS()
	} else {
		return ""
	}
}

func (t *SimpleTerm) haveWildcard() bool {
	if t != nil {
		return false
	} else if t.SingleTerm != nil {
		return t.SingleTerm.haveWildcard()
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.haveWildcard()
	} else {
		return false
	}
}

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

// prefix term: a term is behind of prefix operator symbol ("+" / "-")
type PrefixTerm struct {
	Prefix     *op.PreSymbol `parser:"@@" json:"prefix"`
	SingleTerm *SingleTerm   `parser:"( @@ " json:"simple_term"`
	PhraseTerm *PhraseTerm   `parser:"| @@ " json:"phrase_term"`
	RangeTerm  *RangeTerm    `parser:"| @@)" json:"range_term"`
}

func (t *PrefixTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.Prefix.String() + t.SingleTerm.String()
	} else if t.PhraseTerm != nil {
		return t.Prefix.String() + t.PhraseTerm.String()
	} else if t.RangeTerm != nil {
		return t.Prefix.String() + t.RangeTerm.String()
	} else {
		return ""
	}
}

func (t *PrefixTerm) GetPrefixType() op.PrefixOPType {
	return t.Prefix.GetPrefixType()
}

type WPrefixTerm struct {
	Prefix     *op.PreSymbol `parser:"WHITESPACE @@" json:"prefix"`
	SingleTerm *SingleTerm   `parser:"( @@ " json:"simple_term"`
	PhraseTerm *PhraseTerm   `parser:"| @@ " json:"phrase_term"`
	RangeTerm  *RangeTerm    `parser:"| @@)" json:"range_term"`
}

func (t *WPrefixTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return " " + t.Prefix.String() + t.SingleTerm.String()
	} else if t.PhraseTerm != nil {
		return t.Prefix.String() + t.PhraseTerm.String()
	} else if t.RangeTerm != nil {
		return " " + t.Prefix.String() + t.RangeTerm.String()
	} else {
		return ""
	}
}

func (t *WPrefixTerm) GetPrefixType() op.PrefixOPType {
	return t.Prefix.GetPrefixType()
}

type GroupPrefixTerm struct {
	PrefixTerm  *PrefixTerm    `parser:"LPAREN @@" json:"prefix_term"`
	WPrefixTerm []*WPrefixTerm `parser:"@@+ WHITESPACE RPAREN" json:"prefix_terms`
}

// type OrTerm struct {
// 	AndTerm  *AndTerm   `parser:"@@" json:"and_term"`
// 	OrJTerms []*OrJTerm `parser:"@@+" json:"orj_terms"`
// }

type GroupTerm struct {
	PrefixTerms []PrefixTerm `parser:"LPAREN @@+ RPAREN" json:"pre-bool_terms"`
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
	FuzzySymbol string      `parser:"@FUZZY?  " json:"fuzzy_symbol"`
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

func (t *FuzzyTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.SingleTerm.String() + t.FuzzySymbol
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.String() + t.FuzzySymbol
	} else {
		return ""
	}
}
