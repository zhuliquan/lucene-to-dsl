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

	MinimumShouldMatch int
}

func getMinimumShouldMatch(opType OpType) int {
	if opType&OR == OR {
		return 1
	} else {
		return 0
	}
}

func newDefaultBoolNode(opType OpType) *BoolNode {
	return &BoolNode{
		opNode: opNode{opType: opType},

		MinimumShouldMatch: getMinimumShouldMatch(opType),
	}
}

func NewBoolNode(node AstNode, opType OpType) AstNode {
	boolNode := newDefaultBoolNode(opType)
	switch opType {
	case AND:
		if filterCtxNode, ok := node.(FilterCtxNode); ok && filterCtxNode.GetFilterCtx() {
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
		return boolNodeUnionJoinBoolNode(n, x.(*BoolNode))
	} else {
		return boolNodeUnionJoinLeafNode(n, x)
	}
}

//    (x1 and x2 and (x3 or x4) and not x5 and not x6) or
//    (y1 and y2 and (y3 or y4) and not y5 and not y6)
// -------------------------------------------------------
// =>
// 1. (x1 and x2) or (y1 and y2)
//    (x3  or x4) or (y3  or y4) or
//    (not x5 and not x6) or (not y5 and not y6)
// =>
// 2. (x1 and x2) or (y1 and y2)
//    (x3  or x4) or (y3  or y4) or
//    not (x5 or x6) or not (y5 or y6)
// =>
// 3. (x1 and x2) or (y1 and y2)
//    (x3  or x4) or (y3  or y4) or
//    not ((x5 or x6) and (y5 or y6))
func boolNodeUnionJoinBoolNode(n, o *BoolNode) (AstNode, error) {
	if n.opType == OR && o.opType == OR {
		var tmp AstNode = n
		var err error
		for _, node := range flattenNodes(o.Should) {
			tmp, err = tmp.UnionJoin(node)
			if err != nil {
				return nil, err
			}
		}
		return tmp, nil
	}
	if n.opType == OR {
		n.Should[o.NodeKey()] = append(n.Should[o.NodeKey()], o)
		return n, nil
	}
	if o.opType == OR {
		o.Should[n.NodeKey()] = append(o.Should[n.NodeKey()], n)
		return o, nil
	}

	if n.opType == NOT && o.opType == NOT {
		orNode1 := &BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.MustNot,
			MinimumShouldMatch: 1,
		}
		orNode2 := &BoolNode{
			opNode:             opNode{opType: OR},
			Should:             o.MustNot,
			MinimumShouldMatch: 1,
		}
		andNode, err := orNode1.InterSect(orNode2)
		if err != nil {
			return nil, err
		}
		return andNode.Inverse()
	}

	return &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			OP_KEY: {n, o},
		},
		MinimumShouldMatch: 1,
	}, nil

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
		n.MinimumShouldMatch = 1
	}
	return ReduceAstNode(n), nil
}

func (n *BoolNode) InterSect(x AstNode) (AstNode, error) {
	if x.AstType() == OP_NODE_TYPE {
		return boolNodeIntersectBoolNode(n, x.(*BoolNode))
	} else {
		return boolNodeIntersectLeafNode(n, x)
	}
}

//    (x1 and x2 and (x3 or x4) and not x5 and not x6) and
//    (y1 and y2 and (y3 or y4) and not y5 and not y6)
// -------------------------------------------------------
// =>
//    (x1 and x2 and y1 and y2) and
//    (x3 or x4) and (y3 or y4) and
//    and not x5 and not x6 and not y5 and not y6
func boolNodeIntersectBoolNode(n, o *BoolNode) (AstNode, error) {
	var err error
	var res AstNode = n
	if o.opType&AND == AND {
		if res, err = boolNodeIntersectAndNode(n, o); err != nil {
			return nil, err
		}
		if res.AstType() == LEAF_NODE_TYPE {
			res = NewBoolNode(res, AND).(*BoolNode)
		}
	}
	n = res.(*BoolNode)
	if o.opType&NOT == NOT {
		if res, err = boolNodeIntersectNotNode(n, o); err != nil {
			return nil, err
		}
		if res.AstType() == LEAF_NODE_TYPE {
			res = NewBoolNode(res, AND).(*BoolNode)
		}
	}

	n = res.(*BoolNode)
	if o.opType&OR == OR {
		if res, err = boolNodeIntersectOrNode(n, o); err != nil {
			return nil, err
		}
		if res.AstType() == LEAF_NODE_TYPE {
			res = NewBoolNode(res, AND).(*BoolNode)
		}
	}
	return ReduceAstNode(res), nil
}

// (x1 and x2 and y1 and y2)
func boolNodeIntersectAndNode(n, o *BoolNode) (AstNode, error) {
	if len(o.Must) > 0 {
		if len(n.Must) == 0 {
			n.Must = o.Must
		} else {
			tmp, err := boolNodeIntersectAndNodes(n, o.Must)
			if err != nil {
				return nil, err
			}
			if tmp.AstType() == LEAF_NODE_TYPE {
				n = NewBoolNode(tmp, AND).(*BoolNode)
			}
		}
	}
	if len(o.Filter) > 0 {
		if len(n.Filter) == 0 {
			n.Filter = o.Filter
		} else {
			tmp, err := boolNodeIntersectAndNodes(n, o.Filter)
			if err != nil {
				return nil, err
			}
			if tmp.AstType() == LEAF_NODE_TYPE {
				n = NewBoolNode(tmp, AND).(*BoolNode)
			}
		}
	}

	n.opType |= AND
	return n, nil
}

