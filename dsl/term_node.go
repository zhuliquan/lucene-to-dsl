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
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
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
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
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
	return NewBoolNode(n, NOT), nil
}

func termNodeUnionJoinTermNode(n, o *TermNode) (AstNode, error) {
	if CompareAny(o.value, n.value, n.mType) == 0 {
		return o, nil
	} else {
		return &TermsNode{
			fieldNode: n.fieldNode,
			boostNode: n.boostNode,
			valueType: n.valueType,
			terms:     []LeafValue{n.value, o.value},
		}, nil
	}
}

func termNodeUnionJoinTermsNode(n *TermNode, o *TermsNode) (AstNode, error) {
	var terms = UnionJoinValueLst([]LeafValue{n.value}, o.terms, n.mType)
	if len(terms) == 1 {
		return n, nil
	} else {
		return &TermsNode{
			fieldNode: n.fieldNode,
			boostNode: n.boostNode,
			valueType: n.valueType,
			terms:     terms,
		}, nil
	}
}

func termNodeIntersectTermNode(n, o *TermNode) (AstNode, error) {
	if CompareAny(o.value, n.value, n.mType) == 0 {
		return o, nil
	} else if n.isArrayType() {
		return lfNodeIntersectLfNode(n, o)
	} else {
		return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
	}
}

func termNodeIntersectTermsNode(n *TermNode, o *TermsNode) (AstNode, error) {
	if idx := FindAny(o.terms, n.value, n.mType); idx != -1 {
		if n.isArrayType() {
			terms := append(o.terms[:idx], o.terms[idx+1:]...)
			if len(terms) == 0 {
				return n, nil
			} else if len(terms) == 1 {
				return lfNodeIntersectLfNode(
					n, &TermNode{
						kvNode: kvNode{
							fieldNode: o.fieldNode,
							valueNode: valueNode{value: terms[0], valueType: o.valueType},
						},
						boostNode: o.boostNode,
					},
				)
			} else {
				return lfNodeIntersectLfNode(
					n, &TermsNode{
						fieldNode: o.fieldNode,
						boostNode: o.boostNode,
						valueType: o.valueType,
						terms:     terms,
					},
				)
			}
		} else {
			return n, nil
		}
	} else {
		if n.isArrayType() {
			return lfNodeIntersectLfNode(n, o)
		} else {
			return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", n.ToDSL(), o.ToDSL())
		}
	}
}

func (n *TermNode) ToDSL() DSL {
	return DSL{
		TERM_KEY: DSL{
			n.field: DSL{
				VALUE_KEY: n.toPrintValue(),
				BOOST_KEY: n.getBoost(),
			},
		},
	}
}
