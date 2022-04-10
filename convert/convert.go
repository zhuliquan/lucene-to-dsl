package convert

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/zhuliquan/datemath_parser"
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
		return convertToRange(q.Field, q.Term, property)
	} else if termType|term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
		return convertToSingle(q.Field, q.Term, property)
	} else if termType|term.PHRASE_TERM_TYPE == term.PHRASE_TERM_TYPE {
		return convertToSingle(q.Field, q.Term, property)
	} else if termType|term.GROUP_TERM_TYPE == term.GROUP_TERM_TYPE {
		return convertToGroup(q.Field, q.Term, property)
	} else if termType|term.REGEXP_TERM_TYPE == term.REGEXP_TERM_TYPE {
		return convertToRegexp(q.Field, q.Term, property)
	} else {
		return nil, fmt.Errorf("con't convert term query: %s", q.String())
	}
}

func convertToRange(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.DSLTermValue
		rightValue *dsl.DSLTermValue
		leftCmp    dsl.CompareType
		rightCmp   dsl.CompareType
	)

	if bound.LeftInclude {
		leftCmp = dsl.GTE
	} else {
		leftCmp = dsl.GT
	}

	if bound.RightInclude {
		leftCmp = dsl.LTE
	} else {
		leftCmp = dsl.LT
	}

	switch property.Type {
	case mapping.BYTE_FIELD_TYPE:
		if !bound.LeftValue.IsInf() {
			if li, err := bound.LeftValue.Value(convertToInt64); err != nil {
				return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int8", field, bound.LeftValue.String())
			} else if li64 := li.(int64); li64 < -128 || li64 > 127 {
				return nil, fmt.Errorf("field: %s is byte type, expect left value range [-128, 127]", field)
			} else {
				leftValue = &dsl.DSLTermValue{IntTerm: li64}
			}
		} else {
			leftValue = &dsl.DSLTermValue{IntTerm: -128}
		}
		if !bound.RightValue.IsInf() {
			if ri, err := bound.RightValue.Value(convertToInt64); err != nil {
				return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int8", field, bound.RightValue.String())
			} else if ri64 := ri.(int64); ri64 < -128 || ri64 > 127 {
				return nil, fmt.Errorf("field: %s is byte type, expect right value range [-128, 127]", field)
			} else {
				rightValue = &dsl.DSLTermValue{IntTerm: ri64}
			}
		} else {
			rightValue = &dsl.DSLTermValue{IntTerm: 127}
		}

		return &dsl.RangeNode{
			Field:       field.String(),
			ValueType:   dsl.INT_VALUE,
			LeftValue:   leftValue,
			RightValue:  rightValue,
			LeftCmpSym:  leftCmp,
			RightCmpSym: rightCmp,
			Boost:       termV.Boost(),
		}, nil
	case mapping.SHORT_FIELD_TYPE:
		if !bound.LeftValue.IsInf() {
			if li, err := bound.LeftValue.Value(convertToInt64); err != nil {
				return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int16", field, bound.LeftValue.String())
			} else if li64 := li.(int64); li64 < -32768 || li64 > 32767 {
				return nil, fmt.Errorf("field: %s is byte type, expect left value range [-32768, 32767]", field)
			} else {
				leftValue = &dsl.DSLTermValue{IntTerm: li64}
			}
		} else {
			leftValue = &dsl.DSLTermValue{IntTerm: -32768}
		}
		if !bound.RightValue.IsInf() {
			if ri, err := bound.RightValue.Value(convertToInt64); err != nil {
				return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int16", field, bound.RightValue.String())
			} else if ri64 := ri.(int64); ri64 <= -32768 || ri64 >= 32767 {
				return nil, fmt.Errorf("field: %s is byte type, expect right value range [-32768, 32767]", field)
			} else {
				rightValue = &dsl.DSLTermValue{IntTerm: ri64}
			}
		} else {
			rightValue = &dsl.DSLTermValue{IntTerm: 32767}
		}

		return &dsl.RangeNode{
			Field:       field.String(),
			ValueType:   dsl.INT_VALUE,
			LeftValue:   leftValue,
			RightValue:  rightValue,
			LeftCmpSym:  leftCmp,
			RightCmpSym: rightCmp,
			Boost:       termV.Boost(),
		}, nil

	case mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE:
		if !bound.LeftValue.IsInf() {
			if li, err := bound.LeftValue.Value(convertToInt64); err != nil {
				return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int32", field, bound.LeftValue.String())
			} else if li64 := li.(int64); li64 < -2147483648 || li64 > 2147483647 {
				return nil, fmt.Errorf("field: %s is byte type, expect left value range [-2147483648, 2147483647]", field)
			} else {
				leftValue = &dsl.DSLTermValue{IntTerm: li64}
			}
		} else {
			leftValue = &dsl.DSLTermValue{IntTerm: -2147483648}
		}
		if !bound.RightValue.IsInf() {
			if ri, err := bound.RightValue.Value(convertToInt64); err != nil {
				return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int32", field, bound.RightValue.String())
			} else if ri64 := ri.(int64); ri64 < -2147483648 || ri64 > 2147483647 {
				return nil, fmt.Errorf("field: %s is byte type, expect right value range [-2147483648, 2147483647]", field)
			} else {
				rightValue = &dsl.DSLTermValue{IntTerm: ri64}
			}
		} else {
			rightValue = &dsl.DSLTermValue{IntTerm: 2147483647}
		}

		return &dsl.RangeNode{
			Field:       field.String(),
			ValueType:   dsl.INT_VALUE,
			LeftValue:   leftValue,
			RightValue:  rightValue,
			LeftCmpSym:  leftCmp,
			RightCmpSym: rightCmp,
			Boost:       termV.Boost(),
		}, nil

	case mapping.UNSIGNED_LONG_FIELD_TYPE:

	}

	return nil, nil
}

