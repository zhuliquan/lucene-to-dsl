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
			input: `"dsada 78"^8`,
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Boost: "^8"}},
		},
		{
			name:  "TestTerm03",
			input: `"dsada 78"~8`,
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8"}},
		},
		{
			name:  "TestTerm04",
			input: `"dsada 78"~8^8`,
			want:  &Term{PhraseTerm: &PhraseTerm{Value: `"dsada 78"`, Fuzzy: "~8", Boost: "^8"}},
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
			input: `\/dsada\/\ dasda80980?*\^\^^8`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Boost: `^8`}},
		},
		{
			name:  "TestTerm07",
			input: `\/dsada\/\ dasda80980?*\^\^~8`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`}},
		},
		{
			name:  "TestTerm07",
			input: `\/dsada\/\ dasda80980?*\^\^~8^8`,
			want:  &Term{SimpleTerm: &SimpleTerm{Value: []string{`\/dsada\/\ dasda80980`, `?`, `*`, `\^\^`}, Fuzzy: `~8`, Boost: `^8`}},
		},
		// {},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var out = &Term{}
			if err := termParser.ParseString(tt.input, out); err != nil {
				t.Errorf("failed to parse input: %s", tt.input)
			} else if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("termParser.ParseString( %s ) = %+v, want: %+v", tt.input, out, tt.want)
			}
		})
	}

}
