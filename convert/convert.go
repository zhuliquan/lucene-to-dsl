package convert

import (
	"fmt"
	"net"

	"github.com/zhuliquan/go_tools/ip_tools"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"

	lucene "github.com/zhuliquan/lucene_parser"
	term "github.com/zhuliquan/lucene_parser/term"
)

var convertMapping *mapping.PropertyMapping

func Init(pm *mapping.PropertyMapping) {
	convertMapping = pm
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
	if property, err := convertMapping.GetProperty(q.Field.String()); err != nil {
		return nil, err
	} else {
		var termType = q.Term.GetTermType()
		if termType|term.RANGE_TERM_TYPE == term.RANGE_TERM_TYPE {
			return convertToRange(q.Field, q.Term, property)
		} else if termType|term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
			return convertToSingle(q.Field, q.Term, property)
		} else if termType|term.PHRASE_TERM_TYPE == term.PHRASE_TERM_TYPE {
			return convertToPhrase(q.Field, q.Term, property)
		} else if termType|term.GROUP_TERM_TYPE == term.GROUP_TERM_TYPE {
			return convertToGroup(q.Field, q.Term, property)
		} else if termType|term.REGEXP_TERM_TYPE == term.REGEXP_TERM_TYPE {
			return convertToRegexp(q.Field, q.Term, property)
		} else {
			return nil, fmt.Errorf("con't convert term query: %s", q.String())
		}
	}
}

func convertToRange(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  dsl.LeafValue
		rightValue dsl.LeafValue
		leftCmp    dsl.CompareType
		rightCmp   dsl.CompareType
	)

	if bound.LeftInclude {
		leftCmp = dsl.GTE
	} else {
		leftCmp = dsl.GT
	}

	if bound.RightInclude {
		rightCmp = dsl.LTE
	} else {
		rightCmp = dsl.LT
	}

	if lv, err := termValueToLeafValue(bound.LeftValue, property); err != nil {
		return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s, err: %s",
			field, bound.LeftValue.String(), property.Type, err)
	} else {
		leftValue = lv
	}

	if rv, err := termValueToLeafValue(bound.RightValue, property); err != nil {
		return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s, err: %s",
			field, bound.RightValue.String(), property.Type, err)
	} else {
		rightValue = rv
	}

	var node = &dsl.RangeNode{
		Field:       field.String(),
		ValueType:   property.Type,
		LeftValue:   leftValue,
		RightValue:  rightValue,
		LeftCmpSym:  leftCmp,
		RightCmpSym: rightCmp,
		Boost:       termV.Boost(),
	}
	if err := dsl.CheckValidRangeNode(node); err != nil {
		return nil, fmt.Errorf("field: %s value: %s is invalid, err: %s", field, termV.String(), err)
	} else {
		return node, nil
	}
}

func convertToSingle(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	strVal, _ := termV.Value(convertToString)
	if strVal.(string) == "*" {
		return &dsl.ExistsNode{
			Field: field.String(),
		}, nil
	}
	if property.NullValue == strVal {
		return (&dsl.ExistsNode{
			Field: field.String(),
		}).Inverse()
	}
	return convertToNormal(field, termV, property, strVal.(string))
}

func convertToPhrase(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	strVal, _ := termV.Value(convertToString)
	return convertToNormal(field, termV, property, strVal.(string))
}

