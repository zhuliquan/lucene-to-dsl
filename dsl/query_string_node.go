package dsl

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

// query_string node
type QueryStringNode struct {
	KvNode
	Boost float64
}

func (n *QueryStringNode) getBoost() float64 {
	return n.Boost
}

func (n *QueryStringNode) UnionJoin(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		if b, ok := o.(boostNode); ok {
			if CompareAny(n.getBoost(), b.getBoost(), mapping.DOUBLE_FIELD_TYPE) != 0 {
				return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return &OrNode{
			MinimumShouldMatch: 1,
			Nodes: map[string][]AstNode{
				o.NodeKey(): {n, o},
			},
		}, nil
	}
}

func (n *QueryStringNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.InterSect(n)
	default:
		if b, ok := o.(boostNode); ok {
			if CompareAny(n.getBoost(), b.getBoost(), mapping.DOUBLE_FIELD_TYPE) != 0 {
				return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
			}
		}
		return &AndNode{
			MustNodes: map[string][]AstNode{
				o.NodeKey(): {n, o},
			},
		}, nil
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
	return DSL{"query_string": DSL{"query": n.Value, "default_field": n.Field, "boost": n.Boost}}
}
