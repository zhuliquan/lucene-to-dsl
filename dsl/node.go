package dsl

import "github.com/zhuliquan/lucene-to-dsl/mapping"

// define ast node of dsl
type AstNode interface {
	AstType() AstType
	DslType() DslType
	NodeKey() string

	// union_join / intersect / inverse nodes with same NodeKey (get by NodeKey() func)
	UnionJoin(AstNode) (AstNode, error)
	InterSect(AstNode) (AstNode, error)
	Inverse() (AstNode, error)
	ToDSL() DSL
}

type boostNode interface {
	getBoost() float64
}

type filterNode interface {
	NeedFilter() bool
}

type OpNode struct{}

func (n *OpNode) AstType() AstType {
	return OP_NODE_TYPE
}

// leaf node
type LfNode struct {
	Filter bool
}

func (n *LfNode) NeedFilter() bool {
	return n.Filter
}

func (n *LfNode) AstType() AstType {
	return LEAF_NODE_TYPE
}

// Key value node
type KvNode struct {
	LfNode
	Field string
	Type  mapping.FieldType
	Value LeafValue
}

func (n *KvNode) NodeKey() string {
	return "LEAF:" + n.Field
}
