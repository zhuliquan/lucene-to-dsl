package token

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/stateful"
)

var rules = []stateful.Rule{
	{
		Name:    "EOL",
		Pattern: `\n`,
	},
	{
		Name:    "WHITESPACE",
		Pattern: `[\t\r\f ]+`,
	},
	{
		Name:    "IDENT",
		Pattern: `([^-\!\s:\|\&"\?\*\\\^~\(\)\{\}\[\]\+\/><=]|\.|(\\(\s|:|\&|\||\?|\*|\\|\^|~|\(|\)|\!|\[|\]|\{|\}|\+|-|\/|>|<|=)))+`,
	},
	{
		Name:    "QUOTE",
		Pattern: `"`,
	},
	{
		Name:    "SLASH",
		Pattern: `\/`,
	},
	{
		Name:    "REVERSE",
		Pattern: `\\`,
	},
	{
		Name:    "COLON",
		Pattern: `:`,
	},
	{
		Name:    "COMPARE",
		Pattern: `[<>]=?`,
	},
	{
		Name:    "PLUS",
		Pattern: `\+`,
	},
	{
		Name:    "MINUS",
		Pattern: `-`,
	},
	{
		Name:    "FUZZY",
		Pattern: `~\d*`,
	},
	{
		Name:    "BOOST",
		Pattern: `\^(\d+(\.\d+)?)`,
	},
	{
		Name:    "WILDCARD",
		Pattern: `[\?\*]`,
	},
	{
		Name:    "LPAREN",
		Pattern: `\(`,
	},
	{
		Name:    "RPAREN",
		Pattern: `\)`,
	},
	{
		Name:    "LBRACK",
		Pattern: `\[`,
	},
	{
		Name:    "RBRACK",
		Pattern: `\]`,
	},
	{
		Name:    "LBRACE",
		Pattern: `\{`,
	},
	{
		Name:    "RBRACE",
		Pattern: `\}`,
	},
	{
		Name:    "AND",
		Pattern: `\&`,
	},
	{
		Name:    "SOR",
		Pattern: `\|`,
	},
	{
		Name:    "NOT",
		Pattern: `!`,
	},
}

var Lexer *stateful.Definition
var Scanner *participle.Parser

func Scan(exp string) []*Token {
	var tokens = []*Token{}
	var ch = make(chan *Token, 100)
	if err := Scanner.ParseString(exp, ch); err != nil {
		return nil
	} else {
		for c := range ch {
			tokens = append(tokens, c)
		}
		return tokens
	}
}

func init() {
	var err error
	Lexer, err = stateful.NewSimple(rules)
	if err != nil {
		panic(err.Error())
	}
	Scanner, err = participle.Build(
		&Token{},
		participle.Lexer(Lexer),
	)
	if err != nil {
		panic(err.Error())
	}

}

type Token struct {
	EOL        string `parser:"  @EOL" json:"eol"`
	WHITESPACE string `parser:"| @WHITESPACE" json:"whitespace"`
	IDENT      string `parser:"| @IDENT" json:"ident"`
	QUOTE      string `parser:"| @QUOTE" json:"quote"`
	SLASH      string `parser:"| @SLASH" json:"slash"`
	REVERSE    string `parser:"| @REVERSE" json:"reverse"`
	COLON      string `parser:"| @COLON" json:"colon"`
	COMPARE    string `parser:"| @COMPARE" json:"compare"`
	PLUS       string `parser:"| @PLUS" json:"plus"`
	MINUS      string `parser:"| @MINUS" json:"minus"`
	FUZZY      string `parser:"| @FUZZY" json:"fuzzy"`
	BOOST      string `parser:"| @BOOST" json:"boost"`
	WILDCARD   string `parser:"| @WILDCARD" json:"wildcard"`
	LPAREN     string `parser:"| @LPAREN" json:"lparen"`
	RPAREN     string `parser:"| @RPAREN" json:"rparen"`
	LBRACK     string `parser:"| @LBRACK" json:"lbrack"`
	RBRACK     string `parser:"| @RBRACK" json:"rbrack"`
	LBRACE     string `parser:"| @LBRACE" json:"lbrace"`
	RBRACE     string `parser:"| @RBRACE" json:"rbrace"`
	AND        string `parser:"| @AND" json:"and"`
	SOR        string `parser:"| @SOR" json:"sor"`
	NOT        string `parser:"| @NOT" json:"not"`
}

