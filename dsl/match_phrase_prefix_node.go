package dsl

import "fmt"

// match_phrase node
type MatchPhrasePrefixNode struct {
	kvNode
	slopNode
	boostNode
	expandsNode
	analyzerNode
}

func NewMatchPhrasePrefixNode(kvNode *kvNode, opts ...func(AstNode)) *MatchPhrasePrefixNode {
	var n = &MatchPhrasePrefixNode{
		kvNode:       *kvNode,
		slopNode:     slopNode{slop: 2},
		boostNode:    boostNode{boost: 1.0},
		expandsNode:  expandsNode{maxExpands: 50},
		analyzerNode: analyzerNode{},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *MatchPhrasePrefixNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		if b, ok := o.(BoostNode); ok {
			if compareBoost(n, b) != 0 {
				return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func (n *MatchPhrasePrefixNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	default:
		if b, ok := o.(BoostNode); ok {
			if compareBoost(n, b) != 0 {
				return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	}
}

func (n *MatchPhrasePrefixNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}

func (n *MatchPhrasePrefixNode) DslType() DslType {
	return MATCH_PHRASE_PREFIX_DSL_TYPE
}

func (n *MatchPhrasePrefixNode) ToDSL() DSL {
	d := DSL{
		QUERY_KEY:          n.toPrintValue(),
		SLOP_KEY:           n.getSlop(),
		MAX_EXPANSIONS_KEY: n.getMaxExpands(),
	}
	addValueForDSL(d, ANALYZER_KEY, n.getAnaLyzer())
	return DSL{MATCH_PHRASE_PREFIX_KEY: DSL{n.field: d}}
}
