package convert

import (
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/shopspring/decimal"
	"github.com/x448/float16"
	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/go_tools/ip_tools"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"

	lucene "github.com/zhuliquan/lucene_parser"
	term "github.com/zhuliquan/lucene_parser/term"
)

var fm *mapping.Mapping
var dateParser *datemath_parser.DateMathParser

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
		bound    = termV.GetBound()
		leftCmp  dsl.CompareType
		rightCmp dsl.CompareType
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

	switch property.Type {
	case mapping.BYTE_FIELD_TYPE:
		return convertToInt8Range(field, termV, leftCmp, rightCmp, property)
	case mapping.SHORT_FIELD_TYPE:
		return convertToInt16Range(field, termV, leftCmp, rightCmp, property)
	case mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE:
		return convertToInt32Range(field, termV, leftCmp, rightCmp, property)
	case mapping.LONG_FIELD_TYPE, mapping.LONG_RANGE_FIELD_TYPE:
		return convertToInt64Range(field, termV, leftCmp, rightCmp, property)
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		return convertToUint64Range(field, termV, leftCmp, rightCmp, property)
	case mapping.HALF_FLOAT_FIELD_TYPE:
		return convertToFloat16Range(field, termV, leftCmp, rightCmp, property)
	case mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE:
		return convertToFloat32Range(field, termV, leftCmp, rightCmp, property)
	case mapping.DOUBLE_FIELD_TYPE, mapping.DOUBLE_RANGE_FIELD_TYPE:
		return convertToFloat64Range(field, termV, leftCmp, rightCmp, property)
	case mapping.IP_FIELD_TYPE, mapping.IP_RANGE_FIELD_TYPE:
		return convertToIPRange(field, termV, leftCmp, rightCmp, property)
	case mapping.DATE_FIELD_TYPE, mapping.DATE_RANGE_FIELD_TYPE:
		return convertToDateRange(field, termV, leftCmp, rightCmp, property)
	default:
		return convertToStringRange(field, termV, leftCmp, rightCmp, property)
	}
}