func convertToSingle(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	switch property.Type {
	case mapping.BOOLEAN_FIELD_TYPE:
		if b, err := termV.Value(convertToBool); err != nil {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  dsl.BOOL_VALUE,
					Value: &dsl.DSLTermValue{
						BoolTerm: b.(bool),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.BYTE_FIELD_TYPE:
		if i, err := termV.Value(convertToInt64); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int8", field, termV.String())
		} else if i64 := i.(int64); i64 < -128 || i64 > 127 {
			return nil, fmt.Errorf("field: %s is byte type, expect value range [-128, 127]", field)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
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
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int16", field, termV.String())
		} else if i64 := i.(int64); i64 < -32768 || i64 > 32767 {
			return nil, fmt.Errorf("field: %s is short type, expect value range [-32768, 32767]", field)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
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
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int32", field, termV.String())
		} else if i64 := i.(int64); i64 < -2147483648 || i64 > 2147483647 {
			return nil, fmt.Errorf("field: %s is short type, expect value range [-2147483648, 2147483647]", field)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
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
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int64", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
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
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to uint64", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  dsl.UINT_VALUE,
					Value: &dsl.DSLTermValue{
						UintTerm: i.(uint64),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}

	case mapping.HALF_FLOAT_FIELD_TYPE, mapping.FLOAT_FIELD_TYPE, mapping.SCALED_FLOAT_FIELD_TYPE, mapping.DOUBLE_FIELD_TYPE:
		if f, err := termV.Value(convertToFloat64); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to float", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  dsl.FLOAT_VALUE,
					Value: &dsl.DSLTermValue{
						FloatTerm: f.(float64),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
		if parser, err := datemath_parser.NewDateMathParser(
			datemath_parser.WithFormat(strings.Split(property.Format, "||")),
			datemath_parser.WithTimeZone(time.Local.String()),
		); err != nil {
			return nil, fmt.Errorf("failed to create date math parser, err: %+v", err)
		} else {
			if _, err := termV.Value(convertToDate(parser)); err != nil {
				return nil, fmt.Errorf("field: %s value: %s is invalid, expect to date math expr", field, termV.String())
			} else {
				// return &dsl.RangeNode{
				// 	Field:     field.String(),
				// 	ValeuType: dsl.LEAF_NODE_TYPE,
				// }, nil
				return nil, nil

			}
		}
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		if ip, err := termV.Value(convertToIp); err == nil {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  dsl.IP_VALUE,
					Value: &dsl.DSLTermValue{
						IpTerm: ip.(net.IP),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
		if cidr, err := termV.Value(convertToCidr); err == nil {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  dsl.IP_CIDR_VALUE,
					Value: &dsl.DSLTermValue{
						IpCidrTerm: cidr.(*net.IPNet),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
		return nil, fmt.Errorf("ip value: %s is invalid", termV.String())
	case mapping.KEYWORD_FIELD_TYPE:
		var s, _ = termV.Value(func(s string) (interface{}, error) { return s, nil })
		return &dsl.TermNode{
			EqNode: dsl.EqNode{
				Field: field.String(),
				Type:  dsl.KEYWORD_VALUE,
				Value: &dsl.DSLTermValue{
					StringTerm: s.(string),
				},
			},
			Boost: termV.Boost(),
		}, nil

	case mapping.TEXT_FIELD_TYPE:
		// var s, _ = termV.Value(func(s string) (interface{}, error) { return s, nil })
		// return &dsl.TermNode{
		// 	EqNode: dsl.EqNode{
		// 		Field: field,
		// 		Type:  dsl.KEYWORD_VALUE,
		// 		Value: &dsl.DSLTermValue{
		// 			StringTerm: s.(string),
		// 		},
		// 	},
		// 	Boost: termV.Boost(),
		// }, nil
		return nil, nil
	}

	return nil, nil
}

func convertToRegexp(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	var s, _ = termV.Value(func(x string) (interface{}, error) { return x, nil })
	return &dsl.RegexpNode{EqNode: dsl.EqNode{
		Field: field.String(),
		Type:  dsl.PHRASE_VALUE,
		Value: &dsl.DSLTermValue{
			StringTerm: s.(string),
		},
	}}, nil
}

func convertToGroup(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	return luceneToDSLNode(convertTermGroupToLucene(field, termV.TermGroup))
}

func convertTermGroupToLucene(field *term.Field, termGroup *term.TermGroup) *lucene.Lucene {
	if termGroup == nil {
		return nil
	} else {
		return convertLogicTermGroupToLucene(field, termGroup.LogicTermGroup, termGroup.BoostSymbol)
	}
}

func convertLogicTermGroupToLucene(field *term.Field, termGroup *term.LogicTermGroup, boostSymbol string) *lucene.Lucene {
	if termGroup == nil {
		return nil
	} else {
		var q = &lucene.Lucene{}
		q.OrQuery = convertOrGroupToOrQuery(field, termGroup.OrTermGroup, boostSymbol)
		q.OSQuery = []*lucene.OSQuery{}
		for _, osGroup := range termGroup.OSTermGroup {
			q.OSQuery = append(q.OSQuery, convertOsGroupToOsQuery(field, osGroup, boostSymbol))
		}
		return q
	}
}

func convertOsGroupToOsQuery(field *term.Field, osGroup *term.OSTermGroup, boostSymbol string) *lucene.OSQuery {
	if osGroup == nil {
		return nil
	} else {
		return &lucene.OSQuery{
			OrSymbol: osGroup.OrSymbol,
			OrQuery:  convertOrGroupToOrQuery(field, osGroup.OrTermGroup, boostSymbol),
		}
	}
}

func convertOrGroupToOrQuery(field *term.Field, orGroup *term.OrTermGroup, boostSymbol string) *lucene.OrQuery {
	if orGroup == nil {
		return nil
	} else {
		var q = &lucene.OrQuery{}
		q.AndQuery = convertAndGroupToAndQuery(field, orGroup.AndTermGroup, boostSymbol)
		q.AnSQuery = []*lucene.AnSQuery{}
		for _, ansGroup := range orGroup.AnSTermGroup {
			q.AnSQuery = append(q.AnSQuery, convertAnsGroupToAnsQuery(field, ansGroup, boostSymbol))
		}
		return q
	}
}

func convertAnsGroupToAnsQuery(field *term.Field, ansGroup *term.AnSTermGroup, boostSymbol string) *lucene.AnSQuery {
	if ansGroup == nil {
		return nil
	} else {
		return &lucene.AnSQuery{
			AndSymbol: ansGroup.AndSymbol,
			AndQuery:  convertAndGroupToAndQuery(field, ansGroup.AndTermGroup, boostSymbol),
		}
	}
}

func convertAndGroupToAndQuery(field *term.Field, andGroup *term.AndTermGroup, boostSymbol string) *lucene.AndQuery {
	if andGroup == nil {
		return nil
	} else if andGroup.TermGroupElem != nil {
		return &lucene.AndQuery{
			NotSymbol:  andGroup.NotSymbol,
			FieldQuery: convertGroupElemToFieldTerm(field, andGroup.TermGroupElem, boostSymbol),
		}
	} else if andGroup.ParenTermGroup != nil {
		return &lucene.AndQuery{
			NotSymbol:  andGroup.NotSymbol,
			ParenQuery: convertParenTermGroupToParentTerm(field, andGroup.ParenTermGroup, boostSymbol),
		}
	} else {
		return nil
	}
}

func convertParenTermGroupToParentTerm(field *term.Field, parentTermGroup *term.ParenTermGroup, boostSymbol string) *lucene.ParenQuery {
	return nil
}

func convertGroupElemToFieldTerm(field *term.Field, groupElem *term.TermGroupElem, boostSymbol string) *lucene.FieldQuery {
	if groupElem == nil {
		return nil
	} else if groupElem.SingleTerm != nil {
		return &lucene.FieldQuery{
			Field: field,
			Term: &term.Term{
				FuzzyTerm: &term.FuzzyTerm{
					SingleTerm:  groupElem.SingleTerm,
					BoostSymbol: boostSymbol,
				},
			},
		}
	} else if groupElem.PhraseTerm != nil {
		return &lucene.FieldQuery{
			Field: field,
			Term: &term.Term{
				FuzzyTerm: &term.FuzzyTerm{
					PhraseTerm:  groupElem.PhraseTerm,
					BoostSymbol: boostSymbol,
				},
			},
		}
	} else if groupElem.SRangeTerm != nil {
		return &lucene.FieldQuery{
			Field: field,
			Term: &term.Term{
				RangeTerm: &term.RangeTerm{
					SRangeTerm:  groupElem.SRangeTerm,
					BoostSymbol: boostSymbol,
				},
			},
		}
	} else if groupElem.DRangeTerm != nil {
		return &lucene.FieldQuery{
			Field: field,
			Term: &term.Term{
				RangeTerm: &term.RangeTerm{
					DRangeTerm:  groupElem.DRangeTerm,
					BoostSymbol: boostSymbol,
				},
			},
		}
	} else {
		return nil
	}
}
