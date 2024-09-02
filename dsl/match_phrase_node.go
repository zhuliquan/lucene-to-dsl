package dsl

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
		QUERY_KEY: n.getValue(),
		BOOST_KEY: n.getBoost(),
	}
	addValueForDSL(d, ANALYZER_KEY, n.getAnaLyzer())
	return DSL{MATCH_PHRASE_KEY: DSL{n.field: d}}
}

func (n *MatchPhraseNode) UnionJoin(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.UnionJoin(n)
	}
	switch o.DslType() {
	default:
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func (n *MatchPhraseNode) InterSect(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.InterSect(n)
	}
	switch o.DslType() {
	default:
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	}
}

func (n *MatchPhraseNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}
