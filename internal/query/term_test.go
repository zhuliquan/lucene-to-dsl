package query

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
)

func TestTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&Term{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *Term
	}
	var testCases = []testCase{
		{
			name:  "TestTerm01",
			input: `"dsada 78"`,
			want:  &Term{SRangeTerm: &SRangeTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
		},
		{
			name:  "TestTerm02",
			input: `"dsada 78"^08`,
			want:  &Term{SRangeTerm: &SRangeTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Boost: "^08"}}},
		},
		{
			name:  "TestTerm03",
			input: `"dsada 78"~8`,
			want:  &Term{SRangeTerm: &SRangeTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"}}},
		},
		{
			name:  "TestTerm04",
			input: `"dsada 78"~8^080`,
			want:  &Term{SRangeTerm: &SRangeTerm{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^080"}}},
		},
		{
			name:  "TestTerm05",
			input: `/dsada 78/`,
			want:  &Term{RegexpTerm: &RegexpTerm{Value: `/dsada 78/`}},
		},
		{
			name:  "TestTerm06",
			input: `\/dsada\/\ dasda80980?*`,
			want:  &Term{SRangeTerm: &SRangeTerm{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}}},
		},
		{
			name:  "TestTerm07",
			input: `\/dsada\/\ dasda80980?*\^\^^08`,
			want:  &Term{SRangeTerm: &SRangeTerm{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Boost: `^08`}}},
		},
		{
			name:  "TestTerm08",
			input: `\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &Term{SRangeTerm: &SRangeTerm{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`}}},
		},
		{
			name:  "TestTerm09",
			input: `\/dsada\/\ dasda80980?*\^\^~8^080`,
			want:  &Term{SRangeTerm: &SRangeTerm{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`, Boost: `^080`}}},
		},
		{
			name:  "TestTerm10",
			input: `[1 TO 2]`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "]"},
			},
		},
		{
			name:  "TestTerm11",
			input: `[1 TO 2 }`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "}",
			}},
		},
		{
			name:  `TestTerm12`,
			input: `{ 1 TO 2}`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "}",
			}},
		},
		{
			name:  `TestTerm13`,
			input: `{ 1 TO 2]`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "]",
			}},
		},
		{
			name:  `TestTerm14`,
			input: `[10 TO *]`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"10"}},
				TO:       "TO",
				RValue:   &RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			}},
		},
		{
			name:  `TestTerm15`,
			input: `{* TO 2012-01-01}`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			}},
		},
		{
			name:  `TestTerm16`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &Term{RangeTerm: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{PhraseValue: "\"2012-01-01 09:08:16\""},
				RBRACKET: "}",
			}},
		},
		{
			name:  "TestTerm17",
			input: `<="dsada 78"`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: "<=", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
		},
		{
			name:  "TestTerm18",
			input: `<"dsada 78"^08`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: "<", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Boost: "^08"}}},
		},
		{
			name:  "TestTerm19",
			input: `>="dsada 78"~8`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: ">=", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"}}},
		},
		{
			name:  "TestTerm20",
			input: `>"dsada 78"~8^080`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: ">", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^080"}}},
		},
		{
			name:  "TestTerm21",
			input: `<=dsada\ 78`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: "<=", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}}}},
		},
		{
			name:  "TestTerm22",
			input: `<dsada\ 78^08`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: "<", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}, Boost: "^08"}}},
		},
		{
			name:  "TestTerm23",
			input: `>=dsada\ 78~8`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: ">=", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}, Fuzzy: "~8"}}},
		},
		{
			name:  "TestTerm24",
			input: `>dsada\ 78~8^080`,
			want:  &Term{SRangeTerm: &SRangeTerm{Symbol: ">", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}, Fuzzy: "~8", Boost: "^080"}}},
		},
		{
			name:  "TestTerm25",
			input: `/\d+\d+\.\d+.+/`,
			want:  &Term{RegexpTerm: &RegexpTerm{Value: `/\d+\d+\.\d+.+/`}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestTerm_isRegexp(t *testing.T) {
	var termParser = participle.MustBuild(
		&Term{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  bool
	}

	var testCases = []testCase{
		{
			name:  "TestRegexpTerm01",
			input: `12313\+90`,
			want:  false,
		},
		{
			name:  "TestRegexpTerm02",
			input: `/[1-9]+\.\d+/`,
			want:  true,
		},
		{
			name:  "TestRegexpTerm03",
			input: `"dsad 7089"`,
			want:  false,
		},
		{
			name:  "TestRegexpTerm04",
			input: `[1 TO 454 ]`,
			want:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if out.isRegexp() != tt.want {
				t.Errorf("isRegexp() = %+v, want: %+v", out.isRegexp(), tt.want)
			}
		})
	}

}

func TestTerm_isWildcard(t *testing.T) {

	var termParser = participle.MustBuild(
		&Term{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  bool
	}

	var testCases = []testCase{
		{
			name:  "TestWildcard01",
			input: `12313?`,
			want:  true,
		},
		{
			name:  "TestWildcard02",
			input: `12313\?`,
			want:  false,
		},
		{
			name:  "TestWildcard03",
			input: `12313*`,
			want:  true,
		},
		{
			name:  "TestWildcard04",
			input: `12313\*`,
			want:  false,
		},
		{
			name:  "TestWildcard05",
			input: `/[1-9]+\.\d+/`,
			want:  false,
		},
		{
			name:  "TestWildcard06",
			input: `"dsad?\? 7089*"`,
			want:  true,
		},
		{
			name:  "TestWildcard07",
			input: `"dsadad 789"`,
			want:  false,
		},
		{
			name:  "TestWildcard08",
			input: `[1 TO 2]`,
			want:  false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if out.haveWildcard() != tt.want {
				t.Errorf("haveWildcard() = %+v, want: %+v", out.haveWildcard(), tt.want)
			}
		})
	}
}

func TestTerm_isRange(t *testing.T) {
	var termParser = participle.MustBuild(
		&Term{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  bool
	}

	var testCases = []testCase{
		{
			name:  "TestRangeTerm01",
			input: `12313\+90`,
			want:  false,
		},
		{
			name:  "TestRangeTerm02",
			input: `/[1-9]+\.\d+/`,
			want:  false,
		},
		{
			name:  "TestRangeTerm03",
			input: `"dsad 7089"`,
			want:  false,
		},
		{
			name:  "TestRangeTerm04",
			input: `[1 TO 454 ]`,
			want:  true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if out.isRange() != tt.want {
				t.Errorf("isRange() = %+v, want: %+v", out.isRange(), tt.want)
			}
		})
	}
}

func TestTerm_fuzziness(t *testing.T) {

	var termParser = participle.MustBuild(
		&Term{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  int
	}

	var testCases = []testCase{
		{
			name:  "TestFuzzines01",
			input: `12313\+90`,
			want:  0,
		},
		{
			name:  "TestFuzzines02",
			input: `/[1-9]+\.\d+/`,
			want:  0,
		},
		{
			name:  "TestFuzzines03",
			input: `"dsad 7089"`,
			want:  0,
		},
		{
			name:  "TestFuzzines04",
			input: `[1 TO 454 ]`,
			want:  0,
		},
		{
			name:  "TestFuzzines05",
			input: `12313\+90~3`,
			want:  3,
		},
		{
			name:  "TestFuzzines06",
			input: `"dsad 7089"~3`,
			want:  3,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if out.fuzziness() != tt.want {
				t.Errorf("fuzziness() = %+v, want: %+v", out.fuzziness(), tt.want)
			}
		})
	}

}

func TestTerm_boost(t *testing.T) {

	var termParser = participle.MustBuild(
		&Term{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  float64
	}

	var testCases = []testCase{
		{
			name:  "TestBoost01",
			input: `12313\+90`,
			want:  1.0,
		},
		{
			name:  "TestBoost02",
			input: `/[1-9]+\.\d+/`,
			want:  1.0,
		},
		{
			name:  "TestBoost03",
			input: `"dsad 7089"`,
			want:  1.0,
		},
		{
			name:  "TestBoost04",
			input: `[1 TO 454 ]`,
			want:  1.0,
		},
		{
			name:  "TestBoost05",
			input: `12313\+90^1.2`,
			want:  1.2,
		},
		{
			name:  "TestBoost06",
			input: `12313\+90^0.2`,
			want:  0.2,
		},
		{
			name:  "TestBoost07",
			input: `"dsad 7089"^3.8`,
			want:  3.8,
		},
		{
			name:  "TestBoost07",
			input: `"dsad 7089"^0.8`,
			want:  0.8,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if (out.boost() - tt.want) > 1E-8 {
				t.Errorf("boost() = %+v, want: %+v", out.boost(), tt.want)
			}
		})
	}

}

func TestSimpleTerm(t *testing.T) {
	var simpleTermParser = participle.MustBuild(
		&SimpleTerm{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *SimpleTerm
	}
	var testCases = []testCase{
		{
			name:  "TestSimpleTerm01",
			input: `\/dsada\/\ dasda80980?*`,
			want:  &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}},
		},
		{
			name:  "TestSimpleTerm02",
			input: `\/dsada\/\ dasda80980?*\^\^^08`,
			want:  &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Boost: `^08`},
		},
		{
			name:  "TestSimpleTerm03",
			input: `\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`},
		},
		{
			name:  "TestSimpleTerm04",
			input: `\/dsada\/\ dasda80980?*\^\^~8^080`,
			want:  &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`, Boost: `^080`},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &SimpleTerm{}
			if err := simpleTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("simpleTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestPhraseTerm(t *testing.T) {
	var phraseTermParser = participle.MustBuild(
		&PhraseTerm{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *PhraseTerm
	}
	var testCases = []testCase{
		{
			name:  "PhraseTerm01",
			input: `"dsada 78"`,
			want:  &PhraseTerm{Value: `"dsada 78"`},
		},
		{
			name:  "PhraseTerm02",
			input: `"dsada 78"^08`,
			want:  &PhraseTerm{Value: `"dsada 78"`, Boost: "^08"},
		},
		{
			name:  "PhraseTerm03",
			input: `"dsada 78"~8`,
			want:  &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"},
		},
		{
			name:  "PhraseTerm04",
			input: `"dsada 78"~8^080`,
			want:  &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^080"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &PhraseTerm{}
			if err := phraseTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("phraseTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestRegexpTerm(t *testing.T) {
	var regexpTermParser = participle.MustBuild(
		&RegexpTerm{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *RegexpTerm
	}
	var testCases = []testCase{
		{
			name:  "RegexpTerm01",
			input: `/dsada 78/`,
			want:  &RegexpTerm{Value: `/dsada 78/`},
		},
		{
			name:  "RegexpTerm02",
			input: `/\d+\/\d+\.\d+.+/`,
			want:  &RegexpTerm{Value: `/\d+\/\d+\.\d+.+/`},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &RegexpTerm{}
			if err := regexpTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("regexpTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestRangeTerm(t *testing.T) {
	var rangeTermParser = participle.MustBuild(
		&RangeTerm{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *RangeTerm
	}
	var testCases = []testCase{
		{
			name:  "RangeTerm01",
			input: `[1 TO 2]`,
			want: &RangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "]"},
		},
		{
			name:  "RangeTerm02",
			input: `[1 TO 2 }`,
			want: &RangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `RangeTerm03`,
			input: `{ 1 TO 2}`,
			want: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `RangeTerm04`,
			input: `{ 1 TO 2]`,
			want: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "]",
			},
		},
		{
			name:  `RangeTerm05`,
			input: `[10 TO *]`,
			want: &RangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"10"}},
				TO:       "TO",
				RValue:   &RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			},
		},
		{
			name:  `RangeTerm06`,
			input: `{* TO 2012-01-01}`,
			want: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `RangeTerm07`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &RangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{PhraseValue: "\"2012-01-01 09:08:16\""},
				RBRACKET: "}",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &RangeTerm{}
			if err := rangeTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangeTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestSRangeTerm(t *testing.T) {
	var rangesTermParser = participle.MustBuild(
		&SRangeTerm{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *SRangeTerm
	}
	var testCases = []testCase{
		{
			name:  "SRangeTerm01",
			input: `<="dsada 78"`,
			want:  &SRangeTerm{Symbol: "<=", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
		},
		{
			name:  "SRangeTerm02",
			input: `<"dsada 78"^08`,
			want:  &SRangeTerm{Symbol: "<", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Boost: "^08"}},
		},
		{
			name:  "SRangeTerm03",
			input: `>="dsada 78"~8`,
			want:  &SRangeTerm{Symbol: ">=", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"}},
		},
		{
			name:  "SRangeTerm04",
			input: `>"dsada 78"~8^080`,
			want:  &SRangeTerm{Symbol: ">", PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^080"}},
		},
		{
			name:  "SRangeTerm05",
			input: `<=dsada\ 78`,
			want:  &SRangeTerm{Symbol: "<=", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}}},
		},
		{
			name:  "SRangeTerm06",
			input: `<dsada\ 78^08`,
			want:  &SRangeTerm{Symbol: "<", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}, Boost: "^08"}},
		},
		{
			name:  "SRangeTerm07",
			input: `>=dsada\ 78~8`,
			want:  &SRangeTerm{Symbol: ">=", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}, Fuzzy: "~8"}},
		},
		{
			name:  "SRangeTerm08",
			input: `>dsada\ 78~8^080`,
			want:  &SRangeTerm{Symbol: ">", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}, Fuzzy: "~8", Boost: "^080"}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &SRangeTerm{}
			if err := rangesTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangesTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}
