package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestPrefixNode(t *testing.T) {
	var n1 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

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
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

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
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n2 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

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

func TestPrefixNodeMergeTermsNode(t *testing.T) {
	var n1 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n2 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n3 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, false),
		[]LeafValue{"abc"},
		WithBoost(1.2),
	)

	var n4 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, false),
		[]LeafValue{"abc", "a"},
		WithBoost(1.2),
	)

	var n5 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, false),
		[]LeafValue{"abc", "a", "b"},
		WithBoost(1.2),
	)

	var n6 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, true),
		[]LeafValue{"abc"},
		WithBoost(1.2),
	)

	var n7 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, true),
		[]LeafValue{"abc", "a"},
		WithBoost(1.2),
	)

	var n8 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, true),
		[]LeafValue{"abc", "a", "b"},
		WithBoost(1.2),
	)

	n9, err := n1.UnionJoin(n3)
	assert.Nil(t, err)
	assert.Equal(t, n1, n9)

	n9, err = n1.UnionJoin(n4)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, &TermNode{
				kvNode: kvNode{
					fieldNode: n4.fieldNode,
					valueNode: valueNode{
						valueType: n4.valueType,
						value:     "a",
					},
				},
				boostNode: n4.boostNode,
			}},
		},
		MinimumShouldMatch: 1,
	}, n9)

	n9, err = n1.UnionJoin(n5)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n1, &TermsNode{
				fieldNode: n5.fieldNode,
				valueType: n5.valueType,
				terms:     []LeafValue{"a", "b"},
				boostNode: n5.boostNode,
			}},
		},
		MinimumShouldMatch: 1,
	}, n9)

	n9, err = n2.UnionJoin(n6)
	assert.Nil(t, err)
	assert.Equal(t, n2, n9)

	n9, err = n2.UnionJoin(n7)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n2, &TermNode{
				kvNode: kvNode{
					fieldNode: n7.fieldNode,
					valueNode: valueNode{
						valueType: n7.valueType,
						value:     "a",
					},
				},
				boostNode: n7.boostNode,
			}},
		},
		MinimumShouldMatch: 1,
	}, n9)

	n9, err = n2.UnionJoin(n8)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n2, &TermsNode{
				fieldNode: n8.fieldNode,
				valueType: n8.valueType,
				terms:     []LeafValue{"a", "b"},
				boostNode: n8.boostNode,
			}},
		},
		MinimumShouldMatch: 1,
	}, n9)

	var n10 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, false),
		[]LeafValue{"a"},
		WithBoost(1.2),
	)

	var n11 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, false),
		[]LeafValue{"a", "abc"},
		WithBoost(1.2),
	)

	var n12 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, false),
		[]LeafValue{"a", "abc", "abd"},
		WithBoost(1.2),
	)

	var n13 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, true),
		[]LeafValue{"a"},
		WithBoost(1.2),
	)

	var n14 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, true),
		[]LeafValue{"a", "abc"},
		WithBoost(1.2),
	)

	var n15 = NewTermsNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueType(mapping.TEXT_FIELD_TYPE, true),
		[]LeafValue{"a", "abc", "abd"},
		WithBoost(1.2),
	)

	n9, err = n1.InterSect(n10)
	assert.NotNil(t, err)
	assert.Equal(t, nil, n9)

	n9, err = n1.InterSect(n11)
	assert.Nil(t, err)
	assert.Equal(t, &TermNode{
		kvNode: kvNode{
			fieldNode: n11.fieldNode,
			valueNode: valueNode{
				valueType: n1.valueType,
				value:     "abc",
			},
		},
		boostNode: n11.boostNode,
	}, n9)

	n9, err = n1.InterSect(n12)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: n12.fieldNode,
		valueType: n12.valueType,
		terms:     []LeafValue{"abc", "abd"},
		boostNode: n12.boostNode,
	}, n9)

	n9, err = n1.InterSect(n3)
	assert.Nil(t, err)
	assert.Equal(t, &TermNode{
		kvNode: kvNode{
			fieldNode: n3.fieldNode,
			valueNode: valueNode{
				valueType: n3.valueType,
				value:     "abc",
			},
		},
		boostNode: n3.boostNode,
	}, n9)

	n9, err = n1.InterSect(n4)
	assert.Nil(t, err)
	assert.Equal(t, &TermNode{
		kvNode: kvNode{
			fieldNode: n4.fieldNode,
			valueNode: valueNode{
				valueType: n4.valueType,
				value:     "abc",
			},
		},
		boostNode: n3.boostNode,
	}, n9)

	n9, err = n1.InterSect(n4)
	assert.Nil(t, err)
	assert.Equal(t, &TermNode{
		kvNode: kvNode{
			fieldNode: n4.fieldNode,
			valueNode: valueNode{
				valueType: n4.valueType,
				value:     "abc",
			},
		},
		boostNode: n3.boostNode,
	}, n9)

	/////
	n9, err = n2.InterSect(n13)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				n2, &TermNode{
					kvNode: kvNode{
						fieldNode: n13.fieldNode,
						valueNode: valueNode{
							valueType: n13.valueType,
							value:     "a",
						},
					},
					boostNode: n13.boostNode,
				},
			},
		},
	}, n9)

	n9, err = n2.InterSect(n14)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				n2, &TermNode{
					kvNode: kvNode{
						fieldNode: n14.fieldNode,
						valueNode: valueNode{
							valueType: n14.valueType,
							value:     "a",
						},
					},
					boostNode: n14.boostNode,
				},
			},
		},
	}, n9)

	n9, err = n2.InterSect(n15)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				n2, &TermNode{
					kvNode: kvNode{
						fieldNode: n15.fieldNode,
						valueNode: valueNode{
							valueType: n15.valueType,
							value:     "a",
						},
					},
					boostNode: n15.boostNode,
				},
			},
		},
	}, n9)

	n9, err = n2.InterSect(n6)
	assert.Nil(t, err)
	assert.Equal(t, n2, n9)

	n9, err = n2.InterSect(n7)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				n2, &TermNode{
					kvNode: kvNode{
						fieldNode: n7.fieldNode,
						valueNode: valueNode{
							valueType: n7.valueType,
							value:     "a",
						},
					},
					boostNode: n7.boostNode,
				},
			},
		},
	}, n9)

	n9, err = n2.InterSect(n8)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				n2, &TermsNode{
					fieldNode: n8.fieldNode,
					valueType: n8.valueType,
					terms:     []LeafValue{"a", "b"},
					boostNode: n8.boostNode,
				},
			},
		},
	}, n9)
}

func TestPrefixNodeIntersectPrefixNode(t *testing.T) {
	var n1 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n2 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("abc", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), NewPrefixPattern("abc"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n3 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ed", NewValueType(mapping.TEXT_FIELD_TYPE, true)),
	), NewPrefixPattern("ed"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

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
		MinimumShouldMatch: 1,
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
		MinimumShouldMatch: 1,
	}, n4)

	var n5 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ab", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), NewPrefixPattern("ab"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	var n6 = NewPrefixNode(NewKVNode(
		NewFieldNode(NewLfNode(), "foo"),
		NewValueNode("ed", NewValueType(mapping.TEXT_FIELD_TYPE, false)),
	), NewPrefixPattern("ed"), WithRewrite(CONSTANT_SCORE_BOOLEAN))

	n4, err = n5.UnionJoin(n6)
	assert.Nil(t, err)
	assert.NotNil(t, n4)

	n4, err = n5.InterSect(n6)
	assert.NotNil(t, err)
	assert.Nil(t, nil, n4)
}
