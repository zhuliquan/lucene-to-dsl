package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
)

func TestIdsNode(t *testing.T) {
	var node1 = &IdsNode{
		ids: []string{"1", "2"},
	}
	var node2 = &IdsNode{
		ids: []string{"2", "3"},
	}

	assert.Equal(t, LEAF_NODE_TYPE, node1.AstType())
	assert.Equal(t, IDS_DSL_TYPE, node1.DslType())
	var node3, err = node1.InterSect(node2)
	assert.Equal(t, &IdsNode{ids: []string{"2"}}, node3)
	assert.Nil(t, err)
	node3, err = node1.UnionJoin(node2)
	assert.Equal(t, &IdsNode{ids: []string{"1", "2", "3"}}, node3)
	assert.Nil(t, err)

	var node4 = &MatchNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{
					filterCtxNode: filterCtxNode{
						filterCtx: false,
					},
				},
				field: "foo",
			},
			valueNode: valueNode{
				valueType: valueType{
					aType: false,
					mType: mapping.TEXT_FIELD_TYPE,
				},
				value: "bar",
			},
		},
	}
	node3, err = node1.InterSect(node4)
	assert.NotNil(t, err)
	assert.Nil(t, node3)

	node3, err = node1.UnionJoin(node4)
	assert.NotNil(t, err)
	assert.Nil(t, node3)

	node3, err = node1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"_id": {node1},
		},
	}, node3)

	assert.Equal(t, "_id", node1.NodeKey())
	assert.Equal(t, DSL{"ids": DSL{"values": []string{"1", "2"}}}, node1.ToDSL())

}
