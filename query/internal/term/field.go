package term

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
