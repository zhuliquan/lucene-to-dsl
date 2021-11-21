package query

type Lucene struct {
	OrQuery *OrQuery  `parser:"@@" json:"or_query"`
	OrTerms []*OrTerm `parser:"@@*" json:"or_terms"`
}

type OrTerm struct {
	ORSymbol  *ORSymbol  `parser:"@@" json:"or_symbol"`
	NOTSymbol *NOTSymbol `parser:"@@?" json:"not_symbol"`
	OrQuery   *OrQuery   `parser:"@@" json:"or_query"`
}

type OrQuery struct {
	AndQuery *AndQuery  `parser:"@@" json:"and_query"`
	AndTerms []*AndTerm `parser:"@@*" json:"and_terms" `
}

type AndTerm struct {
	ANDSymbol *ANDSymbol `parser:"@@" json:"and_symbol"`
	NOTSymbol *NOTSymbol `parser:"@@?" json:"not_symbol"`
	AndQuery  *AndQuery  `parser:"@@" json:"and_query"`
}

type AndQuery struct {
	NotQuery   *NotQuery   `parser:"  @@" json:"not_query"`
	ParenQuery *ParenQuery `parser:"| @@" json:"paren_query"`
	FieldQuery *FieldTerm  `parser:"| @@" json:"field_query"`
}

type NotQuery struct {
	NOTSymbol *NOTSymbol `parser:"@@" json:"not_symbol"`
	SubQuery  *Lucene    `parser:"@@" json:"sub_query"`
}

type ParenQuery struct {
	LParen   string  `parser:"@LPAREN" json:"lparen"`
	SubQuery *Lucene `parser:"@@" json:"sub_query"`
	RParen   string  `parser:"@RPAREN" json:"rparen"`
}
