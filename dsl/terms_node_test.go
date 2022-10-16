package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestTermsNode(t *testing.T) {
	var node1 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.INTEGER_FIELD_TYPE, true), []LeafValue{1, 2}, WithBoost(1.2))
	var node2 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.INTEGER_FIELD_TYPE, true), []LeafValue{1, 2}, WithBoost(1.1))

	_, err := node1.InterSect(node2)
	assert.NotNil(t, err)
	_, err = node1.UnionJoin(node2)
	assert.NotNil(t, err)

	assert.Equal(t, DSL{"terms": DSL{"foo": termsToPrintValue(node1.terms, node1.mType), "boost": 1.2}}, node1.ToDSL())
	assert.Equal(t, TERMS_DSL_TYPE, node1.DslType())
	node3, _ := node1.Inverse()
	assert.Equal(t, &NotNode{
		opNode: opNode{filterCtxNode: node1.filterCtxNode},
		Nodes: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: node1.fieldNode,
						valueNode: valueNode{
							valueType: node1.valueType,
							value:     1,
						},
					},
					boostNode: node1.boostNode,
				},
				&TermNode{
					kvNode: kvNode{
						fieldNode: node1.fieldNode,
						valueNode: valueNode{
							valueType: node1.valueType,
							value:     2,
						},
					},
					boostNode: node1.boostNode,
				},
			},
		},
	}, node3)
}

func TestTermsNodeMergeTermsNode(t *testing.T) {
	var node1 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.INTEGER_FIELD_TYPE, true), []LeafValue{1, 2, 3})
	var node2 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.INTEGER_FIELD_TYPE, true), []LeafValue{1, 2, 3})
	var node3 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.INTEGER_FIELD_TYPE, true), []LeafValue{2, 3, 4})
	var node4 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.INTEGER_FIELD_TYPE, true), []LeafValue{4, 5, 6})

	node5, err := node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node5)

	node5, err = node1.InterSect(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node5)

	node5, err = node1.UnionJoin(node3)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: node1.fieldNode,
		boostNode: node1.boostNode,
		valueType: node1.valueType,
		terms:     []LeafValue{1, 2, 3, 4},
	}, node5)

	node5, err = node1.InterSect(node3)
	assert.Nil(t, err)
	assert.Equal(t, &AndNode{
		MustNodes: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: node1.fieldNode,
						valueNode: valueNode{
							valueType: node1.valueType,
							value:     1,
						},
					},
					boostNode: node1.boostNode,
				},
				&TermsNode{
					fieldNode: node1.fieldNode,
					boostNode: node1.boostNode,
					valueType: node1.valueType,
					terms:     []LeafValue{2, 3, 4},
				},
			},
		},
	}, node5)

	node5, err = node1.UnionJoin(node4)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: node1.fieldNode,
		boostNode: node1.boostNode,
		valueType: node1.valueType,
		terms:     []LeafValue{1, 2, 3, 4, 5, 6},
	}, node5)

	node5, err = node1.InterSect(node4)
	assert.Nil(t, err)
	assert.Equal(t, &AndNode{
		MustNodes: map[string][]AstNode{
			"foo": {
				&TermsNode{
					fieldNode: node1.fieldNode,
					boostNode: node1.boostNode,
					valueType: node1.valueType,
					terms:     []LeafValue{1, 2, 3},
				},
				&TermsNode{
					fieldNode: node1.fieldNode,
					boostNode: node1.boostNode,
					valueType: node1.valueType,
					terms:     []LeafValue{4, 5, 6},
				},
			},
		},
	}, node5)
}
