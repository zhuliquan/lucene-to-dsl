package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
)

func TestFuzzyNode(t *testing.T) {
	var node1 = NewFuzzyNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("bar", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
		),
		WithRewrite(CONSTANT_SCORE_BOOLEAN),
		WithMaxExpands(30),
		WithFuzziness("AUTO:1,3"),
		WithPrefixLength(1),
		WithTranspositions(false),
	)

	var node2 = NewExistsNode(NewFieldNode(NewLfNode(), "foo"))

	var node3, _ = node1.UnionJoin(node2)
	assert.Equal(t, node2, node3)

	var node4, _ = node1.InterSect(node2)
	assert.Equal(t, node1, node4)

	var node5, _ = node1.Inverse()
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			node4.NodeKey(): {node1},
		},
	}, node5)

	assert.Equal(t, FUZZY_DSL_TYPE, node1.DslType())
	assert.Equal(t, DSL{
		"fuzzy": DSL{
			"foo": DSL{
				"value":          "bar",
				"fuzziness":      "AUTO:1,3",
				"rewrite":        CONSTANT_SCORE_BOOLEAN,
				"prefix_length":  1,
				"max_expansions": 30,
				"transpositions": false,
			},
		},
	}, node1.ToDSL())

}
