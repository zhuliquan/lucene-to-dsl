package dsl

import (
	"fmt"
)

type TermsNode struct {
	fieldNode
	boostNode
	valueType
	terms []LeafValue
}

func NewTermsNode(fieldNode *fieldNode, valueType *valueType, terms []LeafValue, opts ...func(AstNode)) *TermsNode {
	var n = &TermsNode{
		fieldNode: *fieldNode,
		valueType: *valueType,
		boostNode: boostNode{boost: 1.0},
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
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
			return nil, fmt.Errorf("failed to union join %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case TERMS_DSL_TYPE:
		return termsNodeUnionJoinTermsNode(n, o.(*TermsNode))
	case EXISTS_DSL_TYPE, TERM_DSL_TYPE,
		WILDCARD_DSL_TYPE, PREFIX_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE,
		MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return o.UnionJoin(n)
	case IDS_DSL_TYPE:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	default:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermsNode) InterSect(o AstNode) (AstNode, error) {
	if b, ok := o.(BoostNode); ok {
		if compareBoost(n, b) != 0 {
			return nil, fmt.Errorf("failed to intersect %s and %s, err: boost is conflict", n.ToDSL(), o.ToDSL())
		}
	}
	switch o.DslType() {
	case TERMS_DSL_TYPE:
		return termsNodeIntersectTermsNode(n, o.(*TermsNode))
	case EXISTS_DSL_TYPE, TERM_DSL_TYPE,
		WILDCARD_DSL_TYPE, PREFIX_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE,
		MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return o.InterSect(n)
	case IDS_DSL_TYPE:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermsNode) Inverse() (AstNode, error) {
	var nodes = []AstNode{}
	for _, val := range n.terms {
		nodes = append(
			nodes, &TermNode{
				kvNode: kvNode{
					fieldNode: n.fieldNode,
					valueNode: valueNode{valueType: n.valueType, value: val},
				},
				boostNode: n.boostNode,
			},
		)
	}
	return &NotNode{
		opNode: opNode{filterCtxNode: n.filterCtxNode},
		Nodes:  map[string][]AstNode{n.field: nodes},
	}, nil
}

func (n *TermsNode) ToDSL() DSL {
	return DSL{
		TERMS_KEY: DSL{
			n.field:   termsToPrintValue(n.terms, n.mType),
			BOOST_KEY: n.getBoost(),
		},
	}
}

func termsNodeUnionJoinTermsNode(a, b *TermsNode) (AstNode, error) {
	terms := UnionJoinValueLst(a.terms, b.terms, a.mType)
	if len(terms) == 1 {
		return &TermNode{
			kvNode: kvNode{
				fieldNode: a.fieldNode,
				valueNode: valueNode{
					valueType: a.valueType,
					value:     terms[0],
				},
			},
			boostNode: a.boostNode,
		}, nil
	} else {
		return &TermsNode{
			fieldNode: a.fieldNode,
			boostNode: a.boostNode,
			valueType: a.valueType,
			terms:     terms,
		}, nil
	}
}

func termsNodeIntersectTermsNode(a, b *TermsNode) (AstNode, error) {
	terms := IntersectValueLst(a.terms, b.terms, a.mType)
	if a.isArrayType() {
		diff := DifferenceValueList(a.terms, b.terms, a.mType)
		if len(diff) == 0 {
			return b, nil
		} else if len(diff) == 1 {
			return lfNodeIntersectLfNode(
				&TermNode{
					kvNode: kvNode{
						fieldNode: a.fieldNode,
						valueNode: valueNode{
							valueType: a.valueType,
							value:     diff[0],
						},
					},
					boostNode: a.boostNode,
				}, b,
			)
		} else {
			return lfNodeIntersectLfNode(
				&TermsNode{
					fieldNode: a.fieldNode,
					valueType: a.valueType,
					boostNode: a.boostNode,
					terms:     diff,
				}, b,
			)
		}
	} else {
		if len(terms) == 0 {
			return nil, fmt.Errorf("failed to intersect %v and %v, err: value is conflict", a.ToDSL(), b.ToDSL())
		} else if len(terms) == 1 {
			return &TermNode{
				kvNode: kvNode{
					fieldNode: a.fieldNode,
					valueNode: valueNode{
						valueType: a.valueType,
						value:     terms[0],
					},
				},
				boostNode: a.boostNode,
			}, nil
		} else {
			return &TermsNode{
				fieldNode: a.fieldNode,
				valueType: a.valueType,
				boostNode: a.boostNode,
				terms:     terms,
			}, nil
		}
	}
}
