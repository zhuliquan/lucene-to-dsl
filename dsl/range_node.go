package dsl

import (
	"encoding/json"
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

type RangeNode struct {
	KvNode
	LeftValue   LeafValue
	RightValue  LeafValue
	LeftCmpSym  CompareType
	RightCmpSym CompareType
	Boost       float64
}

func (n *RangeNode) String() string {
	var lv, _ = json.Marshal(n.LeftValue)
	var rv, _ = json.Marshal(n.RightValue)
	var lb = "("
	if n.LeftCmpSym == GTE {
		lb = "["
	}
	var rb = ")"
	if n.RightCmpSym == LTE {
		rb = "]"
	}
	return fmt.Sprintf("%s%s, %s%s", lb, lv, rv, rb)
}

func (n *RangeNode) getBoost() float64 {
	return n.Boost
}

func (n *RangeNode) DslType() DslType {
	return RANGE_DSL_TYPE
}

func (n *RangeNode) NodeKey() string { return "LEAF:" + n.Field }

func (n *RangeNode) UnionJoin(o AstNode) (AstNode, error) {
	if bn, ok := o.(boostNode); ok {
		if CompareAny(bn.getBoost(), n.Boost, n.Type) != 0 {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost value isn't equal", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		return rangeNodeUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return rangeNodeUnionJoinTermsNode(n, o.(*TermsNode))
	case RANGE_DSL_TYPE:
		return rangeNodeUnionJoinRangeNode(n, o.(*RangeNode))
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *RangeNode) InterSect(o AstNode) (AstNode, error) {
	if bn, ok := o.(boostNode); ok {
		if CompareAny(bn.getBoost(), n.Boost, n.Type) != 0 {
			return nil, fmt.Errorf("failed to intersect %s and %s, err: boost value isn't equal", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	case TERM_DSL_TYPE:
		return rangeNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return rangeNodeIntersectTermsNode(n, o.(*TermsNode))
	case RANGE_DSL_TYPE:
		return rangeNodeIntersectRangeNode(n, o.(*RangeNode))
	default:
		return &AndNode{
			MustNodes: map[string][]AstNode{
				n.NodeKey(): {n, o},
			},
		}, nil
	}
}

func (n *RangeNode) Inverse() (AstNode, error) {
	var (
		leftCmpSym  = LT
		rightCmpSym = GT
	)
	if n.LeftCmpSym == GT {
		leftCmpSym = LTE
	}
	if n.RightCmpSym == LT {
		rightCmpSym = GTE
	}
	var (
		isLeftInf  = isMinInf(n.LeftValue, n.Type)
		isRightInf = isMaxInf(n.RightValue, n.Type)
		leftNode   = &RangeNode{
			KvNode:      n.KvNode,
			LeftValue:   minInf[n.Type],
			LeftCmpSym:  GT,
			RightValue:  n.LeftValue,
			RightCmpSym: leftCmpSym,
		}
		rightNode = &RangeNode{
			KvNode:      n.KvNode,
			LeftValue:   n.RightValue,
			RightValue:  maxInf[n.Type],
			LeftCmpSym:  rightCmpSym,
			RightCmpSym: LT,
		}
	)

	if !isLeftInf && !isRightInf {
		return &OrNode{
			MinimumShouldMatch: 1,
			Nodes: map[string][]AstNode{
				n.NodeKey(): {leftNode, rightNode},
			},
		}, nil
	} else if !isLeftInf {
		return leftNode, nil
	} else if !isRightInf {
		return rightNode, nil
	} else {
		return &NotNode{
			Nodes: map[string][]AstNode{
				n.NodeKey(): {
					&ExistsNode{
						KvNode: n.KvNode,
					},
				},
			},
		}, nil
	}
}

func (n *RangeNode) ToDSL() DSL {
	var res = DSL{}
	res[n.LeftCmpSym.String()] = leafValueToPrintValue(n.LeftValue, n.Type)
	res[n.RightCmpSym.String()] = leafValueToPrintValue(n.RightValue, n.Type)
	res["relation"] = "WITHIN"
	res["boost"] = n.Boost
	if mapping.CheckDateType(n.Type) {
		res["format"] = "epoch_millis"
	}
	return DSL{"range": DSL{n.Field: res}}
}

func rangeNodeUnionJoinTermNode(n *RangeNode, t *TermNode) (AstNode, error) {
	if !checkRangeInclude(n, t.Value) {
		if CompareAny(n.LeftValue, t.Value, n.Type) == 0 && n.LeftCmpSym == GT {
			return &RangeNode{
				KvNode:      n.KvNode,
				LeftValue:   n.LeftValue,
				RightValue:  n.RightValue,
				LeftCmpSym:  GTE,
				RightCmpSym: n.RightCmpSym,
				Boost:       n.Boost,
			}, nil
		}
		if CompareAny(n.RightValue, t.Value, n.Type) == 0 && n.RightCmpSym == LT {
			return &RangeNode{
				KvNode:      n.KvNode,
				LeftValue:   n.LeftValue,
				RightValue:  n.RightValue,
				LeftCmpSym:  n.LeftCmpSym,
				RightCmpSym: LTE,
				Boost:       n.Boost,
			}, nil
		}
		return &OrNode{
			MinimumShouldMatch: 1,
			Nodes: map[string][]AstNode{
				n.NodeKey(): {n, t},
			},
		}, nil
	} else {
		return n, nil
	}
}

func rangeNodeUnionJoinTermsNode(n *RangeNode, t *TermsNode) (AstNode, error) {
	var (
		excludeValues = []LeafValue{}
		leftCmpSym    = n.LeftCmpSym
		rightCmpSym   = n.RightCmpSym
	)

	for _, value := range t.Values {
		if !checkRangeInclude(n, value) {
			if CompareAny(n.LeftValue, value, n.Type) == 0 && n.LeftCmpSym == GT {
				leftCmpSym = GTE
			} else if CompareAny(n.RightValue, value, n.Type) == 0 && n.RightCmpSym == LT {
				rightCmpSym = LTE
			} else {
				excludeValues = append(excludeValues, value)
			}
		}
	}
	var rangeNode = &RangeNode{
		KvNode:      n.KvNode,
		LeftValue:   n.LeftValue,
		RightValue:  n.RightValue,
		LeftCmpSym:  leftCmpSym,
		RightCmpSym: rightCmpSym,
		Boost:       n.Boost,
	}
	if len(excludeValues) == 0 {
		return rangeNode, nil
	} else {
		return &OrNode{
			MinimumShouldMatch: 1,
			Nodes: map[string][]AstNode{
				n.NodeKey(): {
					rangeNode,
					&TermsNode{
						KvNode: t.KvNode,
						Values: excludeValues,
						Boost:  t.Boost,
					},
				},
			},
		}, nil

	}
}

func rangeNodeUnionJoinRangeNode(n, t *RangeNode) (AstNode, error) {
	// first check overlap, if no overlap, return or ast node
	if !checkRangeOverlap(n, t) {
		return &OrNode{
			MinimumShouldMatch: 1,
			Nodes: map[string][]AstNode{
				n.NodeKey(): {n, t},
			},
		}, nil
	}
	// compare left value of n and t, and get lower value, and cmp symbol is associate with lower value
	// compare left value of n and t, and get higher value, and cmp symbol is associate with higher value
	var dst = &RangeNode{
		KvNode:    n.KvNode,
		LeftValue: n.LeftValue,
		Boost:     n.Boost,
	}

	unionCmpLeft(n, t, dst)
	unionCmpRight(n, t, dst)
	return dst, nil
}

func unionCmpLeft(n, t, dst *RangeNode) {
	var leftFlag = CompareAny(t.LeftValue, n.LeftValue, n.Type)
	if leftFlag < 0 {
		dst.LeftValue = t.LeftValue
		dst.LeftCmpSym = t.LeftCmpSym
	} else if leftFlag > 0 {
		dst.LeftValue = n.LeftValue
		dst.LeftCmpSym = n.LeftCmpSym
	} else {
		dst.LeftValue = n.LeftValue
		if t.LeftCmpSym == GTE {
			dst.LeftCmpSym = t.LeftCmpSym
		} else {
			dst.LeftCmpSym = n.LeftCmpSym
		}
	}
}

func unionCmpRight(n, t, dst *RangeNode) {
	var rightFlag = CompareAny(t.RightValue, n.RightValue, n.Type)
	if rightFlag > 0 {
		dst.RightValue = t.RightValue
		dst.RightCmpSym = t.RightCmpSym
	} else if rightFlag < 0 {
		dst.RightValue = n.RightValue
		dst.RightCmpSym = n.RightCmpSym
	} else {
		dst.RightValue = n.RightValue
		if t.RightCmpSym == LTE {
			dst.RightCmpSym = t.RightCmpSym
		} else {
			dst.RightCmpSym = n.RightCmpSym
		}
	}
}

func rangeNodeIntersectTermNode(n *RangeNode, t *TermNode) (AstNode, error) {
	if checkRangeInclude(n, t.Value) {
		return t, nil
	} else {
		return nil, fmt.Errorf("range node: %s can't intersect with term node: %s, range doesn't include term value", n.ToDSL(), t.ToDSL())
	}
}

func rangeNodeIntersectTermsNode(n *RangeNode, t *TermsNode) (AstNode, error) {
	var includeValues = []LeafValue{}
	for _, value := range t.Values {
		if checkRangeInclude(n, value) {
			includeValues = append(includeValues, value)
		}
	}
	if len(includeValues) == 0 {
		return nil, fmt.Errorf("failed to intersect %s and %s, err: range doesn't include any term value", n.ToDSL(), t.ToDSL())
	}
	return &TermsNode{
		KvNode: t.KvNode,
		Values: includeValues,
		Boost:  t.Boost,
	}, nil
}

func rangeNodeIntersectRangeNode(n, t *RangeNode) (AstNode, error) {
	// first check have range overlap zone
	if !checkRangeOverlap(n, t) {
		return nil, fmt.Errorf("range node: %s can't intersect with range node: %s, no overlap between two range", n.ToDSL(), t.ToDSL())
	}
	// compare left value of n and t, and get higher value, and cmp symbol is associate with higher value
	// compare left value of n and t, and get lower value, and cmp symbol is associate with lower value
	var dst = &RangeNode{
		KvNode: n.KvNode,
		Boost:  n.Boost,
	}
	intersectCmpLeft(t, n, dst)
	intersectCmpRight(t, n, dst)
	return dst, nil

}

func intersectCmpLeft(n, t, dst *RangeNode) {
	var leftFlag = CompareAny(t.LeftValue, n.LeftValue, n.Type)
	if leftFlag > 0 {
		dst.LeftValue = t.LeftValue
		dst.LeftCmpSym = t.LeftCmpSym
	} else if leftFlag < 0 {
		dst.LeftValue = n.LeftValue
		dst.LeftCmpSym = n.LeftCmpSym
	} else {
		dst.LeftValue = n.LeftValue
		if t.LeftCmpSym == GT {
			dst.LeftCmpSym = t.LeftCmpSym
		} else {
			dst.LeftCmpSym = n.LeftCmpSym
		}
	}
}

func intersectCmpRight(n, t, dst *RangeNode) {
	var rightFlag = CompareAny(t.RightValue, n.RightValue, n.Type)
	if rightFlag < 0 {
		dst.RightValue = t.RightValue
		dst.RightCmpSym = t.RightCmpSym
	} else if rightFlag > 0 {
		dst.RightValue = n.RightValue
		dst.RightCmpSym = n.RightCmpSym
	} else {
		dst.RightValue = n.RightValue
		if t.RightCmpSym == LT {
			dst.RightCmpSym = t.RightCmpSym
		} else {
			dst.RightCmpSym = n.RightCmpSym
		}
	}
}

// check range node include a value
func checkRangeInclude(n *RangeNode, v LeafValue) bool {
	var leftCmpRes = CompareAny(v, n.LeftValue, n.Type)
	var rightCmpRes = CompareAny(v, n.RightValue, n.Type)
	return (leftCmpRes > 0 && rightCmpRes < 0) ||
		(leftCmpRes == 0 && n.LeftCmpSym == GTE) ||
		(rightCmpRes == 0 && n.RightCmpSym == LTE)
}

// check two range overlap
func checkRangeOverlap(a, b *RangeNode) bool {
	var cmpRes1 = CompareAny(a.RightValue, b.LeftValue, a.Type)
	var cmpRes2 = CompareAny(b.RightValue, a.LeftValue, a.Type)
	// two range don't have overlap zone is easy, there are two case:
	// 1. max value of left range is lower than min value of right range
	// 2. two range is exclude type range and max value of left range is equal with min value of right range
	// inverse two case you can check overlap sense.
	return !(cmpRes1 < 0 || cmpRes2 < 0 ||
		(cmpRes1 == 0 && (a.RightCmpSym == LT || b.LeftCmpSym == GT)) ||
		(cmpRes2 == 0 && (b.RightCmpSym == LT || a.LeftCmpSym == GT)))
}
