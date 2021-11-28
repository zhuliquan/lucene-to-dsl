package query

import (
	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
	tm "github.com/zhuliquan/lucene-to-dsl/query/internal/term"
)

type Lucene struct {
	OrQuery *OrQuery  `parser:"@@" json:"or_query"`
	OrTerms []*OrTerm `parser:"@@*" json:"or_terms"`
}

type OrTerm struct {
	ORSymbol  *op.OrSymbol  `parser:"@@" json:"or_symbol"`
	NOTSymbol *op.NotSymbol `parser:"@@?" json:"not_symbol"`
	OrQuery   *OrQuery      `parser:"@@" json:"or_query"`
}

type OrQuery struct {
	AndQuery *AndQuery  `parser:"@@" json:"and_query"`
	AndTerms []*AndTerm `parser:"@@*" json:"and_terms" `
}

type AndTerm struct {
	ANDSymbol *op.AndSymbol `parser:"@@" json:"and_symbol"`
	NOTSymbol *op.NotSymbol `parser:"@@?" json:"not_symbol"`
	AndQuery  *AndQuery     `parser:"@@" json:"and_query"`
}

type AndQuery struct {
	NotQuery   *NotQuery   `parser:"  @@" json:"not_query"`
	ParenQuery *ParenQuery `parser:"| @@" json:"paren_query"`
	FieldQuery *FieldQuery `parser:"| @@" json:"field_query"`
}

type NotQuery struct {
	NOTSymbol *op.NotSymbol `parser:"@@" json:"not_symbol"`
	SubQuery  *Lucene       `parser:"@@" json:"sub_query"`
}

type ParenQuery struct {
	LParen   string  `parser:"@LPAREN" json:"lparen"`
	SubQuery *Lucene `parser:"@@" json:"sub_query"`
	RParen   string  `parser:"@RPAREN" json:"rparen"`
}

type FieldQuery struct {
	Field *tm.Field `parser:"@@ COLON" json:"field"`
	Term  *tm.Term  `parser:"@@" json:"term"`
}

func (f *FieldQuery) String() string {
	if f == nil {
		return ""
	} else if f.Field == nil || f.Term == nil {
		return ""
	} else {
		return f.Field.String() + " : " + f.Term.String()
	}
}
