package dsl

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
		analyzerNode: analyzerNode{},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *MatchPhrasePrefixNode) UnionJoin(AstNode) (AstNode, error) {
	return nil, nil
}

func (n *MatchPhrasePrefixNode) InterSect(AstNode) (AstNode, error) {
	return nil, nil
}

func (n *MatchPhrasePrefixNode) Inverse() (AstNode, error) {
	return nil, nil
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
	if n.getAnaLyzer() != "" {
		d[ANALYZER_KEY] = n.getAnaLyzer()
	}
	return DSL{
		MATCH_PHRASE_PREFIX_KEY: d,
	}
}
