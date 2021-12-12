package convert

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/lucene"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

var fm *mapping.Mapping

func InitConvert(m *mapping.Mapping, covFunc map[string]func(string) (interface{}, error)) error {
	fm = m
	return nil
}

func LuceneToDSLNode(q *lucene.Lucene) (dsl.DSLNode, error) {
	return luceneToDSLNode(q)
}

func luceneToDSLNode(q *lucene.Lucene) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}

	if node, err := orQueryToDSLNode(q.OrQuery); err != nil {
		return nil, err
	} else {
		var nodes = []dsl.DSLNode{node}
		for _, query := range q.OSQuery {
			if node, err := osQueryToDSLNode(query); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, node)
			}
		}
		return &dsl.OrDSLNode{Nodes: nodes}, nil
	}
}

func orQueryToDSLNode(q *lucene.OrQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyOrQuery
	}
	if node, err := andQueryToDSLNode(q.AndQuery); err != nil {
		return nil, err
	} else {
		var nodes = []dsl.DSLNode{node}
		for _, query := range q.AnSQuery {
			if node, err := ansQueryToDSLNode(query); err != nil {
				return nil, err
			} else {
				nodes = append(nodes, node)
			}
		}
		return &dsl.AndDSLNode{MustNodes: nodes}, nil
	}
}

func osQueryToDSLNode(q *lucene.OSQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyOrQuery
	}
	return orQueryToDSLNode(q.OrQuery)
}

func andQueryToDSLNode(q *lucene.AndQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}
	if q.FieldQuery != nil {
		if node, err := fieldQueryToDSLNode(q.FieldQuery); err != nil {
			return nil, err
		} else if q.NotSymbol != nil {
			return &dsl.NotDSLNode{Nodes: []dsl.DSLNode{node}}, nil
		} else {
			return node, nil
		}
	} else {
		if node, err := parenQueryToDSLNode(q.ParenQuery); err != nil {
			return nil, err
		} else if q.NotSymbol != nil {
			return &dsl.NotDSLNode{Nodes: []dsl.DSLNode{node}}, nil
		} else {
			return node, nil
		}
	}
}

func ansQueryToDSLNode(q *lucene.AnSQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}
	return andQueryToDSLNode(q.AndQuery)
}

func parenQueryToDSLNode(q *lucene.ParenQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyParenQuery
	}
	return luceneToDSLNode(q.SubQuery)
}

// very import
func fieldQueryToDSLNode(q *lucene.FieldQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyFieldQuery
	}
	fmt.Println(fm)
	return nil, nil
}
