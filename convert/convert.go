package convert

import (
	"fmt"
	"net"
	"regexp"

	"github.com/zhuliquan/go_tools/ip_tools"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	"github.com/zhuliquan/lucene_parser"
	lucene "github.com/zhuliquan/lucene_parser"
	term "github.com/zhuliquan/lucene_parser/term"
)

type Converter interface {
	LuceneToAstNode(q *lucene.Lucene) (dsl.AstNode, error)
}

func NewConverter(mp *mapping.PropertyMapping) Converter {
	return &converter{
		mp: mp,
	}
}

type converter struct {
	mp *mapping.PropertyMapping
}

func (c *converter) LuceneToAstNode(q *lucene.Lucene) (dsl.AstNode, error) {
	return c.luceneToAstNode(q)
}

func (c *converter) luceneToAstNode(q *lucene.Lucene, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}

	if preNode, err := c.orQueryToAstNode(q.OrQuery, pp...); err != nil {
		return nil, err
	} else {
		var curNode dsl.AstNode
		var convertErr, unionErr error
		for _, osQuery := range q.OSQuery {
			if curNode, convertErr = c.osQueryToAstNode(osQuery, pp...); convertErr != nil {
				return nil, convertErr
			} else if preNode, unionErr = preNode.UnionJoin(curNode); unionErr != nil {
				return nil, unionErr
			}
		}
		return preNode, nil
	}
}

func (c *converter) orQueryToAstNode(q *lucene.OrQuery, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyOrQuery
	}
	if preNode, err := c.andQueryToAstNode(q.AndQuery, pp...); err != nil {
		return nil, err
	} else {
		var curNode dsl.AstNode
		var convertErr, unionErr error
		for _, ansQuery := range q.AnSQuery {
			if curNode, convertErr = c.ansQueryToAstNode(ansQuery, pp...); convertErr != nil {
				return nil, convertErr
			} else if preNode, unionErr = preNode.InterSect(curNode); unionErr != nil {
				return nil, unionErr
			}
		}
		return preNode, nil
	}
}

func (c *converter) osQueryToAstNode(q *lucene.OSQuery, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyOrQuery
	}
	return c.orQueryToAstNode(q.OrQuery, pp...)
}

func (c *converter) andQueryToAstNode(q *lucene.AndQuery, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}
	var (
		node dsl.AstNode
		err  error
	)
	if q.FieldQuery != nil {
		if node, err = c.fieldQueryToAstNode(q.FieldQuery, pp...); err != nil {
			return nil, err
		}
	} else if q.ParenQuery != nil {
		if node, err = c.parenQueryToAstNode(q.ParenQuery, pp...); err != nil {
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

func (c *converter) ansQueryToAstNode(q *lucene.AnSQuery, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyAndQuery
	}
	return c.andQueryToAstNode(q.AndQuery, pp...)
}

func (c *converter) parenQueryToAstNode(q *lucene.ParenQuery, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyParenQuery
	}
	return c.luceneToAstNode(q.SubQuery, pp...)
}

func (c *converter) fieldQueryToAstNode(q *lucene.FieldQuery, pp ...*mapping.Property) (dsl.AstNode, error) {
	if q == nil {
		return nil, ErrEmptyFieldQuery
	} else if q.Field == nil || q.Term == nil {
		return nil, ErrEmptyFieldQuery
	}

	var field = q.Field.String()
	if q.Field.String() == EXIST_FIELD {
		field = q.Term.String()
	}
	if field == "*" && q.Term.String() == "*" {
		return &dsl.MatchAllNode{}, nil
	}

	var props []*mapping.Property
	if len(pp) == 0 {
		props, err := c.mp.GetProperty(field)
		if err != nil {
			return nil, err
		}
		if len(props) == 0 {
			return nil, fmt.Errorf("field: %s don't match any es mapping", field)
		}
	} else {
		props = append(props, pp...)
	}

	var res dsl.AstNode = &dsl.EmptyNode{}
	for _, prop := range props {
		node, err := c.fieldQueryToAstNodeByProp(q, prop)
		if err != nil {
			return nil, err
		}
		if res, err = res.UnionJoin(node); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (c *converter) fieldQueryToAstNodeByProp(q *lucene.FieldQuery, property *mapping.Property) (dsl.AstNode, error) {
	var termType = q.Term.GetTermType()
	if termType&term.RANGE_TERM_TYPE == term.RANGE_TERM_TYPE {
		return c.convertToRange(q.Field, q.Term, property)
	} else if termType&term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
		return c.convertToSingle(q.Field, q.Term, property)
	} else if termType&term.PHRASE_TERM_TYPE == term.PHRASE_TERM_TYPE {
		return c.convertToPhrase(q.Field, q.Term, property)
	} else if termType&term.GROUP_TERM_TYPE == term.GROUP_TERM_TYPE {
		return c.convertToGroup(q.Field, q.Term, property)
	} else if termType&term.REGEXP_TERM_TYPE == term.REGEXP_TERM_TYPE {
		return c.convertToRegexp(q.Field, q.Term, property)
	} else {
		return nil, fmt.Errorf("con't convert term query: %s", q.String())
	}
}

func (c *converter) convertToRange(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
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
			dsl.NewValueType(property.Type, true),
			leftValue, rightValue, leftCmp, rightCmp,
		), dsl.WithBoost(termV.Boost().Float()),
	)
	if err := dsl.CheckValidRangeNode(node); err != nil {
		return nil, fmt.Errorf("field: %s value: %s is invalid, err: %s", field, termV.String(), err)
	} else {
		return node, nil
	}
}

func (c *converter) convertToSingle(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
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
	return c.convertToNormal(field, termV, property, strVal.(string))
}

