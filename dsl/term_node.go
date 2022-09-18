package dsl

import "fmt"

type TermNode struct {
	KvNode
	Boost float64
}

func (n *TermNode) getBoost() float64 {
	return n.Boost
}

func (n *TermNode) DslType() DslType {
	return TERM_DSL_TYPE
}

func (n *TermNode) UnionJoin(o AstNode) (AstNode, error) {
	if n == nil && o == nil {
		return nil, ErrUnionJoinNilNode
	} else if n == nil && o != nil {
		return o, nil
	} else if n != nil && o == nil {
		return n, nil
	}

	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	case TERM_DSL_TYPE:
		var t = o.(*TermNode)
		if CompareAny(n.Value, t.Value, n.Type) == 0 {
			return n, nil
		} else {
			return &TermsNode{
				KvNode: n.KvNode,
				Values: []LeafValue{n.Value, t.Value},
				Boost:  n.Boost,
			}, nil
		}
	case TERMS_DSL_TYPE:
		var t = o.(*TermsNode)
		t.Values = UnionJoinValueLst(t.Values, []LeafValue{n.Value}, t.Type)
		return t, nil

	case RANGE_DSL_TYPE:
		// put logic of compare and collision into range node
		return o.(*RangeNode).UnionJoin(n)
	case QUERY_STRING_DSL_TYPE:
		var t = o.(*QueryStringNode)
		return t, nil
	case IDS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())

	default:
		return nil, fmt.Errorf("failed to union join %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermNode) InterSect(o AstNode) (AstNode, error) {
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
	case TERM_DSL_TYPE:
		var t = o.(*TermNode)
		if CompareAny(t.Value, n.Value, n.Type) == 0 {
			return o, nil
		} else {
			return &QueryStringNode{KvNode: KvNode{Field: n.Field, Value: fmt.Sprintf("%v AND %v", n.Value, t.Value)}, Boost: n.Boost}, nil
		}
	case RANGE_DSL_TYPE:
		return o.(*RangeNode).UnionJoin(n)
	case IDS_DSL_TYPE, TERMS_DSL_TYPE, PREFIX_DSL_TYPE, WILDCARD_DSL_TYPE, FUZZY_DSL_TYPE, REGEXP_DSL_TYPE, MATCH_DSL_TYPE, MATCH_PHRASE_DSL_TYPE, QUERY_STRING_DSL_TYPE:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is conflict", n.ToDSL(), o.ToDSL())
	default:
		return nil, fmt.Errorf("failed to intersect %s and %s, err: term type is unknown", n.ToDSL(), o.ToDSL())
	}
}

func (n *TermNode) Inverse() (AstNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return &NotNode{
		Nodes: map[string][]AstNode{
			n.Field: {n},
		},
	}, nil
}

func (n *TermNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"term": DSL{n.Field: DSL{"value": n.Value, "boost": n.Boost}}}
}
