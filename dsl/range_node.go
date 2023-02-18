package dsl

import (
	"encoding/json"
	"fmt"

	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

type RangeNode struct {
	rgNode
	boostNode
	// format   string // 如果是 date类型直接变为epoch_millis
	relation RelationType
	timeZone string
}

func WithRelation(relation RelationType) func(AstNode) {
	return func(n AstNode) {
		if f, ok := n.(*RangeNode); ok {
			f.relation = relation
		}
	}
}

func WithTimeZone(timeZone string) func(AstNode) {
	return func(n AstNode) {
		if f, ok := n.(*RangeNode); ok {
			f.timeZone = timeZone
		}
	}
}

func NewRangeNode(RgNode *rgNode, opts ...func(AstNode)) *RangeNode {
	var n = &RangeNode{
		rgNode:    *RgNode,
		boostNode: boostNode{boost: 1.0},
		relation:  INTERSECTS,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *RangeNode) String() string {
	var lv, _ = json.Marshal(leafValueToPrintValue(n.lValue, n.mType))
	var rv, _ = json.Marshal(leafValueToPrintValue(n.rValue, n.mType))
	var lb = "("
	if n.lCmpSym == GTE {
		lb = "["
	}
	var rb = ")"
	if n.rCmpSym == LTE {
		rb = "]"
	}
	return fmt.Sprintf("%s%s, %s%s", lb, lv, rv, rb)
}

func (n *RangeNode) DslType() DslType {
	return RANGE_DSL_TYPE
}

func (n *RangeNode) UnionJoin(o AstNode) (AstNode, error) {
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost value isn't equal", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		return rangeNodeUnionJoinTermNode(n, o.(*TermNode))
	case RANGE_DSL_TYPE:
		return rangeNodeUnionJoinRangeNode(n, o.(*RangeNode))
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *RangeNode) InterSect(o AstNode) (AstNode, error) {
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
			return nil, fmt.Errorf("failed to intersect %s and %s, err: boost value isn't equal", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	case TERM_DSL_TYPE:
		return rangeNodeIntersectTermNode(n, o.(*TermNode))
	case RANGE_DSL_TYPE:
		return rangeNodeIntersectRangeNode(n, o.(*RangeNode))
	default:
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *RangeNode) Inverse() (AstNode, error) {
	var (
		lCmpSym = LT
		rCmpSym = GT
	)
	if n.lCmpSym == GT {
		lCmpSym = LTE
	}
	if n.rCmpSym == LT {
		rCmpSym = GTE
	}
	var (
		isLeftInf  = isMinInf(n.lValue, n.mType)
		isRightInf = isMaxInf(n.rValue, n.mType)
		leftNode   = &RangeNode{
			rgNode: rgNode{
				fieldNode: n.fieldNode,
				valueType: n.valueType,
				lValue:    minInf[n.mType],
				rValue:    n.lValue,
				lCmpSym:   GT,
				rCmpSym:   lCmpSym,
			},
			timeZone:  n.timeZone,
			relation:  n.relation,
			boostNode: n.boostNode,
		}
		rightNode = &RangeNode{
			rgNode: rgNode{
				fieldNode: n.fieldNode,
				valueType: n.valueType,
				lValue:    n.rValue,
				rValue:    maxInf[n.mType],
				lCmpSym:   rCmpSym,
				rCmpSym:   LT,
			},
			timeZone:  n.timeZone,
			relation:  n.relation,
			boostNode: n.boostNode,
		}
	)

	if !isLeftInf && !isRightInf {
		return &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				n.NodeKey(): {leftNode, rightNode},
			},
			MinimumShouldMatch: 1,
		}, nil
	} else if !isLeftInf {
		return leftNode, nil
	} else if !isRightInf {
		return rightNode, nil
	} else {
		return NewBoolNode(NewExistsNode(&n.fieldNode), NOT), nil
	}
}

func (n *RangeNode) ToDSL() DSL {
	var res = DSL{
		BOOST_KEY:          n.getBoost(),
		RELATION_KEY:       n.relation,
		n.lCmpSym.String(): leafValueToPrintValue(n.lValue, n.mType),
		n.rCmpSym.String(): leafValueToPrintValue(n.rValue, n.mType),
	}
	if mapping.CheckDateType(n.mType) {
		addValueForDSL(res, FORMAT_KEY, datemath_parser.EPOCH_MILLIS)
	}
	addValueForDSL(res, TIME_ZONE_KEY, n.timeZone)
	return DSL{RANGE_KEY: DSL{n.field: res}}
}

func rangeNodeUnionJoinTermNode(n *RangeNode, t *TermNode) (AstNode, error) {
	if !checkRangeInclude(n, t.value) {
		if CompareAny(n.lValue, t.value, n.mType) == 0 && n.lCmpSym == GT {
			return &RangeNode{
				rgNode: rgNode{
					fieldNode: n.fieldNode,
					valueType: n.valueType,
					lValue:    n.lValue,
					rValue:    n.rValue,
					lCmpSym:   GTE,
					rCmpSym:   n.rCmpSym,
				},
				timeZone:  n.timeZone,
				relation:  n.relation,
				boostNode: n.boostNode,
			}, nil
		}
		if CompareAny(n.rValue, t.value, n.mType) == 0 && n.rCmpSym == LT {
			return &RangeNode{
				rgNode: rgNode{
					fieldNode: n.fieldNode,
					valueType: n.valueType,
					lValue:    n.lValue,
					rValue:    n.rValue,
					lCmpSym:   n.lCmpSym,
					rCmpSym:   LTE,
				},
				timeZone:  n.timeZone,
				relation:  n.relation,
				boostNode: n.boostNode,
			}, nil
		}
		return lfNodeUnionJoinLfNode(n, t)
	} else {
		return n, nil
	}
}

func rangeNodeUnionJoinRangeNode(n, t *RangeNode) (AstNode, error) {
	// first check overlap, if no overlap, return or ast node
	if !checkRangeOverlap(n, t) {
		return &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				n.NodeKey(): {n, t},
			},
			MinimumShouldMatch: 1,
		}, nil
	}
	// compare left value of n and t, and get lower value, and cmp symbol is associate with lower value
	// compare left value of n and t, and get higher value, and cmp symbol is associate with higher value
	var dst = &RangeNode{
		rgNode: rgNode{
			fieldNode: n.fieldNode,
			valueType: n.valueType,
		},
		timeZone:  n.timeZone,
		relation:  n.relation,
		boostNode: n.boostNode,
	}

	unionCmpLeft(n, t, dst)
	unionCmpRight(n, t, dst)
	return dst, nil
}

