package dsl

type EmptyNode struct{}

func (n *EmptyNode) AstType() AstType                   { return EMPTY_NODE_TYPE }
func (n *EmptyNode) DslType() DslType                   { return EMPTY_DSL_TYPE }
func (n *EmptyNode) NodeKey() string                    { return "" }
func (n *EmptyNode) UnionJoin(AstNode) (AstNode, error) { return n, nil }
func (n *EmptyNode) InterSect(AstNode) (AstNode, error) { return n, nil }
func (n *EmptyNode) Inverse() (AstNode, error)          { return n, nil }
func (n *EmptyNode) ToDSL() DSL                         { return EmptyDSL }
