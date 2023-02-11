package dsl

import "fmt"

// match_phrase node
type MatchPhraseNode struct {
	kvNode
	analyzerNode
	boostNode
}

func NewMatchPhraseNode(kvNode *kvNode, opts ...func(AstNode)) *MatchPhraseNode {
	var n = &MatchPhraseNode{
		kvNode:       *kvNode,
		analyzerNode: analyzerNode{},
		boostNode:    boostNode{boost: 1.0},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *MatchPhraseNode) DslType() DslType {
	return MATCH_PHRASE_DSL_TYPE
}

func (n *MatchPhraseNode) ToDSL() DSL {
	d := DSL{
		n.field: n.toPrintValue(),
	}
	if n.getAnaLyzer() != "" {
		d[ANALYZER_KEY] = n.getAnaLyzer()
	}
	return DSL{MATCH_PHRASE_KEY: d}
}

func (n *MatchPhraseNode) UnionJoin(o AstNode) (AstNode, error) {
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

func (n *MatchPhraseNode) InterSect(o AstNode) (AstNode, error) {
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

func (n *MatchPhraseNode) Inverse() (AstNode, error) {
	return &NotNode{
		opNode: opNode{filterCtxNode: n.filterCtxNode},
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
}
