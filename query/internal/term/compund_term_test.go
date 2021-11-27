package term

import (
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
	}
	var testCases = []testCase{
		{
			name:  "TestRangeTerm01",
			input: `<="dsada 78"`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &RangeValue{PhraseValue: `"dsada 78"`}}},
		},
		{
			name:  "TestRangeTerm02",
			input: `<=dsada\ 78`,
			want:  &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: "<=", Value: &RangeValue{SingleValue: []string{`dsada\ 78`}}}},
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
			want:  &PrefixTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm02",
			input: `+"dsada 78"`,
			want:  &PrefixTerm{Symbol: "+", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm03",
			input: `-"dsada 78"`,
			want:  &PrefixTerm{Symbol: "-", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm04",
			input: `\+\/dsada\/\ dasda80980?*`,
			want:  &PrefixTerm{SingleTerm: &SingleTerm{Value: []string{`\+\/dsada\/\ dasda80980`, `?`, `*`}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm05",
			input: `+\/dsada\/\ dasda80980?*`,
			want:  &PrefixTerm{Symbol: "+", SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm06",
			input: `-\-\/dsada\/\ dasda80980?*`,
			want:  &PrefixTerm{Symbol: "-", SingleTerm: &SingleTerm{Value: []string{`\-\/dsada\/\ dasda80980`, `?`, `*`}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm07",
			input: `->890`,
			want:  &PrefixTerm{Symbol: "-", RangeTerm: &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm08",
			input: `>890`,
			want:  &PrefixTerm{RangeTerm: &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm09",
			input: `+>=890`,
			want:  &PrefixTerm{Symbol: "+", RangeTerm: &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: ">=", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm10",
			input: `+[1 TO 2]`,
			want: &PrefixTerm{Symbol: "+", RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				TO:     "TO",
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm11",
			input: `-[1 TO 2]`,
			want: &PrefixTerm{Symbol: "-", RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				TO:     "TO",
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestPrefixTerm12",
			input: `[1 TO 2]`,
			want: &PrefixTerm{RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				TO:     "TO",
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
			want:  &WPrefixTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm02",
			input: `   +"dsada 78"`,
			want:  &WPrefixTerm{Symbol: "+", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm03",
			input: `  -"dsada 78"`,
			want:  &WPrefixTerm{Symbol: "-", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm04",
			input: `  \+\/dsada\/\ dasda80980?*`,
			want:  &WPrefixTerm{SingleTerm: &SingleTerm{Value: []string{`\+\/dsada\/\ dasda80980`, `?`, `*`}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm05",
			input: `  +\/dsada\/\ dasda80980?*`,
			want:  &WPrefixTerm{Symbol: "+", SingleTerm: &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm06",
			input: `  -\-\/dsada\/\ dasda80980?*`,
			want:  &WPrefixTerm{Symbol: "-", SingleTerm: &SingleTerm{Value: []string{`\-\/dsada\/\ dasda80980`, `?`, `*`}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm07",
			input: `  ->890`,
			want:  &WPrefixTerm{Symbol: "-", RangeTerm: &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm08",
			input: `  >890`,
			want:  &WPrefixTerm{RangeTerm: &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: ">", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.SHOULD_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm09",
			input: `  +>=890`,
			want:  &WPrefixTerm{Symbol: "+", RangeTerm: &RangeTerm{SRangeTerm: &SRangeTerm{Symbol: ">=", Value: &RangeValue{SingleValue: []string{`890`}}}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm10",
			input: `   +[1 TO 2]`,
			want: &WPrefixTerm{Symbol: "+", RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				TO:     "TO",
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm11",
			input: `  -[1 TO 2]`,
			want: &WPrefixTerm{Symbol: "-", RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				TO:     "TO",
				RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
			}}},
			oType: op.MUST_NOT_PREFIX_TYPE,
		},
		{
			name:  "TestWPrefixTerm12",
			input: `  [1 TO 2]`,
			want: &WPrefixTerm{RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
				LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
				TO:     "TO",
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
			input: `( 8908  "dsada 78" +"89080  xxx" -"xx yyyy" +\+dsada\ 7897 -\-\-dsada\-7897  )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{SingleTerm: &SingleTerm{Value: []string{"8908"}}},
				PrefixTerms: []*WPrefixTerm{
					{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
					{Symbol: "+", PhraseTerm: &PhraseTerm{Value: `"89080  xxx"`}},
					{Symbol: "-", PhraseTerm: &PhraseTerm{Value: `"xx yyyy"`}},
					{Symbol: "+", SingleTerm: &SingleTerm{Value: []string{`\+dsada\ 7897`}}},
					{Symbol: "-", SingleTerm: &SingleTerm{Value: []string{`\-\-dsada\-7897`}}},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup02",
			input: `( 8908 )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{SingleTerm: &SingleTerm{Value: []string{"8908"}}},
			},
		},
		{
			name:  "TestPrefixTermGroup03",
			input: `( 8908 [ -1 TO 3]  )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{SingleTerm: &SingleTerm{Value: []string{"8908"}}},
				PrefixTerms: []*WPrefixTerm{
					{
						RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
							LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
							TO:     "TO",
							RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
						}},
					},
				},
			},
		},
		{
			name:  "TestPrefixTermGroup04",
			input: `( +>2021-11-04 +<2021-11-11 )`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{Symbol: "+", RangeTerm: &RangeTerm{
					SRangeTerm: &SRangeTerm{
						Symbol: ">",
						Value:  &RangeValue{SingleValue: []string{`2021`, "-", "11", "-", "04"}},
					},
				}},
				PrefixTerms: []*WPrefixTerm{
					{Symbol: "+", RangeTerm: &RangeTerm{
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
				PrefixTerm: &PrefixTerm{RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
					LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
					TO:     "TO",
					RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
				}}},
			},
		},
		{
			name:  "TestPrefixTermGroup06",
			input: `( [-1 TO 3] [1 TO 2] +[5 TO 10}  -{8 TO 90])`,
			want: &PrefixTermGroup{
				PrefixTerm: &PrefixTerm{RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
					LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"-", "1"}},
					TO:     "TO",
					RValue: &RangeValue{SingleValue: []string{"3"}}, RBRACKET: "]",
				}}},
				PrefixTerms: []*WPrefixTerm{
					{RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
						LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"1"}},
						TO:     "TO",
						RValue: &RangeValue{SingleValue: []string{"2"}}, RBRACKET: "]",
					}}},
					{Symbol: "+", RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
						LBRACKET: "[", LValue: &RangeValue{SingleValue: []string{"5"}},
						TO:     "TO",
						RValue: &RangeValue{SingleValue: []string{"10"}}, RBRACKET: "}",
					}}},
					{Symbol: "-", RangeTerm: &RangeTerm{DRangeTerm: &DRangeTerm{
						LBRACKET: "{", LValue: &RangeValue{SingleValue: []string{"8"}},
						TO:     "TO",
						RValue: &RangeValue{SingleValue: []string{"90"}}, RBRACKET: "]",
					}}},
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
