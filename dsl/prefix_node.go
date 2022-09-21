package dsl

import "strings"

type PrefixNode struct {
	KvNode
}

func (n *PrefixNode) DslType() DslType {
	return PREFIX_DSL_TYPE
}

func (n *PrefixNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case PREFIX_DSL_TYPE:
		return prefixNodeUnionJoinPrefixNode(n, o.(*PrefixNode))
	// case RANGE_DSL_TYPE:
	// 	return prefixNodeUnionJoinRangeNode(n, o.(*RangeNode))
	default:
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *PrefixNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	case PREFIX_DSL_TYPE:
		return prefixNodeIntersectPrefixNode(n, o.(*PrefixNode))
	default:
		return &AndNode{
			MustNodes: map[string][]AstNode{
				n.NodeKey(): {n, o},
			},
		}, nil
	}
}

func (n *PrefixNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
}

func (n *PrefixNode) ToDSL() DSL {
	return DSL{"prefix": DSL{n.Field: n.Value}}
}

func prefixNodeUnionJoinPrefixNode(n, o *PrefixNode) (AstNode, error) {
	if strings.HasPrefix(n.Value.(string), o.Value.(string)) {
		return o, nil
	} else if strings.HasPrefix(o.Value.(string), n.Value.(string)) {
		return n, nil
	} else {
		return &OrNode{
			MinimumShouldMatch: 1,
			Nodes: map[string][]AstNode{
				n.NodeKey(): {n, o},
			},
		}, nil
	}
}

// func prefixNodeUnionJoinRangeNode(n *PrefixNode, o *RangeNode) (AstNode, error) {
// 	var prefix = n.Value.(string)
// 	var leftVal = o.LeftValue.(string)
// 	var rightVal = o.RightValue.(string)
// 	if prefix < leftVal {
// 		if strings.HasPrefix(leftVal, prefix) {
// 			return n, nil
// 		}

// 	}
// 	if prefix == leftVal {

// 	}
// 	if prefix > leftVal && prefix < rightVal {

// 	}

// 	if prefix == rightVal {

// 	}

// 	if prefix > rightVal {

// 	}

// 	return nil, nil
// }

func prefixNodeIntersectPrefixNode(n, o *PrefixNode) (AstNode, error) {
	var prefixN = n.Value.(string)
	var prefixO = o.Value.(string)
	if strings.HasPrefix(prefixN, prefixO) {
		return n, nil
	} else if strings.HasPrefix(prefixO, prefixN) {
		return o, nil
	} else {
		return &AndNode{
			MustNodes: map[string][]AstNode{
				n.NodeKey(): {n, o},
			},
		}, nil
	}
}

// func prefixNodeIntersectPrefixNode(n, o *RangeNode) (AstNode, error) {
// 	if strings.HasPrefix(n.)
// }
