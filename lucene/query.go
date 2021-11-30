package lucene

import (
	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	tk "github.com/zhuliquan/lucene-to-dsl/lucene/internal/token"
)

type Query interface {
	String() string
	ToASTNode() (dsl.ASTNode, error)
}

var QueryParser *participle.Parser

func init() {
	var err error
	QueryParser, err = participle.Build(
		&Lucene{},
		participle.Lexer(tk.Lexer),
	)
	if err != nil {
		panic(err)
	}
}