func (c *converter) convertToPhrase(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	strVal, _ := termV.Value(convertToString)
	return c.convertToNormal(field, termV, property, strVal.(string))
}

func (c *converter) convertToNormal(field *term.Field, termV *term.Term, property *mapping.Property, strVal string) (dsl.AstNode, error) {
	// trick for id
	if field.String() == ID_FIELD {
		strLst, _ := termV.Value(toStrLst)
		return dsl.NewIdsNode(
			dsl.NewLfNode(), strLst.([]string),
		), nil
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
				dsl.WithBoost(termV.Boost().Float()),
			), nil
		}

	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE, mapping.DATE_NANOS_FIELD_TYPE:
		if dr, err := termV.Value(convertToDateRange(property)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to date math expr", field, termV.String())
		} else {
			var dateRange = dr.(*dateRange)
			return dsl.NewRangeNode(
				dsl.NewRgNode(
					dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
					dsl.NewValueType(property.Type, true),
					dateRange.from, dateRange.to, dsl.GTE, dsl.LTE,
				),
				dsl.WithBoost(termV.Boost().Float()),
			), nil
		}
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		if ip, err := termV.Value(convertToIp); err == nil {
			return dsl.NewTermNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(ip, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost().Float()),
			), nil
		}
		if ip1, ip2, err := ip_tools.GetRangeIpByIpCidr(termV.String()); err == nil {
			return dsl.NewRangeNode(dsl.NewRgNode(
				dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
				dsl.NewValueType(property.Type, true),
				net.IP(ip1), net.IP(ip2), dsl.GTE, dsl.LTE,
			), dsl.WithBoost(termV.Boost().Float())), nil
		}
		return nil, fmt.Errorf("field: %s value: %s is invalid, type: %s",
			field, termV.String(), property.Type)

	case mapping.TEXT_FIELD_TYPE, mapping.MATCH_ONLY_TEXT_FIELD_TYPE:
		if termV.GetTermType()&term.SINGLE_TERM_TYPE == term.SINGLE_TERM_TYPE {
			return dsl.NewQueryStringNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(strVal, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost().Float()),
			), nil
		} else {
			return dsl.NewMatchPhraseNode(
				dsl.NewKVNode(dsl.NewFieldNode(dsl.NewLfNode(), field.String()), dsl.NewValueNode(strVal, dsl.NewValueType(property.Type, true))),
				dsl.WithBoost(termV.Boost().Float()),
			), nil
		}
	default:
		return nil, fmt.Errorf("field: %s mapping type: %s don't support lucene", field, property.Type)
	}
}

func (c *converter) convertToRegexp(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	if !mapping.CheckStringType(property.Type) {
		return nil, fmt.Errorf("type: %s, don't support regex query, expect text", property.Type)
	}
	var valStr, _ = termV.Value(convertToString)
	if pattern, err := regexp.Compile(valStr.(string)); err != nil {
		return nil, fmt.Errorf("regexp str: %+v is invalid, err: %+v", valStr, err)
	} else {
		return dsl.NewRegexpNode(
			dsl.NewKVNode(
				dsl.NewFieldNode(dsl.NewLfNode(), field.String()),
				dsl.NewValueNode(valStr, dsl.NewValueType(property.Type, true)),
			),
			pattern,
		), nil
	}

}

func (c *converter) convertToGroup(field *term.Field, termV *term.Term, property *mapping.Property) (dsl.AstNode, error) {
	return c.luceneToAstNode(lucene_parser.TermGroupToLucene(field, termV.TermGroup), property)
}

func termValueToLeafValue(termV termValue, property *mapping.Property) (dsl.LeafValue, error) {
	switch typ := property.Type; typ {
	case mapping.BOOLEAN_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return false, nil
			} else if termR.IsInf(1) {
				return true, nil
			} else {
				return termR.Value(convertToBool)
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
			} else if termR.IsInf(1) {
				return dsl.MaxInt[bits], nil
			} else {
				return termR.Value(convertToInt(bits))
			}
		} else {
			return termV.Value(convertToInt(bits))
		}
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		var bits = fieldTypeBits[typ]
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinUint, nil
			} else if termR.IsInf(1) {
				return dsl.MaxUint[bits], nil
			} else {
				return termR.Value(convertToUInt(bits))
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
			} else if termR.IsInf(1) {
				return dsl.MaxFloat[bits], nil
			} else {
				return termR.Value(convertToFloat(bits, property.ScalingFactor))
			}
		} else {
			return termV.Value(convertToFloat(bits, property.ScalingFactor))
		}
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinIP, nil
			} else if termR.IsInf(1) {
				return dsl.MaxIP, nil
			} else {
				return termR.Value(convertToIp)
			}
		} else {
			return termV.Value(convertToIp)
		}
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE, mapping.DATE_NANOS_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinTime, nil
			} else if termR.IsInf(1) {
				return dsl.MaxTime, nil
			} else {
				return termR.Value(convertToDate(property))
			}
		} else {
			return termV.Value(convertToDate(property))
		}
	case mapping.VERSION_FIELD_TYPE:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinVersion, nil
			} else if termR.IsInf(1) {
				return dsl.MaxVersion, nil
			} else {
				return termR.Value(convertToVersion)
			}
		} else {
			return termV.Value(convertToVersion)
		}
	default:
		if termR, ok := termV.(rangeValue); ok {
			if termR.IsInf(-1) {
				return dsl.MinString, nil
			} else if termR.IsInf(1) {
				return dsl.MaxString, nil
			} else {
				return termR.Value(convertToString)
			}
		} else {
			return termV.Value(convertToString)
		}
	}
}

// func termValueToFuzzyNode(field *term.Field, termV *term.Term, property *mapping.Property)
