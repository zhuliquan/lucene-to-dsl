package dsl

type PrefixNode struct {
	KvNode
}

func (n *PrefixNode) DslType() DslType {
	return PREFIX_DSL_TYPE
}

func (n *PrefixNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return &OrNode{
			Nodes: map[string][]AstNode{
				n.NodeKey(): {n, o},
			},
		}, nil
	}
}

func (n *PrefixNode) InterSect(o AstNode) (AstNode, error) {
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

func (n *PrefixNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
}

func (n *PrefixNode) ToDSL() DSL {
	return DSL{"prefix": DSL{n.Field: DSL{"values": n.Value}}}
}
