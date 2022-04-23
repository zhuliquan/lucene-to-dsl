package convert

import (
	"fmt"
	"net"
	"strings"

	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/go_tools/ip_tools"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"

	lucene "github.com/zhuliquan/lucene_parser"
	term "github.com/zhuliquan/lucene_parser/term"
)

var fm *mapping.Mapping

type termValue interface {
	Value(func(string) (interface{}, error)) (interface{}, error)
}

type rangeValue interface {
	termValue
	IsInf(int) bool
}

var convertToString = func(s string) (interface{}, error) { return s, nil }

var fieldTypeBits = map[mapping.FieldType]int{
	mapping.BYTE_FIELD_TYPE:           8,
	mapping.SHORT_FIELD_TYPE:          16,
	mapping.INTEGER_FIELD_TYPE:        32,
	mapping.INTERGER_RANGE_FIELD_TYPE: 32,
	mapping.LONG_FIELD_TYPE:           64,
	mapping.LONG_RANGE_FIELD_TYPE:     64,
	mapping.UNSIGNED_LONG_FIELD_TYPE:  64,
	mapping.HALF_FLOAT_FIELD_TYPE:     16,
	mapping.FLOAT_FIELD_TYPE:          32,
	mapping.FLOAT_RANGE_FIELD_TYPE:    32,
	mapping.DOUBLE_FIELD_TYPE:         64,
	mapping.DOUBLE_RANGE_FIELD_TYPE:   64,
	mapping.SCALED_FLOAT_FIELD_TYPE:   128,
}

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

	return &dsl.RangeNode{
		Field:       field.String(),
		ValueType:   property.Type,
		LeftValue:   leftValue,
		RightValue:  rightValue,
		LeftCmpSym:  leftCmp,
		RightCmpSym: rightCmp,
		Boost:       termV.Boost(),
	}, nil
}

func convertToSingle(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.DSLNode, error) {
	switch property.Type {
	case mapping.BOOLEAN_FIELD_TYPE,
		mapping.BYTE_FIELD_TYPE, mapping.SHORT_FIELD_TYPE,
		mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE,
		mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE, mapping.UNSIGNED_LONG_FIELD_TYPE,
		mapping.HALF_FLOAT_FIELD_TYPE, mapping.SCALED_FLOAT_FIELD_TYPE,
		mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE,
		mapping.VERSION_FIELD_TYPE,
		mapping.KEYWORD_FIELD_TYPE, mapping.CONSTANT_KEYWORD_FIELD_TYPE:
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
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
		var dateParser = getDateParserFromMapping(property)
		if _, err := termV.Value(convertToDate(dateParser)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to date math expr", field, termV.String())
		} else {
			// return &dsl.RangeNode{
			// 	Field:     field.String(),
			// 	ValeuType: dsl.LEAF_NODE_TYPE,
			// }, nil
			return nil, nil

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
		return nil, fmt.Errorf("ip value: %s is invalid", termV.String())

	case mapping.TEXT_FIELD_TYPE, mapping.WILDCARD_FIELD_TYPE:
		// var s, _ = termV.Value(func(s string) (interface{}, error) { return s, nil })
		// return &dsl.MatchNode{
		// 	EqNode: dsl.EqNode{
		// 		Field: field,
		// 		Type:  dsl.KEYWORD_VALUE,
		// 		Value: &dsl.LeafValue{
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
		Type:  property.Type,
		Value: s,
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
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
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

func getDateParserFromMapping(property *mapping.Property) *datemath_parser.DateMathParser {
	var opts = []datemath_parser.DateMathParserOption{}
	if property.Format != "" {
		opts = append(opts, datemath_parser.WithFormat(strings.Split(property.Format, "||")))
	}
	if dp, err := datemath_parser.NewDateMathParser(opts...); err != nil {
		panic(err)
	} else {
		return dp
	}
}
