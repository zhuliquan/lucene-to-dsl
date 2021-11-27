package term

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

func TestRangeTerm(t *testing.T) {
	var rangesTermParser = participle.MustBuild(
		&RangeTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *RangeTerm
	}
	var testCases = []testCase{
		{
			name:  "TestRangeTerm01",
			input: `<="dsada 78"`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
		},
		{
			name:  "TestRangeTerm02",
			input: `<=dsada\ 78`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", SingleTerm: &SingleTerm{Value: []string{`dsada\ 78`}}}},
		},
		{
			name:  "TestRangeTerm03",
			input: `[1 TO 2]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}},
		},
		{
			name:  "TestRangeTerm04",
			input: `[1 TO 2 }`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}},
		},
		{
			name:  `TestRangeTerm05`,
			input: `{ 1 TO 2}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}},
		},
		{
			name:  `TestRangeTerm06`,
			input: `{ 1 TO 2]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}},
		},
		{
			name:  `TestRangeTerm07`,
			input: `[10 TO *]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"10"}},
				TO:       "TO",
				RValue:   &RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			}},
		},
		{
			name:  `TestRangeTerm08`,
			input: `{* TO 2012-01-01}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			}},
		},
		{
			name:  `TestRangeTerm09`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{PhraseValue: "\"2012-01-01 09:08:16\""},
				RBRACKET: "}",
			}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &RangeTerm{}
			if err := rangesTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangesTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestFuzzyTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&FuzzyTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *FuzzyTerm
	}
	var testCases = []testCase{
		{
			name:  "TestFuzzyTerm01",
			input: `"dsada 78"`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
		},
		{
			name:  "TestFuzzyTerm02",
			input: `"dsada 78"^08`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, BoostSymbol: "^08"},
		},
		{
			name:  "TestFuzzyTerm03",
			input: `"dsada 78"~8`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, FuzzySymbol: "~8"},
		},
		{
			name:  "TestFuzzyTerm04",
			input: `"dsada 78"~`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, FuzzySymbol: "~"},
		},
		{
			name:  "TestFuzzyTerm05",
			input: `\/dsada\/\ dasda80980?*`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}},
		},
		{
			name:  "TestFuzzyTerm06",
			input: `\/dsada\/\ dasda80980?*\^\^^08`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, BoostSymbol: `^08`},
		},
		{
			name:  "TestFuzzyTerm07",
			input: `\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~8`},
		},
		{
			name:  "TestFuzzyTerm08",
			input: `\/dsada\/\ dasda80980?*\^\^~`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~`},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &FuzzyTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("fuzzyTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestPrefixTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&FuzzyTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *FuzzyTerm
	}
	var testCases = []testCase{
		{
			name:  "TestFuzzyTerm01",
			input: `"dsada 78"`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
		},
		{
			name:  "TestFuzzyTerm02",
			input: `"dsada 78"^08`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, BoostSymbol: "^08"},
		},
		{
			name:  "TestFuzzyTerm03",
			input: `"dsada 78"~8`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, FuzzySymbol: "~8"},
		},
		{
			name:  "TestFuzzyTerm04",
			input: `"dsada 78"~`,
			want:  &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, FuzzySymbol: "~"},
		},
		{
			name:  "TestFuzzyTerm05",
			input: `\/dsada\/\ dasda80980?*`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}},
		},
		{
			name:  "TestFuzzyTerm06",
			input: `\/dsada\/\ dasda80980?*\^\^^08`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, BoostSymbol: `^08`},
		},
		{
			name:  "TestFuzzyTerm07",
			input: `\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~8`},
		},
		{
			name:  "TestFuzzyTerm08",
			input: `\/dsada\/\ dasda80980?*\^\^~`,
			want:  &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~`},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &FuzzyTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("fuzzyTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}
