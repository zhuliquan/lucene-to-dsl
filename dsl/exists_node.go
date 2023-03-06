package dsl

type ExistsNode struct {
	fieldNode
}

func (n *ExistsNode) DslType() DslType {
	return EXISTS_DSL_TYPE
}

func NewExistsNode(fieldNode *fieldNode) *ExistsNode {
	return &ExistsNode{
		fieldNode: *fieldNode,
	}
}

// if union same field node, you can return exist node, for example {"exists": {"field" : "x"}} union {"match": {"x": "foo bar"}}
// "exists": {"field": "x"} > "match": {"x": "foo bar"}
func (n *ExistsNode) UnionJoin(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.UnionJoin(n)
	}
	switch o.DslType() {
	default:
		return n, nil
	}
}

func (n *ExistsNode) InterSect(o AstNode) (AstNode, error) {
	if checkCommonDslType(o.DslType()) {
		return o.InterSect(n)
	}
	switch o.DslType() {
	default:
		return o, nil
	}
}

func (n *ExistsNode) Inverse() (AstNode, error) {
	return inverseNode(n), nil
}

func (n *ExistsNode) ToDSL() DSL {
	return DSL{
		EXISTS_KEY: DSL{
			FIELD_KEY: n.field,
		},
	}
}
