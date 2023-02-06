package dsl

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
	return n, nil
}

func (n *MatchNode) InterSect(o AstNode) (AstNode, error) {
	return n, nil
}

func (n *MatchNode) Inverse() (AstNode, error) {
	return n, nil
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
