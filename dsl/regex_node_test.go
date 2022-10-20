package dsl

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestRegexNode(t *testing.T) {
	pattern := regexp.MustCompile("^[1-5]{1,9}")
	var node1 = NewRegexNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("^[1-5]{1,9}", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
		),
		pattern,
		WithRewrite(SCORING_BOOLEAN),
		WithMaxDeterminizedStates(10),
		WithFlags(COMPLEMENT_FLAG),
	)

	assert.Equal(t, REGEXP_DSL_TYPE, node1.DslType())
	assert.Equal(t, DSL{"regexp": DSL{
		"foo": DSL{
			"value":                   "^[1-5]{1,9}",
			"rewrite":                 SCORING_BOOLEAN,
			"max_determinized_states": 10,
			"flags":                   COMPLEMENT_FLAG,
		},
	}}, node1.ToDSL())
	node2, _ := node1.Inverse()
	assert.Equal(t, &NotNode{
		opNode: opNode{filterCtxNode: node1.filterCtxNode},
		Nodes: map[string][]AstNode{
			"foo": {node1},
		},
	}, node2)
}

// func TestRegexNode(t *testing.T) {
// 	var node1 = &RegexpNode{
// 		KvNode: KvNode{
// 			Field: "foo",
// 			Type:  mapping.TEXT_FIELD_TYPE,
// 			Value: "^[1-5]{1,9}",
// 		},
// 	}
// 	var node2 = &ExistsNode{
// 		KvNode: KvNode{
// 			Field: "foo",
// 			Type:  mapping.TEXT_FIELD_TYPE,
// 		},
// 	}
// 	assert.Equal(t, REGEXP_DSL_TYPE, node1.DslType())
// 	assert.Equal(t, LEAF_NODE_TYPE, node1.AstType())
// 	var node3, err = node1.UnionJoin(node2)
// 	assert.Nil(t, err)
// 	assert.Equal(t, node2, node3)

// 	node3, err = node1.InterSect(node2)
// 	assert.Nil(t, err)
// 	assert.Equal(t, node1, node3)

// 	var node4 = &MatchNode{
// 		KvNode: KvNode{
// 			Field: "foo",
// 			Type:  mapping.TEXT_FIELD_TYPE,
// 			Value: "bar",
// 		},
// 	}

// 	node3, err = node1.UnionJoin(node4)
// 	assert.Nil(t, err)
// 	assert.Equal(t, &OrNode{
// 		MinimumShouldMatch: 1,
// 		Nodes: map[string][]AstNode{
// 			node1.NodeKey(): {node1, node4},
// 		},
// 	}, node3)

// 	node3, err = node1.InterSect(node4)
// 	assert.Nil(t, err)
// 	assert.Equal(t, &AndNode{
// 		MustNodes: map[string][]AstNode{
// 			node1.NodeKey(): {node1, node4},
// 		},
// 	}, node3)

// 	node3, err = node1.Inverse()
// 	assert.Nil(t, err)
// 	assert.Equal(t, &NotNode{
// 		Nodes: map[string][]AstNode{
// 			node1.NodeKey(): {node1},
// 		},
// 	}, node3)

// 	assert.Equal(t, DSL{"regexp": DSL{"foo": DSL{"value": "^[1-5]{1,9}"}}}, node1.ToDSL())
// 	assert.Equal(t, "LEAF:foo", node1.NodeKey())

// 	var node5 = &RegexpNode{
// 		KvNode:  KvNode{Field: "foo", Value: "^[1-5]{1,9}"},
// 		Rewrite: "constant_score",
// 	}
// 	assert.Equal(t, DSL{"regexp": DSL{"foo": DSL{"value": "^[1-5]{1,9}", "rewrite": "constant_score"}}}, node5.ToDSL())
// }
