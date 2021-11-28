package term

import (
	"strconv"
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
		return t.SRangeTerm.String() + t.BoostSymbol
	} else if t.DRangeTerm != nil {
		return t.DRangeTerm.String() + t.BoostSymbol
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
