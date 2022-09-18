package dsl

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

type IdsNode struct {
	LfNode
	Values []string
}

func (n *IdsNode) DslType() DslType {
	return IDS_DSL_TYPE
}

func (n *IdsNode) UnionJoin(o AstNode) (AstNode, error) {
	if o.DslType() == IDS_DSL_TYPE {
		var t = o.(*IdsNode)
		return &IdsNode{
			Values: ValueLstToStrLst(
				UnionJoinValueLst(
					StrLstToValueLst(n.Values),
					StrLstToValueLst(t.Values),
					mapping.KEYWORD_FIELD_TYPE,
				),
			),
		}, nil
	} else {
		return nil, fmt.Errorf("failed to union join %v and %v, err: id dsl only support union join with id dsl", n, o)
	}
}

func (n *IdsNode) InterSect(o AstNode) (AstNode, error) {
	if o.DslType() == IDS_DSL_TYPE {
		var t = o.(*IdsNode)
		return &IdsNode{
			Values: ValueLstToStrLst(
				IntersectValueLst(
					StrLstToValueLst(n.Values),
					StrLstToValueLst(t.Values),
					mapping.KEYWORD_FIELD_TYPE,
				),
			),
		}, nil
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: id dsl only support intersect with id dsl", n, o)
	}
}

func (n *IdsNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.NodeKey(): {n},
		},
	}, nil
}

func (n *IdsNode) NodeKey() string {
	return "LEAF:_id"
}

func (n *IdsNode) ToDSL() DSL {
	return DSL{"ids": DSL{"values": n.Values}}
}
