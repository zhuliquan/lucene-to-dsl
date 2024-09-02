package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
)

func TestMatchNode(t *testing.T) {
	var node1 = NewMatchNode(&kvNode{
		fieldNode: fieldNode{
			lfNode: lfNode{},
			field:  "foo",
		},
		valueNode: valueNode{valueType: valueType{mType: mapping.TEXT_FIELD_TYPE}, value: "bar2"},
	},
		WithBoost(1.4),
	)
	var node3 = NewQueryStringNode(&kvNode{
		fieldNode: fieldNode{
			lfNode: lfNode{},
			field:  "foo",
		},
		valueNode: valueNode{valueType: valueType{mType: mapping.TEXT_FIELD_TYPE}, value: "this AND that"},
	},
		WithBoost(1.4),
	)
	var node4 = &ExistsNode{
		fieldNode: fieldNode{
			lfNode: lfNode{},
			field:  "foo",
		},
	}

	var node5, err = node1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo": {node1},
		},
	}, node5)

	node5, err = node1.InterSect(node3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {node1, node3},
		},
	}, node5)

	node5, err = node1.InterSect(node4)
	assert.Nil(t, err)
	assert.Equal(t, node1, node5)

	node5, err = node1.UnionJoin(node3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {node1, node3},
		},
		MinimumShouldMatch: 1,
	}, node5)

	node5, err = node1.UnionJoin(node4)
	assert.Nil(t, err)
	assert.Equal(t, node4, node5)

	assert.Equal(t, "foo", node1.NodeKey())
	assert.Equal(t, MATCH_DSL_TYPE, node1.DslType())
	assert.Equal(t, LEAF_NODE_TYPE, node1.AstType())
	assert.Equal(t, DSL{
		"match": DSL{
			node1.field: DSL{
				"query":          node1.getValue(),
				"boost":          node1.getBoost(),
				"max_expansions": 50,
			},
		},
	}, node1.ToDSL())

}
