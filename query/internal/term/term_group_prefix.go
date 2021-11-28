package term

import (
	"strconv"
	"strings"

	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
)

// term group elem
type TermGroupElem struct {
	SingleTerm *SingleTerm `parser:"  @@" json:"single_term"`
	PhraseTerm *PhraseTerm `parser:"| @@" json:"phrase_term"`
	SRangeTerm *SRangeTerm `parser:"| @@" json:"single_range_term"`
	DRangeTerm *DRangeTerm `parser:"| @@" json:"double_range_term"`
}

func (t *TermGroupElem) String() string {
	if t == nil {
		return ""
	} else if t.SingleTerm != nil {
		return t.SingleTerm.String()
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.String()
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.String()
	} else if t.DRangeTerm != nil {
		return t.DRangeTerm.String()
	} else {
		return ""
	}
}

// prefix term: a term is behind of prefix operator symbol ("+" / "-")
type PrefixTerm struct {
	Symbol string         `parser:"@( PLUS | MINUS)?" json:"symbol"`
	Elem   *TermGroupElem `parser:"@@" json:"elem"`
}

func (t *PrefixTerm) String() string {
	if t == nil {
		return ""
	} else if t.Elem != nil {
		return t.Symbol + t.Elem.String()
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
	Symbol string         `parser:"WHITESPACE @(PLUS|MINUS)?" json:"symbol"`
	Elem   *TermGroupElem `parser:"@@" json:"elem"`
}

func (t *WPrefixTerm) String() string {
	if t == nil {
		return ""
	} else if t.Elem != nil {
		return " " + t.Symbol + t.Elem.String()
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
	PrefixTerm  *PrefixTerm    `parser:"LPAREN WHITESPACE* @@ " json:"prefix_term"`
	PrefixTerms []*WPrefixTerm `parser:"@@*  WHITESPACE* RPAREN" json:"prefix_terms"`
	BoostSymbol string         `parser:"@BOOST?" json:"boost_symbol"`
}

func (t *PrefixTermGroup) String() string {
	if t == nil {
		return ""
	} else if t.PrefixTerm != nil {
		var sl = []string{t.PrefixTerm.String()}
		for _, x := range t.PrefixTerms {
			sl = append(sl, x.String())
		}
		return strings.Join(sl, "")
	} else {
		return ""
	}
}

func (t *PrefixTermGroup) Boost() float64 {
	if t == nil {
		return 0.0
	} else if len(t.BoostSymbol) == 0 {
		return 1.0
	} else {
		var res, _ = strconv.ParseFloat(t.BoostSymbol[1:], 64)
		return res
	}
}
