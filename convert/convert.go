package convert

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/zhuliquan/datemath_parser"
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

func LuceneToAstNode(q *lucene.Lucene) (dsl.AstNode, error) {
	return luceneToAstNode(q)
}

func luceneToAstNode(q *lucene.Lucene) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}

	if node, err := orQueryToAstNode(q.OrQuery); err != nil {
		return nil, err
	} else {
		var nodes = map[string][]dsl.AstNode{node.NodeKey(): {node}}
		for _, query := range q.OSQuery {
			if curNode, err := osQueryToAstNode(query); err != nil {
				return nil, err
			} else {
				if preNode, ok := nodes[curNode.NodeKey()]; ok {
					if curNode.DslType() == dsl.AND_DSL_TYPE ||
						curNode.DslType() == dsl.NOT_DSL_TYPE {
						nodes[curNode.NodeKey()] = append(nodes[curNode.NodeKey()], curNode)
					} else {
						if node, err := preNode[0].UnionJoin(curNode); err != nil {
							return nil, err
						} else {
							delete(nodes, curNode.NodeKey())
							nodes[node.NodeKey()] = []dsl.AstNode{node}
						}
					}
				} else {
					nodes[curNode.NodeKey()] = []dsl.AstNode{curNode}
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
		return &dsl.OrNode{Nodes: nodes}, nil
	}
}

func orQueryToAstNode(q *lucene.OrQuery) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyOrQuery
	}
	if node, err := andQueryToAstNode(q.AndQuery); err != nil {
		return nil, err
	} else {
		var nodes = map[string][]dsl.AstNode{node.NodeKey(): {node}}
		for _, query := range q.AnSQuery {
			if curNode, err := ansQueryToAstNode(query); err != nil {
				return nil, err
			} else {
				if preNode, ok := nodes[curNode.NodeKey()]; ok {
					if curNode.DslType() == dsl.OR_DSL_TYPE {
						nodes[curNode.NodeKey()] = append(nodes[curNode.NodeKey()], curNode)
					} else {
						if node, err := preNode[0].InterSect(curNode); err != nil {
							return nil, err
						} else {
							delete(nodes, curNode.NodeKey())
							nodes[node.NodeKey()] = []dsl.AstNode{node}
						}
					}
				} else {
					nodes[curNode.NodeKey()] = []dsl.AstNode{curNode}
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
		return &dsl.AndNode{MustNodes: nodes}, nil
	}
}

func osQueryToAstNode(q *lucene.OSQuery) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyOrQuery
	}
	return orQueryToAstNode(q.OrQuery)
}

func andQueryToAstNode(q *lucene.AndQuery) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}
	var (
		node dsl.AstNode
		err  error
	)
	if q.FieldQuery != nil {
		if node, err = fieldQueryToAstNode(q.FieldQuery); err != nil {
			return nil, err
		}
	} else if q.ParenQuery != nil {
		if node, err = parenQueryToAstNode(q.ParenQuery); err != nil {
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

func ansQueryToAstNode(q *lucene.AnSQuery) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}
	return andQueryToAstNode(q.AndQuery)
}

func parenQueryToAstNode(q *lucene.ParenQuery) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyParenQuery
	}
	return luceneToAstNode(q.SubQuery)
}

