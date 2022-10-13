package dsl

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

// query_string node
type QueryStringNode struct {
	kvNode
	boostNode
}

func NewQueryStringNode(kvNode *kvNode, opts ...func(AstNode)) *QueryStringNode {
	var n = &QueryStringNode{
		kvNode:    *kvNode,
		boostNode: boostNode{boost: 1.0},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *QueryStringNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		if b, ok := o.(BoostNode); ok {
			if CompareAny(n.getBoost(), b.getBoost(), mapping.DOUBLE_FIELD_TYPE) != 0 {
				return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return lfNodeUnionJoinLfNode(n, o)
	}
}

func (n *QueryStringNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	default:
		if b, ok := o.(BoostNode); ok {
			if CompareAny(n.getBoost(), b.getBoost(), mapping.DOUBLE_FIELD_TYPE) != 0 {
				return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return lfNodeIntersectLfNode(n, o)
	}
}

func (n *QueryStringNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
}

func (n *QueryStringNode) DslType() DslType {
	return QUERY_STRING_DSL_TYPE
}

func (n *QueryStringNode) ToDSL() DSL {
	return DSL{
		QUERY_STRING_KEY: DSL{
			QUERY_KEY:         n.toPrintValue(),
			BOOST_KEY:         n.getBoost(),
			DEFAULT_FIELD_KEY: n.field,
		},
	}
}
