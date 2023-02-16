package dsl

import "fmt"

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

func NewBoolNode(node AstNode, opType OpType) AstNode {
	boolNode := &BoolNode{
		opNode: opNode{opType: opType},
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

func (n *BoolNode) UnionJoin(x AstNode) (AstNode, error) {
	if x.AstType() == OP_NODE_TYPE {
		o := x.(*BoolNode)
		n.opType |= o.opType
		return n, nil
	} else {
		return boolNodeUnionJoinLeafNode(n, x)
	}
}

func boolNodeUnionJoinLeafNode(n *BoolNode, x AstNode) (AstNode, error) {
	n.opType |= OR
	if n.Should == nil {
		n.Should = map[string][]AstNode{}
	}
	nodes := n.Should[x.NodeKey()]
	errs := []error{}
	merge := false
	for i, node := range nodes {
		n3, err := node.UnionJoin(x)
		if err == nil {
			if n3.AstType() != OP_NODE_TYPE { // union join two nodes into a single node as soon as possible
				nodes[i] = n3
				merge = true
				break
			}
		} else {
			errs = append(errs, err)
		}
	}
	if !merge {
		if len(errs) == len(nodes) && len(nodes) != 0 {
			return nil, fmt.Errorf("failed to union node: %+v, errs: %+v", x, errs)
		} else {
			nodes = append(nodes, x)
		}
	}
	n.Should[x.NodeKey()] = nodes
	return n, nil
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
		res[MINIMUM_SHOULD_MATCH_KEY] = n.MinimumShouldMatch
	}
	if nodes := flattenNodes(n.MustNot); nodes != nil {
		res[MUST_NOT_KEY] = nodes
	}
	if len(res) == 0 {
		return EmptyDSL
	}
	res[BOOST_KEY] = n.getBoost()
	return DSL{BOOL_KEY: res}
}

// or
// var nodes = map[string][]dsl.AstNode{node.NodeKey(): {node}}
// 		for _, query := range q.OSQuery {
// 			if curNode, err := osQueryToAstNode(query); err != nil {
// 				return nil, err
// 			} else {
// 				if preNode, ok := nodes[curNode.NodeKey()]; ok {
// 					if curNode.DslType() == dsl.AND_DSL_TYPE ||
// 						curNode.DslType() == dsl.NOT_DSL_TYPE {
// 						nodes[curNode.NodeKey()] = append(nodes[curNode.NodeKey()], curNode)
// 					} else {
// 						if node, err := preNode[0].UnionJoin(curNode); err != nil {
// 							return nil, err
// 						} else {
// 							delete(nodes, curNode.NodeKey())
// 							nodes[node.NodeKey()] = []dsl.AstNode{node}
// 						}
// 					}
// 				} else {
// 					nodes[curNode.NodeKey()] = []dsl.AstNode{curNode}
// 				}
// 			}
// 		}
// 		if len(nodes) == 1 {
// 			for _, ns := range nodes {
// 				if len(ns) == 1 {
// 					return ns[0], nil
// 				}
// 			}
// 		}
// 		return &dsl.OrNode{Nodes: nodes}, nil

// and
// var nodes = map[string][]dsl.AstNode{node.NodeKey(): {node}}
// 		for _, query := range q.AnSQuery {
// 			if curNode, err := ansQueryToAstNode(query); err != nil {
// 				return nil, err
// 			} else {

// 				if preNode, ok := nodes[curNode.NodeKey()]; ok {
// 					if curNode.DslType() == dsl.OR_DSL_TYPE {
// 						nodes[curNode.NodeKey()] = append(nodes[curNode.NodeKey()], curNode)
// 					} else {
// 						if node, err := preNode[0].InterSect(curNode); err != nil {
// 							return nil, err
// 						} else {
// 							delete(nodes, curNode.NodeKey())
// 							nodes[node.NodeKey()] = []dsl.AstNode{node}
// 						}
// 					}
// 				} else {
// 					nodes[curNode.NodeKey()] = []dsl.AstNode{curNode}
// 				}
// 			}
// 		}
// 		if len(nodes) == 1 {
// 			for _, ns := range nodes {
// 				if len(ns) == 1 {
// 					return ns[0], nil
// 				}
// 			}
// 		}
// 		return &dsl.AndNode{MustNodes: nodes}, nil

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
