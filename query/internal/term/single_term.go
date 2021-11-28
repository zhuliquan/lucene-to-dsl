package term

import (
	"fmt"
	"strings"
)

// simple term: is a single term without escape char and whitespace
type SingleTerm struct {
	Value []string `parser:"@(IDENT|WILDCARD)+" json:"value"`
}

func (t *SingleTerm) ValueS() string {
	if t == nil {
		return ""
	} else {
		return strings.Join(t.Value, "")
	}
}

func (t *SingleTerm) String() string {
	if t == nil {
		return ""
	} else {
		return strings.Join(t.Value, "")
	}
}

func (t *SingleTerm) haveWildcard() bool {
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

// phrase term: a series of terms be surrounded with quotation, for instance "foo bar".
type PhraseTerm struct {
	Value string `parser:"@STRING" json:"value"`
}

func (t *PhraseTerm) ValueS() string {
	if t == nil || len(t.Value) == 0 {
		return ""
	} else {
		return t.Value[1 : len(t.Value)-1]
	}
}

func (t *PhraseTerm) String() string {
	if t == nil {
		return ""
	} else {
		return t.Value
	}
}

func (t *PhraseTerm) haveWildcard() bool {
	if t == nil {
		return false
	}
	for i := 1; i < len(t.Value)-1; i++ {
		if i == 1 && (t.Value[1] == '?' || t.Value[1] == '*') {
			return true
		}
		if i > 1 && ((t.Value[i] == '?' || t.Value[i] == '*') && t.Value[i-1] != '\\') {
			return true
		}

	}
	return false
}

// a regexp term is surrounded be slash, for instance /\d+\.?\d+/ in here if you want present '/' you should type '\/'
type RegexpTerm struct {
	Value string `parser:"@REGEXP" json:"value"`
}

func (t *RegexpTerm) ValuesS() string {
	if t == nil || len(t.Value) == 0 {
		return ""
	} else {
		return t.Value[1 : len(t.Value)-1]
	}
}

func (t *RegexpTerm) String() string {
	if t == nil {
		return ""
	} else {
		return t.Value
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
	SingleValue []string `parser:"| @(IDENT|PLUS|MINUS)+" json:"simple_value"`
}

func (v *RangeValue) String() string {
	if v == nil {
		return ""
	} else if len(v.PhraseValue) != 0 {
		return v.PhraseValue
	} else if len(v.InfinityVal) != 0 {
		return v.InfinityVal
	} else if len(v.SingleValue) != 0 {
		return strings.Join(v.SingleValue, "")
	} else {
		return ""
	}
}

//double side of range term: a term is surrounded by brace / bracket, for instance [1 TO 2] / [1 TO 2} / {1 TO 2] / {1 TO 2}
type DRangeTerm struct {
	LBRACKET string      `parser:"@(LBRACE|LBRACK) WHITESPACE*" json:"left_bracket"`
	LValue   *RangeValue `parser:"@@ WHITESPACE+ 'TO'" json:"left_value"`
	RValue   *RangeValue `parser:"WHITESPACE+ @@" json:"right_value"`
	RBRACKET string      `parser:"WHITESPACE* @(RBRACK|RBRACE)" json:"right_bracket"`
}

func (t *DRangeTerm) ToBound() *Bound {
	if t == nil {
		return nil
	} else if t.LBRACKET == "[" || t.RBRACKET == "]" {
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

func (t *DRangeTerm) String() string {
	if t == nil {
		return ""
	} else {
		return fmt.Sprintf("%s %s TO %s %s", t.LBRACKET, t.LValue.String(), t.RValue.String(), t.RBRACKET)
	}
}

// single side of range term: a term is behind of symbol ('>' / '<' / '>=' / '<=')
type SRangeTerm struct {
	Symbol string      `parser:"@COMPARE" json:"symbol"`
	Value  *RangeValue `parser:"@@" json:"value"`
}

func (t *SRangeTerm) toDRangeTerm() *DRangeTerm {
	if t == nil {
		return nil
	} else if t.Symbol == ">" && t.Value != nil {
		return &DRangeTerm{LBRACKET: "{", LValue: t.Value, RValue: &RangeValue{InfinityVal: "*"}, RBRACKET: "}"}
	} else if t.Symbol == ">=" && t.Value != nil {
		return &DRangeTerm{LBRACKET: "[", LValue: t.Value, RValue: &RangeValue{InfinityVal: "*"}, RBRACKET: "}"}
	} else if t.Symbol == "<" && t.Value != nil {
		return &DRangeTerm{LBRACKET: "{", LValue: &RangeValue{InfinityVal: "*"}, RValue: t.Value, RBRACKET: "}"}
	} else if t.Symbol == "<=" && t.Value != nil {
		return &DRangeTerm{LBRACKET: "{", LValue: &RangeValue{InfinityVal: "*"}, RValue: t.Value, RBRACKET: "]"}
	} else {
		return nil
	}
}

func (t *SRangeTerm) ToBound() *Bound {
	return t.toDRangeTerm().ToBound()
}

func (t *SRangeTerm) String() string {
	return t.toDRangeTerm().String()
}
