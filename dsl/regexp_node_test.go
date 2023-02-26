package dsl

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestRegexpNode(t *testing.T) {
	pattern := regexp.MustCompile("^[1-5]{1,9}")
	var node1 = NewRegexpNode(
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
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo": {node1},
		},
	}, node2)
}

func TestRegexpNodeMergeTermNode(t *testing.T) {
	var pattern = regexp.MustCompile(`^a\wb\w`)
	var n1 = NewRegexpNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("^a\\wb\\w", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), pattern)

	var n2 = NewRegexpNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("^a\\wb\\w", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), pattern)

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

func TestRegexpNodeMergeRegexpNode(t *testing.T) {
	p1 := regexp.MustCompile(`aab\w+`)
	p2 := regexp.MustCompile(`a\wc`)
	var n1 = NewRegexpNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("aab\\w+", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), p1)

	var n2 = NewRegexpNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a\\wc", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), p2)

	var n3 = NewRegexpNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("aab\\w+", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), p1)

	var n4 = NewRegexpNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("a\\wc", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
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
		minimumShouldMatch: 1,
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
		minimumShouldMatch: 1,
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
