package lucene

import "fmt"

var (
	EMPTY_FIELD_QUERY_ERR = fmt.Errorf("empty field query")
	EMPTY_NOT_QUERY_ERR   = fmt.Errorf("empty not query")
	EMPTY_PAREN_QUERY_ERR = fmt.Errorf("empty paren query")
	EMPTY_AND_QUERY_ERR   = fmt.Errorf("empty and query")
	EMPTY_OR_QUERY_ERR    = fmt.Errorf("empty or query")
)
