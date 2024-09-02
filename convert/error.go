package convert

import "fmt"

var (
	ErrEmptyLuceneQuery = fmt.Errorf("empty lucene query")
	ErrEmptyFieldQuery  = fmt.Errorf("empty field query")
	ErrEmptyNotQuery    = fmt.Errorf("empty not query")
	ErrEmptyParenQuery  = fmt.Errorf("empty paren query")
	ErrEmptyAndQuery    = fmt.Errorf("empty and query")
	ErrEmptyOrQuery     = fmt.Errorf("empty or query")
)