func convertToStringRange(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)

	if !bound.LeftValue.IsInf() {
		leftValue = &dsl.LeafValue{String: bound.LeftValue.String()}
	} else {
		leftValue = &dsl.LeafValue{String: ""}
	}
	if !bound.RightValue.IsInf() {
		rightValue = &dsl.LeafValue{String: bound.RightValue.String()}
	} else {
		rightValue = &dsl.LeafValue{String: "~"}
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

func convertToDateRange(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)

	if !bound.LeftValue.IsInf() {
		if ld, err := bound.LeftValue.Value(convertToDate(dateParser)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to date", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{DateTime: ld.(time.Time)}
		}
	} else {
		leftValue = &dsl.LeafValue{DateTime: dsl.MinTime}
	}
	if !bound.RightValue.IsInf() {
		if rd, err := bound.RightValue.Value(convertToDate(dateParser)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to date", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{DateTime: rd.(time.Time)}
		}
	} else {
		rightValue = &dsl.LeafValue{DateTime: dsl.MaxTime}
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

func convertToInt8Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)

	if !bound.LeftValue.IsInf() {
		if li, err := bound.LeftValue.Value(convertToInt(8)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int8", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{TinyInt: int16(li.(int64))}
		}
	} else {
		leftValue = &dsl.LeafValue{TinyInt: math.MinInt8}
	}
	if !bound.RightValue.IsInf() {
		if ri, err := bound.RightValue.Value(convertToInt(8)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int8", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{TinyInt: int16(ri.(int64))}
		}
	} else {
		rightValue = &dsl.LeafValue{TinyInt: math.MaxInt8}
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

func convertToInt16Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if li, err := bound.LeftValue.Value(convertToInt(16)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int16", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{TinyInt: int16(li.(int64))}
		}
	} else {
		leftValue = &dsl.LeafValue{TinyInt: math.MinInt16}
	}
	if !bound.RightValue.IsInf() {
		if ri, err := bound.RightValue.Value(convertToInt(16)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int16", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{TinyInt: int16(ri.(int64))}
		}
	} else {
		rightValue = &dsl.LeafValue{TinyInt: math.MaxInt16}
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

func convertToInt32Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if li, err := bound.LeftValue.Value(convertToInt(32)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int32", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{Decimal: decimal.NewFromInt32(int32(li.(int64)))}
		}
	} else {
		leftValue = &dsl.LeafValue{Decimal: decimal.NewFromInt32(math.MinInt32)}
	}
	if !bound.RightValue.IsInf() {
		if ri, err := bound.RightValue.Value(convertToInt(32)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int32", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{Decimal: decimal.NewFromInt32(int32(ri.(int64)))}
		}
	} else {
		rightValue = &dsl.LeafValue{Decimal: decimal.NewFromInt32(math.MaxInt32)}
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

func convertToInt64Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if li, err := bound.LeftValue.Value(convertToInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to int64", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{Decimal: decimal.NewFromInt(li.(int64))}
		}
	} else {
		leftValue = &dsl.LeafValue{Decimal: decimal.NewFromInt(math.MinInt64)}
	}
	if !bound.RightValue.IsInf() {
		if ri, err := bound.RightValue.Value(convertToInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to int64", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{Decimal: decimal.NewFromInt(ri.(int64))}
		}
	} else {
		rightValue = &dsl.LeafValue{Decimal: decimal.NewFromInt(math.MaxInt64)}
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

func convertToUint64Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if li, err := bound.LeftValue.Value(convertToUInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to uint64", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{LongInt: li.(uint64)}
		}
	} else {
		leftValue = &dsl.LeafValue{LongInt: 0}
	}
	if !bound.RightValue.IsInf() {
		if ri, err := bound.RightValue.Value(convertToUInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to uint64", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{LongInt: ri.(uint64)}
		}
	} else {
		rightValue = &dsl.LeafValue{LongInt: math.MaxUint64}
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

func convertToFloat16Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if lf, err := bound.LeftValue.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to float16", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{Float16: float16.Fromfloat32(float32(lf.(float64)))}
		}
	} else {
		leftValue = &dsl.LeafValue{Float16: dsl.MinFloat16}
	}

	if !bound.RightValue.IsInf() {
		if lf, err := bound.LeftValue.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to float16", field, bound.LeftValue.String())
		} else {
			rightValue = &dsl.LeafValue{Float16: float16.Fromfloat32(float32(lf.(float64)))}
		}

	} else {
		rightValue = &dsl.LeafValue{Float16: dsl.MaxFloat16}
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

func convertToFloat32Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if lf, err := bound.LeftValue.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to float32", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat32(float32(lf.(float64)))}
		}

	} else {
		leftValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat32(-math.MaxFloat32)}
	}

	if !bound.RightValue.IsInf() {
		if rf, err := bound.RightValue.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to float32", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat32(float32(rf.(float64)))}
		}
	} else {
		rightValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat32(math.MaxFloat32)}
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

