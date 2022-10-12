package dsl

// match node
type MatchNode struct {
	kvNode
	boostNode
	analyzerNode
}

func NewMatchNode(kvNode *kvNode, opts ...func(AstNode)) *MatchNode {
	var n = &MatchNode{
		kvNode:    *kvNode,
		boostNode: boostNode{boost: 1.0},
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
	return DSL{
		"match": DSL{
			n.field: DSL{
				"query": n.toPrintValue(),
				"boost": n.getBoost(),
			},
		},
	}
}
