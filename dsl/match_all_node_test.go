package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchAllNode(t *testing.T) {
	var matchAll = &MatchAllNode{}
	assert.Equal(t, LEAF_NODE_TYPE, matchAll.AstType())
	assert.Equal(t, MATCH_ALL_DSL_TYPE, matchAll.DslType())
	var otherNode1, _ = matchAll.UnionJoin(&EmptyNode{})
	assert.Equal(t, matchAll, otherNode1)
	var otherNode2, _ = matchAll.InterSect(&EmptyNode{})
	assert.Equal(t, &EmptyNode{}, otherNode2)
	var otherNode3, _ = matchAll.Inverse()
	assert.Equal(t, &EmptyNode{}, otherNode3)
	assert.Equal(t, "*", matchAll.NodeKey())
	assert.Equal(t, DSL{"match_all": DSL{}}, matchAll.ToDSL())
}
