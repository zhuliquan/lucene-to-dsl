package dsl

type NotNode struct {
	OpNode
	Nodes map[string][]AstNode
}

func (n *NotNode) DslType() DslType {
	return NOT_DSL_TYPE
}

func (n *NotNode) NodeKey() string {
	return NOT_OP_KEY
}

func (n *NotNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	var res = []DSL{}
	for _, nodes := range n.Nodes {
		for _, node := range nodes {
			res = append(res, node.ToDSL())
		}
	}
	if len(res) == 1 {
		return DSL{"bool": DSL{"must_not": res[0]}}
	} else {
		return DSL{"bool": DSL{"must_not": res}}
	}
}

func (n *NotNode) UnionJoin(AstNode) (AstNode, error) {
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
