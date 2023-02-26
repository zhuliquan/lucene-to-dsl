package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	"github.com/zhuliquan/lucene-to-dsl/utils"
)

func TestPrefixNode(t *testing.T) {
	var n1 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), utils.NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	assert.Equal(t, PREFIX_DSL_TYPE, n1.DslType())
	assert.Equal(t, DSL{"prefix": DSL{"foo": DSL{"value": "ab", "rewrite": CONSTANT_SCORE_BOOLEAN}}}, n1.ToDSL())
	n2, _ := n1.Inverse()
	assert.Equal(t, &BoolNode{
		opNode:  opNode{opType: NOT},
		MustNot: map[string][]AstNode{"foo": {n1}},
	}, n2)
}

func TestPrefixNodeMergeExistNode(t *testing.T) {
	var node1 = NewExistsNode(NewFieldNode(NewLfNode(), "foo"))
	var node2 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), utils.NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	node3, err := node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node3)

	node3, err = node1.InterSect(node2)
	assert.Nil(t, err)
	assert.Equal(t, node2, node3)
}

func TestPrefixNodeMergeTermNode(t *testing.T) {
	var n1 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), utils.NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n2 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), utils.NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n3 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), WithBoost(1.2))

	var n4 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("abc", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), WithBoost(1.2))

	var n5 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), WithBoost(1.2))

	var n6 = NewTermNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("abc", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), WithBoost(1.2))

	n7, err := n1.UnionJoin(n3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, n3},
		},
		minimumShouldMatch: 1,
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
		minimumShouldMatch: 1,
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

func TestPrefixNodeIntersectPrefixNode(t *testing.T) {
	var n1 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), utils.NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n2 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("abc", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), utils.NewPrefixPattern("abc"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n3 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ed", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), utils.NewPrefixPattern("ed"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	n4, err := n1.UnionJoin(n2)
	assert.Nil(t, err)
	assert.Equal(t, n1, n4)

	n4, err = n2.UnionJoin(n1)
	assert.Nil(t, err)
	assert.Equal(t, n1, n4)

	n4, err = n1.UnionJoin(n3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, n3},
		},
		minimumShouldMatch: 1,
	}, n4)

	n4, err = n1.InterSect(n2)
	assert.Nil(t, err)
	assert.Equal(t, n2, n4)

	n4, err = n2.UnionJoin(n1)
	assert.Nil(t, err)
	assert.Equal(t, n1, n4)

	n4, err = n1.UnionJoin(n3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, n3},
		},
		minimumShouldMatch: 1,
	}, n4)

	var n5 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), utils.NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n6 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ed", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), utils.NewPrefixPattern("ed"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	n4, err = n5.UnionJoin(n6)
	assert.Nil(t, err)
	assert.NotNil(t, n4)

	n4, err = n5.InterSect(n6)
	assert.NotNil(t, err)
	assert.Nil(t, nil, n4)
}
