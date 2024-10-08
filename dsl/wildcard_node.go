package dsl

import (
	"github.com/zhuliquan/lucene-to-dsl/utils"
)

type WildCardNode struct {
	kvNode
	boostNode
	rewriteNode
	patternNode
}

func NewWildCardNode(kvNode *kvNode, pattern utils.PatternMatcher, opts ...func(AstNode)) *WildCardNode {
	var n = &WildCardNode{
		kvNode:      *kvNode,
		boostNode:   boostNode{boost: 1.0},
		rewriteNode: rewriteNode{rewrite: CONSTANT_SCORE},
		patternNode: patternNode{matcher: pattern},
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
		return patternNodeUnionJoinTermNode(n, o.(*TermNode))
	case WILDCARD_DSL_TYPE:
		return valueNodeUnionJoinValueNode(n, o)
	default:
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func (n *WildCardNode) InterSect(o AstNode) (AstNode, error) {
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
		return patternNodeIntersectTermNode(n, o.(*TermNode))
	case WILDCARD_DSL_TYPE:
		return valueNodeIntersectValueNode(n, o)
	default:
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	}
}

func (n *WildCardNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}

func (n *WildCardNode) ToDSL() DSL {
	return DSL{
		WILDCARD_KEY: DSL{
			n.field: DSL{
				VALUE_KEY:   n.toPrintValue(),
				BOOST_KEY:   n.getBoost(),
				REWRITE_KEY: n.getRewrite(),
			},
		},
	}
}
