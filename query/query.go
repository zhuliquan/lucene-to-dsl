package query

import (
	"github.com/alecthomas/participle"
	"github.com/zhuliquan/lucene-to-dsl/query/internal/lucene"
	tk "github.com/zhuliquan/lucene-to-dsl/query/internal/token"
)

var QueryParser *participle.Parser

func init() {
	var err error
	QueryParser, err = participle.Build(
		&lucene.Lucene{},
		participle.Lexer(tk.Lexer),
	)
	if err != nil {
		panic(err)
	}
}