func convertToFloat64Range(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if lf, err := bound.LeftValue.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect to float64", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat(lf.(float64))}
		}
	} else {
		leftValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat(-math.MaxFloat64)}
	}
	if !bound.RightValue.IsInf() {
		if rf, err := bound.RightValue.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect to float64", field, bound.RightValue.String())
		} else {
			rightValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat(rf.(float64))}
		}
	} else {
		rightValue = &dsl.LeafValue{Decimal: decimal.NewFromFloat(math.MaxFloat64)}
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

func convertToIPRange(field *term.Field, termV *term.Term,
	leftCmp, rightCmp dsl.CompareType, property *mapping.Property) (dsl.DSLNode, error) {
	var (
		bound      = termV.GetBound()
		leftValue  *dsl.LeafValue
		rightValue *dsl.LeafValue
	)
	if !bound.LeftValue.IsInf() {
		if li, err := bound.LeftValue.Value(convertToIp); err != nil {
			return nil, fmt.Errorf("field: %s left value: %s is invalid, expect ip", field, bound.LeftValue.String())
		} else {
			leftValue = &dsl.LeafValue{IpValue: li.(net.IP)}
		}
	} else {
		leftValue = &dsl.LeafValue{IpValue: dsl.MinIP}
	}

	if !bound.RightValue.IsInf() {
		if ri, err := bound.RightValue.Value(convertToIp); err != nil {
			return nil, fmt.Errorf("field: %s right value: %s is invalid, expect ip", field, bound.RightValue.String())
		} else {
			leftValue = &dsl.LeafValue{IpValue: ri.(net.IP)}
		}
	} else {
		leftValue = &dsl.LeafValue{IpValue: dsl.MaxIP}
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
	case mapping.BOOLEAN_FIELD_TYPE:
		if b, err := termV.Value(convertToBool); err != nil {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Boolean: b.(bool),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.BYTE_FIELD_TYPE:
		if i, err := termV.Value(convertToInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int8", field, termV.String())
		} else if i64 := i.(int64); i64 < math.MinInt8 || i64 > math.MaxInt8 {
			return nil, fmt.Errorf("field: %s is byte type, expect value range [%d, %d]", field, math.MinInt8, math.MaxInt8)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						TinyInt: int16(i.(int64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.SHORT_FIELD_TYPE:
		if i, err := termV.Value(convertToInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int16", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						TinyInt: int16(i.(int64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.INTEGER_FIELD_TYPE, mapping.INTERGER_RANGE_FIELD_TYPE:
		if i, err := termV.Value(convertToInt(32)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int32", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Decimal: decimal.NewFromInt(i.(int64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.LONG_FIELD_TYPE:
		if i, err := termV.Value(convertToInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to int64", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Decimal: decimal.NewFromInt(i.(int64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.UNSIGNED_LONG_FIELD_TYPE:
		if i, err := termV.Value(convertToUInt(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to uint64", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						LongInt: i.(uint64),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}

	case mapping.HALF_FLOAT_FIELD_TYPE:
		if f, err := termV.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to float16", field, termV.String())
		} else if f64 := f.(float64); f64 < math.MinInt32 || f64 > math.MaxInt32 {
			return nil, fmt.Errorf("field: %s is float16 type, expect value range [%d, %d]", field, math.MinInt32, math.MaxInt32)
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Float16: float16.Fromfloat32(float32(f64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.FLOAT_FIELD_TYPE, mapping.FLOAT_RANGE_FIELD_TYPE:
		if f, err := termV.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to float32", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Decimal: decimal.NewFromFloat(f.(float64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
	case mapping.DOUBLE_FIELD_TYPE, mapping.DOUBLE_RANGE_FIELD_TYPE:
		if f, err := termV.Value(convertToFloat(64)); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to float64", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Decimal: decimal.NewFromFloat(f.(float64)),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}

	case mapping.SCALED_FLOAT_FIELD_TYPE:
		if f, err := decimal.NewFromString(termV.String()); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to float128", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Value: &dsl.LeafValue{
						Decimal: f,
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
					Type:  property.Type,
					Value: &dsl.LeafValue{
						IpValue: ip.(net.IP),
					},
				},
				Boost: termV.Boost(),
			}, nil
		}
		if ip1, ip2, err := ip_tools.GetRangeIpByIpCidr(termV.String()); err == nil {
			return &dsl.RangeNode{
				Field:     field.String(),
				ValueType: property.Type,
				LeftValue: &dsl.LeafValue{
					IpValue: net.IP(ip1),
				},
				RightValue: &dsl.LeafValue{
					IpValue: net.IP(ip2),
				},
				LeftCmpSym:  dsl.GTE,
				RightCmpSym: dsl.LTE,
				Boost:       termV.Boost(),
			}, nil
		}
		return nil, fmt.Errorf("ip value: %s is invalid", termV.String())
	case mapping.VERSION_FIELD_TYPE:
		if v, err := termV.Value(convertToVersion); err != nil {
			return nil, fmt.Errorf("field: %s value: %s is invalid, expect to version", field, termV.String())
		} else {
			return &dsl.TermNode{
				EqNode: dsl.EqNode{
					Field: field.String(),
					Type:  property.Type,
					Value: &dsl.LeafValue{
						Version: v.(*version.Version),
					},
				},
				Boost: termV.Boost(),
			}, nil

		}

	case mapping.KEYWORD_FIELD_TYPE, mapping.CONSTANT_KEYWORD_FIELD_TYPE:
		var s, _ = termV.Value(func(s string) (interface{}, error) { return s, nil })
		return &dsl.TermNode{
			EqNode: dsl.EqNode{
				Field: field.String(),
				Type:  property.Type,
				Value: &dsl.LeafValue{
					String: s.(string),
				},
			},
			Boost: termV.Boost(),
		}, nil

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
		Value: &dsl.LeafValue{
			String: s.(string),
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
