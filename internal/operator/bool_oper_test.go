package operator

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/internal/token"
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
			want:  &ANDSymbol{LAND: "AND"},
		},
		{
			name:  "TestAndSymbol02",
			input: ` and `,
			want:  &ANDSymbol{LAND: "and"},
		},
		{
			name:  "TestAndSymbol03",
			input: ` && `,
			want:  &ANDSymbol{SAND: "&&"},
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
			want:  &ORSymbol{LOR: "OR"},
		},
		{
			name:  "TestOrSymbol02",
			input: ` or  `,
			want:  &ORSymbol{LOR: "or"},
		},
		{
			name:  "TestOrSymbol03",
			input: ` ||  `,
			want:  &ORSymbol{SOR: "||"},
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
			want:  &NOTSymbol{LNOT: "NOT"},
		},
		{
			name:  "TestNotSymbol02",
			input: `not `,
			want:  &NOTSymbol{LNOT: "not"},
		},
		{
			name:  "TestNotSymbol03",
			input: `! `,
			want:  &NOTSymbol{SNOT: "!"},
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

func TestPreSymbol(t *testing.T) {
	var operatorParser = participle.MustBuild(
		&PreSymbol{},
		participle.Lexer(token.Lexer),
	)
	type testCase struct {
		name  string
		input string
		want  *PreSymbol
	}
	var testCases = []testCase{
		{
			name:  "TestPreSymbol01",
			input: `-`,
			want:  &PreSymbol{MustNOT: "-"},
		},
		{
			name:  "TestPreSymbol02",
			input: `+`,
			want:  &PreSymbol{Must: "+"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var symbol = &PreSymbol{}
			if err := operatorParser.ParseString(tt.input, symbol); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(symbol, tt.want) {
				t.Errorf("ParseString( %s ) = %+v, want: %+v", tt.input, symbol, tt.want)
			}
		})
	}
}
