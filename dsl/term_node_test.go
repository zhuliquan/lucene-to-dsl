package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestTermNode(t *testing.T) {
	var node0 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar3", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))))
	assert.Equal(t, TERM_DSL_TYPE, node0.DslType())
	assert.Equal(t, DSL{"term": DSL{"foo": DSL{"value": "bar3", "boost": 1.0}}}, node0.ToDSL())
	var node1 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))), WithBoost(1.2))
	var node2 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar2", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))), WithBoost(1.1))
	_, err := node1.InterSect(node2)
	assert.NotNil(t, err)

	_, err = node1.UnionJoin(node2)
	assert.NotNil(t, err)

}
func TestTermNodeMergeTermNode(t *testing.T) {
	var node1 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))))
	var node2 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))))
	var node3 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar2", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))))

	var node4 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, true))))
	var node5 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, true))))
	var node6 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar2", NewValueType(mapping.KEYWORD_FIELD_TYPE, true))))

	node8, err := node1.InterSect(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node8)

	node8, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node8)

	node8, err = node1.InterSect(node3)
	assert.NotNil(t, err)
	assert.Equal(t, nil, node8)

	node8, err = node1.UnionJoin(node3)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: node1.fieldNode,
		boostNode: node1.boostNode,
		valueType: node1.valueType,
		terms:     []LeafValue{"bar1", "bar2"},
	}, node8)

	node8, err = node4.InterSect(node5)
	assert.Nil(t, err)
	assert.Equal(t, node4, node8)

	node8, err = node4.UnionJoin(node5)
	assert.Nil(t, err)
	assert.Equal(t, node4, node8)

	node8, err = node4.InterSect(node6)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			node4.NodeKey(): {node4, node6},
		},
	}, node8)

	node8, err = node4.UnionJoin(node6)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: node4.fieldNode,
		boostNode: node4.boostNode,
		valueType: node4.valueType,
		terms:     []LeafValue{"bar1", "bar2"},
	}, node8)
}

func TestTermNodeMergeTermsNode(t *testing.T) {
	var node1 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false))))
	var node2 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, false), []LeafValue{"bar1"})
	var node3 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, false), []LeafValue{"bar2", "bar1"})
	var node4 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, false), []LeafValue{"bar2", "bar1", "bar3"})
	var node5 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, false), []LeafValue{"bar2", "bar3"})

	var node6 = NewTermNode(NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, true))))
	var node7 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, true), []LeafValue{"bar1"})
	var node8 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, true), []LeafValue{"bar2", "bar1"})
	var node9 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, true), []LeafValue{"bar2", "bar1", "bar3"})
	var node10 = NewTermsNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, true), []LeafValue{"bar2", "bar3"})

	node11, err := node1.InterSect(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node11)

	node11, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, node1, node11)

	node11, err = node1.InterSect(node3)
	assert.Nil(t, err)
	assert.Equal(t, node1, node11)

	node11, err = node1.UnionJoin(node3)
	assert.Nil(t, err)
	assert.Equal(t, node3, node11)

	node11, err = node1.InterSect(node4)
	assert.Nil(t, err)
	assert.Equal(t, node1, node11)

	node11, err = node1.UnionJoin(node4)
	assert.Nil(t, err)
	assert.Equal(t, node4, node11)

	node11, err = node1.InterSect(node5)
	assert.NotNil(t, err)
	assert.Equal(t, nil, node11)

	node11, err = node1.UnionJoin(node5)
	assert.Nil(t, err)
	assert.Equal(t, node4, node11)

	node11, err = node6.InterSect(node7)
	assert.Nil(t, err)
	assert.Equal(t, node6, node11)

	node11, err = node6.UnionJoin(node7)
	assert.Nil(t, err)
	assert.Equal(t, node6, node11)

	node11, err = node6.InterSect(node8)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				node6, &TermNode{
					kvNode: kvNode{
						fieldNode: node6.fieldNode,
						valueNode: valueNode{value: "bar2", valueType: node6.valueType},
					},
					boostNode: node6.boostNode,
				},
			},
		},
	}, node11)

	node11, err = node6.UnionJoin(node8)
	assert.Nil(t, err)
	assert.Equal(t, node8, node11)

	node11, err = node6.InterSect(node9)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				node6, &TermsNode{
					fieldNode: node6.fieldNode,
					valueType: node6.valueType,
					boostNode: node6.boostNode,
					terms:     []LeafValue{"bar2", "bar3"},
				},
			},
		},
	}, node11)

	node11, err = node6.UnionJoin(node9)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: node6.fieldNode,
		valueType: node6.valueType,
		boostNode: node6.boostNode,
		terms:     []LeafValue{"bar1", "bar2", "bar3"},
	}, node11)

	node11, err = node6.InterSect(node10)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {
				node6, &TermsNode{
					fieldNode: node6.fieldNode,
					valueType: node6.valueType,
					boostNode: node6.boostNode,
					terms:     []LeafValue{"bar2", "bar3"},
				},
			},
		},
	}, node11)

	node11, err = node6.UnionJoin(node10)
	assert.Nil(t, err)
	assert.Equal(t, &TermsNode{
		fieldNode: node6.fieldNode,
		valueType: node6.valueType,
		boostNode: node6.boostNode,
		terms:     []LeafValue{"bar1", "bar2", "bar3"},
	}, node11)
}

func TestDifferenceValueList(t *testing.T) {
	assert.Equal(t, []LeafValue{1, 2}, DifferenceValueList(
		[]LeafValue{1, 1, 2}, []LeafValue{3, 4}, mapping.INTEGER_FIELD_TYPE,
	))
	assert.Equal(t, []LeafValue{1}, DifferenceValueList(
		[]LeafValue{1, 1, 2, 2}, []LeafValue{2, 3, 4}, mapping.INTEGER_FIELD_TYPE,
	))
	assert.Equal(t, []LeafValue{}, DifferenceValueList(
		[]LeafValue{1, 1, 2}, []LeafValue{1, 2}, mapping.INTEGER_FIELD_TYPE,
	))
}
