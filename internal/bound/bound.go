package bound

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// range bound like this [1, 2] [1, 2) (1, 2] (1, 2)
type Bound struct {
	LeftValue    *RangeValue `json:"left_value,omitempty"`
	RightValue   *RangeValue `json:"right_Value,omitempty"`
	LeftInclude  bool        `json:"left_include,omitempty"`
	RightInclude bool        `json:"right_include,omitempty"`
}

func (n *Bound) GetBoundType() BoundType {
	if n == nil {
		return UNKNOWN_BOUND_TYPE
	} else if n.LeftInclude && n.RightInclude {
		return LEFT_INCLUDE_RIGHT_INCLUDE
	} else if n.LeftInclude && !n.RightInclude {
		return LEFT_INCLUDE_RIGHT_EXCLUDE
	} else if !n.LeftInclude && n.RightInclude {
		return LEFT_EXCLUDE_RIGHT_INCLUDE
	} else {
		return LEFT_EXCLUDE_RIGHT_EXCLUDE
	}
}

type RangeValue struct {
	InfinityVal string   `parser:"  @('*')" json:"infinity_val"`
	PhraseValue []string `parser:"| QUOTE @( REVERSE QUOTE | !QUOTE )* QUOTE" json:"phrase_value"`
	SingleValue []string `parser:"| @(IDENT|NUMBER|'.'|'+'|'-'|'|'|'/'|':')+" json:"simple_value"`
}

func (v *RangeValue) String() string {
	if v == nil {
		return ""
	} else if len(v.PhraseValue) != 0 {
		return strings.Join(v.PhraseValue, "")
	} else if len(v.InfinityVal) != 0 {
		return v.InfinityVal
	} else if len(v.SingleValue) != 0 {
		return strings.Join(v.SingleValue, "")
	} else {
		return ""
	}
}

func (v *RangeValue) Int() (int, error) {
	if v == nil {
		return 0, ErrEmptyValue
	} else {
		return strconv.Atoi(v.String())
	}
}

func (v *RangeValue) Float() (float64, error) {
	if v == nil {
		return 0.0, ErrEmptyValue
	} else {
		return strconv.ParseFloat(v.String(), 64)
	}
}

func (v *RangeValue) Time(format []string) (*time.Time, error) {
	if v == nil {
		return nil, ErrEmptyValue
	} else {
		var s = v.String()
		if s == "" {
			return nil, ErrEmptyValue
		}
		var sv = strings.Split(s, "||")
		var timePart = sv[0]
		var res *time.Time
		if timePart == "now" {
			var t = time.Now()
			res = &t
		} else {
			for _, f := range format {
				if t, err := time.Parse(f, timePart); err == nil {
					res = &t
					break
				}
			}
		}
		if res == nil {
			return nil, fmt.Errorf("failed to parse date: '%s' according to format: %v", s, format)
		}
		if len(sv) > 1 {
			var dv = strings.Split(sv[1], "/")
			var durationPart = dv[0]
			if offsetDuration, err := ParseDuration(durationPart); err != nil {
				return nil, fmt.Errorf("failed to parse offset duration: '%s', err: %+v", durationPart, err)
			} else {
				res.Add(offsetDuration)
			}
			if len(dv) > 1 {
				var roundPart = dv[1]
				if roundDuration, err := ParseDuration(roundPart); err != nil {
					return nil, fmt.Errorf("failed to parse round duration: '%s', err: %+v", durationPart, err)
				} else {
					res.Round(roundDuration)
				}
			}
		}
		return res, nil
	}

}

func (v *RangeValue) IsInf() bool {
	return v != nil && len(v.InfinityVal) != 0
}
