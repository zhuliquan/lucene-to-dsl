package token

import (
	"reflect"
	"testing"

	"github.com/alecthomas/participle"
)

func TestLexer(t *testing.T) {
	var scan = func(scanner *participle.Parser, exp string) []*Token {
		var tokens = []*Token{}
		var ch = make(chan *Token, 100)
		scanner.ParseString(exp, ch)
		for c := range ch {
			tokens = append(tokens, c)
		}
		return tokens
	}
	var scanner *participle.Parser
	var err error
	scanner, err = participle.Build(
		&Token{},
		participle.Lexer(Lexer),
	)
	if err != nil {
		panic(err)
	}

	type testCase struct {
		name  string
		input string
		want  []*Token
	}

	var testCases = []testCase{
		{
			name:  "TestScan01",
			input: `\ \ \:7:>8908 8+9 x:>=90`,
			want: []*Token{
				{IDENT: `\ \ \:7`},
				{COLON: ":"},
				{COMPARE: ">"},
				{IDENT: "8908"},
				{WHITESPACE: " "},
				{IDENT: "8"},
				{PLUS: "+"},
				{IDENT: "9"},
				{WHITESPACE: " "},
				{IDENT: "x"},
				{COLON: ":"},
				{COMPARE: ">="},
				{IDENT: "90"},
			},
		},
		{
			name:  "TestScan02",
			input: `now-8d x:/[\d\s]+/ y:"dasda 8\ : +"`,
			want: []*Token{
				{IDENT: "now"},
				{MINUS: "-"},
				{IDENT: "8d"},
				{WHITESPACE: " "},
				{IDENT: "x"},
				{COLON: ":"},
				{REGEXP: `/[\d\s]+/`},
				{WHITESPACE: " "},
				{IDENT: "y"},
				{COLON: ":"},
				{STRING: `"dasda 8\ : +"`},
			},
		},
		{
			name:  "TestScan03",
			input: `\!\:.\ \\:<=<(you OR !& \!\&*\** [{ you\[\]+ you?}])^090~9~ouo |`,
			want: []*Token{
				{IDENT: `\!\:.\ \\`},
				{COLON: ":"},
				{COMPARE: "<="},
				{COMPARE: "<"},
				{LPAREN: "("},
				{IDENT: "you"},
				{WHITESPACE: " "},
				{IDENT: "OR"},
				{WHITESPACE: " "},
				{NOT: "!"},
				{AND: "&"},
				{WHITESPACE: " "},
				{IDENT: `\!\&`},
				{WILDCARD: "*"},
				{IDENT: `\*`},
				{WILDCARD: "*"},
				{WHITESPACE: " "},
				{LBRACK: "["},
				{LBRACE: "{"},
				{WHITESPACE: " "},
				{IDENT: `you\[\]`},
				{PLUS: `+`},
				{WHITESPACE: " "},
				{IDENT: "you"},
				{WILDCARD: "?"},
				{RBRACE: "}"},
				{RBRACK: "]"},
				{RPAREN: ")"},
				{BOOST: `^090`},
				{FUZZY: `~9`},
				{FUZZY: `~`},
				{IDENT: "ouo"},
				{WHITESPACE: " "},
				{SOR: "|"},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if out := scan(scanner, tt.input); !reflect.DeepEqual(out, tt.want) {
				t.Errorf("Scan ( %+v ) = %+v, but want: %+v", tt.input, out, tt.want)
			}
		})

	}
}
