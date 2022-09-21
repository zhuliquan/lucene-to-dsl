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
	var res = DSL{}
	if nodes := getMustNodes(n); nodes != nil {
		res["must"] = nodes
	}
	if nodes := getFilterNodes(n); nodes != nil {
		res["filter"] = nodes
	}
	if len(res) == 0 {
		return EmptyDSL
	}
	return DSL{"bool": res}
}

func getMustNodes(n *AndNode) interface{} {
	var mustDSLes = []DSL{}
	for _, nodes := range n.MustNodes {
		for _, node := range nodes {
			mustDSLes = append(mustDSLes, node.ToDSL())
		}
	}
	if len(mustDSLes) == 0 {
		return nil
	} else if len(mustDSLes) == 1 {
		return mustDSLes[0]
	} else {
		return mustDSLes
	}
}

func getFilterNodes(n *AndNode) interface{} {
	var filterDSLes = []DSL{}
	for _, nodes := range n.FilterNodes {
		for _, node := range nodes {
			filterDSLes = append(filterDSLes, node.ToDSL())
		}
	}
	if len(filterDSLes) == 0 {
		return nil
	} else if len(filterDSLes) == 1 {
		return filterDSLes[0]
	} else {
		return filterDSLes
	}
}
