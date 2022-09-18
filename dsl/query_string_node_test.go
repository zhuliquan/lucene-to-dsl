package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestQueryStringNode(t *testing.T) {
	var node1 = &QueryStringNode{
		KvNode: KvNode{Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "this AND that"},
		Boost:  1.4,
	}
	var node2 = &MatchNode{
		KvNode: KvNode{Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1"},
		Boost:  1.3,
	}
	var node3 = &MatchNode{
		KvNode: KvNode{Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar2"},
		Boost:  1.4,
	}
	var node4 = &ExistsNode{KvNode: KvNode{Field: "foo"}}

	var node5, err = node1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &NotNode{Nodes: map[string][]AstNode{node1.NodeKey(): {node1}}}, node5)

	node5, err = node1.InterSect(node2)
	assert.NotNil(t, err)
	assert.Nil(t, node5)

	node5, err = node1.InterSect(node3)
	assert.Nil(t, err)
	assert.Equal(t, &AndNode{MustNodes: map[string][]AstNode{node1.NodeKey(): {node1, node3}}}, node5)

	node5, err = node1.InterSect(node4)
	assert.Nil(t, err)
	assert.Equal(t, node1, node5)

	node5, err = node1.UnionJoin(node2)
	assert.NotNil(t, err)
	assert.Nil(t, node5)

	node5, err = node1.UnionJoin(node3)
	assert.Nil(t, err)
	assert.Equal(t, &OrNode{MinimumShouldMatch: 1, Nodes: map[string][]AstNode{node1.NodeKey(): {node1, node3}}}, node5)

	node5, err = node1.UnionJoin(node4)
	assert.Nil(t, err)
	assert.Equal(t, node4, node5)

	assert.Equal(t, QUERY_STRING_DSL_TYPE, node1.DslType())
	assert.Equal(t, LEAF_NODE_TYPE, node1.AstType())
	assert.Equal(t, DSL{
		"query_string": DSL{
			"query":         node1.Value,
			"default_field": node1.Field,
			"boost":         node1.Boost,
		},
	}, node1.ToDSL())
}