func convertToNormal(field *term.Field, termV *term.Term, property *mapping.Property, strVal string) (dsl.DSLNode, error) {
	switch property.Type {
	case mapping.BOOLEAN_FIELD_TYPE,
		mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE,
		mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE,
		mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE, mapping.UNSIGNED_LONG_FIELD_TYPE,
		mapping.HALF_FLOAT_FIELD_TYPE, mapping.SCALED_FLOAT_FIELD_TYPE,
		mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE,
		mapping.VERSION_FIELD_TYPE,
		mapping.KEYWORD_FIELD_TYPE, mapping.CONSTANT_KEYWORD_FIELD_TYPE, mapping.WILDCARD_FIELD_TYPE:
		if val, err := termValueToLeafValue(termV, property); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s, err: %s",
				field, termV.String(), property.Type, err)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: val,
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE, mapping.DATE_NANOS_FIELD_TYPE:
		var dateParser = getDateParserFromMapping(property)
		if _, err := termV.Value(convertToDate(dateParser)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to date math expr", field, termV.String())
		} else {
			// TODO: 需要考虑如何解决如何处理 日缺失想查一个月的锁有天的情况，月缺失想查整年的情况, 即：2019-02 / 2019。
			// var lowerDate, upperDate = getDateRange(d.(time.Time))
			// if reflect.DeepEqual(lowerDate, upperDate) {
			// 	return &dsl.TermNode{
			// 		EqNode: dsl.EqNode{
			// 			Field: field.String(),
			// 			Type:  mapping.KEYWORD_FIELD_TYPE,
			// 			Value: strVal,
			// 		},
			// 		Boost: termV.Boost(),
			// 	}, nil
			// } else {
			// 	return &dsl.RangeNode{
			// 		Field:       field.String(),
			// 		ValueType:   property.Type,
			// 		LeftValue:   lowerDate,
			// 		RightValue:  upperDate,
			// 		LeftCmpSym:  dsl.GTE,
			// 		RightCmpSym: dsl.LTE,
			// 		Boost:       termV.Boost(),
			// 	}, nil
			// }
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  mapping.KEYWORD_FIELD_TYPE,
					Value: strVal,
				},
				Boost: termV.Boost(),
			}, nil

		}
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		if ip, err := termV.Value(convertToIp); err == nil {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: ip,
				},
				Boost: termV.Boost(),
			}, nil
		}
		if ip1, ip2, err := ip_tools.GetRangeIpByIpCidr(termV.String()); err == nil {
			return &dsl.RangeNode{
				Field:       field.String(),
				ValueType:   property.Type,
				LeftValue:   net.IP(ip1),
				RightValue:  net.IP(ip2),
				LeftCmpSym:  dsl.GTE,
				RightCmpSym: dsl.LTE,
				Boost:       termV.Boost(),
			}, nil

		}
		return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s",
			field, termV.String(), property.Type)

	case mapping.TEXT_FIELD_TYPE, mapping.MATCH_ONLY_TEXT_FIELD_TYPE:
		if termV.GetTermType()|term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
			return &dsl.QueryStringNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: strVal,
				},
				Boost: termV.Boost(),
			}, nil
		} else {
			return &dsl.MatchPhraseNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: strVal,
				},
				Boost: termV.Boost(),
			}, nil
		}
	default:
		return nil, fmt.Errorf("field: %s mapping type: %s don't support lucene", field, property.Type)
	}
}

