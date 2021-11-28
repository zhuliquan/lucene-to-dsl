package term

import (
	"fmt"
	"strings"

	bnd "github.com/zhuliquan/lucene-to-dsl/query/internal/bound"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

// simple term: is a single term without escape char and whitespace
type SingleTerm struct {
	Value    []string `parser:"@(IDENT|WILDCARD)+" json:"value"`
	wildcard int8
}

func (t *SingleTerm) GetTermType() TermType {
	return SINGLE_TERM_TYPE
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
	} else if t.wildcard == -1 {
		return false
	} else if t.wildcard == 1 {
		return true
	} else {
		for i := 0; i < len(t.Value); i++ {
			if t.Value[i] == "?" || t.Value[i] == "*" {
				t.wildcard = 1
				return true
			}
		}
		t.wildcard = -1
		return false
	}

}

// phrase term: a series of terms be surrounded with quotation, for instance "foo bar".
type PhraseTerm struct {
	Value    string `parser:"@STRING" json:"value"`
	wildcard int8
}

func (t *PhraseTerm) GetTermType() TermType {
	return PHRASE_TERM_TYPE
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
	} else if t.wildcard == -1 {
		return false
	} else if t.wildcard == 1 {
		return true
	} else {
		for _, x := range token.Scan(t.Value[1 : len(t.Value)-1]) {
			if x.GetTokenType() == token.WILDCARD_TOKEN_TYPE {
				t.wildcard = 1
				return true
			}
		}
		t.wildcard = -1
		return false
	}

}

// a regexp term is surrounded be slash, for instance /\d+\.?\d+/ in here if you want present '/' you should type '\/'
type RegexpTerm struct {
	Value string `parser:"@REGEXP" json:"value"`
}

func (t *RegexpTerm) GetTermType() TermType {
	return REGEXP_TERM_TYPE
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

//double side of range term: a term is surrounded by brace / bracket, for instance [1 TO 2] / [1 TO 2} / {1 TO 2] / {1 TO 2}
type DRangeTerm struct {
	LBRACKET string          `parser:"@(LBRACE|LBRACK) WHITESPACE*" json:"left_bracket"`
	LValue   *bnd.RangeValue `parser:"@@ WHITESPACE+ 'TO'" json:"left_value"`
	RValue   *bnd.RangeValue `parser:"WHITESPACE+ @@" json:"right_value"`
	RBRACKET string          `parser:"WHITESPACE* @(RBRACK|RBRACE)" json:"right_bracket"`
}

func (t *DRangeTerm) GetTermType() TermType {
	return RANGE_TERM_TYPE
}

func (t *DRangeTerm) GetBound() *bnd.Bound {
	if t == nil {
		return nil
	} else if t.LBRACKET == "[" && t.RBRACKET == "]" {
		return &bnd.Bound{LeftInclude: t.LValue, RightInclude: t.RValue}
	} else if t.LBRACKET == "[" && t.RBRACKET == "}" {
		return &bnd.Bound{LeftInclude: t.LValue, RightExclude: t.RValue}
	} else if t.LBRACKET == "{" && t.RBRACKET == "]" {
		return &bnd.Bound{LeftExclude: t.LValue, RightInclude: t.RValue}
	} else if t.LBRACKET == "{" && t.RBRACKET == "}" {
		return &bnd.Bound{LeftExclude: t.LValue, RightExclude: t.RValue}
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
	Symbol string          `parser:"@COMPARE" json:"symbol"`
	Value  *bnd.RangeValue `parser:"@@" json:"value"`
	drange *DRangeTerm
}

func (t *SRangeTerm) GetTermType() TermType {
	return RANGE_TERM_TYPE
}

func (t *SRangeTerm) toDRangeTerm() *DRangeTerm {
	if t == nil {
		return nil
	} else if t.drange != nil {
		return t.drange
	} else {
		if t.Symbol == ">" && t.Value != nil {
			t.drange = &DRangeTerm{LBRACKET: "{", LValue: t.Value, RValue: &bnd.RangeValue{InfinityVal: "*"}, RBRACKET: "}"}
		} else if t.Symbol == ">=" && t.Value != nil {
			t.drange = &DRangeTerm{LBRACKET: "[", LValue: t.Value, RValue: &bnd.RangeValue{InfinityVal: "*"}, RBRACKET: "}"}
		} else if t.Symbol == "<" && t.Value != nil {
			t.drange = &DRangeTerm{LBRACKET: "{", LValue: &bnd.RangeValue{InfinityVal: "*"}, RValue: t.Value, RBRACKET: "}"}
		} else if t.Symbol == "<=" && t.Value != nil {
			t.drange = &DRangeTerm{LBRACKET: "{", LValue: &bnd.RangeValue{InfinityVal: "*"}, RValue: t.Value, RBRACKET: "]"}
		}
	}
	return t.drange
}

func (t *SRangeTerm) GetBound() *bnd.Bound {
	return t.toDRangeTerm().GetBound()
}

func (t *SRangeTerm) String() string {
	return t.toDRangeTerm().String()
}
