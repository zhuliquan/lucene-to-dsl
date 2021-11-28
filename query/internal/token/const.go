package token

type TokenType uint32

const (
	UNKNOWN_TOKEN_TYPE    TokenType = 0
	EOL_TOKEN_TYPE        TokenType = 1
	WHITESPACE_TOKEN_TYPE TokenType = 2
	IDENT_TOKEN_TYPE      TokenType = 3
	STRING_TOKEN_TYPE     TokenType = 4
	REGEXP_TOKEN_TYPE     TokenType = 5
	COLON_TOKEN_TYPE      TokenType = 6
	PLUS_TOKEN_TYPE       TokenType = 7
	COMPARE_TOKEN_TYPE    TokenType = 8
	MINUS_TOKEN_TYPE      TokenType = 9
	FUZZY_TOKEN_TYPE      TokenType = 10
	BOOST_TOKEN_TYPE      TokenType = 11
	WILDCARD_TOKEN_TYPE   TokenType = 12
	LPAREN_TOKEN_TYPE     TokenType = 13
	RPAREN_TOKEN_TYPE     TokenType = 14
	LBRACK_TOKEN_TYPE     TokenType = 15
	RBRACK_TOKEN_TYPE     TokenType = 16
	LBRACE_TOKEN_TYPE     TokenType = 17
	RBRACE_TOKEN_TYPE     TokenType = 18
	AND_TOKEN_TYPE        TokenType = 19
	SOR_TOKEN_TYPE        TokenType = 20
	NOT_TOKEN_TYPE        TokenType = 21
)
