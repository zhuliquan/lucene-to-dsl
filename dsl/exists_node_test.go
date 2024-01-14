package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
)

func TestExistNode(t *testing.T) {
	var node1 = NewExistsNode(NewFieldNode(NewLfNode(), "foo"))
	var node2 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))), WithBoost(1.2))
	assert.Equal(t, LEAF_NODE_TYPE, node1.AstType())
	assert.Equal(t, EXISTS_DSL_TYPE, node1.DslType())

	node3, err := node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node3)

	node3, err = node1.InterSect(node2)
	assert.Nil(t, err)
	assert.Equal(t, node2, node3)

	node3, err = node1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo": {node1},
		},
	}, node3)

	assert.Equal(t, DSL{"exists": DSL{"field": "foo"}}, node1.ToDSL())

}
