package dsl

import "fmt"

type WildCardNode struct {
	kvNode
	boostNode
	rewriteNode
	statesNode
	flags string
}

func NewWildCardNode(kvNode *kvNode, opts ...func(AstNode)) *WildCardNode {
	var n = &WildCardNode{
		kvNode:      *kvNode,
		boostNode:   boostNode{boost: 1.0},
		rewriteNode: rewriteNode{rewrite: CONSTANT_SCORE},
		statesNode:  statesNode{maxDeterminizedStates: 10000},
		flags:       ALL_FLAG,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *WildCardNode) DslType() DslType {
	return WILDCARD_DSL_TYPE
}

func (n *WildCardNode) UnionJoin(o AstNode) (AstNode, error) {
	if n == nil && o == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && o != nil {
		return o, nil
	} else if n != nil && o == nil {
		return n, nil
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to union join wildcard node")
	}
}

func (n *WildCardNode) InterSect(o AstNode) (AstNode, error) {
	if n == nil || o == nil {
		return nil, ErrIntersectNilNode
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to intersect wildcard node")
	}
}

func (n *WildCardNode) Inverse() (AstNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("failed to inverse wildcard node")
}

func (n *WildCardNode) ToDSL() DSL {
	return DSL{
		WILDCARD_KEY: DSL{
			n.field: DSL{
				VALUE_KEY:   n.toPrintValue(),
				BOOST_KEY:   n.getBoost(),
				FLAGS_KEY:   n.flags,
				REWRITE_KEY: n.getRewrite(),

				MAX_DETERMINIZED_STATES_KEY: n.getMaxDeterminizedStates(),
			},
		},
	}

}
