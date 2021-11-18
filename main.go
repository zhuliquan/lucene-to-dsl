package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/stateful"
)

type FieldToken struct {
	Field string `parser:"@WORD"  json:"field"`
	Colon string `parser:"@COLON" json:"colon"`
	TERM  *Term  `parser:"@@" json:"term"`
}

type Field struct {
}

type Term struct {
	PhraseTerm *PhraseTerm `parser:"  @@" json:"phrase_term"`
	RegexpTerm *RegexpTerm `parser:"| @@" json:"regexp_term"`
	SimpleTerm *SimpleTerm `parser:"| @@" json:"simple_term"`
}

func (t *Term) String() string {
	if t == nil {
		return ""
	} else if t.PhraseTerm != nil {
		return t.PhraseTerm.String()
	} else if t.RegexpTerm != nil {
		return t.RegexpTerm.String()
	} else if t.SimpleTerm != nil {
		return t.SimpleTerm.String()
	} else {
		return ""
	}
}

type RegexpTerm struct {
	Value string `parser:"@REGEXP" json:"value"`
}

func (t *RegexpTerm) String() string {
	if t == nil {
		return ""
	} else if t.Value != "" {
		return t.Value[1 : len(t.Value)-1]
	} else {
		return ""
	}
}

type PhraseTerm struct {
	Value string `parser:"@STRING" json:"value"`
	Fuzzy string `parser:"@FUZZY?" json:"fuzzy"`
	Boost string `parser:"@BOOST?" json:"boost"`
}

func (t *PhraseTerm) String() string {
	if t == nil {
		return ""
	} else if t.Value != "" {
		return t.Value[1 : len(t.Value)-1]
	} else {
		return ""
	}
}

func (t *PhraseTerm) isWildCard() bool {
	for i := 1; i < len(t.Value)-1; i++ {
		if i > 1 && (t.Value[i] == '?' || t.Value[i] == '*' && t.Value[i-1] != '\\') {
			return true
		}
		if i == 1 && (t.Value[i] == '?' || t.Value[i] == '*') {
			return true
		}
	}
	return false
}

type SimpleTerm struct {
	Value []string `parser:"@(IDENT|WILDCARD|MATH_OPERATOR)+" json:"value"`
	Fuzzy string   `parser:"@FUZZY?" json:"fuzzy"`
	Boost string   `parser:"@BOOST?" json:"boost"`
}

func (t *SimpleTerm) String() string {
	return strings.Join(t.Value, "")
}

func (t *SimpleTerm) isWildCard() bool {
	for i := 0; i < len(t.Value); i++ {
		if t.Value[i] == "?" || t.Value[i] == "*" {
			return true
		}
	}
	return false
}

// func (t *SimpleTerm) getFuzzy() (float64, bool) {
// 	if t.Fuzzy == "" {
// 		return 0.0, false
// 	} else {
// 		return strconv.ParseFloat()
// 	}
// }

type Token struct {
	STRING        string `  @STRING`
	REGEXP        string `| @REGEXP`
	COLON         string `| @COLON`
	IDENT         string `| @IDENT`
	EOL           string `| @EOL`
	WHITESPACE    string `| @WHITESPACE`
	FUZZY         string `| @FUZZY`
	BOOST         string `| @BOOST`
	WILDCARD      string `| @WILDCARD`
	BRACKET       string `| @BRACKET`
	BOOL_OPERATOR string `| @BOOL_OPERATOR`
	MATH_OPERATOR string `| @MATH_OPERATOR`
}

func (t *Token) String() string {
	if t == nil {
		return ""
	} else if t.STRING != "" {
		return t.STRING
	} else if t.REGEXP != "" {
		return t.REGEXP
	} else if t.COLON != "" {
		return t.COLON
	} else if t.IDENT != "" {
		return t.IDENT
	} else if t.EOL != "" {
		return t.EOL
	} else if t.WHITESPACE != "" {
		return t.WHITESPACE
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

func main() {
	var _lexer, err1 = stateful.NewSimple(
		[]stateful.Rule{
			{
				Name:    "STRING",
				Pattern: `"(\\"|[^"])*"`,
			},
			{
				Name:    "REGEXP",
				Pattern: `\/([^"\/]|\\"|\\\/)+\/`,
			},
			{
				Name:    "IDENT",
				Pattern: `([^\!\s:\|\&"\?\*\\\^~\(\)\{\}\[\]\+-\/]|(\\(\s|:|\&|\||\?|\*|\\|\^|~|\(|\)|\!|\[|\]|\{|\}|\+|-|\/)))+`,
			},
			{
				Name:    "MATH_OPERATOR",
				Pattern: `[-\+\*\/]`,
			},
			{
				Name:    "COLON",
				Pattern: `:([<|>]=?)?`,
			},
			{
				Name:    "WHITESPACE",
				Pattern: `[\t\r\f ]+`,
			},
			{
				Name:    "EOL",
				Pattern: `\n`,
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
				Name:    "BOOL_OPERATOR",
				Pattern: `[!\&\|]`,
			},
			{
				Name:    "BRACKET",
				Pattern: `[\(\)\[\]\{\}]`,
			},
		},
	)

	fmt.Println(_lexer, err1)
	var parser, err2 = participle.Build(
		&Token{},
		participle.Lexer(_lexer),
	)
	fmt.Println(parser, err2)

	tokens := make(chan *Token, 128)
	s := `(     x:>"10\ 90"   && Y:/[0-9]+/   ||  Z:you^91+2  X:y?Sou*\ \ ^9 ! ypu:90)`
	var err = parser.ParseString(s, tokens)
	fmt.Println(err)
	for token := range tokens {
		fmt.Printf("%s\n", token)
	}

	var termParser, err3 = participle.Build(
		&Term{},
		participle.Lexer(_lexer),
	)
	fmt.Println(termParser, err3)
	var term = &Term{}

	err = termParser.ParseString(`?uiouio*\\\ 8980~0.9^`, term)
	fmt.Println(term.String(), err)
	err = termParser.ParseString(`/?uiouio*\\\ 8980~0.9^/`, term)
	fmt.Println(term.String(), err)
	err = termParser.ParseString(`"?uiouio*\\\ 8980"~0.9^`, term)
	fmt.Println(term.String(), err)
}
