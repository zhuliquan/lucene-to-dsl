package operator

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

func TestAndSymbol(t *testing.T) {
	var operatorParser = participle.MustBuild(
		&ANDSymbol{},
		participle.Lexer(token.Lexer),
	)
	type testCase struct {
		name  string
		input string
		want  *ANDSymbol
	}
	var testCases = []testCase{
		{
			name:  "TestAndSymbol01",
			input: ` AND   `,
			want:  &ANDSymbol{Symbol: "AND"},
		},
		{
			name:  "TestAndSymbol02",
			input: ` and `,
			want:  &ANDSymbol{Symbol: "and"},
		},
		{
			name:  "TestAndSymbol03",
			input: ` && `,
			want:  &ANDSymbol{Symbol: "&&"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var symbol = &ANDSymbol{}
			if err := operatorParser.ParseString(tt.input, symbol); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(symbol, tt.want) {
				t.Errorf("ParseString( %s ) = %+v, want: %+v", tt.input, symbol, tt.want)
			}
		})
	}
}

func TestOrSymbol(t *testing.T) {
	var operatorParser = participle.MustBuild(
		&ORSymbol{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *ORSymbol
	}
	var testCases = []testCase{
		{
			name:  "TestOrSymbol01",
			input: ` OR  `,
			want:  &ORSymbol{Symbol: "OR"},
		},
		{
			name:  "TestOrSymbol02",
			input: ` or  `,
			want:  &ORSymbol{Symbol: "or"},
		},
		{
			name:  "TestOrSymbol03",
			input: ` ||  `,
			want:  &ORSymbol{Symbol: "||"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var symbol = &ORSymbol{}
			if err := operatorParser.ParseString(tt.input, symbol); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(symbol, tt.want) {
				t.Errorf("ParseString( %s ) = %+v, want: %+v", tt.input, symbol, tt.want)
			}
		})
	}
}

func TestNotSymbol(t *testing.T) {
	var operatorParser = participle.MustBuild(
		&NOTSymbol{},
		participle.Lexer(token.Lexer),
	)
	type testCase struct {
		name  string
		input string
		want  *NOTSymbol
	}
	var testCases = []testCase{
		{
			name:  "TestNotSymbol01",
			input: `NOT `,
			want:  &NOTSymbol{Symbol: "NOT"},
		},
		{
			name:  "TestNotSymbol02",
			input: `not `,
			want:  &NOTSymbol{Symbol: "not"},
		},
		{
			name:  "TestNotSymbol03",
			input: `! `,
			want:  &NOTSymbol{Symbol: "!"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var symbol = &NOTSymbol{}
			if err := operatorParser.ParseString(tt.input, symbol); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(symbol, tt.want) {
				t.Errorf("ParseString( %s ) = %+v, want: %+v", tt.input, symbol, tt.want)
			}
		})
	}
}
