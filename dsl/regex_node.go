package dsl

type RegexpNode struct {
	KvNode
}

func (n *RegexpNode) DslType() DslType {
	return REGEXP_DSL_TYPE
}

func (n *RegexpNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *RegexpNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	default:
		return &AndNode{
			MustNodes: map[string][]AstNode{
				n.NodeKey(): {n, o},
			},
		}, nil
	}
}

func (n *RegexpNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
}

func (n *RegexpNode) ToDSL() DSL {
	return DSL{"regexp": DSL{n.Field: DSL{"value": n.Value}}}
}
