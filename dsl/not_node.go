package dsl

type NotNode struct {
	opNode
	Nodes map[string][]AstNode
}

func (n *NotNode) DslType() DslType {
	return NOT_DSL_TYPE
}

func (n *NotNode) NodeKey() string {
	return NOT_OP_KEY
}

func (n *NotNode) ToDSL() DSL {
	if nodes := flattenNodes(n.Nodes); nodes != nil {
		return DSL{
			BOOL_KEY: DSL{
				MUST_NOT_KEY: nodes,
			},
		}
	} else {
		return EmptyDSL
	}
}

func (n *NotNode) UnionJoin(o AstNode) (AstNode, error) {
	// var t = o.(*OrNode)
	// var nNodes = n.Nodes
	// var tNodes = t.Nodes

	return nil, nil
}

func (n *NotNode) InterSect(o AstNode) (AstNode, error) {
	if o == nil || n == nil {
		return nil, ErrIntersectNilNode
	}
	var t = o.(*NotNode)
	for key, curNodes := range t.Nodes {
		if preNodes, ok := n.Nodes[key]; ok {
			if key == OR_OP_KEY {
				n.Nodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.Nodes, key)
					n.Nodes[key] = []AstNode{newNode}
				}
			}
		} else {
			n.Nodes[key] = curNodes
		}
	}
	return n, nil
}

// 全部都不是的反例是至少有一个:
func (n *NotNode) Inverse() (AstNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return &OrNode{
		Nodes:              n.Nodes,
		MinimumShouldMatch: 1,
	}, nil
}
