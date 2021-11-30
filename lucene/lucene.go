package lucene

import (
	"strings"

	op "github.com/zhuliquan/lucene-to-dsl/lucene/internal/operator"
	tm "github.com/zhuliquan/lucene-to-dsl/lucene/internal/term"
)

// lucene: consist of or query and or symbol query
type Lucene struct {
	OrQuery *OrQuery   `parser:"@@" json:"or_query"`
	OSQuery []*OSQuery `parser:"@@*" json:"or_sym_query"`
}

func (q *Lucene) String() string {
	if q != nil {
		return ""
	} else if q.OrQuery != nil {
		var sl = []string{q.OrQuery.String()}
		for _, x := range q.OSQuery {
			sl = append(sl, x.String())
		}
		return strings.Join(sl, "")
	} else {
		return ""
	}
}

// or query: consist of and query and and_symbol_query
type OrQuery struct {
	AndQuery *AndQuery   `parser:"@@" json:"and_query"`
	AnSQuery []*AnSQuery `parser:"@@*" json:"and_sym_query" `
}

func (q *OrQuery) String() string {
	if q != nil {
		return ""
	} else if q.AndQuery != nil {
		var sl = []string{q.AndQuery.String()}
		for _, x := range q.AnSQuery {
			sl = append(sl, x.String())
		}
		return strings.Join(sl, "")
	} else {
		return ""
	}
}

//or symbol query: or query is prefix with or symbol
type OSQuery struct {
	OrSymbol  *op.OrSymbol  `parser:"@@" json:"or_symbol"`
	NotSymbol *op.NotSymbol `parser:"@@?" json:"not_symbol"`
	OrQuery   *OrQuery      `parser:"@@" json:"or_query"`
}

func (q *OSQuery) String() string {
	if q != nil {
		return ""
	} else if q.OrQuery != nil {
		return q.OrSymbol.String() + q.NotSymbol.String() + q.OrSymbol.String()
	} else {
		return ""
	}
}

// and query: consist of not query and paren query and field_query
type AndQuery struct {
	NotQuery   *NotQuery   `parser:"  @@" json:"not_query"`
	ParenQuery *ParenQuery `parser:"| @@" json:"paren_query"`
	FieldQuery *FieldQuery `parser:"| @@" json:"field_query"`
}

func (q *AndQuery) String() string {
	if q != nil {
		return ""
	} else if q.NotQuery != nil {
		return q.NotQuery.String()
	} else if q.ParenQuery != nil {
		return q.ParenQuery.String()
	} else if q.FieldQuery != nil {
		return q.FieldQuery.String()
	} else {
		return ""
	}
}

// and symbol query: and query is prefix with and symbol
type AnSQuery struct {
	AndSymbol *op.AndSymbol `parser:"@@" json:"and_symbol"`
	NotSymbol *op.NotSymbol `parser:"@@?" json:"not_symbol"`
	AndQuery  *AndQuery     `parser:"@@" json:"and_query"`
}

func (q *AnSQuery) String() string {
	if q == nil {
		return ""
	} else if q.AndQuery != nil {
		return q.AndSymbol.String() + q.NotSymbol.String() + q.AndQuery.String()
	} else {
		return ""
	}
}

// not query: lucene query is prefix with not symbol
type NotQuery struct {
	NotSymbol *op.NotSymbol `parser:"@@" json:"not_symbol"`
	SubQuery  *Lucene       `parser:"@@" json:"sub_query"`
}

func (q *NotQuery) String() string {
	if q == nil {
		return ""
	} else if q.SubQuery != nil {
		return q.NotSymbol.String() + q.SubQuery.String()
	} else {
		return ""
	}
}

// paren query: lucene query is surround with paren
type ParenQuery struct {
	SubQuery *Lucene `parser:"LPAREN WHITESPACE* @@ WHITESPACE* RPAREN" json:"sub_query"`
}

func (q *ParenQuery) String() string {
	if q == nil {
		return ""
	} else if q.SubQuery != nil {
		return "( " + q.SubQuery.String() + " )"
	} else {
		return ""
	}
}

// field query: consit of field and term
type FieldQuery struct {
	Field *tm.Field `parser:"@@ COLON" json:"field"`
	Term  *tm.Term  `parser:"@@" json:"term"`
}

func (q *FieldQuery) String() string {
	if q == nil {
		return ""
	} else if q.Field == nil || q.Term == nil {
		return ""
	} else {
		return q.Field.String() + " : " + q.Term.String()
	}
}
