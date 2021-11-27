package term

import (
	"strconv"

	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
)

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
	Symbol     string      `parser:"@( PLUS | MINUS)?" json:"symbol"`
	SingleTerm *SingleTerm `parser:"( @@ " json:"single_term"`
	PhraseTerm *PhraseTerm `parser:"| @@ " json:"phrase_term"`
	RangeTerm  *RangeTerm  `parser:"| @@)" json:"range_term"`
}

func (t *PrefixTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.Symbol + t.SingleTerm.String()
	} else if t.PhraseTerm != nil {
		return t.Symbol + t.PhraseTerm.String()
	} else if t.RangeTerm != nil {
		return t.Symbol + t.RangeTerm.String()
	} else {
		return ""
	}
}

func (t *PrefixTerm) GetPrefixType() op.PrefixOPType {
	if t == nil {
		return op.UNKNOWN_PREFIX_TYPE
	} else if t.Symbol == "+" {
		return op.MUST_PREFIX_TYPE
	} else if t.Symbol == "-" {
		return op.MUST_NOT_PREFIX_TYPE
	} else {
		return op.SHOULD_PREFIX_TYPE
	}
}

// whitespace is prefix with prefix term
type WPrefixTerm struct {
	Symbol     string      `parser:"WHITESPACE @(PLUS|MINUS)?" json:"symbol"`
	SingleTerm *SingleTerm `parser:"( @@ " json:"single_term"`
	PhraseTerm *PhraseTerm `parser:"| @@ " json:"phrase_term"`
	RangeTerm  *RangeTerm  `parser:"| @@)" json:"range_term"`
}

func (t *WPrefixTerm) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return " " + t.Symbol + t.SingleTerm.String()
	} else if t.PhraseTerm != nil {
		return " " + t.Symbol + t.PhraseTerm.String()
	} else if t.RangeTerm != nil {
		return " " + t.Symbol + t.RangeTerm.String()
	} else {
		return ""
	}
}

func (t *WPrefixTerm) GetPrefixType() op.PrefixOPType {
	if t == nil {
		return op.UNKNOWN_PREFIX_TYPE
	} else if t.Symbol == "+" {
		return op.MUST_PREFIX_TYPE
	} else if t.Symbol == "-" {
		return op.MUST_NOT_PREFIX_TYPE
	} else {
		return op.SHOULD_PREFIX_TYPE
	}
}

type PrefixTermGroup struct {
	PrefixTerm  *PrefixTerm    `parser:"LPAREN WHITESPACE* @@" json:"prefix_term"`
	PrefixTerms []*WPrefixTerm `parser:"@@* WHITESPACE* RPAREN" json:"prefix_terms`
}

// type OrTerm struct {
// 	AndTerm  *AndTerm   `parser:"@@" json:"and_term"`
// 	OrJTerms []*OrJTerm `parser:"@@+" json:"orj_terms"`
// }

type TermGroup struct {
	PrefixTerms []PrefixTerm `parser:"LPAREN @@+ RPAREN" json:"pre-bool_terms"`
}

// a term with boost symbol like this ( foo bar )^2 / [1 TO 2]^2 or nothing (default 1.0)
type BoostTerm struct {
	RangeTerm   *RangeTerm       `parser:"( @@  " json:"range_term"`
	GroupTerm   *PrefixTermGroup `parser:"| @@ )" json:"group_term"`
	BoostSymbol string           `parser:"@BOOST?" json:"boost_symbol"`
}

// fuzzy term: term can by suffix with fuzzy or boost like this foo^2 / "foo bar"^2 / foo~ / "foo bar"~2
type FuzzyTerm struct {
	SingleTerm  *SingleTerm `parser:"( @@ " json:"single_term"`
	PhraseTerm  *PhraseTerm `parser:"| @@)" json:"phrase_term"`
	FuzzySymbol string      `parser:"( @FUZZY " json:"fuzzy_symbol"`
	BoostSymbol string      `parser:"| @BOOST )?" json:"boost_symbol`
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
