package dsl

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestEmptyNode(t *testing.T) {
// 	var emptyNode = &EmptyNode{}
// 	assert.Equal(t, EMPTY_NODE_TYPE, emptyNode.AstType())
// 	assert.Equal(t, EMPTY_DSL_TYPE, emptyNode.DslType())
// 	var otherNode1, _ = emptyNode.UnionJoin(emptyNode)
// 	assert.Equal(t, emptyNode, otherNode1)
// 	var otherNode2, _ = emptyNode.InterSect(emptyNode)
// 	assert.Equal(t, emptyNode, otherNode2)
// 	var otherNode3, _ = emptyNode.Inverse()
// 	assert.Equal(t, emptyNode, otherNode3)
// 	assert.Equal(t, "", emptyNode.NodeKey())
// 	assert.Equal(t, EmptyDSL, emptyNode.ToDSL())
// }
