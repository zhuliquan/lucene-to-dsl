package term

import (
	"math"
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	bnd "github.com/zhuliquan/lucene-to-dsl/lucene/internal/bound"
	"github.com/zhuliquan/lucene-to-dsl/lucene/internal/token"
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
		boost float64
		bound *bnd.Bound
	}
	var testCases = []testCase{
		{
			name:  "TestRangeTerm01",
			input: `<="dsada 78"`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &bnd.RangeValue{PhraseValue: []string{`dsada`, ` `, `78`}}}},
			boost: 1.0,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{InfinityVal: "*"}, RightInclude: &bnd.RangeValue{PhraseValue: []string{`dsada`, ` `, `78`}}},
		},
		{
			name:  "TestRangeTerm02",
			input: `<="dsada 78"^8.9`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &bnd.RangeValue{PhraseValue: []string{`dsada`, ` `, `78`}}}, BoostSymbol: "^8.9"},
			boost: 8.9,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{InfinityVal: "*"}, RightInclude: &bnd.RangeValue{PhraseValue: []string{`dsada`, ` `, `78`}}},
		},
		{
			name:  "TestRangeTerm03",
			input: `<=dsada\ 78`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}}}},
			boost: 1.0,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{InfinityVal: "*"}, RightInclude: &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}}},
		},
		{
			name:  "TestRangeTerm04",
			input: `<=dsada\ 78^0.5`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}}}, BoostSymbol: "^0.5"},
			boost: 0.5,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{InfinityVal: "*"}, RightInclude: &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}}},
		},
		{
			name:  "TestRangeTerm05",
			input: `[1 TO 2]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}},
			boost: 1.0,
			bound: &bnd.Bound{LeftInclude: &bnd.RangeValue{SingleValue: []string{`1`}}, RightInclude: &bnd.RangeValue{SingleValue: []string{"2"}}},
		},
		{
			name:  "TestRangeTerm06",
			input: `[1 TO 2]^0.7`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}, BoostSymbol: "^0.7"},
			boost: 0.7,
			bound: &bnd.Bound{LeftInclude: &bnd.RangeValue{SingleValue: []string{`1`}}, RightInclude: &bnd.RangeValue{SingleValue: []string{"2"}}},
		},
		{
			name:  "TestRangeTerm07",
			input: `[1 TO 2 }`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}},
			boost: 1.0,
			bound: &bnd.Bound{LeftInclude: &bnd.RangeValue{SingleValue: []string{`1`}}, RightExclude: &bnd.RangeValue{SingleValue: []string{"2"}}},
		},
		{
			name:  "TestRangeTerm08",
			input: `[1 TO 2 }^0.9`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}, BoostSymbol: "^0.9"},
			boost: 0.9,
			bound: &bnd.Bound{LeftInclude: &bnd.RangeValue{SingleValue: []string{`1`}}, RightExclude: &bnd.RangeValue{SingleValue: []string{"2"}}},
		},
		{
			name:  `TestRangeTerm09`,
			input: `{ 1 TO 2}^7`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}, BoostSymbol: "^7"},
			boost: 7.0,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{SingleValue: []string{`1`}}, RightExclude: &bnd.RangeValue{SingleValue: []string{"2"}}},
		},
		{
			name:  `TestRangeTerm10`,
			input: `{ 1 TO 2]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}},
			boost: 1.0,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{SingleValue: []string{`1`}}, RightInclude: &bnd.RangeValue{SingleValue: []string{"2"}}},
		},
		{
			name:  `TestRangeTerm11`,
			input: `[10 TO *]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"10"}},
				RValue:   &bnd.RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			}},
			boost: 1.0,
			bound: &bnd.Bound{LeftInclude: &bnd.RangeValue{SingleValue: []string{`10`}}, RightInclude: &bnd.RangeValue{InfinityVal: "*"}},
		},
		{
			name:  `TestRangeTerm12`,
			input: `{* TO 2012-01-01}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{InfinityVal: "*"},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			}},
			boost: 1.0,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{InfinityVal: "*"}, RightExclude: &bnd.RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}}},
		},
		{
			name:  `TestRangeTerm13`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{InfinityVal: "*"},
				RValue:   &bnd.RangeValue{PhraseValue: []string{"2012", "-", "01", "-", "01", " ", "09", ":", "08", ":", "16"}},
				RBRACKET: "}",
			}},
			boost: 1.0,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{InfinityVal: "*"}, RightExclude: &bnd.RangeValue{PhraseValue: []string{"2012", "-", "01", "-", "01", " ", "09", ":", "08", ":", "16"}}},
		},
		{
			name:  `TestRangeTerm14`,
			input: `>2012-01-01^9.8`,
			want: &RangeTerm{SRangeTerm: &SRangeTerm{
				Symbol: ">",
				Value:  &bnd.RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
			}, BoostSymbol: "^9.8"},
			boost: 9.8,
			bound: &bnd.Bound{LeftExclude: &bnd.RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}}, RightExclude: &bnd.RangeValue{InfinityVal: "*"}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &RangeTerm{}
			if err := rangesTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangesTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if math.Abs(tt.boost-out.Boost()) > 1E-6 {
				t.Errorf("expect get boost: %f, but get boost: %f", tt.boost, out.Boost())
			} else if !reflect.DeepEqual(out.GetBound(), tt.bound) {
				t.Errorf("expect get bound: %+v, but get bound: %+v", tt.bound, out.GetBound())
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
		name     string
		input    string
		want     *FuzzyTerm
		valueS   string
		wildcard bool
		boost    float64
		fuzzy    int
	}
	var testCases = []testCase{
		{
			name:     "TestFuzzyTerm01",
			input:    `"dsada\* 78"`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada\*`, ` `, `78`}}},
			valueS:   `dsada\* 78`,
			wildcard: false,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm02",
			input:    `"dsada* 78"`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada`, `*`, ` `, `78`}}},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm03",
			input:    `"dsada\* 78"^08`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada\*`, ` `, `78`}}, BoostSymbol: "^08"},
			valueS:   `dsada\* 78`,
			wildcard: false,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm04",
			input:    `"dsada* 78"^08`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada`, `*`, ` `, `78`}}, BoostSymbol: "^08"},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm05",
			input:    `"dsada\* 78"~8`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada\*`, ` `, `78`}}, FuzzySymbol: "~8"},
			valueS:   `dsada\* 78`,
			wildcard: false,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm06",
			input:    `"dsada* 78"~8`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada`, `*`, ` `, `78`}}, FuzzySymbol: "~8"},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm07",
			input:    `"dsada 78"~`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada`, ` `, `78`}}, FuzzySymbol: "~"},
			valueS:   `dsada 78`,
			wildcard: false,
			fuzzy:    1,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm08",
			input:    `"dsada* 78"~`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: []string{`dsada`, `*`, ` `, `78`}}, FuzzySymbol: "~"},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    1,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm09",
			input:    `\/dsada\/\ dasda80980?*`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `?`, `*`}}},
			valueS:   `\/dsada\/\ dasda80980?*`,
			wildcard: true,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm10",
			input:    `\/dsada\/\ dasda80980?*\^\^^08`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `?`, `*`, `\^\^`}}, BoostSymbol: `^08`},
			valueS:   `\/dsada\/\ dasda80980?*\^\^`,
			wildcard: true,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm11",
			input:    `\/dsada\/\ dasda80980?*\^\^~8`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~8`},
			valueS:   `\/dsada\/\ dasda80980?*\^\^`,
			wildcard: true,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm12",
			input:    `\/dsada\/\ dasda80980?*\^\^~`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~`},
			valueS:   `\/dsada\/\ dasda80980?*\^\^`,
			wildcard: true,
			fuzzy:    1,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm13",
			input:    `\/dsada\/\ dasda80980\?\*`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `\?\*`}}},
			valueS:   `\/dsada\/\ dasda80980\?\*`,
			wildcard: false,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm14",
			input:    `\/dsada\/\ dasda80980\?\*\^\^^08`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `\?\*\^\^`}}, BoostSymbol: `^08`},
			valueS:   `\/dsada\/\ dasda80980\?\*\^\^`,
			wildcard: false,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm15",
			input:    `\/dsada\/\ dasda80980\?\*\^\^~8`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `\?\*\^\^`}}, FuzzySymbol: `~8`},
			valueS:   `\/dsada\/\ dasda80980\?\*\^\^`,
			wildcard: false,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm16",
			input:    `\/dsada\/\ dasda80980\?\*\^\^~`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `\?\*\^\^`}}, FuzzySymbol: `~`},
			valueS:   `\/dsada\/\ dasda80980\?\*\^\^`,
			wildcard: false,
			fuzzy:    1,
			boost:    1.0,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &FuzzyTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("fuzzyTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if out.ValueS() != tt.valueS {
				t.Errorf("expect get values: %s but get values: %s", tt.valueS, out.ValueS())
			} else if out.haveWildcard() != tt.wildcard {
				t.Errorf("expect wildcard: %+v, but wildcard: %+v", tt.wildcard, out.haveWildcard())
			} else if out.Fuzziness() != tt.fuzzy {
				t.Errorf("expect get fuzzy: %d, but get fuzzy: %d", tt.fuzzy, out.Fuzziness())
			} else if math.Abs(out.Boost()-tt.boost) > 1E-6 {
				t.Errorf("expect get boost: %f, but get boost: %f", tt.boost, out.Boost())
			}
		})
	}
}
