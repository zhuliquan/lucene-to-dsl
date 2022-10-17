package dsl

import "fmt"

type WildCardNode struct {
	kvNode
	boostNode
	rewriteNode
	patternNode
}

type wildcardPattern struct {
	pattern []rune
}

func NewWildCardPattern(pattern string) PatternMatcher {
	return &wildcardPattern{pattern: []rune(pattern)}
}

func (w *wildcardPattern) Match(text []byte) bool {
	return wildcardMatch([]rune(string(text)), w.pattern)
}

func NewWildCardNode(kvNode *kvNode, pattern PatternMatcher, opts ...func(AstNode)) *WildCardNode {
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
	switch o.DslType() {
	case TERM_DSL_TYPE:
		return patternNodeUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return patternNodeUnionJoinTermsNode(n, o.(*TermsNode))
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to union join wildcard node")
	}
}

func (n *WildCardNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case TERM_DSL_TYPE:
		return patternNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return patternNodeIntersectTermsNode(n, o.(*TermsNode))
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to intersect wildcard node")
	}
}

func (n *WildCardNode) Inverse() (AstNode, error) {
	return &NotNode{
		opNode: opNode{filterCtxNode: n.filterCtxNode},
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
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
