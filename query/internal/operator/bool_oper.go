package operator

// and operator ("AND" / "and" / "&&")
type ANDSymbol struct {
	Symbol string `parser:"  WHITESPACE+ @( AND AND | 'AND' | 'and' ) WHITESPACE+" json:"symbol"`
}

func (o *ANDSymbol) String() string {
	if o == nil {
		return ""
	} else if o.Symbol != "" {
		return o.Symbol
	} else {
		return ""
	}
}

func (o *ANDSymbol) GetLogicType() LogicOPType {
	if o == nil {
		return UNKNOWN_LOGIC_TYPE
	} else {
		return AND_LOGIC_TYPE
	}
}

// or operator ("OR" / "or" / "||")
type ORSymbol struct {
	Symbol string `parser:"  WHITESPACE+ @( OR OR | 'OR' | 'or' ) WHITESPACE+" json:"symbol"`
}

func (o *ORSymbol) String() string {
	if o == nil {
		return ""
	} else if o.Symbol != "" {
		return o.Symbol
	} else {
		return ""
	}
}

func (o *ORSymbol) GetLogicType() LogicOPType {
	if o == nil {
		return UNKNOWN_LOGIC_TYPE
	} else {
		return OR_LOGIC_TYPE
	}
}

// not operator ("NOT" / "not" / "!")
type NOTSymbol struct {
	Symbol string `parser:"@( NOT | 'NOT' | 'not') WHITESPACE+" json:"symbol"`
}

func (o *NOTSymbol) String() string {
	if o == nil {
		return ""
	} else if o.Symbol != "" {
		return o.Symbol
	} else {
		return ""
	}
}

func (o *NOTSymbol) GetLogicType() LogicOPType {
	if o == nil {
		return UNKNOWN_LOGIC_TYPE
	} else {
		return NOT_LOGIC_TYPE
	}
}

// (" " ? ( "+" / "-")? )
type PreSymbol struct {
	Should  string `parser:"@WHITESPACE*" json:"should"`
	MustNOT string `parser:"( @MINUS " json:"must_not"`
	Must    string `parser:"| @PLUS )?" json:"must"`
}

func (o *PreSymbol) String() string {
	if o == nil {
		return ""
	} else if len(o.MustNOT) != 0 {
		return o.MustNOT
	} else if len(o.Must) != 0 {
		return o.Must
	} else if len(o.Should) != 0 {
		return " "
	} else {
		return ""
	}
}

func (o *PreSymbol) GetPrefixType() PrefixOPType {
	if o == nil {
		return UNKNOWN_PREFIX_TYPE
	} else if len(o.MustNOT) != 0 {
		return MUST_NOT_PREFIX_TYPE
	} else if len(o.Must) != 0 {
		return MUST_PREFIX_TYPE
	} else {
		return SHOULD_PREFIX_TYPE
	}
}
