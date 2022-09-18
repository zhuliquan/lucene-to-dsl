package dsl

type AndNode struct {
	OpNode
	MustNodes   map[string][]AstNode
	FilterNodes map[string][]AstNode
}

func (n *AndNode) DslType() DslType {
	return AND_DSL_TYPE
}

func (n *AndNode) NodeKey() string {
	return AND_OP_KEY
}

func (n *AndNode) UnionJoin(AstNode) (AstNode, error) {
	return nil, nil
}

func (n *AndNode) InterSect(o AstNode) (AstNode, error) {
	var t *AndNode = o.(*AndNode)
	for key, curNodes := range t.MustNodes {
		if preNodes, ok := n.MustNodes[key]; ok {
			if key == OR_OP_KEY {
				n.MustNodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.MustNodes, key)
					n.MustNodes[key] = []AstNode{newNode}
				}
			}

		} else {
			n.MustNodes[key] = curNodes
		}
	}

	for key, curNodes := range t.FilterNodes {
		if preNodes, ok := n.FilterNodes[key]; ok {
			if key == OR_OP_KEY {
				n.FilterNodes[key] = append(preNodes, curNodes...)
			} else {
				if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
					return nil, err
				} else {
					delete(n.FilterNodes, key)
					n.FilterNodes[key] = []AstNode{newNode}
				}
			}
		} else {
			n.FilterNodes[key] = curNodes
		}
	}

	return n, nil
}

func (n *AndNode) Inverse() (AstNode, error) {
	var resNodes = make(map[string][]AstNode)
	for key, nodes := range n.MustNodes {
		resNodes[key] = nodes
	}
	for key, nodes := range n.FilterNodes {
		resNodes[key] = nodes
	}
	return &OrNode{Nodes: resNodes, MinimumShouldMatch: -1}, nil
}

func (n *AndNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	var FRes = []DSL{}
	var MRes = []DSL{}
	for _, nodes := range n.MustNodes {
		for _, node := range nodes {
			MRes = append(MRes, node.ToDSL())
		}
	}
	for _, nodes := range n.FilterNodes {
		for _, node := range nodes {
			FRes = append(FRes, node.ToDSL())
		}
	}

	if len(FRes) == 1 && len(n.MustNodes) == 0 {
		return DSL{"bool": DSL{"filter": FRes[0]}}
	} else if len(FRes) == 1 && len(n.MustNodes) == 1 {
		return DSL{"bool": DSL{"must": MRes[0], "filter": FRes[0]}}
	} else if len(FRes) == 1 && len(n.MustNodes) > 1 {
		return DSL{"bool": DSL{"must": MRes, "filter": FRes[0]}}
	} else if len(FRes) == 0 && len(n.MustNodes) == 1 {
		return MRes[0]
	} else if len(FRes) == 0 && len(n.MustNodes) > 1 {
		return DSL{"bool": DSL{"must": MRes}}
	} else if len(FRes) > 1 && len(n.MustNodes) == 0 {
		return DSL{"bool": DSL{"filter": FRes}}
	} else if len(FRes) > 1 && len(n.MustNodes) == 1 {
		return DSL{"bool": DSL{"must": MRes[0], "filter": FRes}}
	} else {
		return DSL{"bool": DSL{"must": MRes, "filter": FRes}}
	}
}
