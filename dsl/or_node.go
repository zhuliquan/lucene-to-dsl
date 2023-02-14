package dsl

// type OrNode struct {
// 	opNode
// 	Should map[string][]AstNode
// 	// minimum node should match
// 	MinimumShouldMatch int
// }

// func NewOrNode(node AstNode) *OrNode {
// 	return &OrNode{
// 		opNode: opNode{},
// 		Should: map[string][]AstNode{
// 			node.NodeKey(): {node},
// 		},
// 		MinimumShouldMatch: 1,
// 	}
// }

// func (n *ShouldBoolNode) DslType() DslType {
// 	return OR_DSL_TYPE
// }

// func (n *ShouldBoolNode) NodeKey() string {
// 	return OR_OP_KEY
// }

// func (n *ShouldBoolNode) UnionJoin(o AstNode) (AstNode, error) {
// 	if o == nil && n == nil {
// 		return nil, ErrUnionJoinNilNode
// 	} else if o == nil && n != nil {
// 		return n, nil
// 	} else if o != nil && n == nil {
// 		return o, nil
// 	}
// 	switch o.DslType() {
// 	case OR_DSL_TYPE:
// 		return nil, nil
// 	case AND_DSL_TYPE:
// 		return nil, nil
// 	default:

// 	}
// 	var t = o.(*OrNode)
// 	for key, curNodes := range t.Nodes {
// 		if preNodes, ok := n.Nodes[key]; ok {
// 			// if key == AND_OP_KEY || key == NOT_OP_KEY {
// 			// 	n.Nodes[key] = append(preNodes, curNodes...)
// 			// } else {
// 			// 	if newNode, err := preNodes[0].UnionJoin(curNodes[0]); err != nil {
// 			// 		return nil, err
// 			// 	} else {
// 			// 		delete(n.Nodes, key)
// 			// 		n.Nodes[key] = []DSLNode{newNode}
// 			// 	}
// 			// }
// 			n.Nodes[key] = append(preNodes, curNodes...)

// 		} else {
// 			n.Nodes[key] = curNodes
// 		}
// 	}
// 	return n, nil
// }

// func (n *OrNode) InterSect(AstNode) (AstNode, error) {
// 	return nil, nil
// }

// func (n *OrNode) Inverse() (AstNode, error) {
// 	if n == nil {
// 		return nil, ErrInverseNilNode
// 	}
// 	return &NotNode{
// 		Nodes: n.Nodes,
// 	}, nil
// }

// func (n *OrNode) ToDSL() DSL {
// 	if n == nil {
// 		return EmptyDSL
// 	}
// 	var res = []DSL{}
// 	for _, nodes := range n.Nodes {
// 		for _, node := range nodes {
// 			res = append(res, node.ToDSL())
// 		}
// 	}
// 	if len(res) == 1 {
// 		return res[0]
// 	} else {
// 		var shouldMatch = 1
// 		if n.MinimumShouldMatch != 0 {
// 			shouldMatch = n.MinimumShouldMatch
// 		}
// 		return DSL{
// 			BOOL_KEY: DSL{
// 				SHOULD_KEY: res,
// 			},
// 			MINIMUM_SHOULD_MATCH_KEY: shouldMatch,
// 		}
// 	}
// }
