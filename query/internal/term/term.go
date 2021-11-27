package term

type Term struct {
	RegexpTerm *RegexpTerm `parser:"  @@" json:"regexp_term"`
	FuzzyTerm  *FuzzyTerm  `parser:"| @@" json:"ranges_term"`
	BoostTerm  *BoostTerm  `parser:"| @@" json:"range_term"`
}

func (t *Term) String() string {
	if t == nil {
		return ""
	} else if t.RegexpTerm != nil {
		return t.RegexpTerm.String()
	} else if t.FuzzyTerm != nil {
		return t.FuzzyTerm.String()
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
	} else if t.FuzzyTerm != nil {
		return t.FuzzyTerm.Fuzziness()
	} else {
		return 0
	}
}

func (t *Term) boost() float64 {
	if t == nil {
		return 0.0
	} else if t.FuzzyTerm != nil {
		return t.FuzzyTerm.Boost()
	} else {
		return 0.0
	}
}