func unionCmpLeft(n, t, dst *RangeNode) {
	var leftFlag = CompareAny(t.lValue, n.lValue, n.mType)
	if leftFlag < 0 {
		dst.lValue = t.lValue
		dst.lCmpSym = t.lCmpSym
	} else if leftFlag > 0 {
		dst.lValue = n.lValue
		dst.lCmpSym = n.lCmpSym
	} else {
		dst.lValue = n.lValue
		if t.lCmpSym == GTE {
			dst.lCmpSym = t.lCmpSym
		} else {
			dst.lCmpSym = n.lCmpSym
		}
	}
}

func unionCmpRight(n, t, dst *RangeNode) {
	var rightFlag = CompareAny(t.rValue, n.rValue, n.mType)
	if rightFlag > 0 {
		dst.rValue = t.rValue
		dst.rCmpSym = t.rCmpSym
	} else if rightFlag < 0 {
		dst.rValue = n.rValue
		dst.rCmpSym = n.rCmpSym
	} else {
		dst.rValue = n.rValue
		if t.rCmpSym == LTE {
			dst.rCmpSym = t.rCmpSym
		} else {
			dst.rCmpSym = n.rCmpSym
		}
	}
}

func rangeNodeIntersectTermNode(n *RangeNode, t *TermNode) (AstNode, error) {
	if checkRangeInclude(n, t.value) {
		return t, nil
	} else if n.isArrayType() {
		return lfNodeIntersectLfNode(n, t)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), t.ToDSL())
	}
}

func rangeNodeIntersectRangeNode(n, t *RangeNode) (AstNode, error) {
	// first check have range overlap zone
	if !checkRangeOverlap(n, t) {
		if n.isArrayType() {
			return lfNodeIntersectLfNode(n, t)
		} else {
			return nil, fmt.Errorf("range node: %s can't intersect with range node: %s, no overlap between two range", n.ToDSL(), t.ToDSL())
		}
	}
	// compare left value of n and t, and get higher value, and cmp symbol is associate with higher value
	// compare left value of n and t, and get lower value, and cmp symbol is associate with lower value
	var dst = &RangeNode{
		rgNode: rgNode{
			fieldNode: n.fieldNode,
			valueType: n.valueType,
		},
		timeZone:  n.timeZone,
		relation:  n.relation,
		boostNode: n.boostNode,
	}
	intersectCmpLeft(t, n, dst)
	intersectCmpRight(t, n, dst)
	return dst, nil

}

func intersectCmpLeft(n, t, dst *RangeNode) {
	var leftFlag = CompareAny(t.lValue, n.lValue, n.mType)
	if leftFlag > 0 {
		dst.lValue = t.lValue
		dst.lCmpSym = t.lCmpSym
	} else if leftFlag < 0 {
		dst.lValue = n.lValue
		dst.lCmpSym = n.lCmpSym
	} else {
		dst.lValue = n.lValue
		if t.lCmpSym == GT {
			dst.lCmpSym = t.lCmpSym
		} else {
			dst.lCmpSym = n.lCmpSym
		}
	}
}

func intersectCmpRight(n, t, dst *RangeNode) {
	var rightFlag = CompareAny(t.rValue, n.rValue, n.mType)
	if rightFlag < 0 {
		dst.rValue = t.rValue
		dst.rCmpSym = t.rCmpSym
	} else if rightFlag > 0 {
		dst.rValue = n.rValue
		dst.rCmpSym = n.rCmpSym
	} else {
		dst.rValue = n.rValue
		if t.rCmpSym == LT {
			dst.rCmpSym = t.rCmpSym
		} else {
			dst.rCmpSym = n.rCmpSym
		}
	}
}

// check range node include a value
func checkRangeInclude(n *RangeNode, v LeafValue) bool {
	var leftCmpRes = CompareAny(v, n.lValue, n.mType)
	var rightCmpRes = CompareAny(v, n.rValue, n.mType)
	return (leftCmpRes > 0 && rightCmpRes < 0) ||
		(leftCmpRes == 0 && n.lCmpSym == GTE) ||
		(rightCmpRes == 0 && n.rCmpSym == LTE)
}

// check two range overlap
func checkRangeOverlap(a, b *RangeNode) bool {
	var cmpRes1 = CompareAny(a.rValue, b.lValue, a.mType)
	var cmpRes2 = CompareAny(b.rValue, a.lValue, a.mType)
	// two range don't have overlap zone is easy, there are two case:
	// 1. max value of left range is lower than min value of right range
	// 2. two range is exclude type range and max value of left range is equal with min value of right range
	// inverse two case you can check overlap sense.
	return !(cmpRes1 < 0 || cmpRes2 < 0 ||
		(cmpRes1 == 0 && (a.rCmpSym == LT || b.lCmpSym == GT)) ||
		(cmpRes2 == 0 && (b.rCmpSym == LT || a.lCmpSym == GT)))
}
