package dsl

import "fmt"

type WildCardNode struct {
	LfNode
	KvNode
	Boost float64
}

func (n *WildCardNode) getBoost() float64 {
	return n.Boost
}

func (n *WildCardNode) DslType() DslType {
	return WILDCARD_DSL_TYPE
}

func (n *WildCardNode) UnionJoin(o AstNode) (AstNode, error) {
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
	default:
		return nil, fmt.Errorf("failed to union join wildcard node")
	}
}

func (n *WildCardNode) InterSect(o AstNode) (AstNode, error) {
	if n == nil || o == nil {
		return nil, ErrIntersectNilNode
	}
	switch o.DslType() {
	case EXISTS_DSL_TYPE:
		return o.UnionJoin(n)
	default:
		return nil, fmt.Errorf("failed to intersect wildcard node")
	}
}

func (n *WildCardNode) Inverse() (AstNode, error) {
	if n == nil {
		return nil, ErrInverseNilNode
	}
	return nil, fmt.Errorf("failed to inverse wildcard node")
}

func (n *WildCardNode) ToDSL() DSL {
	if n == nil {
		return EmptyDSL
	}
	return DSL{"wildcard": DSL{n.Field: DSL{"values": n.Value, "boost": n.Boost}}}

}
