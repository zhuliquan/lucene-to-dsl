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
	return DSL{
		"match_phrase": DSL{
			n.field: n.toPrintValue(),
		},
	}
}

func (n *MatchPhraseNode) UnionJoin(AstNode) (AstNode, error) {
	return nil, nil
}

func (n *MatchPhraseNode) InterSect(AstNode) (AstNode, error) {
	return nil, nil
}

func (n *MatchPhraseNode) Inverse() (AstNode, error) {
	return nil, nil
}
