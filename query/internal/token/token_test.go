package token

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	var err error
	if err != nil {
		panic(err)
	}

	type testCase struct {
		name  string
		input string
		want  []*Token
		typeS []TokenType
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
			typeS: []TokenType{
				IDENT_TOKEN_TYPE,
				COLON_TOKEN_TYPE,
				COMPARE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				PLUS_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				COLON_TOKEN_TYPE,
				COMPARE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
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
			typeS: []TokenType{
				IDENT_TOKEN_TYPE,
				MINUS_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				COLON_TOKEN_TYPE,
				REGEXP_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				COLON_TOKEN_TYPE,
				STRING_TOKEN_TYPE,
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
			typeS: []TokenType{
				IDENT_TOKEN_TYPE,
				COLON_TOKEN_TYPE,
				COMPARE_TOKEN_TYPE,
				COMPARE_TOKEN_TYPE,
				LPAREN_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				NOT_TOKEN_TYPE,
				AND_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WILDCARD_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WILDCARD_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				LBRACK_TOKEN_TYPE,
				LBRACE_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				PLUS_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WILDCARD_TOKEN_TYPE,
				RBRACE_TOKEN_TYPE,
				RBRACK_TOKEN_TYPE,
				RPAREN_TOKEN_TYPE,
				BOOST_TOKEN_TYPE,
				FUZZY_TOKEN_TYPE,
				FUZZY_TOKEN_TYPE,
				IDENT_TOKEN_TYPE,
				WHITESPACE_TOKEN_TYPE,
				SOR_TOKEN_TYPE,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if out := Scan(tt.input); !reflect.DeepEqual(out, tt.want) {
				t.Errorf("Scan ( %+v ) = %+v, but want: %+v", tt.input, out, tt.want)
			} else {
				for i := 0; i < len(out); i++ {
					if out[i].GetTokenType() != tt.typeS[i] {
						t.Errorf("expect get type: %+v, but get type: %+v", tt.typeS[i], out[i].GetTokenType())
					}
				}
			}
		})

	}
}
