package convert

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	lucene "github.com/zhuliquan/lucene_parser"
	term "github.com/zhuliquan/lucene_parser/term"
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
					if curNode.GetDSLType() == dsl.AND_DSL_TYPE ||
						curNode.GetDSLType() == dsl.NOT_DSL_TYPE {
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
		if len(nodes) == 1 {
			for _, ns := range nodes {
				if len(ns) == 1 {
					return ns[0], nil
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
					if curNode.GetDSLType() == dsl.OR_DSL_TYPE {
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
		if len(nodes) == 1 {
			for _, ns := range nodes {
				if len(ns) == 1 {
					return ns[0], nil
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

	var property, _ = mapping.GetProperty(q.Field.String())
	if property.NullValue == q.Term.String() || "\""+property.NullValue+"\"" == q.Term.String() {
		var d = &dsl.ExistsNode{Field: q.Field.String()}
		return d.Inverse()
	}
	var termType = q.Term.GetTermType()
	if termType|term.RANGE_TERM_TYPE == term.RANGE_TERM_TYPE {
		return convertToRange(q.Field.String(), q.Term, property)
	} else if termType|term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
		return convertToSingle(q.Field.String(), q.Term, property)
	} else if termType|term.PHRASE_TERM_TYPE == term.PHRASE_TERM_TYPE {
		return convertToSingle(q.Field.String(), q.Term, property)
	} else if termType|term.GROUP_TERM_TYPE == term.GROUP_TERM_TYPE {
		return convertToGroup(q.Field.String(), q.Term, property)
	} else if termType|term.REGEXP_TERM_TYPE == term.REGEXP_TERM_TYPE {
		return convertToRegexp(q.Field.String(), q.Term, property)
	} else {
		return nil, fmt.Errorf("con't convert term query: %s", q.String())
	}
}

func convertToRange(field string, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	return nil, nil
}

func convertToSingle(field string, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	switch property.Type {
	case mapping.BOOLEAN_FIELD_TYPE:
		if termV.String() == "true" || termV.String() == "false" || termV.String() == "\"true\"" || termV.String() == "\"false\"" {
			return &dsl.TermNode{}, nil
		}
	case mapping.BYTE_FIELD_TYPE:
		if i, err := termV.Value(convertToInt64); err != nil {
			return nil, err
		} else if i64 := i.(int64); i64 < -128 || i64 > 127 {
			return nil, fmt.Errorf("field: %s is byte type, expect value range [-128, 127]", field)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field,
					Type:  dsl.INT_VALUE,
					Value: &dsl.DSLTermValue{
						IntTerm: i64,
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.SHORT_FIELD_TYPE:
		if i, err := termV.Value(convertToInt64); err != nil {
			return nil, err
		} else if i64 := i.(int64); i64 < -32768 || i64 > 32767 {
			return nil, fmt.Errorf("field: %s is short type, expect value range [-32768, 32767]", field)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field,
					Type:  dsl.INT_VALUE,
					Value: &dsl.DSLTermValue{
						IntTerm: i64,
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE:
		if i, err := termV.Value(convertToInt64); err != nil {
			return nil, err
		} else if i64 := i.(int64); i64 < -2147483648 || i64 > 2147483647 {
			return nil, fmt.Errorf("field: %s is short type, expect value range [-2147483648, 2147483647]", field)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field,
					Type:  dsl.INT_VALUE,
					Value: &dsl.DSLTermValue{
						IntTerm: i64,
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.LONG_FIELD_TYPE:
		if i, err := termV.Value(convertToInt64); err != nil {
			return nil, err
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field,
					Type:  dsl.INT_VALUE,
					Value: &dsl.DSLTermValue{
						IntTerm: i.(int64),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		if i, err := termV.Value(convertToUInt64); err != nil {
			return nil, err
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field,
					Type:  dsl.UINT_VALUE,
					Value: &dsl.DSLTermValue{
						UintTerm: i.(uint64),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}

	}

	return nil, nil
}

func convertToRegexp(field string, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	var s, _ = termV.Value(func(x string) (interface{}, error) { return x, nil })
	return &dsl.RegexpNode{EqNode: dsl.EqNode{
		Field: field,
		Type:  dsl.PHRASE_VALUE,
		Value: &dsl.DSLTermValue{
			StringTerm: s.(string),
		},
	}}, nil
}

func convertToGroup(field string, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	return nil, nil
}
