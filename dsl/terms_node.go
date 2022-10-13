package dsl

import (
	"fmt"

	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

type TermsNode struct {
	fieldNode
	boostNode
	mType mapping.FieldType
	terms []LeafValue
}

func NewTermsNode(fieldNode *fieldNode, mType mapping.FieldType, terms []LeafValue, opts ...func(AstNode)) *TermsNode {
	var n = &TermsNode{
		fieldNode: *fieldNode,
		boostNode: boostNode{boost: 1.0},
		mType:     mType,
		terms:     terms,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *TermsNode) DslType() DslType {
	return TERMS_DSL_TYPE
}

func (n *TermsNode) UnionJoin(o AstNode) (AstNode, error) {
	if bn, ok := o.(BoostNode); ok {
		if CompareAny(bn.getBoost(), n.getBoost(), mapping.DOUBLE_FIELD_TYPE) != 0 {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		return o.UnionJoin(n)
	case TERMS_DSL_TYPE:
		var t = o.(*TermsNode)
		t.terms = UnionJoinValueLst(t.terms, n.terms, n.mType)
		return t, nil
	case RANGE_DSL_TYPE:
		return o.(*RangeNode).UnionJoin(n)
	case QUERY_STRING_DSL_TYPE:
		var t = o.(*QueryStringNode)
		// var s = ""
		// for _, val := range n.Values {
		// TODO: 需要 %s 修改
		// s += fmt.Sprintf(" OR %s", val)
		// }
		// t.Value = fmt.Sprintf("%s%s", t.Value, s)
		return t, nil
	case IDS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())

	default:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())

	}
}

func (n *TermsNode) InterSect(o AstNode) (AstNode, error) {
	if n == nil || o == nil {
		return nil, ErrIntersectNilNode
	}
	if v, ok := o.(BoostNode); ok {
		if CompareAny(v.getBoost(), n.getBoost(), mapping.DOUBLE_FIELD_TYPE) != 0 {
			return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case RANGE_DSL_TYPE:
		return o.(*RangeNode).UnionJoin(n)
	case TERM_DSL_TYPE:
		// TODO: 如果 values 存在还行/ 不存在就要怎么办？
		return nil, nil
	case TERMS_DSL_TYPE:
		return nil, nil
	case IDS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermsNode) Inverse() (AstNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	var nodes = []AstNode{}
	for _, val := range n.terms {
		nodes = append(
			nodes, &TermNode{
				kvNode: kvNode{
					fieldNode: n.fieldNode,
					valueNode: valueNode{mType: n.mType, value: val},
				},
				boostNode: n.boostNode,
			},
		)
	}
	return &NotNode{Nodes: map[string][]AstNode{n.field: nodes}}, nil
}

func (n *TermsNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{
		TERMS_KEY: DSL{
			n.field:   n.terms,
			BOOST_KEY: n.getBoost(),
		},
	}
}
