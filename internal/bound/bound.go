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

func (b *Bound) adjust() {
	if b.LeftInclude != nil && len(b.LeftInclude.InfinityVal) != 0 {
		b.LeftExclude = &RangeValue{InfinityVal: "*"}
		b.LeftInclude = nil
	}
	if b.RightInclude != nil && len(b.RightInclude.InfinityVal) != 0 {
		b.RightExclude = &RangeValue{InfinityVal: "*"}
		b.RightInclude = nil
	}
}

type PhraseValue struct {
	Value []string `parser:"" json:"value"`
}

type RangeValue struct {
	InfinityVal string   `parser:"  @('*')" json:"infinity_val"`
	PhraseValue []string `parser:"| QUOTE @( REVERSE QUOTE | !QUOTE )* QUOTE" json:"phrase_value"`
	SingleValue []string `parser:"| @(IDENT|NUMBER|DOT|PLUS|MINUS)+" json:"simple_value"`
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
