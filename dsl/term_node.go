package dsl

import "fmt"

type TermNode struct {
	kvNode
	boostNode
}

func NewTermNode(kvNode *kvNode, opts ...func(AstNode)) *TermNode {
	var n = &TermNode{
		kvNode:    *kvNode,
		boostNode: boostNode{boost: 1.0},
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *TermNode) DslType() DslType {
	return TERM_DSL_TYPE
}

func (n *TermNode) UnionJoin(o AstNode) (AstNode, error) {
	if bn, ok := o.(BoostNode); ok {
		if CompareAny(bn.getBoost(), n.getBoost(), n.mType) != 0 {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case TERM_DSL_TYPE:
		return termNodeUnionJoinTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return termNodeUnionJoinTermsNode(n, o.(*TermsNode))
	case EXISTS_DSL_TYPE, RANGE_DSL_TYPE, PREFIX_DSL_TYPE, REGEXP_DSL_TYPE:
		return o.UnionJoin(n)
	case IDS_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())

	default:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermNode) InterSect(o AstNode) (AstNode, error) {
	if bn, ok := o.(BoostNode); ok {
		if CompareAny(bn.getBoost(), n.getBoost(), n.mType) != 0 {
			return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case TERM_DSL_TYPE:
		return termNodeIntersectTermNode(n, o.(*TermNode))
	case TERMS_DSL_TYPE:
		return termNodeIntersectTermsNode(n, o.(*TermsNode))
	case EXISTS_DSL_TYPE, RANGE_DSL_TYPE, PREFIX_DSL_TYPE, REGEXP_DSL_TYPE:
		return o.InterSect(n)
	case IDS_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermNode) Inverse() (AstNode, error) {
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.field: {n},
		},
	}, nil
}

func termNodeUnionJoinTermNode(n, o *TermNode) (AstNode, error) {
	if CompareAny(o.value, n.value, n.mType) == 0 {
		return o, nil
	} else {
		return &TermsNode{
			fieldNode: n.fieldNode,
			boostNode: n.boostNode,
			mType:     n.mType,
			terms:     []LeafValue{n.value, o.value},
		}, nil
	}
}

func termNodeUnionJoinTermsNode(n *TermNode, o *TermsNode) (AstNode, error) {
	return &TermsNode{
		fieldNode: n.fieldNode,
		boostNode: n.boostNode,
		mType:     n.mType,
		terms:     UnionJoinValueLst([]LeafValue{n.value}, o.terms, n.mType),
	}, nil
}

func termNodeIntersectTermNode(n, o *TermNode) (AstNode, error) {
	if CompareAny(o.value, n.value, n.mType) == 0 {
		return o, nil
	} else {
		return lfNodeIntersectLfNode(n, o)
	}
}

func termNodeIntersectTermsNode(n *TermNode, o *TermsNode) (AstNode, error) {
	// var values = IntersectValueLst([]LeafValue{n.Value}, o.Values, n.Type)
	// if len(o.Values) == len(values) {
	// 	n.
	// }
	return nil, nil
}

func (n *TermNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{
		"term": DSL{
			n.field: DSL{
				"value": n.toPrintValue(),
				"boost": n.getBoost(),
			},
		},
	}
}
