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

	Must    map[string][]AstNode // must
	MustNot map[string][]AstNode // must_not
	Filter  map[string][]AstNode // filter
	Should  map[string][]AstNode // should

	MinimumShouldMatch int
}

func newDefaultBoolNode(opType OpType) *BoolNode {
	minimumShouldMatch := 0
	if opType == OR {
		minimumShouldMatch = 1
	}
	return &BoolNode{
		opNode: opNode{opType: opType},

		MinimumShouldMatch: minimumShouldMatch,
	}
}

func NewBoolNode(node AstNode, opType OpType, opts ...func(AstNode)) AstNode {
	boolNode := newDefaultBoolNode(opType)

	for _, opt := range opts {
		opt(boolNode)
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
	if nodes, err := reduceNodes(append(n.Should[key], x), UNION_JOIN, UnionJoin); err != nil {
		return nil, err
	} else {
		n.Should[key] = nodes
		n.MinimumShouldMatch = 1
	}
	return n, nil
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
	if nodes, err := reduceNodes(append(n.Filter[key], x), INTERSECT, Intersect); err != nil {
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
	if nodes, err := reduceNodes(append(n.Must[key], x), INTERSECT, Intersect); err != nil {
		return nil, err
	} else {
		n.Must[key] = nodes
	}
	return n, nil
}

func reduceNodes(nodes []AstNode, mergeMethodName string, mergeMethodFunc MergeMethodFunc) ([]AstNode, error) {
	for before, first := nodes, true; ; first = false {
		rest := before[:len(before)-1]
		node := before[len(before)-1]
		errs := []error{}
		join := false

		for i, n1 := range rest {
			if n2, err := mergeMethodFunc(n1, node); err == nil {
				if n2.AstType() != OP_NODE_TYPE { // merge two nodes into a single node as soon as possible
					rest[i] = n2
					join = true
					break
				}
			} else {
				errs = append(errs, err)
			}
		}

		if !join {
			if first {
				if len(errs) == len(rest) && len(errs) > 0 { // all error for nodes merge with n0
					return nil, fmt.Errorf("failed to %s node: %+v, errs: %+v", mergeMethodName, node, errs)
				} else {
					rest = append(rest, node)
				}
			} else {
				rest = append(rest, node)
			}
			return rest, nil
		}

		before = rest // loop find any other node which can be merge with n0
	}
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
	var onlyMust = true
	if nodes := flattenNodes(n.Filter); nodes != nil {
		res[FILTER_KEY] = nodes
		if _, ok := nodes.([]DSL); ok {
			onlyMust = false
		}
	}
	if nodes := flattenNodes(n.Should); nodes != nil {
		res[SHOULD_KEY] = nodes
		onlyMust = false
	}
	if nodes := flattenNodes(n.MustNot); nodes != nil {
		res[MUST_NOT_KEY] = nodes
		onlyMust = false
	}
	if len(res) == 0 {
		return EmptyDSL
	}
	if onlyMust {
		return res[FILTER_KEY].(DSL)
	} else {
		addValueForDSL(res, MINIMUM_SHOULD_MATCH_KEY, n.MinimumShouldMatch)
		return DSL{BOOL_KEY: res}
	}
}
