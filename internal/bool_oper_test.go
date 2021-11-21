package internal

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
)

func TestBoolOperator(t *testing.T) {

	var boolOperParser = participle.MustBuild(
		&BoolOper{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *BoolOper
	}

	var testCases = []testCase{
		{
			name:  "TestBoolOper01",
			input: ` AND   `,
			want:  &BoolOper{AndIdent: "AND"},
		},
		{
			name:  "TestBoolOper02",
			input: ` and `,
			want:  &BoolOper{AndIdent: "and"},
		},
		{
			name:  "TestBoolOper03",
			input: ` && `,
			want:  &BoolOper{AndSymbol: "&"},
		},
		{
			name:  "TestBoolOper04",
			input: ` OR  `,
			want:  &BoolOper{OrIdent: "OR"},
		},
		{
			name:  "TestBoolOper05",
			input: ` or  `,
			want:  &BoolOper{OrIdent: "or"},
		},
		{
			name:  "TestBoolOper06",
			input: ` ||  `,
			want:  &BoolOper{OrSymbol: "|"},
		},
		{
			name:  "TestBoolOper07",
			input: `NOT `,
			want:  &BoolOper{NotIdent: "NOT"},
		},
		{
			name:  "TestBoolOper08",
			input: `not `,
			want:  &BoolOper{NotIdent: "not"},
		},
		{
			name:  "TestBoolOper09",
			input: `! `,
			want:  &BoolOper{NotSymbol: "!"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &BoolOper{}
			if err := boolOperParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(out, tt.want) {
				t.Errorf("ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}

}
