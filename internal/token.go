package internal

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
		Pattern: `([^\!\s:\|\&"\?\*\\\^~\(\)\{\}\[\]\+-\/]|\.|(\\(\s|:|\&|\||\?|\*|\\|\^|~|\(|\)|\!|\[|\]|\{|\}|\+|-|\/)))+`,
	},
	{
		Name:    "STRING",
		Pattern: `"(\\"|[^"])*"`,
	},
	{
		Name:    "REGEXP",
		Pattern: `\/([^"\/]|\\"|\\\/)+\/`,
	},
	{
		Name:    "COLON",
		Pattern: `:([<|>]=?)?`,
	},
	{
		Name:    "FUZZY",
		Pattern: `~((0\.)?\d+)?`,
	},
	{
		Name:    "BOOST",
		Pattern: `\^((0\.)?\d+)?`,
	},
	{
		Name:    "WILDCARD",
		Pattern: `[\?\*]`,
	},
	{
		Name:    "BRACKET",
		Pattern: `[\(\)\[\]\{\}]`,
	},
	{
		Name:    "BOOL_OPERATOR",
		Pattern: `[!\&\|]`,
	},
	{
		Name:    "MATH_OPERATOR",
		Pattern: `[-\+\*\/]`,
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
	EOL           string `parser:"  @EOL" json:"eol"`
	WHITESPACE    string `parser:"| @WHITESPACE" json:"whitespace"`
	IDENT         string `parser:"| @IDENT" json:"ident"`
	STRING        string `parser:"| @STRING" json:"string"`
	REGEXP        string `parser:"| @REGEXP" json:"regexp"`
	COLON         string `parser:"| @COLON" json:"colon"`
	FUZZY         string `parser:"| @FUZZY" json:"fuzzy"`
	BOOST         string `parser:"| @BOOST" json:"boost"`
	WILDCARD      string `parser:"| @WILDCARD" json:"wildcard"`
	BRACKET       string `parser:"| @BRACKET" json:"bracket"`
	BOOL_OPERATOR string `parser:"| @BOOL_OPERATOR" json:"bool_operator"`
	MATH_OPERATOR string `parser:"| @MATH_OPERATOR" json:"math_operator"`
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
	} else if t.FUZZY != "" {
		return t.FUZZY
	} else if t.BOOST != "" {
		return t.BOOST
	} else if t.WILDCARD != "" {
		return t.WILDCARD
	} else if t.BRACKET != "" {
		return t.BRACKET
	} else if t.BOOL_OPERATOR != "" {
		return t.BOOL_OPERATOR
	} else if t.MATH_OPERATOR != "" {
		return t.MATH_OPERATOR
	} else {
		return ""
	}
}
