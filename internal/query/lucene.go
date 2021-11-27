package query

import (
	op "github.com/zhuliquan/lucene-to-dsl/internal/operator"
	"github.com/zhuliquan/lucene-to-dsl/internal/term"
)

type Lucene struct {
	OrQuery *OrQuery  `parser:"@@" json:"or_query"`
	OrTerms []*OrTerm `parser:"@@*" json:"or_terms"`
}

type OrTerm struct {
	ORSymbol  *op.ORSymbol  `parser:"@@" json:"or_symbol"`
	NOTSymbol *op.NOTSymbol `parser:"@@?" json:"not_symbol"`
	OrQuery   *OrQuery      `parser:"@@" json:"or_query"`
}

type OrQuery struct {
	AndQuery *AndQuery  `parser:"@@" json:"and_query"`
	AndTerms []*AndTerm `parser:"@@*" json:"and_terms" `
}

type AndTerm struct {
	ANDSymbol *op.ANDSymbol `parser:"@@" json:"and_symbol"`
	NOTSymbol *op.NOTSymbol `parser:"@@?" json:"not_symbol"`
	AndQuery  *AndQuery     `parser:"@@" json:"and_query"`
}

type AndQuery struct {
	NotQuery   *NotQuery   `parser:"  @@" json:"not_query"`
	ParenQuery *ParenQuery `parser:"| @@" json:"paren_query"`
	FieldQuery *FieldQuery `parser:"| @@" json:"field_query"`
}

type NotQuery struct {
	NOTSymbol *op.NOTSymbol `parser:"@@" json:"not_symbol"`
	SubQuery  *Lucene       `parser:"@@" json:"sub_query"`
}

type ParenQuery struct {
	LParen   string  `parser:"@LPAREN" json:"lparen"`
	SubQuery *Lucene `parser:"@@" json:"sub_query"`
	RParen   string  `parser:"@RPAREN" json:"rparen"`
}

type FieldQuery struct {
	Field *term.Field `parser:"@@ COLON" json:"field"`
	Term  *term.Term  `parser:"@@" json:"term"`
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
