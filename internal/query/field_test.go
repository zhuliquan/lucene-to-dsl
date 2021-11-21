package query

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
)

func TestFieldTerm(t *testing.T) {
	var fieldTermParser = participle.MustBuild(
		&FieldTerm{},
		participle.Lexer(Lexer),
	)

	type testCase struct {
		name  string
		input string
		want  *FieldTerm
	}
	var testCases = []testCase{
		{
			name:  "TestFieldTerm01",
			input: `x:"dsada 78"`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`}}},
		},
		{
			name:  "TestFieldTerm02",
			input: `x:"dsada 78"^08`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Boost: "^08"}}},
		},
		{
			name:  "TestFieldTerm03",
			input: `x:"dsada 78"~8`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"}}},
		},
		{
			name:  "TestFieldTerm04",
			input: `x:"dsada 78"~8^080`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^080"}}},
		},
		{
			name:  "TestFieldTerm05",
			input: `x-y:/dsada 78/`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x", "-", "y"}}, Colon: ":", Term: &Term{RegexpTerm: &RegexpTerm{Value: `/dsada 78/`}}},
		},
		{
			name:  "TestFieldTerm06",
			input: `x.z-y:\/dsada\/\ dasda80980?*`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x.z", "-", "y"}}, Colon: ":", Term: &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`}}}},
		},
		{
			name:  "TestFieldTerm07",
			input: `x:\/dsada\/\ dasda80980?*\^\^^08`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Boost: `^08`}}},
		},
		{
			name:  "TestFieldTerm08",
			input: `x:\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`}}},
		},
		{
			name:  "TestFieldTerm09",
			input: `x:\/dsada\/\ dasda80980?*\^\^~8^080`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`, Boost: `^080`}}},
		},
		{
			name:  "TestFieldTerm10",
			input: `x:[1 TO 2]`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{RangeTerm: &RangeTerm{LBRACKET: "[", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "]"}}},
		},
		{
			name:  "TestFieldTerm11",
			input: `x:[1 TO 2 }`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{RangeTerm: &RangeTerm{LBRACKET: "[", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "}"}}},
		},
		{
			name:  `TestFieldTerm12`,
			input: `x:{ 1 TO 2}`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{RangeTerm: &RangeTerm{LBRACKET: "{", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "}"}}},
		},
		{
			name:  `TestFieldTerm13`,
			input: `x:{ 1 TO 2]`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{RangeTerm: &RangeTerm{LBRACKET: "{", LValue: []string{"1"}, TO: "TO", RValue: []string{"2"}, RBRACKET: "]"}}},
		},
		{
			name:  `TestFieldTerm14`,
			input: `x:[10 TO *]`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{RangeTerm: &RangeTerm{LBRACKET: "[", LValue: []string{"10"}, TO: "TO", RValue: []string{"*"}, RBRACKET: "]"}}},
		},
		{
			name:  `TestFieldTerm15`,
			input: `x:{* TO 2012-01-01}`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":", Term: &Term{RangeTerm: &RangeTerm{LBRACKET: "{", LValue: []string{"*"}, TO: "TO", RValue: []string{"2012", "-", "01", "-", "01"}, RBRACKET: "}"}}},
		},
		{
			name:  `TestFieldTerm16`,
			input: `x:>89`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":>", Term: &Term{SimpleTerm: &SimpleTerm{Value: []string{"89"}}}},
		},
		{
			name:  `TestFieldTerm17`,
			input: `x:>=89`,
			want:  &FieldTerm{Field: &Field{Value: []string{"x"}}, Colon: ":>=", Term: &Term{SimpleTerm: &SimpleTerm{Value: []string{"89"}}}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &FieldTerm{}
			if err := fieldTermParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s, err: %+v", tt.input, err)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}
}