func boolNodeIntersectAndNodes(n *BoolNode, nodesMap map[string][]AstNode) (AstNode, error) {
	var tmp AstNode = n
	var err error
	for _, node := range flattenNodes(nodesMap) {
		tmp, err = tmp.InterSect(node)
		if err != nil {
			return nil, err
		}
	}
	return tmp, nil
}

// and not x1 and not x2 and not y1 and not y2
func boolNodeIntersectNotNode(n, o *BoolNode) (AstNode, error) {
	if n.opType&NOT != NOT {
		n.MustNot = o.MustNot
		n.opType |= NOT
		return n, nil
	} else {
		orNode1 := &BoolNode{
			opNode:             opNode{opType: OR},
			Should:             n.MustNot,
			MinimumShouldMatch: 1,
		}
		orNode2 := &BoolNode{
			opNode:             opNode{opType: OR},
			Should:             o.MustNot,
			MinimumShouldMatch: 1,
		}
		orNode, err := orNode1.UnionJoin(orNode2)
		if err != nil {
			return nil, err
		}
		notNode, err := orNode.Inverse()
		if err != nil {
			return nil, err
		}
		andNode := &BoolNode{
			opNode: opNode{opType: n.opType & ^NOT | AND},
			Must:   n.Must,
			Filter: n.Filter,
			Should: n.Should,
		}
		if andNode.Must == nil {
			andNode.Must = map[string][]AstNode{}
		}
		andNode.Must[andNode.NodeKey()] = append(andNode.Must[andNode.NodeKey()], notNode)
		return andNode, nil
	}
}

func boolNodeIntersectOrNode(n, o *BoolNode) (AstNode, error) {
	if n.opType&OR != OR {
		n.Should = o.Should
		n.opType |= OR
		n.MinimumShouldMatch = 1
		return n, nil
	} else {
		orNode1 := ReduceAstNode(&BoolNode{
			opNode: opNode{opType: OR},
			Should: n.Should,

			MinimumShouldMatch: 1,
		})
		orNode2 := ReduceAstNode(&BoolNode{
			opNode: opNode{opType: OR},
			Should: o.Should,

			MinimumShouldMatch: 1,
		})
		node := &BoolNode{
			opNode:  opNode{opType: n.opType & ^OR},
			Must:    n.Must,
			Filter:  n.Filter,
			MustNot: n.MustNot,
		}
		if node.Must == nil {
			node.Must = make(map[string][]AstNode)
		}
		node.Must[orNode1.NodeKey()] = append(node.Must[orNode1.NodeKey()], orNode1)
		node.Must[orNode2.NodeKey()] = append(node.Must[orNode2.NodeKey()], orNode2)
		node.opType |= AND
		return node, nil
	}
}

func boolNodeIntersectLeafNode(n *BoolNode, x AstNode) (AstNode, error) {
	n.opType |= AND
	if filterNode, ok := x.(FilterCtxNode); ok && filterNode.GetFilterCtx() {
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
	case AND, AND | OR, AND | OR | NOT:
		// case1:   not (x1 and x2) 
		//       => #*:* -(#x1 #x2)
		// case2:   not ((x1 and x2) and (x3 or x4))
		//       => #*:* -(#x1 #x2 x3 x4)
		// case3:   not ((x1 and x2) and (x3 or x4) and (not x5 and not x6))
		//       => #*:* - (#1 #2 x3 x4 (#*:* -x5 -x6))
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
		// case1:   not (not x1 and not x2)
		//       => not not x1 or not not x2
		//       => x or y => should clause query
		// case2:   not (not x1) 
		//       => x1
		return ReduceAstNode(&BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,

			MinimumShouldMatch: 1,
		}), nil
	case AND | NOT:
		//    not (x1 and x2 and not x3 and not x4)
		// => not (x1 and x2 and not (x3 or x4))
		// => not x1 or not x2 or x3 or x4
		// => not (x1 and x2) or x3 or x4
		notNode, err := (&BoolNode{
			opNode: opNode{opType: AND},
			Must:   n.Must,
			Filter: n.Filter,
		}).Inverse()
		if err != nil {
			return nil, err
		}
		notNode = ReduceAstNode(notNode)

		orNode := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,

			MinimumShouldMatch: 1,
		}
		return orNode.UnionJoin(notNode)
	case OR | NOT:
		//    not ((x1 or x2) and not x3 and not x4)
		// => not ((x1 or x2) and not (x3 or x4))
		// => not (x1 or x2) or (x3 or x4)
		// => not (x1 or x2) or x3 or x4
		notNode, err := (&BoolNode{
			opNode: opNode{opType: OR},
			Should: n.Should,

			MinimumShouldMatch: 1,
		}).Inverse()
		if err != nil {
			return nil, err
		}
		notNode = ReduceAstNode(notNode)
		orNode := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,

			MinimumShouldMatch: 1,
		}
		return orNode.UnionJoin(notNode)
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
	res[MINIMUM_SHOULD_MATCH_KEY] = n.MinimumShouldMatch
	return DSL{BOOL_KEY: res}
}
