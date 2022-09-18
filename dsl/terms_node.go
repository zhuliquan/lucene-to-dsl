package dsl

import (
	"fmt"
)

type TermsNode struct {
	KvNode
	Values []LeafValue
	Boost  float64
}

func (n *TermsNode) getBoost() float64 {
	return n.Boost
}

func (n *TermsNode) DslType() DslType {
	return TERMS_DSL_TYPE
}

func (n *TermsNode) UnionJoin(o AstNode) (AstNode, error) {
	if bn, ok := o.(boostNode); ok {
		if CompareAny(bn.getBoost(), n.Boost, n.Type) != 0 {
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
		t.Values = UnionJoinValueLst(t.Values, n.Values, n.Type)
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
	if bn, ok := o.(boostNode); ok {
		if CompareAny(bn.getBoost(), n.Boost, n.Type) != 0 {
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
	for _, val := range n.Values {
		nodes = append(nodes, &TermNode{KvNode: KvNode{Field: n.Field, Value: val}, Boost: n.Boost})
	}
	return &NotNode{Nodes: map[string][]AstNode{n.Field: nodes}}, nil
}

func (n *TermsNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"terms": DSL{n.Field: n.Values, "boost": n.Boost}}
}
