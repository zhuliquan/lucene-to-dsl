package dsl

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/zhuliquan/lucene-to-dsl/mapping"
// )

// func TestNode(t *testing.T) {
// 	var opNode = NewOpNode(true)
// 	assert.Equal(t, OP_NODE_TYPE, opNode.AstType())
// 	assert.Equal(t, true, opNode.Is())
// 	var lfNode = lfNode{true, false}
// 	assert.Equal(t, LEAF_NODE_TYPE, lfNode.AstType())
// 	assert.Equal(t, true, lfNode.NeedFilter())
// 	assert.Equal(t, false, lfNode.IsArray())
// 	var value = NewValue("bar", mapping.KEYWORD_FIELD_TYPE)
// 	var kvNode = NewKVNode("foo", WithValue(value), WithFilter(false), WithIsList(false))
// 	assert.Equal(t, LEAF_NODE_TYPE, kvNode.AstType())
// 	assert.Equal(t, "foo", kvNode.NodeKey())
// }
