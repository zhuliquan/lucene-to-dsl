package bound

import (
	"strings"
)

// range bound like this [1, 2] [1, 2) (1, 2] (1, 2)
type Bound struct {
	LeftInclude  *RangeValue `json:"left_include"`
	LeftExclude  *RangeValue `json:"left_exclude"`
	RightInclude *RangeValue `json:"right_include"`
	RightExclude *RangeValue `json:"right_exclude"`
}

func (b *Bound) AdjustInf() {
	if b.LeftInclude.IsInf() {
		b.LeftExclude = Inf
		b.LeftInclude = nil
	}
	if b.RightInclude.IsInf() {
		b.RightExclude = Inf
		b.RightInclude = nil
	}
}

// func (b *Bound) ToDSL(field string ) (dsl.DSL, bool) {
// 	b.adjustInf()
// 	if b.LeftExclude.IsInf() && b.RightExclude.IsInf() {
// 		// 变为exists标识
// 		return nil, false
// 	} else if b.LeftExclude.IsInf() && !b.RightExclude.IsInf() {
// 		return DSL{}
// 	}

// }

type PhraseValue struct {
	Value []string `parser:"" json:"value"`
}

type RangeValue struct {
	InfinityVal string   `parser:"  @('*')" json:"infinity_val"`
	PhraseValue []string `parser:"| QUOTE @( REVERSE QUOTE | !QUOTE )* QUOTE" json:"phrase_value"`
	SingleValue []string `parser:"| @(IDENT|NUMBER|DOT|PLUS|MINUS)+" json:"simple_value"`
}

func (v *RangeValue) IsInf() bool {
	return v != nil && len(v.InfinityVal) != 0
}

func (v *RangeValue) String() string {
	if v == nil {
		return ""
	} else if v.PhraseValue != nil {
		return strings.Join(v.PhraseValue, "")
	} else if len(v.InfinityVal) != 0 {
		return v.InfinityVal
	} else if len(v.SingleValue) != 0 {
		return strings.Join(v.SingleValue, "")
	} else {
		return ""
	}
}
