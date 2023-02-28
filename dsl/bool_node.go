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
		return boolNodeUnionJoinBoolNode(n, x.(*BoolNode))
	} else {
		return boolNodeUnionJoinLeafNode(n, x)
	}
}

//    (x1 and x2 and (x3 or x4) and not x5 and not x6) or
//    (y1 and y2 and (y3 or y4) and not y5 and not y6)
// -------------------------------------------------------
// 1、 (x3 or x4) or (y3 or y4)
// 2、 (not x5 and not x6) or (not y5 and not y6) => not (x5 or x6) or not (y5 or y6) => not ((x5 or x6) and (y5 or y6))
func boolNodeUnionJoinBoolNode(n, o *BoolNode) (AstNode, error) {
	if n.opType == o.opType && n.opType == OR {
		var t AstNode = n
		var err error
		for _, node := range flattenNodes(o.Should) {
			t, err = t.UnionJoin(node)
			if err != nil {
				return nil, err
			}
		}
		return t, nil
	}
	if n.opType == o.opType && n.opType == NOT {
		orNode1 := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,
		}
		orNode2 := &BoolNode{
			opNode: opNode{opType: OR},
			Should: o.MustNot,
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
		n.minimumShouldMatch = 1
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
// => (x1 and x2 and y1 and y2) and
// => (x3 or x4) and (y3 or y4) and
// => and not x5 and not x6 and not y5 and not y6
func boolNodeIntersectBoolNode(n, o *BoolNode) (AstNode, error) {
	var err error
	if o.opType|AND == AND {
		if n, err = boolNodeIntersectAndNode(n, o); err != nil {
			return nil, err
		}
	}
	if n.opType|NOT == NOT {
		if n, err = boolNodeIntersectNotNode(n, o); err != nil {
			return nil, err
		}
	}

	if n.opType|OR == OR {
		if n, err = boolNodeIntersectOrNode(n, o); err != nil {
			return nil, err
		}
	}
	return n, nil
}

// (x1 and x2 and y1 and y2)
func boolNodeIntersectAndNode(n, o *BoolNode) (*BoolNode, error) {
	var err error
	if len(o.Must) > 0 && len(n.Must) > 0 {
		var t AstNode = n
		for _, node := range flattenNodes(o.Must) {
			t, err = t.InterSect(node)
			if err != nil {
				return nil, err
			}
		}
		n = t.(*BoolNode)
	}
	if len(o.Must) > 0 && len(n.Must) == 0 {
		n.Must = o.Must
	}
	if len(o.Filter) > 0 && len(n.Filter) > 0 {
		var t AstNode = n
		for _, node := range flattenNodes(o.Filter) {
			t, err = t.InterSect(node)
			if err != nil {
				return nil, err
			}
		}
		n = t.(*BoolNode)
	}
	if len(o.Filter) > 0 && len(n.Filter) == 0 {
		n.Filter = o.Filter
	}
	n.opType |= AND
	return n, nil
}

// and not x1 and not x2 and not y1 and not y2
func boolNodeIntersectNotNode(n, o *BoolNode) (*BoolNode, error) {
	if n.opType|NOT != NOT {
		n.MustNot = o.MustNot
		n.opType |= NOT
		return n, nil
	} else {
		orNode1 := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,
		}
		orNode2 := &BoolNode{
			opNode: opNode{opType: OR},
			Should: o.MustNot,
		}
		orNode, err := orNode1.UnionJoin(orNode2)
		if err != nil {
			return nil, err
		}
		n.MustNot = nil
		n.opType ^= NOT
		notNode, err := orNode.Inverse()
		if err != nil {
			return nil, err
		}
		x, err := n.InterSect(notNode)
		if err != nil {
			return nil, err
		}
		n = x.(*BoolNode)
		return n, nil
	}
}

func boolNodeIntersectOrNode(n, o *BoolNode) (*BoolNode, error) {
	if n.opType|OR != OR {
		n.Should = o.Should
		n.opType |= OR
		return n, nil
	} else {
		orNode1 := ReduceAstNode(&BoolNode{
			opNode: opNode{opType: OR},
			Should: n.Should,

			minimumShouldMatch: 1,
		})
		orNode2 := ReduceAstNode(&BoolNode{
			opNode: opNode{opType: OR},
			Should: n.Should,

			minimumShouldMatch: 1,
		})
		node := &BoolNode{
			opNode:  opNode{opType: n.opType ^ OR},
			Must:    n.Must,
			Filter:  n.Filter,
			MustNot: n.MustNot,
		}
		node.Must[orNode1.NodeKey()] = append(node.Must[orNode1.NodeKey()], orNode1)
		node.Must[orNode2.NodeKey()] = append(node.Must[orNode2.NodeKey()], orNode2)
		return node, nil
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
			opNode: opNode{opType: OR},
			Should: n.MustNot,

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
			opNode: opNode{opType: OR},
			Should: n.Should,

			minimumShouldMatch: 1,
		}).Inverse()
		notNode = ReduceAstNode(notNode)
		orNode := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,

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
			opNode: opNode{opType: OR},
			Should: n.Should,

			minimumShouldMatch: 1,
		}).Inverse()
		notNode2 = ReduceAstNode(notNode2)
		orNode := &BoolNode{
			opNode: opNode{opType: OR},
			Should: n.MustNot,

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