func fieldQueryToAstNode(q *lucene.FieldQuery) (dsl.AstNode, error) {
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

func convertToRange(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
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

	var node = dsl.NewRangeNode(
		dsl.NewRgNode(
			dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
			property.Type, leftValue, rightValue, leftCmp, rightCmp,
		), dsl.WithBoost(termV.Boost()),
	)
	if err := dsl.CheckValidRangeNode(node); err != nil {
		return nil, fmt.Errorf("field: %s value: %s is invalid, err: %s", field, termV.String(), err)
	} else {
		return node, nil
	}
}

func convertToSingle(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	strVal, _ := termV.Value(convertToString)
	if strVal.(string) == "*" {
		return dsl.NewExistsNode(
			dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
		), nil
	}
	if property.NullValue == strVal {
		return dsl.NewExistsNode(
			dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
		).Inverse()
	}
	return convertToNormal(field, termV, property, strVal.(string))
}

func convertToPhrase(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	strVal, _ := termV.Value(convertToString)
	return convertToNormal(field, termV, property, strVal.(string))
}

func convertToNormal(field *term.Field, termV *term.Term, property *mapping.Property, strVal string) (dsl.AstNode, error) {
	// trick for id
	if field.String() == "_id" {
		if strLst, err := termV.Value(toStrLst); err != nil {
			return nil, err
		} else {
			return dsl.NewIdsNode(
				dsl.NewLfNode(), strLst.([]string),
			), nil
		}
	}
	switch property.Type {
	case mapping.BOOLEAN_FIELD_TYPE,
		mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE,
		mapping.INTEGER_FIELD_TYPE, mapping.INTEGER_RANGE_FIELD_TYPE,
		mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE, mapping.UNSIGNED_LONG_FIELD_TYPE,
		mapping.HALF_FLOAT_FIELD_TYPE, mapping.SCALED_FLOAT_FIELD_TYPE,
		mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE,
		mapping.VERSION_FIELD_TYPE,
		mapping.KEYWORD_FIELD_TYPE, mapping.CONSTANT_KEYWORD_FIELD_TYPE, mapping.WILDCARD_FIELD_TYPE:
		if val, err := termValueToLeafValue(termV, property); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s, err: %s",
				field, termV.String(), property.Type, err)
		} else {
			return dsl.NewTermNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(val, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost()),
			), nil
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE, mapping.DATE_NANOS_FIELD_TYPE:
		var dateParser = getDateParserFromMapping(property)
		if d, err := termV.Value(convertToDate(dateParser)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to date math expr", field, termV.String())
		} else {
			// TODO: 需要考虑如何解决如何处理 日缺失想查一个月的锁有天的情况，月缺失想查整年的情况, 即：2019-02 / 2019。
			var lowerDate, upperDate = getDateRange(d.(time.Time))
			return dsl.NewRangeNode(
				dsl.NewRgNode(
					dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
					property.Type,
					lowerDate, upperDate,
					dsl.GTE, dsl.LTE,
				),
				dsl.WithBoost(termV.Boost()),
				dsl.WithFormat(datemath_parser.EPOCH_MILLIS),
			), nil
		}
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		if ip, err := termV.Value(convertToIp); err == nil {
			return dsl.NewTermNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(ip, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost()),
			), nil
		}
		if ip1, ip2, err := ip_tools.GetRangeIpByIpCidr(termV.String()); err == nil {
			return dsl.NewRangeNode(dsl.NewRgNode(
				dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
				property.Type, net.IP(ip1), net.IP(ip2), dsl.GTE, dsl.LTE,
			), dsl.WithBoost(termV.Boost())), nil
		}
		return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s",
			field, termV.String(), property.Type)

	case mapping.TEXT_FIELD_TYPE, mapping.MATCH_ONLY_TEXT_FIELD_TYPE:
		if termV.GetTermType()|term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
			return dsl.NewQueryStringNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(strVal, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost()),
			), nil
		} else {
			return dsl.NewMatchPhraseNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(strVal, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost()),
			), nil
		}
	default:
		return nil, fmt.Errorf("field: %s mapping type: %s don't support lucene", field, property.Type)
	}
}

func convertToRegexp(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	if !mapping.CheckStringType(property.Type) {
		return nil, fmt.Errorf("type: %s, don't support regex query, expect text", property.Type)
	}
	var valStr, _ = termV.Value(convertToString)
	if pattern, err := regexp.Compile(valStr.(string)); err != nil {
		return nil, fmt.Errorf("regexp str: %+v is invalid, err: %+v", valStr, err)
	} else {
		return dsl.NewRegexNode(
			dsl.NewKVNode(
				dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
				dsl.NewValueNode(valStr, dsl.NewValueType(property.Type, true)),
			),
			pattern,
		), nil
	}

}

func convertToGroup(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	return luceneToAstNode(convertTermGroupToLucene(field, termV.TermGroup))
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
		mapping.INTEGER_FIELD_TYPE, mapping.INTEGER_RANGE_FIELD_TYPE,
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
				return termR.Value(convertToFloat(bits, property.ScalingFactor))
			} else {
				return dsl.MaxFloat[bits], nil
			}
		} else {
			return termV.Value(convertToFloat(bits, property.ScalingFactor))
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

// func termValueToFuzzyNode(field *term.Field, termV *term.Term, property *mapping.Property)
