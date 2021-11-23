package query

import (
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
		Name:    "STRING",
		Pattern: `"(\\"|[^"])*"`,
	},
	{
		Name:    "REGEXP",
		Pattern: `\/([^\/\\]|\\\\|\\\/)+\/`,
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
		Pattern: `~(0*[1-9][0-9]*)?`,
	},
	{
		Name:    "BOOST",
		Pattern: `\^(\d+\.?\d+)`,
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
		Name:    "OR",
		Pattern: `\|`,
	},
	{
		Name:    "NOT",
		Pattern: `!`,
	},
}

var Lexer *stateful.Definition

func init() {
	var err error
	Lexer, err = stateful.NewSimple(rules)
	if err != nil {
		panic(err.Error())
	}

}

type Token struct {
	EOL        string `parser:"  @EOL" json:"eol"`
	WHITESPACE string `parser:"| @WHITESPACE" json:"whitespace"`
	IDENT      string `parser:"| @IDENT" json:"ident"`
	STRING     string `parser:"| @STRING" json:"string"`
	REGEXP     string `parser:"| @REGEXP" json:"regexp"`
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
	SOR        string `parser:"| @OR" json:"sor"`
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
	} else if t.STRING != "" {
		return t.STRING
	} else if t.REGEXP != "" {
		return t.REGEXP
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
