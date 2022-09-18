package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	var opNode = &OpNode{}
	assert.Equal(t, OP_NODE_TYPE, opNode.AstType())
	var lfNode = &LfNode{}
	assert.Equal(t, LEAF_NODE_TYPE, lfNode.AstType())
	var kvNode = &KvNode{Field: "foo "}
	assert.Equal(t, LEAF_NODE_TYPE, kvNode.AstType())
	assert.Equal(t, "LEAF:foo", kvNode.NodeKey())
}
