package internal

// bool operator ("AND" / "OR" / "NOT") or ("and" / "or" / "not") or ("&&" / "||" / "!")
type BoolOper struct {
	AndSymbol string `parser:"  WHITESPACE+ @AND AND WHITESPACE+" json:"and_symbol"`
	OrSymbol  string `parser:"| WHITESPACE+ @OR OR WHITESPACE+" json:"or_symbol"`
	NotSymbol string `parser:"| @NOT WHITESPACE+" json:"not_symbol"`
	AndIdent  string `parser:"| WHITESPACE+ @('AND' | 'and') WHITESPACE+" json:"and_ident"`
	OrIdent   string `parser:"| WHITESPACE+ @('OR' | 'or') WHITESPACE+" json:"or_ident"`
	NotIdent  string `parser:"| @('NOT' | 'not') WHITESPACE+" json:"not_ident"`
}

func (o *BoolOper) String() string {
	if o == nil {
		return ""
	} else if o.AndSymbol != "" {
		return "AND"
	} else if o.AndIdent != "" {
		return "AND"
	} else if o.OrIdent != "" {
		return "OR"
	} else if o.OrSymbol != "" {
		return "OR"
	} else if o.NotIdent != "" {
		return "NOT"
	} else if o.NotSymbol != "" {
		return "NOT"
	} else {
		return ""
	}
}
