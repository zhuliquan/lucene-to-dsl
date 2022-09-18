package dsl

// match_phrase node
type MatchPhraseNode struct {
	KvNode
	Boost float64
}

func (n *MatchPhraseNode) getBoost() float64 {
	return n.Boost
}

func (n *MatchPhraseNode) DslType() DslType {
	return MATCH_PHRASE_DSL_TYPE
}

func (n *MatchPhraseNode) ToDSL() DSL {
	return DSL{"match_phrase": DSL{n.Field: n.Value}}
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
