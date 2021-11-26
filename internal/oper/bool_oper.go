package query

// bool operator ("AND" / "OR" / "NOT") or ("and" / "or" / "not") or ("&&" / "||" / "!")
type ANDSymbol struct {
	SAND string `parser:"  WHITESPACE+ @(AND AND) WHITESPACE+" json:"sand"`
	LAND string `parser:"| WHITESPACE+ @('AND' | 'and') WHITESPACE+" json:"land"`
}

func (o *ANDSymbol) String() string {
	if o == nil {
		return ""
	} else if o.SAND != "" {
		return o.SAND
	} else if o.LAND != "" {
		return o.LAND
	} else {
		return ""
	}
}

type ORSymbol struct {
	SOR string `parser:"  WHITESPACE+ @(OR OR) WHITESPACE+" json:"sor"`
	LOR string `parser:"| WHITESPACE+ @('OR' | 'or') WHITESPACE+" json:"lor"`
}

func (o *ORSymbol) String() string {
	if o == nil {
		return ""
	} else if o.SOR != "" {
		return o.SOR
	} else if o.LOR != "" {
		return o.LOR
	} else {
		return ""
	}
}

type NOTSymbol struct {
	SNOT string `parser:"  @NOT WHITESPACE+" json:"snot"`
	LNOT string `parser:"| @('NOT' | 'not') WHITESPACE+" json:"lnot"`
}

func (o *NOTSymbol) String() string {
	if o == nil {
		return ""
	} else if o.SNOT != "" {
		return "!"
	} else if o.LNOT != "" {
		return o.LNOT
	} else {
		return ""
	}
}

// ("+" / "-")
type PreSymbol struct {
	MustNOT string `parser:"  @MINUS" json:"must_not"`
	Must    string `parser:"| @PLUS" json:"must"`
}

func (o *PreSymbol) String() string {
	if o == nil {
		return ""
	} else if len(o.MustNOT) != 0 {
		return o.MustNOT
	} else if len(o.Must) != 0 {
		return o.Must
	} else {
		return ""
	}
}
