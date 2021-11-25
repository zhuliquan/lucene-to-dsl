package term

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/internal/token"
)

func TestSimpleTerm(t *testing.T) {
	var simpleTermParser = participle.MustBuild(
		&SimpleTerm{},
		participle.Lexer(token.Lexer),
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
		participle.Lexer(token.Lexer),
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
		participle.Lexer(token.Lexer),
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

func TestDRangeTerm(t *testing.T) {
	var rangeTermParser = participle.MustBuild(
		&DRangeTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *DRangeTerm
	}
	var testCases = []testCase{
		{
			name:  "DRangeTerm01",
			input: `[1 TO 2]`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "]"},
		},
		{
			name:  "DRangeTerm02",
			input: `[1 TO 2 }`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `DRangeTerm03`,
			input: `{ 1 TO 2}`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `DRangeTerm04`,
			input: `{ 1 TO 2]`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SimpleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2"}},
				RBRACKET: "]",
			},
		},
		{
			name:  `DRangeTerm05`,
			input: `[10 TO *]`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SimpleValue: []string{"10"}},
				TO:       "TO",
				RValue:   &RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			},
		},
		{
			name:  `DRangeTerm06`,
			input: `{* TO 2012-01-01}`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{InfinityVal: "*"},
				TO:       "TO",
				RValue:   &RangeValue{SimpleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `DRangeTerm07`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &DRangeTerm{
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
			var out = &DRangeTerm{}
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
		participle.Lexer(token.Lexer),
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
			name:  "SRangeTerm05",
			input: `<=dsada\ 78`,
			want:  &SRangeTerm{Symbol: "<=", SimpleTerm: &SimpleTerm{Value: []string{`dsada\ 78`}}},
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
