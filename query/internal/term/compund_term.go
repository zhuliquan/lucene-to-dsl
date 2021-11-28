package term

import (
	"strconv"
	"strings"

	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
)

// single side range term or double side range and with boost like this [1 TO 2]^2
type RangeTerm struct {
	SRangeTerm  *SRangeTerm `parser:"( @@ " json:"s_range_term"`
	DRangeTerm  *DRangeTerm `parser:"| @@)" json:"d_range_term"`
	BoostSymbol string      `parser:"@BOOST?" json:"boost_symbol"`
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

func (t *RangeTerm) Boost() float64 {
	if t == nil {
		return 0.0
	} else if len(t.BoostSymbol) == 0 {
		return 1.0
	} else {
		var res, _ = strconv.ParseFloat(t.BoostSymbol[1:], 64)
		return res
	}
}

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
	Elem   *TermGroupElem `parser:"@@" json:"elem`
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
	PrefixTerm  *PrefixTerm    `parser:"@@ " json:"prefix_term"`
	PrefixTerms []*WPrefixTerm `parser:"@@*" json:"prefix_terms`
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

// term group: join sum prefix term group together
type TermGroup struct {
	OrTermGroup *OrTermGroup   `parser:"LPAREN WHITESPACE* @@ " json:"or_term_group"`
	OSTermGroup []*OSTermGroup `parser:"@@* WHITESPACE* RPAREN" json:"or_symbol_term_group"`
	BoostSymbol string         `parser:"@BOOST?" json:"boost_symbol"`
}

func (t *TermGroup) String() string {
	if t != nil {
		return ""
	} else if t.OrTermGroup != nil {
		var sl = []string{t.OrTermGroup.String()}
		for _, x := range t.OSTermGroup {
			sl = append(sl, x.String())
		}
		return strings.Join(sl, "")
	} else {
		return ""
	}
}

func (t *TermGroup) Boost() float64 {
	if t == nil {
		return 0.0
	} else if len(t.BoostSymbol) == 0 {
		return 1.0
	} else {
		var res, _ = strconv.ParseFloat(t.BoostSymbol[1:], 64)
		return res
	}
}

type OrTermGroup struct {
	AndTermGroup *AndTermGroup   `parser:"@@ " json:"and_term_group"`
	AnSTermGroup []*AnSTermGroup `parser:"@@*" json:"and_symbol_term_group"`
}

func (t *OrTermGroup) String() string {
	if t != nil {
		return ""
	} else if t.AndTermGroup != nil {
		var sl = []string{t.AndTermGroup.String()}
		for _, x := range t.AnSTermGroup {
			sl = append(sl, x.String())
		}
		return strings.Join(sl, "")
	} else {
		return ""
	}
}

type OSTermGroup struct {
	OrSymbol    *op.OrSymbol  `parser:"@@ " json:"or_symbol"`
	NotSymbol   *op.NotSymbol `parser:"@@?" json:"not_symbol"`
	OrTermGroup *OrTermGroup  `parser:"@@ " json:"or_term_group"`
}

func (t *OSTermGroup) String() string {
	if t != nil {
		return ""
	} else if t.OrTermGroup != nil {
		return t.OrSymbol.String() + t.NotSymbol.String() + t.OrTermGroup.String()
	} else {
		return ""
	}
}

type AndTermGroup struct {
	NotTermGroup   *NotTermGroup    `parser:"  @@" json:"not_term_group"`
	ParenTermGroup *ParenTermGroup  `parser:"| @@" json:"paren_term_group"`
	TermGroupElem  *PrefixTermGroup `parser:"| @@" json:"term_group_elem"`
}

func (t *AndTermGroup) String() string {
	if t != nil {
		return ""
	} else if t.NotTermGroup != nil {
		return t.NotTermGroup.String()
	} else if t.ParenTermGroup != nil {
		return t.ParenTermGroup.String()
	} else if t.TermGroupElem != nil {
		return t.TermGroupElem.String()
	} else {
		return ""
	}
}

type AnSTermGroup struct {
	AndSymbol    op.AndSymbol  `parser:"@@" json:"and_symbol"`
	NotSymbol    op.AndSymbol  `parser:"@@?" json:"not_symbol`
	AndTermGroup *AndTermGroup `parser:"@@" json:"and_term_group"`
}

func (t *AnSTermGroup) String() string {
	if t != nil {
		return ""
	} else if t.AndTermGroup != nil {
		return t.AndSymbol.String() + t.NotSymbol.String() + t.AndTermGroup.String()
	} else {
		return ""
	}
}

type NotTermGroup struct {
	NotSymbol    op.NotSymbol `parser:"@@" json:"not_symbol"`
	SubTermGroup *TermGroup   `parser:"@@" json:"sub_term_group"`
}

func (t *NotTermGroup) String() string {
	if t != nil {
		return ""
	} else if t.SubTermGroup != nil {
		return t.NotSymbol.String() + t.SubTermGroup.String()
	} else {
		return ""
	}
}

type ParenTermGroup struct {
	SubTermGroup *TermGroup `parser:"LPAREN WHITESPACE* @@ WHITESPACE* RPAREN" json:"sub_term_group"`
}

func (t *ParenTermGroup) String() string {
	if t != nil {
		return ""
	} else if t.SubTermGroup != nil {
		return "( " + t.SubTermGroup.String() + " )"
	} else {
		return ""
	}
}

// a term with boost symbol like this ( foo bar )^2 / [1 TO 2]^2 or nothing (default 1.0)
type BoostTerm struct {
	RangeTerm   *RangeTerm       `parser:"( @@  " json:"range_term"`
	GroupTerm   *PrefixTermGroup `parser:"| @@ )" json:"group_term"`
	BoostSymbol string           `parser:"@BOOST?" json:"boost_symbol"`
}

func (t *BoostTerm) String() string {
	if t == nil {
		return ""
	} else if t.RangeTerm != nil {
		return t.RangeTerm.String() + t.BoostSymbol
	} else if t.GroupTerm != nil {
		return t.GroupTerm.String() + t.BoostSymbol
	} else {
		return ""
	}
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

func (t *FuzzyTerm) haveWildcard() bool {
	if t == nil {
		return false
	} else if t.SingleTerm != nil {
		return t.SingleTerm.haveWildcard()
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.haveWildcard()
	} else {
		return false
	}
}

func (t *FuzzyTerm) ValueS() string {
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
