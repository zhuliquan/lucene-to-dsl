package dsl

// bool node
// must = [a, b]
// must_not = [c, d]
// should = [e, f, g]
// we can get bool expr below:
// a && b && !c && !d && (e || f || g)
type BoolNode struct {
	opNode
	boostNode

	Must    map[string][]AstNode // must
	MustNot map[string][]AstNode // must_not
	Filter  map[string][]AstNode // filter
	Should  map[string][]AstNode // should

	MinimumShouldMatch int
}

func NewBoolNode(node AstNode, opType OpType) *BoolNode {
	boolNode := &BoolNode{
		opNode: opNode{
			opType: opType,
		},
	}
	switch opType {
	case AND:
		if filterCtxNode, ok := node.(FilterCtxNode); ok && filterCtxNode.getFilterCtx() {
			boolNode.Filter = map[string][]AstNode{node.NodeKey(): {node}}
		} else {
			boolNode.Must = map[string][]AstNode{node.NodeKey(): {node}}
		}
	case OR:
		boolNode.Should = map[string][]AstNode{node.NodeKey(): {node}}
		boolNode.MinimumShouldMatch = 1
	case NOT:
		boolNode.MustNot = map[string][]AstNode{node.NodeKey(): {node}}
	}

	return boolNode
}

func (n *BoolNode) DslType() DslType {
	return BOOL_DSL_TYPE
}

func (n *BoolNode) NodeKey() string {
	return OP_KEY
}

func (n *BoolNode) UnionJoin(AstNode) (AstNode, error) {
	return nil, nil
}

func (n *BoolNode) InterSect(o AstNode) (AstNode, error) {
	// var t *AndNode = o.(*AndNode)
	// for key, curNodes := range t.MustNodes {
	// 	if preNodes, ok := n.MustNodes[key]; ok {
	// 		if key == OR_OP_KEY {
	// 			n.MustNodes[key] = append(preNodes, curNodes...)
	// 		} else {
	// 			if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
	// 				return nil, err
	// 			} else {
	// 				delete(n.MustNodes, key)
	// 				n.MustNodes[key] = []AstNode{newNode}
	// 			}
	// 		}

	// 	} else {
	// 		n.MustNodes[key] = curNodes
	// 	}
	// }

	// for key, curNodes := range t.FilterNodes {
	// 	if preNodes, ok := n.FilterNodes[key]; ok {
	// 		if key == OR_OP_KEY {
	// 			n.FilterNodes[key] = append(preNodes, curNodes...)
	// 		} else {
	// 			if newNode, err := preNodes[0].InterSect(curNodes[0]); err != nil {
	// 				return nil, err
	// 			} else {
	// 				delete(n.FilterNodes, key)
	// 				n.FilterNodes[key] = []AstNode{newNode}
	// 			}
	// 		}
	// 	} else {
	// 		n.FilterNodes[key] = curNodes
	// 	}
	// }

	return nil, nil
}

// not 全部都不是的反例是至少有一个:
func (n *BoolNode) Inverse() (AstNode, error) {
	// var resNodes = make(map[string][]AstNode)
	// for key, nodes := range n.MustNodes {
	// 	resNodes[key] = nodes
	// }
	// for key, nodes := range n.FilterNodes {
	// 	resNodes[key] = nodes
	// }
	// return &OrNode{Nodes: resNodes, MinimumShouldMatch: -1}, nil
	return nil, nil
}

func (n *BoolNode) ToDSL() DSL {
	var res = DSL{}
	if nodes := flattenNodes(n.Must); nodes != nil {
		res[MUST_KEY] = nodes
	}
	if nodes := flattenNodes(n.Filter); nodes != nil {
		res[FILTER_KEY] = nodes
	}
	if nodes := flattenNodes(n.Should); nodes != nil {
		res[SHOULD_KEY] = nodes
	}
	if nodes := flattenNodes(n.MustNot); nodes != nil {
		res[MUST_NOT_KEY] = nodes
	}
	if len(res) == 0 {
		return EmptyDSL
	}
	return DSL{BOOL_KEY: res}
}
