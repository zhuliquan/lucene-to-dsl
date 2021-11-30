package term

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
	bnd "github.com/zhuliquan/lucene-to-dsl/query/internal/bound"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

func TestSingleTerm(t *testing.T) {
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
			want:     &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `?`}},
			values:   `\/dsada\/\ dasda80980?`,
			wildward: true,
		},
		{
			name:     "TestSimpleTerm02",
			input:    `\/dsada\/\ dasda80980*`,
			want:     &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `80980`, `*`}},
			values:   `\/dsada\/\ dasda80980*`,
			wildward: true,
		},
		{
			name:     "TestSimpleTerm03",
			input:    `\/dsada\/\ dasda8\?0980\*`,
			want:     &SingleTerm{Value: []string{`\/dsada\/\ dasda`, `8`, `\?`, `0980`, `\*`}},
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
				t.Errorf("expect get wildcard: %+v, but get wildcard: %+v", tt.wildward, out.haveWildcard())
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
			want:     &PhraseTerm{Value: []string{`dsada`, ` `, `78`}},
			values:   `dsada 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm02",
			input:    `"*dsada 78"`,
			want:     &PhraseTerm{Value: []string{`*`, `dsada`, ` `, `78`}},
			values:   `*dsada 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm03",
			input:    `"?dsada 78"`,
			want:     &PhraseTerm{Value: []string{`?`, `dsada`, ` `, `78`}},
			values:   `?dsada 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm04",
			input:    `"dsada* 78"`,
			want:     &PhraseTerm{Value: []string{`dsada`, `*`, ` `, `78`}},
			values:   `dsada* 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm05",
			input:    `"dsada? 78"`,
			want:     &PhraseTerm{Value: []string{`dsada`, `?`, ` `, `78`}},
			values:   `dsada? 78`,
			wildward: true,
		},
		{
			name:     "TestPhraseTerm06",
			input:    `"dsada\* 78"`,
			want:     &PhraseTerm{Value: []string{`dsada\*`, ` `, `78`}},
			values:   `dsada\* 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm07",
			input:    `"dsada\? 78"`,
			want:     &PhraseTerm{Value: []string{`dsada\?`, ` `, `78`}},
			values:   `dsada\? 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm09",
			input:    `"\*dsada 78"`,
			want:     &PhraseTerm{Value: []string{`\*dsada`, ` `, `78`}},
			values:   `\*dsada 78`,
			wildward: false,
		},
		{
			name:     "TestPhraseTerm10",
			input:    `"\?dsada 78"`,
			want:     &PhraseTerm{Value: []string{`\?dsada`, ` `, `78`}},
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
				t.Errorf("expect get wildcard: %+v, but get wildcard: %+v", tt.wildward, out.haveWildcard())
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
			want:  &RegexpTerm{Value: []string{`dsada`, ` `, `78`}},
		},
		{
			name:  "RegexpTerm02",
			input: `/\d+\/\d+\.\d+.+/`,
			want:  &RegexpTerm{Value: []string{`\`, `d`, `+`, `\/`, `\`, `d`, `+`, `\`, `.`, `\`, `d`, `+`, `.`, `+`}},
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
		bound *bnd.Bound
	}
	var testCases = []testCase{
		{
			name:  "DRangeTerm01",
			input: `[1 TO 2]`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			},
			bound: &bnd.Bound{
				LeftInclude:  &bnd.RangeValue{SingleValue: []string{"1"}},
				RightInclude: &bnd.RangeValue{SingleValue: []string{"2"}},
			},
		},
		{
			name:  "DRangeTerm02",
			input: `[1 TO 2 }`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			},
			bound: &bnd.Bound{
				LeftInclude:  &bnd.RangeValue{SingleValue: []string{"1"}},
				RightExclude: &bnd.RangeValue{SingleValue: []string{"2"}},
			},
		},
		{
			name:  `DRangeTerm03`,
			input: `{ 1 TO 2}`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "}",
			},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{SingleValue: []string{"1"}},
				RightExclude: &bnd.RangeValue{SingleValue: []string{"2"}},
			},
		},
		{
			name:  `DRangeTerm04`,
			input: `{ 1 TO 2]`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{SingleValue: []string{"1"}},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2"}},
				RBRACKET: "]",
			},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{SingleValue: []string{"1"}},
				RightInclude: &bnd.RangeValue{SingleValue: []string{"2"}},
			},
		},
		{
			name:  `DRangeTerm05`,
			input: `[10 TO *]`,
			want: &DRangeTerm{
				LBRACKET: "[",
				LValue:   &bnd.RangeValue{SingleValue: []string{"10"}},
				RValue:   &bnd.RangeValue{InfinityVal: "*"},
				RBRACKET: "]",
			},
			bound: &bnd.Bound{
				LeftInclude:  &bnd.RangeValue{SingleValue: []string{"10"}},
				RightInclude: &bnd.RangeValue{InfinityVal: "*"},
			},
		},
		{
			name:  `DRangeTerm06`,
			input: `{* TO 2012-01-01}`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{InfinityVal: "*"},
				RValue:   &bnd.RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
				RBRACKET: "}",
			},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{InfinityVal: "*"},
				RightExclude: &bnd.RangeValue{SingleValue: []string{"2012", "-", "01", "-", "01"}},
			},
		},
		{
			name:  `DRangeTerm07`,
			input: `{* TO "2012-01-01 09:08:16"}`,
			want: &DRangeTerm{
				LBRACKET: "{",
				LValue:   &bnd.RangeValue{InfinityVal: "*"},
				RValue:   &bnd.RangeValue{PhraseValue: []string{"2012", "-", "01", "-", "01", " ", "09", ":", "08", ":", "16"}},
				RBRACKET: "}",
			},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{InfinityVal: "*"},
				RightExclude: &bnd.RangeValue{PhraseValue: []string{"2012", "-", "01", "-", "01", " ", "09", ":", "08", ":", "16"}},
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
			} else if !reflect.DeepEqual(tt.bound, out.GetBound()) {
				t.Errorf("expect get bound: %+v, but get bound: %+v", tt.bound, out.GetBound())
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
		bound *bnd.Bound
	}
	var testCases = []testCase{
		{
			name:  "SRangeTerm01",
			input: `<="dsada\455 78"`,
			want:  &SRangeTerm{Symbol: "<=", Value: &bnd.RangeValue{PhraseValue: []string{`dsada`, `\`, `455`, ` `, `78`}}},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{InfinityVal: "*"},
				RightInclude: &bnd.RangeValue{PhraseValue: []string{`dsada`, `\`, `455`, ` `, `78`}},
			},
		},
		{
			name:  "SRangeTerm02",
			input: `<"dsada\\ 78"`,
			want:  &SRangeTerm{Symbol: "<", Value: &bnd.RangeValue{PhraseValue: []string{`dsada\\`, ` `, `78`}}},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{InfinityVal: "*"},
				RightExclude: &bnd.RangeValue{PhraseValue: []string{`dsada\\`, ` `, `78`}},
			},
		},
		{
			name:  "SRangeTerm03",
			input: `>=dsada\ 78`,
			want:  &SRangeTerm{Symbol: ">=", Value: &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}}},
			bound: &bnd.Bound{
				LeftInclude:  &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}},
				RightExclude: &bnd.RangeValue{InfinityVal: "*"},
			},
		},
		{
			name:  "SRangeTerm04",
			input: `>dsada\ 78`,
			want:  &SRangeTerm{Symbol: ">", Value: &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}}},
			bound: &bnd.Bound{
				LeftExclude:  &bnd.RangeValue{SingleValue: []string{`dsada\ `, `78`}},
				RightExclude: &bnd.RangeValue{InfinityVal: "*"},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &SRangeTerm{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("rangesTermParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			} else if !reflect.DeepEqual(tt.bound, out.GetBound()) {
				t.Errorf("expect get bound: %+v, but get bound: %+v", tt.bound, out.GetBound())
			}
		})
	}
}
