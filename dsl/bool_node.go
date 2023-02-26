package dsl

// bool node
// must = [a, b]
// must_not = [c, d]
// should = [e, f, g]
// we can get bool expr below:
// a && b && !c && !d && (e || f || g)
type BoolNode struct {
	opNode

	Must    map[string][]AstNode // must
	MustNot map[string][]AstNode // must_not
	Filter  map[string][]AstNode // filter
	Should  map[string][]AstNode // should

	minimumShouldMatch int
}

func getMinimumShouldMatch(opType OpType) int {
	if opType|OR == OR {
		return 1
	} else {
		return 0
	}
}

func newDefaultBoolNode(opType OpType) *BoolNode {
	return &BoolNode{
		opNode: opNode{opType: opType},

		minimumShouldMatch: getMinimumShouldMatch(opType),
	}
}

func NewBoolNode(node AstNode, opType OpType) AstNode {
	boolNode := newDefaultBoolNode(opType)
	switch opType {
	case AND:
		if filterCtxNode, ok := node.(FilterCtxNode); ok && filterCtxNode.getFilterCtx() {
			boolNode.Filter = map[string][]AstNode{node.NodeKey(): {node}}
		} else {
			boolNode.Must = map[string][]AstNode{node.NodeKey(): {node}}
		}
	case OR:
		boolNode.Should = map[string][]AstNode{node.NodeKey(): {node}}
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
	} else {
		return boolNodeUnionJoinLeafNode(n, x)
	}
}

func boolNodeUnionJoinLeafNode(n *BoolNode, x AstNode) (AstNode, error) {
	n.opType |= OR
	if n.Should == nil {
		n.Should = make(map[string][]AstNode, 0)
	}
	key := x.NodeKey()
	if nodes, err := reduceAstNodes(append(n.Should[key], x), UNION_JOIN, UnionJoin); err != nil {
		return nil, err
	} else {
		n.Should[key] = nodes
		n.minimumShouldMatch = 1
	}
	return ReduceAstNode(n), nil
}

func (n *BoolNode) InterSect(x AstNode) (AstNode, error) {
	if x.AstType() == OP_NODE_TYPE {
		o := x.(*BoolNode)
		n.opType |= o.opType
		return n, nil
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
	} else {
		return boolNodeIntersectLeafNode(n, x)
	}
}

func boolNodeIntersectLeafNode(n *BoolNode, x AstNode) (AstNode, error) {
	n.opType |= AND
	if filterNode, ok := x.(FilterCtxNode); ok && filterNode.getFilterCtx() {
		return boolNodeIntersectFilterLeafNode(n, x)
	} else {
		return boolNodeIntersectMustLeafNode(n, x)
	}
}

func boolNodeIntersectFilterLeafNode(n *BoolNode, x AstNode) (AstNode, error) {
	if n.Filter == nil {
		n.Filter = make(map[string][]AstNode, 0)
	}
	key := x.NodeKey()
	if nodes, err := reduceAstNodes(append(n.Filter[key], x), INTERSECT, Intersect); err != nil {
		return nil, err
	} else {
		n.Filter[key] = nodes
	}
	return n, nil
}

func boolNodeIntersectMustLeafNode(n *BoolNode, x AstNode) (AstNode, error) {
	if n.Must == nil {
		n.Must = make(map[string][]AstNode, 0)
	}
	key := x.NodeKey()
	if nodes, err := reduceAstNodes(append(n.Must[key], x), INTERSECT, Intersect); err != nil {
		return nil, err
	} else {
		n.Must[key] = nodes
	}
	return ReduceAstNode(n), nil
}

