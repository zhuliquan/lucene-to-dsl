package dsl

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

const _ID = "_id"

type IdsNode struct {
	lfNode
	ids []string
}

func NewIdsNode(lfNode *lfNode, ids []string) *IdsNode {
	return &IdsNode{lfNode: *lfNode, ids: ids}
}

func (n *IdsNode) DslType() DslType {
	return IDS_DSL_TYPE
}

func (n *IdsNode) UnionJoin(o AstNode) (AstNode, error) {
	if o.DslType() == IDS_DSL_TYPE {
		var t = o.(*IdsNode)
		return &IdsNode{
			ids: ValueLstToStrLst(
				UnionJoinValueLst(
					StrLstToValueLst(n.ids),
					StrLstToValueLst(t.ids),
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
			ids: ValueLstToStrLst(
				IntersectValueLst(
					StrLstToValueLst(n.ids),
					StrLstToValueLst(t.ids),
					mapping.KEYWORD_FIELD_TYPE,
				),
			),
		}, nil
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: id dsl only support intersect with id dsl", n, o)
	}
}

func (n *IdsNode) Inverse() (AstNode, error) {
	return NewBoolNode(n, NOT), nil
}

func (n *IdsNode) NodeKey() string {
	return _ID
}

func (n *IdsNode) ToDSL() DSL {
	return DSL{
		IDS_KEY: DSL{
			VALUES_KEY: n.ids,
		},
	}
}
