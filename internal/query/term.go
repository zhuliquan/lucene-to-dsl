package query

import (
	"fmt"
	"strconv"
	"strings"
)

type Term struct {
	RegexpTerm *RegexpTerm `parser:"  @@" json:"regexp_term"`
	SRangeTerm *SRangeTerm `parser:"| @@" json:"ranges_term"`
	RangeTerm  *RangeTerm  `parser:"| @@" json:"range_term"`
}

func (t *Term) String() string {
	if t == nil {
		return ""
	} else if t.RegexpTerm != nil {
		return t.RegexpTerm.String()
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.String()
	} else if t.RangeTerm != nil {
		return t.RangeTerm.String()
	} else {
		return ""
	}
}

func (t *Term) isRegexp() bool {
	return t != nil && t.RegexpTerm != nil
}

func (t *Term) haveWildcard() bool {
	if t == nil {
		return false
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.haveWildcard()
	} else {
		return false
	}
}

func (t *Term) isRange() bool {
	return t != nil && (t.RangeTerm != nil || t.SRangeTerm.isRange())
}

func (t *Term) fuzziness() int {
	if t == nil {
		return 0
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.fuzziness()
	} else {
		return 0
	}
}

func (t *Term) boost() float64 {
	if t == nil {
		return 0.0
	} else if t.SRangeTerm != nil {
		return t.SRangeTerm.boost()
	} else {
		return 1.0
	}
}

// side range term: a term is behind of symbol ('>' / '<' / '>=' / '<=')
type SRangeTerm struct {
	Symbol     string      `parser:"@COMPARE?" json:"symbol"`
	SimpleTerm *SimpleTerm `parser:"( @@ " json:"simple_term"`
	PhraseTerm *PhraseTerm `parser:"| @@)" json:"phrase_term"`
}

func (t *SRangeTerm) String() string {
	if t == nil {
		return ""
	} else if t.PhraseTerm != nil {
		return t.Symbol + t.PhraseTerm.String()
	} else if t.SimpleTerm != nil {
		return t.Symbol + t.SimpleTerm.String()
	} else {
		return ""
	}
}

func (t *SRangeTerm) isRange() bool {
	return t != nil && len(t.Symbol) != 0
}

func (t *SRangeTerm) haveWildcard() bool {
	if t == nil {
		return false
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.haveWildcard()
	} else if t.SimpleTerm != nil {
		return t.SimpleTerm.haveWildcard()
	} else {
		return false
	}
}

func (t *SRangeTerm) fuzziness() int {
	if t == nil {
		return 0
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.fuzziness()
	} else if t.SimpleTerm != nil {
		return t.SimpleTerm.fuzziness()
	} else {
		return 0
	}
}

func (t *SRangeTerm) boost() float64 {
	if t == nil {
		return 0.0
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.boost()
	} else if t.SimpleTerm != nil {
		return t.SimpleTerm.boost()
	} else {
		return 1.0
	}
}

// a regexp term is surrounded be slash, for instance /\d+\.?\d+/
type RegexpTerm struct {
	Value string `parser:"@REGEXP" json:"value"`
}

func (t *RegexpTerm) String() string {
	if t == nil {
		return ""
	} else if t.Value != "" {
		return t.Value
	} else {
		return ""
	}
}

// phrase term: a series of terms be surrounded with quotation, for instance "foo bar".
type PhraseTerm struct {
	Value string `parser:"@STRING" json:"value"`
	Fuzzy string `parser:"@FUZZY?" json:"fuzzy"`
	Boost string `parser:"@BOOST?" json:"boost"`
}

func (t *PhraseTerm) String() string {
	if t == nil {
		return ""
	} else if t.Value != "" {
		var res = t.Value[1 : len(t.Value)-1]
		if t.Fuzzy != "" {
			res += t.Fuzzy
		}
		if t.Boost != "" {
			res += t.Boost
		}
		return res
	} else {
		return ""
	}
}

func (t *PhraseTerm) haveWildcard() bool {
	if t == nil {
		return false
	}
	for i := 1; i < len(t.Value)-1; i++ {
		if i > 1 && (t.Value[i] == '?' || t.Value[i] == '*' && t.Value[i-1] != '\\') {
			return true
		}
		if i == 1 && (t.Value[i] == '?' || t.Value[i] == '*') {
			return true
		}
	}
	return false
}

func (t *PhraseTerm) fuzziness() int {
	if t == nil || len(t.Fuzzy) == 0 {
		return 0
	} else if t.Fuzzy == "~" {
		return 1
	} else {
		var v, _ = strconv.Atoi(t.Fuzzy[1:])
		return v
	}
}

func (t *PhraseTerm) boost() float64 {
	if t == nil {
		return 0.0
	} else if len(t.Boost) != 0 {
		var v, _ = strconv.ParseFloat(t.Boost[1:], 64)
		return v
	} else {
		return 1.0
	}
}

// simple term: is a single term without escape char and whitespace
type SimpleTerm struct {
	Value []string `parser:"@(IDENT|WILDCARD)+" json:"value"`
	Fuzzy string   `parser:"@FUZZY?" json:"fuzzy"`
	Boost string   `parser:"@BOOST?" json:"boost"`
}