func convertToRegexp(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	var valStr, _ = termV.Value(convertToString)
	return &dsl.RegexpNode{EqNode: dsl.EqNode{
		Field: field.String(),
		Type:  property.Type,
		Value: valStr,
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
		if q.OrQuery == nil {
			return nil
		}
		for _, osGroup := range termGroup.OSTermGroup {
			if t := convertOsGroupToOsQuery(field, osGroup, boostSymbol); t != nil {
				q.OSQuery = append(q.OSQuery, t)
			}
		}
		return q
	}
}

func convertOsGroupToOsQuery(field *term.Field, osGroup *term.OSTermGroup, boostSymbol string) *lucene.OSQuery {
	if osGroup == nil {
		return nil
	} else {
		if t := convertOrGroupToOrQuery(field, osGroup.OrTermGroup, boostSymbol); t != nil {
			return &lucene.OSQuery{
				OrSymbol: osGroup.OrSymbol,
				OrQuery:  t,
			}
		} else {
			return nil
		}
	}
}

func convertOrGroupToOrQuery(field *term.Field, orGroup *term.OrTermGroup, boostSymbol string) *lucene.OrQuery {
	if orGroup == nil {
		return nil
	} else {
		var q = &lucene.OrQuery{}
		q.AndQuery = convertAndGroupToAndQuery(field, orGroup.AndTermGroup, boostSymbol)
		if q.AndQuery == nil {
			return nil
		}

		for _, ansGroup := range orGroup.AnSTermGroup {
			if t := convertAnsGroupToAnsQuery(field, ansGroup, boostSymbol); t != nil {
				q.AnSQuery = append(q.AnSQuery, t)
			}
		}
		return q
	}
}

func convertAnsGroupToAnsQuery(field *term.Field, ansGroup *term.AnSTermGroup, boostSymbol string) *lucene.AnSQuery {
	if ansGroup == nil {
		return nil
	} else {
		if t := convertAndGroupToAndQuery(field, ansGroup.AndTermGroup, boostSymbol); t != nil {
			return &lucene.AnSQuery{
				AndSymbol: ansGroup.AndSymbol,
				AndQuery:  t,
			}
		} else {
			return nil
		}
	}
}

func convertAndGroupToAndQuery(field *term.Field, andGroup *term.AndTermGroup, boostSymbol string) *lucene.AndQuery {
	if andGroup == nil {
		return nil
	} else if andGroup.TermGroupElem != nil {
		if t := convertGroupElemToFieldTerm(field, andGroup.TermGroupElem, boostSymbol); t != nil {
			return &lucene.AndQuery{
				NotSymbol:  andGroup.NotSymbol,
				FieldQuery: t,
			}
		} else {
			return nil
		}
	} else if andGroup.ParenTermGroup != nil {
		if t := convertParenTermGroupToParentTerm(field, andGroup.ParenTermGroup, boostSymbol); t != nil {
			return &lucene.AndQuery{
				NotSymbol:  andGroup.NotSymbol,
				ParenQuery: t,
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func convertParenTermGroupToParentTerm(field *term.Field, parentTermGroup *term.ParenTermGroup, boostSymbol string) *lucene.ParenQuery {
	if t := convertLogicTermGroupToLucene(field, parentTermGroup.SubTermGroup, boostSymbol); t != nil {
		return &lucene.ParenQuery{SubQuery: t}
	} else {
		return nil
	}
}

func convertGroupElemToFieldTerm(field *term.Field, groupElem *term.TermGroupElem, boostSymbol string) *lucene.FieldQuery {
	// it's impossible for groupElem is nil
	if groupElem.SingleTerm != nil {
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

func termValueToLeafValue(termV termValue, property *mapping.Property) (dsl.LeafValue, error) {
	switch typ := property.Type; typ {
	case mapping.BOOLEAN_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return false, nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToBool)
			} else {
				return true, nil
			}
		} else {
			return termV.Value(convertToBool)
		}
	case mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE,
		mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE,
		mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE:
		var bits = fieldTypeBits[typ]
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinInt[bits], nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToInt(bits))
			} else {
				return dsl.MaxInt[bits], nil
			}
		} else {
			return termV.Value(convertToInt(bits))
		}
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		var bits = fieldTypeBits[typ]
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinUint, nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToUInt(bits))
			} else {
				return dsl.MaxUint[bits], nil
			}
		} else {
			return termV.Value(convertToUInt(bits))
		}
	case mapping.HALF_FLOAT_FIELD_TYPE, mapping.SCALED_FLOAT_FIELD_TYPE,
		mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE,
		mapping.DOUBLE_FIELD_TYPE, mapping.DOUBLE_RANGE_FIELD_TYPE:
		var bits = fieldTypeBits[typ]
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinFloat[bits], nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToFloat(bits))
			} else {
				return dsl.MaxFloat[bits], nil
			}
		} else {
			return termV.Value(convertToFloat(bits))
		}
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinIP, nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToIp)
			} else {
				return dsl.MaxIP, nil
			}
		} else {
			return termV.Value(convertToIp)
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE, mapping.DATE_NANOS_FIELD_TYPE:
		var dateParser = getDateParserFromMapping(property)
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinTime, nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToDate(dateParser))
			} else {
				return dsl.MaxTime, nil
			}
		} else {
			return termV.Value(convertToDate(dateParser))
		}
	case mapping.VERSION_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinVersion, nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToVersion)
			} else {
				return dsl.MaxVersion, nil
			}
		} else {
			return termV.Value(convertToVersion)
		}
	default:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinString, nil
			} else if termR.IsInf(0) {
				return termR.Value(convertToString)
			} else {
				return dsl.MaxString, nil
			}
		} else {
			return termV.Value(convertToString)
		}
	}
}
