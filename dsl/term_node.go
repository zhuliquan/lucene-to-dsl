package dsl

import "fmt"

type TermNode struct {
	kvNode
	boostNode
}

func NewTermNode(kvNode *kvNode, opts ...func(AstNode)) *TermNode {
	var n = &TermNode{
		kvNode:    *kvNode,
		boostNode: boostNode{boost: 1.0},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *TermNode) DslType() DslType {
	return TERM_DSL_TYPE
}

func (n *TermNode) UnionJoin(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.UnionJoin(n)
	}
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
			return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
		}
	}
	switch o.DslType() {
	case TERM_DSL_TYPE:
		return termNodeUnionJoinTermNode(n, o.(*TermNode))
	case RANGE_DSL_TYPE, PREFIX_DSL_TYPE, REGEXP_DSL_TYPE, WILDCARD_DSL_TYPE, IDS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func (n *TermNode) InterSect(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.InterSect(n)
	}
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
			return lfNodeIntersectLfNode(n.NodeKey(), n, o)
		}
	}
	switch o.DslType() {
	case TERM_DSL_TYPE:
		return termNodeIntersectTermNode(n, o.(*TermNode))
	case RANGE_DSL_TYPE, PREFIX_DSL_TYPE, REGEXP_DSL_TYPE, WILDCARD_DSL_TYPE, IDS_DSL_TYPE:
		return o.InterSect(n)
	default:
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	}
}

func (n *TermNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}

func termNodeUnionJoinTermNode(n, o *TermNode) (AstNode, error) {
	if CompareAny(o.value, n.value, n.mType) == 0 {
		return o, nil
	} else {
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func termNodeIntersectTermNode(n, o *TermNode) (AstNode, error) {
	if CompareAny(o.value, n.value, n.mType) == 0 {
		return o, nil
	} else if n.IsArrayType() {
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermNode) ToDSL() DSL {
	return DSL{
		TERM_KEY: DSL{
			n.field: DSL{
				VALUE_KEY: n.toPrintValue(),
				BOOST_KEY: n.getBoost(),
			},
		},
	}
}
