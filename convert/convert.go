package convert

import (
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/lucene"
	"github.com/zhuliquan/lucene-to-dsl/lucene/term"
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
		var nodes = map[string][]dsl.DSLNode{node.GetId(): {node}}
		for _, query := range q.OSQuery {
			if curNode, err := osQueryToDSLNode(query); err != nil {
				return nil, err
			} else {
				if preNode, ok := nodes[curNode.GetId()]; ok {
					if curNode.GetNodeType() == dsl.OP_NODE_TYPE {
						nodes[curNode.GetId()] = append(nodes[curNode.GetId()], curNode)
					} else {
						if node, err := preNode[0].UnionJoin(curNode); err != nil {
							return nil, err
						} else {
							delete(nodes, curNode.GetId())
							nodes[node.GetId()] = []dsl.DSLNode{node}
						}
					}

				} else {
					nodes[curNode.GetId()] = []dsl.DSLNode{curNode}
				}
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
		var nodes = map[string][]dsl.DSLNode{node.GetId(): {node}}
		for _, query := range q.AnSQuery {
			if curNode, err := ansQueryToDSLNode(query); err != nil {
				return nil, err
			} else {
				if preNode, ok := nodes[curNode.GetId()]; ok {
					if curNode.GetNodeType() == dsl.OP_NODE_TYPE {
						nodes[curNode.GetId()] = append(nodes[curNode.GetId()], curNode)
					} else {
						if node, err := preNode[0].InterSect(curNode); err != nil {
							return nil, err
						} else {
							delete(nodes, curNode.GetId())
							nodes[node.GetId()] = []dsl.DSLNode{node}
						}
					}
				} else {
					nodes[curNode.GetId()] = []dsl.DSLNode{curNode}
				}
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
	var (
		node dsl.DSLNode
		err  error
	)
	if q.FieldQuery != nil {
		if node, err = fieldQueryToDSLNode(q.FieldQuery); err != nil {
			return nil, err
		}
	} else if q.ParenQuery != nil {
		if node, err = parenQueryToDSLNode(q.ParenQuery); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrEmptyAndQuery
	}

	if q.NotSymbol != nil {
		return node.Inverse()
	} else {
		return node, nil
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

func fieldQueryToDSLNode(q *lucene.FieldQuery) (dsl.DSLNode, error) {
	if q == nil {
		return nil, ErrEmptyFieldQuery
	} else if q.Field == nil || q.Term == nil {
		return nil, ErrEmptyFieldQuery
	}
	var termType = q.Term.GetTermType()

	if termType|term.RANGE_TERM_TYPE == term.RANGE_TERM_TYPE {

	}

	switch q.Term.GetTermType() {
	case term.RANGE_TERM_TYPE:
		return &dsl.RangeNode{}, nil
	case term.REGEXP_TERM_TYPE:
	}
	return nil, nil
}

func convertToRange() (dsl.DSLNode, error) {

}
