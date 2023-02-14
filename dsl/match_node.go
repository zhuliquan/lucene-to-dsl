package dsl

import "fmt"

// match node
type MatchNode struct {
	kvNode
	boostNode
	expandsNode
	analyzerNode
}

func NewMatchNode(kvNode *kvNode, opts ...func(AstNode)) *MatchNode {
	var n = &MatchNode{
		kvNode:      *kvNode,
		boostNode:   boostNode{boost: 1.0},
		expandsNode: expandsNode{maxExpands: 50},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *MatchNode) DslType() DslType {
	return MATCH_DSL_TYPE
}

func (n *MatchNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		if b, ok := o.(BoostNode); ok {
			if compareBoost(n, b) != 0 {
				return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *MatchNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	default:
		if b, ok := o.(BoostNode); ok {
			if compareBoost(n, b) != 0 {
				return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *MatchNode) Inverse() (AstNode, error) {
	return NewBoolNode(n, NOT), nil
}

func (n *MatchNode) ToDSL() DSL {
	d := DSL{
		QUERY_KEY:          n.toPrintValue(),
		BOOST_KEY:          n.getBoost(),
		MAX_EXPANSIONS_KEY: n.getMaxExpands(),
	}
	if n.getAnaLyzer() != "" {
		d[ANALYZER_KEY] = n.getAnaLyzer()
	}
	return DSL{MATCH_KEY: DSL{n.field: d}}
}
