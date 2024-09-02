package dsl

import (
	"fmt"

	mapping "github.com/zhuliquan/es-mapping"
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
	if checkCommonDslType(o.DslType()) {
		return o.UnionJoin(n)
	}
	switch o.DslType() {
	case IDS_DSL_TYPE:
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
	default:
		return nil, fmt.Errorf("failed to union join %v and %v, err: id dsl only support union join with id dsl", n, o)
	}
}

func (n *IdsNode) InterSect(o AstNode) (AstNode, error) {
	switch o.DslType() {
	case IDS_DSL_TYPE:
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
	case BOOL_DSL_TYPE, MATCH_ALL_DSL_TYPE, EMPTY_DSL_TYPE:
		return o.InterSect(n)
	default:
		return nil, fmt.Errorf("failed to intersect %v and %v, err: id dsl only support intersect with id dsl", n, o)
	}
}

func (n *IdsNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
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