func (t *Token) String() string {
	if t == nil {
		return ""
	} else if t.EOL != "" {
		return t.EOL
	} else if t.WHITESPACE != "" {
		return t.WHITESPACE
	} else if t.IDENT != "" {
		return t.IDENT
	} else if t.QUOTE != "" {
		return t.QUOTE
	} else if t.SLASH != "" {
		return t.SLASH
	} else if t.REVERSE != "" {
		return t.REVERSE
	} else if t.COLON != "" {
		return t.COLON
	} else if t.COMPARE != "" {
		return t.COMPARE
	} else if t.PLUS != "" {
		return t.PLUS
	} else if t.MINUS != "" {
		return t.MINUS
	} else if t.FUZZY != "" {
		return t.FUZZY
	} else if t.BOOST != "" {
		return t.BOOST
	} else if t.WILDCARD != "" {
		return t.WILDCARD
	} else if t.LPAREN != "" {
		return t.LPAREN
	} else if t.RPAREN != "" {
		return t.RPAREN
	} else if t.LBRACK != "" {
		return t.LBRACK
	} else if t.RBRACK != "" {
		return t.RBRACK
	} else if t.LBRACE != "" {
		return t.LBRACE
	} else if t.RBRACE != "" {
		return t.RBRACE
	} else if t.AND != "" {
		return t.AND
	} else if t.SOR != "" {
		return t.SOR
	} else if t.NOT != "" {
		return t.NOT
	} else {
		return ""
	}
}

func (t *Token) GetTokenType() TokenType {
	if t == nil {
		return UNKNOWN_TOKEN_TYPE
	} else if t.EOL != "" {
		return EOL_TOKEN_TYPE
	} else if t.WHITESPACE != "" {
		return WHITESPACE_TOKEN_TYPE
	} else if t.IDENT != "" {
		return IDENT_TOKEN_TYPE
	} else if t.QUOTE != "" {
		return QUOTE_TOKEN_TYPE
	} else if t.SLASH != "" {
		return SLASH_TOKEN_TYPE
	} else if t.REVERSE != "" {
		return REVERSE_TOKEN_TYPE
	} else if t.COLON != "" {
		return COLON_TOKEN_TYPE
	} else if t.COMPARE != "" {
		return COMPARE_TOKEN_TYPE
	} else if t.PLUS != "" {
		return PLUS_TOKEN_TYPE
	} else if t.MINUS != "" {
		return MINUS_TOKEN_TYPE
	} else if t.FUZZY != "" {
		return FUZZY_TOKEN_TYPE
	} else if t.BOOST != "" {
		return BOOST_TOKEN_TYPE
	} else if t.WILDCARD != "" {
		return WILDCARD_TOKEN_TYPE
	} else if t.LPAREN != "" {
		return LPAREN_TOKEN_TYPE
	} else if t.RPAREN != "" {
		return RPAREN_TOKEN_TYPE
	} else if t.LBRACK != "" {
		return LBRACK_TOKEN_TYPE
	} else if t.RBRACK != "" {
		return RBRACK_TOKEN_TYPE
	} else if t.LBRACE != "" {
		return LBRACE_TOKEN_TYPE
	} else if t.RBRACE != "" {
		return RBRACE_TOKEN_TYPE
	} else if t.AND != "" {
		return AND_TOKEN_TYPE
	} else if t.SOR != "" {
		return SOR_TOKEN_TYPE
	} else if t.NOT != "" {
		return NOT_TOKEN_TYPE
	} else {
		return UNKNOWN_TOKEN_TYPE
	}
}
