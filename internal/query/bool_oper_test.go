package query

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
)

func TestBoolOperator(t *testing.T) {

	var andOperParser = participle.MustBuild(
		&ANDSymbol{},
		participle.Lexer(Lexer),
	)
	var orOperParser = participle.MustBuild(
		&ORSymbol{},
		participle.Lexer(Lexer),
	)
	var notOperParser = participle.MustBuild(
		&NOTSymbol{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name    string
		parser  interface{}
		boolObj interface{}
		input   string
		want    interface{}
	}

	var testCases = []testCase{
		{
			name:    "TestBoolOper01",
			parser:  andOperParser,
			boolObj: &ANDSymbol{},
			input:   ` AND   `,
			want:    &ANDSymbol{LAND: "AND"},
		},
		{
			name:    "TestBoolOper02",
			parser:  andOperParser,
			boolObj: &ANDSymbol{},
			input:   ` and `,
			want:    &ANDSymbol{LAND: "and"},
		},
		{
			name:    "TestBoolOper03",
			parser:  andOperParser,
			boolObj: &ANDSymbol{},
			input:   ` && `,
			want:    &ANDSymbol{SAND: "&&"},
		},
		{
			name:    "TestBoolOper04",
			parser:  orOperParser,
			boolObj: &ORSymbol{},
			input:   ` OR  `,
			want:    &ORSymbol{LOR: "OR"},
		},
		{
			name:    "TestBoolOper05",
			parser:  orOperParser,
			boolObj: &ORSymbol{},
			input:   ` or  `,
			want:    &ORSymbol{LOR: "or"},
		},
		{
			name:    "TestBoolOper06",
			parser:  orOperParser,
			boolObj: &ORSymbol{},
			input:   ` ||  `,
			want:    &ORSymbol{SOR: "||"},
		},
		{
			name:    "TestBoolOper07",
			parser:  notOperParser,
			boolObj: &NOTSymbol{},
			input:   `NOT `,
			want:    &NOTSymbol{LNOT: "NOT"},
		},
		{
			name:    "TestBoolOper08",
			parser:  notOperParser,
			boolObj: &NOTSymbol{},
			input:   `not `,
			want:    &NOTSymbol{LNOT: "not"},
		},
		{
			name:    "TestBoolOper09",
			parser:  notOperParser,
			boolObj: &NOTSymbol{},
			input:   `! `,
			want:    &NOTSymbol{SNOT: "!"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.parser.(*participle.Parser).ParseString(tt.input, tt.boolObj); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.boolObj, tt.want) {
				t.Errorf("ParseString( %s ) = %+v, want: %+v", tt.input, tt.boolObj, tt.want)
			}
		})
	}

}
