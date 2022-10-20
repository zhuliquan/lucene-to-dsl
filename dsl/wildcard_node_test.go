package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestWildcardNode(t *testing.T) {
	pattern := NewWildCardPattern("a?b*")
	node1 := NewWildCardNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("a?b*", NewValueType(mapping.TEXT_FIELD_TYPE, false))),
		pattern, WithBoost(1.2),
	)
	node2, _ := node1.Inverse()
	assert.Equal(t, &NotNode{
		opNode: opNode{filterCtxNode: node1.filterCtxNode},
		Nodes: map[string][]AstNode{
			"foo": {node1},
		},
	}, node2)
	assert.Equal(t, WILDCARD_DSL_TYPE, node1.DslType())
	assert.Equal(t, DSL{"wildcard": DSL{"foo": DSL{"value": "a?b*", "boost": 1.2, "rewrite": CONSTANT_SCORE}}}, node1.ToDSL())
}
