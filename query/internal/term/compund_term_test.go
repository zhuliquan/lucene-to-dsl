package term

import (
	"math"
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	op "github.com/zhuliquan/lucene-to-dsl/query/internal/operator"
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
		boost float64
	}
	var testCases = []testCase{
		{
			name:  "TestRangeTerm01",
			input: `<="dsada 78"`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &RangeValue{PhraseValue: `"dsada 78"`}}},
			boost: 1.0,
		},
		{
			name:  "TestRangeTerm02",
			input: `<="dsada 78"^8.9`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &RangeValue{PhraseValue: `"dsada 78"`}}, BoostSymbol: "^8.9"},
			boost: 8.9,
		},
		{
			name:  "TestRangeTerm03",
			input: `<=dsada\ 78`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &RangeValue{SingleValue: []string{`dsada\ 78`}}}},
			boost: 1.0,
		},
		{
			name:  "TestRangeTerm04",
			input: `<=dsada\ 78^0.5`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &RangeValue{SingleValue: []string{`dsada\ 78`}}}, BoostSymbol: "^0.5"},
			boost: 0.5,
		},
		{
			name:  "TestRangeTerm05",
			input: `[1 TO 2]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}},
			boost: 1.0,
		},
		{
			name:  "TestRangeTerm06",
			input: `[1 TO 2]^0.7`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}, BoostSymbol: "^0.7"},
			boost: 0.7,
		},
		{
			name:  "TestRangeTerm07",
			input: `[1 TO 2 }`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}},
			boost: 1.0,
		},
		{
			name:  "TestRangeTerm08",
			input: `[1 TO 2 }^0.9`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}, BoostSymbol: "^0.9"},
			boost: 0.9,
		},
		{
			name:  `TestRangeTerm09`,
			input: `{ 1 TO 2}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			}},
			boost: 1.0,
		},
		{
			name:  `TestRangeTerm10`,
			input: `{ 1 TO 2]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			}},
			boost: 1.0,
		},
		{
			name:  `TestRangeTerm11`,
			input: `[10 TO *]`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"10"}},
				RValue:   &RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			}},
			boost: 1.0,
		},
		{
			name:  `TestRangeTerm12`,
			input: `{* TO 2012-01-01}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				RValue:   &RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			}},
			boost: 1.0,
		},
		{
			name:  `TestRangeTerm13`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				RValue:   &RangeValue{PhraseValue: "\"2012-01-01 09:08:16\""},
				RBRACKET: "}",
			}},
			boost: 1.0,
		},
		{
			name:  `TestRangeTerm14`,
			input: `>2012-01-01^9.8`,
			want: &RangeTerm{SRangeTerm: &SRangeTerm{
				Symbol: ">",
				Value:  &RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
			}, BoostSymbol: "^9.8"},
			boost: 9.8,
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
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada\* 78"`}},
			valueS:   `dsada\* 78`,
			wildcard: false,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm02",
			input:    `"dsada* 78"`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada* 78"`}},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm03",
			input:    `"dsada\* 78"^08`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada\* 78"`}, BoostSymbol: "^08"},
			valueS:   `dsada\* 78`,
			wildcard: false,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm04",
			input:    `"dsada* 78"^08`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada* 78"`}, BoostSymbol: "^08"},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm05",
			input:    `"dsada\* 78"~8`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada\* 78"`}, FuzzySymbol: "~8"},
			valueS:   `dsada\* 78`,
			wildcard: false,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm06",
			input:    `"dsada* 78"~8`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada* 78"`}, FuzzySymbol: "~8"},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm07",
			input:    `"dsada 78"~`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}, FuzzySymbol: "~"},
			valueS:   `dsada 78`,
			wildcard: false,
			fuzzy:    1,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm08",
			input:    `"dsada* 78"~`,
			want:     &FuzzyTerm{PhraseTerm: &PhraseTerm{Value: `"dsada* 78"`}, FuzzySymbol: "~"},
			valueS:   `dsada* 78`,
			wildcard: true,
			fuzzy:    1,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm05",
			input:    `\/dsada\/\ dasda80980?*`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}},
			valueS:   `\/dsada\/\ dasda80980?*`,
			wildcard: true,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm06",
			input:    `\/dsada\/\ dasda80980?*\^\^^08`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, BoostSymbol: `^08`},
			valueS:   `\/dsada\/\ dasda80980?*\^\^`,
			wildcard: true,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm07",
			input:    `\/dsada\/\ dasda80980?*\^\^~8`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~8`},
			valueS:   `\/dsada\/\ dasda80980?*\^\^`,
			wildcard: true,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm08",
			input:    `\/dsada\/\ dasda80980?*\^\^~`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}}, FuzzySymbol: `~`},
			valueS:   `\/dsada\/\ dasda80980?*\^\^`,
			wildcard: true,
			fuzzy:    1,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm09",
			input:    `\/dsada\/\ dasda80980\?\*`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980\?\*`}}},
			valueS:   `\/dsada\/\ dasda80980\?\*`,
			wildcard: false,
			fuzzy:    0,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm10",
			input:    `\/dsada\/\ dasda80980\?\*\^\^^08`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980\?\*\^\^`}}, BoostSymbol: `^08`},
			valueS:   `\/dsada\/\ dasda80980\?\*\^\^`,
			wildcard: false,
			fuzzy:    0,
			boost:    8.0,
		},
		{
			name:     "TestFuzzyTerm11",
			input:    `\/dsada\/\ dasda80980\?\*\^\^~8`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980\?\*\^\^`}}, FuzzySymbol: `~8`},
			valueS:   `\/dsada\/\ dasda80980\?\*\^\^`,
			wildcard: false,
			fuzzy:    8,
			boost:    1.0,
		},
		{
			name:     "TestFuzzyTerm12",
			input:    `\/dsada\/\ dasda80980\?\*\^\^~`,
			want:     &FuzzyTerm{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980\?\*\^\^`}}, FuzzySymbol: `~`},
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

func TestPrefixTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&PrefixTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *PrefixTerm
		oType op.PrefixOPType
	}
	var testCases = []testCase{
		{
			name:  "TestPrefixTerm01",
			input: `"dsada 78"`,
			want:  &PrefixTerm{Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm02",
			input: `+"dsada 78"`,
			want:  &PrefixTerm{Symbol: "+", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm03",
			input: `-"dsada 78"`,
			want:  &PrefixTerm{Symbol: "-", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm04",
			input: `\+\/dsada\/\ dasda80980?*`,
			want:  &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\+\/dsada\/\ dasda80980`, `?`, `*`}}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm05",
			input: `+\/dsada\/\ dasda80980?*`,
			want: &PrefixTerm{Symbol: "+", Elem: &TermGroupElem{
				SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}},
			}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm06",
			input: `-\-\/dsada\/\ dasda80980?*`,
			want: &PrefixTerm{Symbol: "-", Elem: &TermGroupElem{
				SingleTerm: &SingleTerm{Value: []string{`\-\/dsada\/\ dasda80980`, `?`, `*`}},
			}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm07",
			input: `->890`,
			want: &PrefixTerm{Symbol: "-", Elem: &TermGroupElem{
				SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}},
			}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm08",
			input: `>890`,
			want:  &PrefixTerm{Elem: &TermGroupElem{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm09",
			input: `+>=890`,
			want:  &PrefixTerm{Symbol: "+", Elem: &TermGroupElem{SRangeTerm: &SRangeTerm{Symbol: ">=", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm10",
			input: `+[1 TO 2]`,
			want: &PrefixTerm{Symbol: "+", Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm11",
			input: `-[1 TO 2]`,
			want: &PrefixTerm{Symbol: "-", Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm12",
			input: `[1 TO 2]`,
			want: &PrefixTerm{Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &PrefixTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if out.GetPrefixType() != tt.oType {
				t.Errorf("expect get type: %+v, but get type: %+v", tt.oType, out.GetPrefixType())
			}
		})
	}
}

func TestWPrefixTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&WPrefixTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *WPrefixTerm
		oType op.PrefixOPType
	}
	var testCases = []testCase{
		{
			name:  "TestWPrefixTerm01",
			input: `  "dsada 78"`,
			want:  &WPrefixTerm{Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm02",
			input: `   +"dsada 78"`,
			want:  &WPrefixTerm{Symbol: "+", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm03",
			input: `  -"dsada 78"`,
			want:  &WPrefixTerm{Symbol: "-", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm04",
			input: `  \+\/dsada\/\ dasda80980?*`,
			want:  &WPrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\+\/dsada\/\ dasda80980`, `?`, `*`}}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm05",
			input: `  +\/dsada\/\ dasda80980?*`,
			want:  &WPrefixTerm{Symbol: "+", Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm06",
			input: `  -\-\/dsada\/\ dasda80980?*`,
			want:  &WPrefixTerm{Symbol: "-", Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\-\/dsada\/\ dasda80980`, `?`, `*`}}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm07",
			input: `  ->890`,
			want:  &WPrefixTerm{Symbol: "-", Elem: &TermGroupElem{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm08",
			input: `  >890`,
			want:  &WPrefixTerm{Elem: &TermGroupElem{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm09",
			input: `  +>=890`,
			want:  &WPrefixTerm{Symbol: "+", Elem: &TermGroupElem{SRangeTerm: &SRangeTerm{Symbol: ">=", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm10",
			input: `   +[1 TO 2]`,
			want: &WPrefixTerm{Symbol: "+", Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm11",
			input: `  -[1 TO 2]`,
			want: &WPrefixTerm{Symbol: "-", Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm12",
			input: `  [1 TO 2]`,
			want: &WPrefixTerm{Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &WPrefixTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if out.GetPrefixType() != tt.oType {
				t.Errorf("expect get type: %+v, but get type: %+v", tt.oType, out.GetPrefixType())
			}
		})
	}
}

func TestPrefixTermGroup(t *testing.T) {
	var termParser = participle.MustBuild(
		&PrefixTermGroup{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *PrefixTermGroup
	}
	var testCases = []testCase{
		{
			name:  "TestPrefixTermGroup01",
			input: `8908  "dsada 78" +"89080  xxx" -"xx yyyy" +\+dsada\ 7897 -\-\-dsada\-7897`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{"8908"}}}},
				PrefixTerms: []*WPrefixTerm{
					{Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
					{Symbol: "+", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"89080  xxx"`}}},
					{Symbol: "-", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"xx yyyy"`}}},
					{Symbol: "+", Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\+dsada\ 7897`}}}},
					{Symbol: "-", Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\-\-dsada\-7897`}}}},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup02",
			input: `8908`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{"8908"}}}},
			},
		},
		{
			name:  "TestPrefixTermGroup03",
			input: `8908 [ -1 TO 3]`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{"8908"}}}},
				PrefixTerms: []*WPrefixTerm{
					{
						Elem: &TermGroupElem{DRangeTerm: &DRangeTerm{
							LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
							RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
						}},
					},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup04",
			input: `+>2021-11-04 +<2021-11-11`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Symbol: "+", Elem: &TermGroupElem{
					SRangeTerm: &SRangeTerm{
						Symbol: ">",
						Value:  &RangeValue{SingleValue: []string{`2021`, "-", "11", "-", "04"}},
					},
				}},
				PrefixTerms: []*WPrefixTerm{
					{Symbol: "+", Elem: &TermGroupElem{
						SRangeTerm: &SRangeTerm{
							Symbol: "<",
							Value:  &RangeValue{SingleValue: []string{`2021`, "-", "11", "-", "11"}},
						},
					}},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup05",
			input: `[-1 TO 3]`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{
					DRangeTerm: &DRangeTerm{
						LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
						RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
					}}},
			},
		},
		{
			name:  "TestPrefixTermGroup06",
			input: `[-1 TO 3] [1 TO 2] +[5 TO 10}  -{8 TO 90]`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{
					DRangeTerm: &DRangeTerm{
						LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
						RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
					}}},
				PrefixTerms: []*WPrefixTerm{
					{Elem: &TermGroupElem{
						DRangeTerm: &DRangeTerm{
							LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
							RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
						},
					}},
					{Symbol: "+", Elem: &TermGroupElem{
						DRangeTerm: &DRangeTerm{
							LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"5"}},
							RValue: &RangeValue{SingleValue: []string{"10"}}, RBRACKET: "}",
						},
					}},
					{Symbol: "-", Elem: &TermGroupElem{
						DRangeTerm: &DRangeTerm{
							LBRACKET: "{", LValue: &RangeValue{SingleValue: []string{"8"}},
							RValue: &RangeValue{SingleValue: []string{"90"}}, RBRACKET: "]",
						},
					}},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &PrefixTermGroup{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestTermGroup(t *testing.T) {
	var termParser = participle.MustBuild(
		&PrefixTermGroup{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *PrefixTermGroup
	}
	var testCases = []testCase{
		{
			name:  "TestPrefixTermGroup01",
			input: `( 8908  "dsada 78" +"89080  xxx" -"xx yyyy" +\+dsada\ 7897 -\-\-dsada\-7897  )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{"8908"}}}},
				PrefixTerms: []*WPrefixTerm{
					{Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
					{Symbol: "+", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"89080  xxx"`}}},
					{Symbol: "-", Elem: &TermGroupElem{PhraseTerm: &PhraseTerm{Value: `"xx yyyy"`}}},
					{Symbol: "+", Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\+dsada\ 7897`}}}},
					{Symbol: "-", Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{`\-\-dsada\-7897`}}}},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup02",
			input: `( 8908 )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{"8908"}}}},
			},
		},
		{
			name:  "TestPrefixTermGroup03",
			input: `( 8908 [ -1 TO 3]  )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{SingleTerm: &SingleTerm{Value: []string{"8908"}}}},
				PrefixTerms: []*WPrefixTerm{
					{
						Elem: &TermGroupElem{
							DRangeTerm: &DRangeTerm{
								LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
								RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
							},
						},
					},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup04",
			input: `( +>2021-11-04 +<2021-11-11 )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Symbol: "+", Elem: &TermGroupElem{
					SRangeTerm: &SRangeTerm{
						Symbol: ">",
						Value:  &RangeValue{SingleValue: []string{`2021`, "-", "11", "-", "04"}},
					},
				}},
				PrefixTerms: []*WPrefixTerm{
					{Symbol: "+", Elem: &TermGroupElem{
						SRangeTerm: &SRangeTerm{
							Symbol: "<",
							Value:  &RangeValue{SingleValue: []string{`2021`, "-", "11", "-", "11"}},
						},
					}},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup05",
			input: `( [-1 TO 3]  )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{
					DRangeTerm: &DRangeTerm{
						LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
						RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
					},
				}},
			},
		},
		{
			name:  "TestPrefixTermGroup06",
			input: `( [-1 TO 3] [1 TO 2] +[5 TO 10}  -{8 TO 90])`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Elem: &TermGroupElem{
					DRangeTerm: &DRangeTerm{
						LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
						RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
					},
				}},
				PrefixTerms: []*WPrefixTerm{
					{Elem: &TermGroupElem{
						DRangeTerm: &DRangeTerm{
							LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
							RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
						},
					}},
					{Symbol: "+", Elem: &TermGroupElem{
						DRangeTerm: &DRangeTerm{
							LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"5"}},
							RValue: &RangeValue{SingleValue: []string{"10"}}, RBRACKET: "}",
						},
					}},
					{Symbol: "-", Elem: &TermGroupElem{
						DRangeTerm: &DRangeTerm{
							LBRACKET: "{", LValue: &RangeValue{SingleValue: []string{"8"}},
							RValue: &RangeValue{SingleValue: []string{"90"}}, RBRACKET: "]",
						},
					}},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &PrefixTermGroup{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}
