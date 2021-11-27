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
