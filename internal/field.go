package internal

import "strings"

type Field struct {
	Value []string `parser:"@(IDENT|MINUS)+"`
}

func (f *Field) String() string {
	if f == nil {
		return ""
	} else {
		return strings.Join(f.Value, "")
	}
}

type FieldTerm struct {
	Field *Field `parser:"@@" json:"field"`
	Colon string `parser:"@COLON" json:"colon"`
	Term  *Term  `parser:"@@" json:"term"`
}

func (f *FieldTerm) String() string {
	if f == nil {
		return ""
	} else if f.Field == nil || f.Term == nil {
		return ""
	} else {
		return f.Field.String() + " " + f.Colon + " " + f.Term.String()
	}
}
