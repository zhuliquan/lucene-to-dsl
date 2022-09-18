package dsl

// match node
type MatchNode struct {
	KvNode
	Boost float64
}

func (n *MatchNode) getBoost() float64 {
	return n.Boost
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
	if n == nil {
		return EmptyDSL
	}
	return DSL{"match": DSL{n.Field: DSL{"value": n.Value, "boost": n.Boost}}}
}