func (t *SimpleTerm) String() string {
	if t == nil {
		return ""
	} else if len(t.Value) != 0 {
		var res = strings.Join(t.Value, "")
		if len(t.Fuzzy) != 0 {
			res += t.Fuzzy
		}
		if len(t.Boost) != 0 {
			res += t.Boost
		}
		return res
	} else {
		return ""
	}
}

func (t *SimpleTerm) haveWildcard() bool {
	if t == nil {
		return false
	}
	for i := 0; i < len(t.Value); i++ {
		if t.Value[i] == "?" || t.Value[i] == "*" {
			return true
		}
	}
	return false
}

func (t *SimpleTerm) fuzziness() int {
	if t == nil || len(t.Fuzzy) == 0 {
		return 0
	} else if t.Fuzzy == "~" {
		return 1
	} else {
		var v, _ = strconv.Atoi(t.Fuzzy[1:])
		return v
	}
}

func (t *SimpleTerm) boost() float64 {
	if t == nil {
		return 0.0
	} else if len(t.Boost) != 0 {
		var v, _ = strconv.ParseFloat(t.Boost[1:], 64)
		return v
	} else {
		return 1.0
	}
}

// range bound like this [1, 2] [1, 2) (1, 2] (1, 2)
type Bound struct {
	LeftInclude  *RangeValue `json:"left_include"`
	LeftExclude  *RangeValue `json:"left_exclude"`
	RightInclude *RangeValue `json:"right_include"`
	RightExclude *RangeValue `json:"right_exclude"`
}

type RangeValue struct {
	PhraseValue string   `parser:"  @STRING" json:"phrase_value"`
	InfinityVal string   `parser:"| @('*')" json:"infinity_val"`
	SimpleValue []string `parser:"| @(IDENT|PLUS|MINUS)+" json:"simple_value"`
}

func (v *RangeValue) String() string {
	if v == nil {
		return ""
	} else if len(v.PhraseValue) != 0 {
		return v.PhraseValue
	} else if len(v.InfinityVal) != 0 {
		return v.InfinityVal
	} else if len(v.SimpleValue) != 0 {
		return strings.Join(v.SimpleValue, "")
	} else {
		return ""
	}
}

//range term: a term is surrounded by brace / bracket, for instance [1 TO 2] / [1 TO 2} / {1 TO 2] / {1 TO 2}
type RangeTerm struct {
	LBRACKET string      `parser:"@(LBRACE|LBRACK) WHITESPACE*" json:"left_bracket"`
	LValue   *RangeValue `parser:"@@" json:"left_value"`
	TO       string      `parser:"WHITESPACE+ @(\"TO\") WHITESPACE+"`
	RValue   *RangeValue `parser:"@@" json:"right_value"`
	RBRACKET string      `parser:"WHITESPACE* @(RBRACK|RBRACE)" json:"right_bracket"`
}

func (t *RangeTerm) ToBound() *Bound {
	if t == nil {
		return nil
	} else {
		if t.LBRACKET == "[" || t.RBRACKET == "]" {
			return &Bound{LeftInclude: t.LValue, RightInclude: t.RValue}
		} else if t.LBRACKET == "[" || t.RBRACKET == "}" {
			return &Bound{LeftInclude: t.LValue, RightExclude: t.RValue}
		} else if t.LBRACKET == "{" || t.RBRACKET == "]" {
			return &Bound{LeftExclude: t.LValue, RightInclude: t.RValue}
		} else if t.LBRACKET == "{" || t.RBRACKET == "}" {
			return &Bound{LeftExclude: t.LValue, RightExclude: t.RValue}
		} else {
			return nil
		}
	}
}

func (t *RangeTerm) String() string {
	if t == nil {
		return ""
	} else {
		return fmt.Sprintf("%s %s TO %s %s", t.LBRACKET, t.LValue.String(), t.RValue.String(), t.RBRACKET)
	}
}

// bool term: a term is behind of symbol ("+" / "-" / "!")
type BoolTerm struct {
	BoolSymbol string      `parser:"@(MINUS|PLUS|NOT)" json:"prefix"`
	PhraseTerm *PhraseTerm `parser:"( @@ " json:"phrase_term"`
	SimpleTerm *SimpleTerm `parser:"| @@ " json:"simple_term"`
	RegexpTerm *RegexpTerm `parser:"| @@ " json:"regexp_term"`
	RangeTerm  *RangeTerm  `parser:"| @@)" json:"range_term"`
}

// type GroupTerm struct {
// 	LParen 	`parser:"@"`
// }

// type GroupElemT struct {
// 	SimpleTerm *SimpleTerm `parser:"  @@" json:"simple_term"`
// 	PhraseTerm *PhraseTerm `parser:"| @@" json:"phrase_term"`
// }

// type GroupJoin struct {

// }

// type GroupElemS struct {
// 	WHITESPACE string `parser:"@@" json:""`

// }
