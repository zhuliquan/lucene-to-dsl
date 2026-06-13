package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
	"github.com/zhuliquan/lucene-to-dsl/utils"
)

func TestWildcardNode(t *testing.T) {
	pattern := utils.NewWildCardPattern("a?b*")
	node1 := NewWildCardNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("a?b*", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
		),
		pattern, WithBoost(1.2),
	)
	node2 := NewWildCardNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("a?b*", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
		),
		pattern, WithBoost(1.5),
	)

	t.Run("test_wildcard_node", func(t *testing.T) {
		assert.Equal(t, WILDCARD_DSL_TYPE, node1.DslType())
		assert.Equal(t, DSL{
			WILDCARD_KEY: DSL{
				"foo": DSL{
					VALUE_KEY:   "a?b*",
					BOOST_KEY:   1.2,
					REWRITE_KEY: CONSTANT_SCORE,
				},
			},
		}, node1.ToDSL())
	})

	t.Run("test_inverse_wildcard_node", func(t *testing.T) {
		res, err := node1.Inverse()
		assert.Nil(t, err)
		assert.Equal(t, &BoolNode{
			opNode: opNode{opType: NOT},
			MustNot: map[string][]AstNode{
				"foo": {node1},
			},
		}, res)
	})

	t.Run("test_intersect_wildcard_node", func(t *testing.T) {
		res, err := node1.InterSect(node2)
		assert.Nil(t, err)
		assert.Equal(t, &BoolNode{
			opNode: opNode{opType: AND},
			Must: map[string][]AstNode{
				"foo": {node1, node2},
			},
		}, res)
	})

	t.Run("test_union_join_wildcard_node", func(t *testing.T) {
		res, err := node1.UnionJoin(node2)
		assert.Nil(t, err)
		assert.Equal(t, &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {node1, node2},
			},
			MinimumShouldMatch: 1,
		}, res)
	})
}

func TestWildcardNodeMergeTermNode(t *testing.T) {
	var pattern = utils.NewWildCardPattern("a?b*")
	var n1 = NewWildCardNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("a?b*", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
		),
		pattern,
	)

	var n2 = NewWildCardNode(
		NewKVNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueNode("a?b*", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
		),
		pattern,
	)

	var n3 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	))

	var n4 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("abbb", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	))

	var n5 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	))

	var n6 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("abbb", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	))

	n7, err := n1.UnionJoin(n3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, n3},
		},
		MinimumShouldMatch: 1,
	}, n7)

	n7, err = n1.UnionJoin(n4)
	assert.Nil(t, err)
	assert.Equal(t, n1, n7)

	n7, err = n1.InterSect(n3)
	assert.NotNil(t, err)
	assert.Equal(t, nil, n7)

	n7, err = n1.InterSect(n4)
	assert.Nil(t, err)
	assert.Equal(t, n4, n7)

	n7, err = n2.UnionJoin(n5)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n2, n5},
		},
		MinimumShouldMatch: 1,
	}, n7)

	n7, err = n2.UnionJoin(n6)
	assert.Nil(t, err)
	assert.Equal(t, n2, n7)

	n7, err = n2.InterSect(n5)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {n2, n5},
		},
	}, n7)

	n7, err = n2.InterSect(n6)
	assert.Nil(t, err)
	assert.Equal(t, n6, n7)
}

func TestWildcardNodeMergeWildcardNode(t *testing.T) {
	p1 := utils.NewWildCardPattern("aab*")
	p2 := utils.NewWildCardPattern("a?c")
	var n1 = NewWildCardNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("aab*", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), p1)

	var n2 = NewWildCardNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a?c", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), p2)

	var n3 = NewWildCardNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("aab*", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), p1)

	var n4 = NewWildCardNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a?c", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), p2)

	n5, err := n1.UnionJoin(n1)
	assert.Nil(t, err)
	assert.Equal(t, n1, n5)

	n5, err = n1.UnionJoin(n2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, n2},
		},
		MinimumShouldMatch: 1,
	}, n5)

	n5, err = n1.InterSect(n1)
	assert.Nil(t, err)
	assert.Equal(t, n1, n5)

	n5, err = n1.InterSect(n2)
	assert.NotNil(t, err)
	assert.Nil(t, n5)

	n5, err = n3.UnionJoin(n3)
	assert.Nil(t, err)
	assert.Equal(t, n3, n5)

	n5, err = n3.UnionJoin(n4)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n3, n4},
		},
		MinimumShouldMatch: 1,
	}, n5)

	n5, err = n3.InterSect(n3)
	assert.Nil(t, err)
	assert.Equal(t, n3, n5)

	n5, err = n3.InterSect(n4)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {n3, n4},
		},
	}, n5)
}