func (n *BoolNode) Inverse() (AstNode, error) {
	// rule1: make must_not clause fewer
	switch n.opType {
	case AND:
		// not (x1 and x2) => #*:* -(#x1 #x2)
		return inverseNode(n), nil
	case OR:
		//    not (x1 or x2)
		// => not x1 and not x2
		// => #*:* -x1 -x2 => must_not clause query
		return &BoolNode{
			opNode:  opNode{opType: NOT},
			MustNot: n.Should,
		}, nil
	case NOT:
		// case1: not (not x1 and not x2)
		//     => not not x1 or not not x2
		//     => x or y => should clause query
		// not (not x1) => x1
		return ReduceAstNode(&BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.MustNot,
			minimumShouldMatch: 1,
		}), nil
	case AND | NOT:
		//    not (x1 and x2 and not x3 and not x4)
		// => not (x1 and x2 and not (x3 or x4))
		// => not x1 or not x2 or x3 or x4
		// => not (x1 and x2) or x3 or x4
		notNode, _ := (&BoolNode{
			opNode: opNode{opType: AND},
			Must:   n.Must,
			Filter: n.Filter,
		}).Inverse()
		notNode = ReduceAstNode(notNode)
		orNode := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.Should,

			minimumShouldMatch: 1,
		}
		orNode.Should[notNode.NodeKey()] = append(orNode.Should[notNode.NodeKey()], notNode)
		return orNode, nil
	case AND | OR:
		return inverseNode(n), nil
	case OR | NOT:
		//    not ((x1 or x2) and not x3 and not x4)
		// => not ((x1 or x2) and not (x3 or x4))
		// => not (x1 or x2) or (x3 or x4)
		// => not (x1 or x2) or x3 or x4
		notNode, _ := (&BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.Should,
			minimumShouldMatch: 1,
		}).Inverse()
		notNode = ReduceAstNode(notNode)
		orNode := &BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.MustNot,
			minimumShouldMatch: 1,
		}
		orNode.Should[notNode.NodeKey()] = append(orNode.Should[notNode.NodeKey()], notNode)
		return orNode, nil
	case AND | OR | NOT:
		//    not ((x1 and x2) and (x3 or x4) and not x5 and not x6)
		// => not (x1 and x2 and (x3 or x4) and not (x3 or x6))
		// => not (x1 and x2) or not (x3 or x4) or x3 or x6
		notNode1, _ := (&BoolNode{
			opNode: opNode{opType: AND},
			Must:   n.Must,
			Filter: n.Filter,
		}).Inverse()
		notNode1 = ReduceAstNode(notNode1)
		notNode2, _ := (&BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.Should,
			minimumShouldMatch: 1,
		}).Inverse()
		notNode2 = ReduceAstNode(notNode2)
		orNode := &BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.MustNot,
			minimumShouldMatch: 1,
		}
		orNode.Should[notNode1.NodeKey()] = append(orNode.Should[notNode1.NodeKey()], notNode1)
		orNode.Should[notNode2.NodeKey()] = append(orNode.Should[notNode2.NodeKey()], notNode2)
		return orNode, nil
	default:
		return nil, ErrInverseNilNode
	}
}

func (n *BoolNode) ToDSL() DSL {
	var res = DSL{}
	if nodes := flattenNodes(n.Must); len(nodes) != 0 {
		res[MUST_KEY] = reduceDSLList(nodesToDSLList(nodes))
	}
	if nodes := flattenNodes(n.Filter); len(nodes) != 0 {
		res[FILTER_KEY] = reduceDSLList(nodesToDSLList(nodes))
	}
	if nodes := flattenNodes(n.Should); len(nodes) != 0 {
		res[SHOULD_KEY] = reduceDSLList(nodesToDSLList(nodes))
	}
	if nodes := flattenNodes(n.MustNot); len(nodes) != 0 {
		res[MUST_NOT_KEY] = reduceDSLList(nodesToDSLList(nodes))
	}
	if len(res) == 0 {
		return EmptyDSL
	}
	res[MINIMUM_SHOULD_MATCH_KEY] = n.minimumShouldMatch
	return DSL{BOOL_KEY: res}
}
