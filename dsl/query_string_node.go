package dsl

// query_string node
type QueryStringNode struct {
	kvNode
	boostNode
	rewriteNode
	analyzerNode
}

func NewQueryStringNode(kvNode *kvNode, opts ...func(AstNode)) *QueryStringNode {
	var n = &QueryStringNode{
		kvNode:       *kvNode,
		boostNode:    boostNode{boost: 1.0},
		rewriteNode:  rewriteNode{rewrite: CONSTANT_SCORE},
		analyzerNode: analyzerNode{},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *QueryStringNode) UnionJoin(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.UnionJoin(n)
	}
	switch o.DslType() {
	default:
		return lfNodeUnionJoinLfNode(n.NodeKey(), n, o)
	}
}

func (n *QueryStringNode) InterSect(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.InterSect(n)
	}
	switch o.DslType() {
	default:
		return lfNodeIntersectLfNode(n.NodeKey(), n, o)
	}
}

func (n *QueryStringNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}

func (n *QueryStringNode) DslType() DslType {
	return QUERY_STRING_DSL_TYPE
}

func (n *QueryStringNode) ToDSL() DSL {
	d := DSL{
		QUERY_KEY:         n.toPrintValue(),
		BOOST_KEY:         n.getBoost(),
		REWRITE_KEY:       n.getRewrite(),
		DEFAULT_FIELD_KEY: n.field,
	}
	addValueForDSL(d, ANALYZER_KEY, n.getAnaLyzer())
	return DSL{QUERY_STRING_KEY: d}
}
