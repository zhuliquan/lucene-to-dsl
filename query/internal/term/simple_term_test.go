package term

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

func TestSimpleTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&SingleTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name     string
		input    string
		want     *SingleTerm
		values   string
		wildward bool
	}
	var testCases = []testCase{
		{
			name:     "TestSimpleTerm01",
			input:    `\/dsada\/\ dasda80980?`,
			want:     &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`}},
			values:   `\/dsada\/\ dasda80980?`,
			wildward: true,
		},
		{
			name:     "TestSimpleTerm02",
			input:    `\/dsada\/\ dasda80980*`,
			want:     &SingleTerm{Value: []string{`\/dsada\/\ dasda80980`, `*`}},
			values:   `\/dsada\/\ dasda80980*`,
			wildward: true,
		},
		{
			name:     "TestSimpleTerm03",
			input:    `\/dsada\/\ dasda8\?0980\*`,
			want:     &SingleTerm{Value: []string{`\/dsada\/\ dasda8\?0980\*`}},
			values:   `\/dsada\/\ dasda8\?0980\*`,
			wildward: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &SingleTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if tt.values != out.ValueS() {
				t.Errorf("expect get values: %s, but get values: %+v", tt.values, out.ValueS())
			} else if tt.wildward != out.haveWildcard() {
				t.Errorf("expect get wildcard: %s, but get wildcard: %+v", tt.wildward, out.haveWildcard())
			}
		})
	}
}

func TestPhraseTerm(t *testing.T) {
	var termParser = participle.MustBuild(
		&PhraseTerm{},
		participle.Lexer(token.Lexer),
	)

	type testCase struct {
		name     string
		input    string
		want     *PhraseTerm
		values   string
		wildward bool
	}
	var testCases = []testCase{
		{
			name:     "TestPhraseTerm01",
			input:    `"dsada 78"`,
			want:     &PhraseTerm{Value: `"dsada 78"`},
			values:   `dsada 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm02",
			input:    `"*dsada 78"`,
			want:     &PhraseTerm{Value: `"*dsada 78"`},
			values:   `*dsada 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm03",
			input:    `"?dsada 78"`,
			want:     &PhraseTerm{Value: `"?dsada 78"`},
			values:   `?dsada 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm04",
			input:    `"dsada* 78"`,
			want:     &PhraseTerm{Value: `"dsada* 78"`},
			values:   `dsada* 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm05",
			input:    `"dsada? 78"`,
			want:     &PhraseTerm{Value: `"dsada? 78"`},
			values:   `dsada? 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm06",
			input:    `"dsada\* 78"`,
			want:     &PhraseTerm{Value: `"dsada\* 78"`},
			values:   `dsada\* 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm07",
			input:    `"dsada\? 78"`,
			want:     &PhraseTerm{Value: `"dsada\? 78"`},
			values:   `dsada\? 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm09",
			input:    `"\*dsada 78"`,
			want:     &PhraseTerm{Value: `"\*dsada 78"`},
			values:   `\*dsada 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm10",
			input:    `"\?dsada 78"`,
			want:     &PhraseTerm{Value: `"\?dsada 78"`},
			values:   `\?dsada 78`,
			wildward: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &PhraseTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("phraseTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if tt.values != out.ValueS() {
				t.Errorf("expect get values: %s, but get values: %+v", tt.values, out.ValueS())
			} else if tt.wildward != out.haveWildcard() {
				t.Errorf("expect get wildcard: %s, but get wildcard: %+v", tt.wildward, out.haveWildcard())
			}
		})
	}
}

func TestRegexpTerm(t *testing.T) {
	var termParser = participle.MustBuild(
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
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("regexpTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestDRangeTerm(t *testing.T) {
	var termParser = participle.MustBuild(
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
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]"},
		},
		{
			name:  "DRangeTerm02",
			input: `[1 TO 2 }`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `DRangeTerm03`,
			input: `{ 1 TO 2}`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			},
		},
		{
			name:  `DRangeTerm04`,
			input: `{ 1 TO 2]`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &RangeValue{SingleValue: []string{"1"}},
				TO:       "TO",
				RValue:   &RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			},
		},
		{
			name:  `DRangeTerm05`,
			input: `[10 TO *]`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &RangeValue{SingleValue: []string{"10"}},
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
				RValue:   &RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
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
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangeTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}

func TestSRangeTerm(t *testing.T) {
	var termParser = participle.MustBuild(
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
			want:  &SRangeTerm{Symbol: "<=", SingleTerm: &SingleTerm{Value: []string{`dsada\ 78`}}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &SRangeTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangesTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}
