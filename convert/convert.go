package convert

import (
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/lucene"
)

type Convert func(lucene.Query) (dsl.DSLNode, error)

func LuceneToDSL(q *lucene.Lucene) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_AND_QUERY_ERR
	}

	return nil, nil
}

func OrQueryToDSL(q *lucene.OrQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_OR_QUERY_ERR
	}
	return nil, nil
}

func OsQueryToDSL(q *lucene.OSQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_OR_QUERY_ERR
	}
	return nil, nil
}

func AndQueryToDSL(q *lucene.AndQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_AND_QUERY_ERR
	}
	return nil, nil
}

func AnsQueryToDSL(q *lucene.AnSQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_AND_QUERY_ERR
	}
	return nil, nil
}

func ParenQueryToDSL(q *lucene.AnSQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_PAREN_QUERY_ERR
	}
	return nil, nil
}

func FieldQueryToDSL(q *lucene.FieldQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, EMPTY_FIELD_QUERY_ERR
	}
	return nil, nil
}
