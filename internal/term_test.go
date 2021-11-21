package internal

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
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}},
		},
		{
			name:  "TestTerm02",
			input: `"dsada 78"^08`,
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Boost: "^08"}},
		},
		{
			name:  "TestTerm03",
			input: `"dsada 78"~8`,
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"}},
		},
		{
			name:  "TestTerm04",
			input: `"dsada 78"~8^080`,
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^080"}},
		},
		{
			name:  "TestTerm05",
			input: `/dsada 78/`,
			want:  &Term{RegexpTerm: &RegexpTerm{Value: `/dsada 78/`}},
		},
		{
			name:  "TestTerm06",
			input: `\/dsada\/\ dasda80980?*`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}},
		},
		{
			name:  "TestTerm07",
			input: `\/dsada\/\ dasda80980?*\^\^^08`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Boost: `^08`}},
		},
		{
			name:  "TestTerm08",
			input: `\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`}},
		},
		{
			name:  "TestTerm09",
			input: `\/dsada\/\ dasda80980?*\^\^~8^080`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`, Boost: `^080`}},
		},
		{
			name:  "TestTerm10",
			input: `[1 TO 2]`,
			want:  &Term{RangeTerm: &RangeTerm{LBRACKET: "[", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "]"}},
		},
		{
			name:  "TestTerm11",
			input: `[1 TO 2 }`,
			want:  &Term{RangeTerm: &RangeTerm{LBRACKET: "[", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "}"}},
		},
		{
			name:  `TestTerm12`,
			input: `{ 1 TO 2}`,
			want:  &Term{RangeTerm: &RangeTerm{LBRACKET: "{", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "}"}},
		},
		{
			name:  `TestTerm13`,
			input: `{ 1 TO 2]`,
			want:  &Term{RangeTerm: &RangeTerm{LBRACKET: "{", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "]"}},
		},
		{
			name:  `TestTerm14`,
			input: `[10 TO *]`,
			want:  &Term{RangeTerm: &RangeTerm{LBRACKET: "[", LValue: []string{"10"}, TO: "TO", RValue: []string{"*"}, RBRACKET: "]"}},
		},
		{
			name:  `TestTerm15`,
			input: `{* TO 2012-01-01}`,
			want:  &Term{RangeTerm: &RangeTerm{LBRACKET: "{", LValue: []string{"*"}, TO: "TO", RValue: []string{"2012", "-", "01", "-", "01"}, RBRACKET: "}"}},
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
