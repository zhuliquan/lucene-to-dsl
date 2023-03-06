package dsl

type MatchAllNode struct {
}

func (n *MatchAllNode) UnionJoin(x AstNode) (AstNode, error) {
	return n, nil
}

func (n *MatchAllNode) InterSect(x AstNode) (AstNode, error) {
	return x, nil
}

func (n *MatchAllNode) Inverse() (AstNode, error) {
	return &EmptyNode{}, nil
}

func (n *MatchAllNode) NodeKey() string {
	return "*"
}

func (n *MatchAllNode) DslType() DslType {
	return MATCH_ALL_DSL_TYPE
}

func (n *MatchAllNode) AstType() AstType {
	return LEAF_NODE_TYPE
}

func (n *MatchAllNode) ToDSL() DSL {
	return DSL{
		"match_all": DSL{},
	}
}
