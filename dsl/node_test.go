package dsl

import (
	"net"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
)

func TestOpNode(t *testing.T) {
	var n = NewOpNode(OR)
	assert.Equal(t, OP_NODE_TYPE, n.AstType())
	assert.Equal(t, OR, n.OpType())
	WithFilterCtx(true)(n)
	assert.Equal(t, true, n.GetFilterCtx())
	WithFilterCtx(false)(n)
	assert.Equal(t, false, n.GetFilterCtx())
}

func TestLeafNode(t *testing.T) {
	var n = NewLfNode()
	assert.Equal(t, LEAF_NODE_TYPE, n.AstType())
	WithFilterCtx(true)(n)
	assert.Equal(t, true, n.GetFilterCtx())
	WithFilterCtx(false)(n)
	assert.Equal(t, false, n.GetFilterCtx())
}

func TestFieldNode(t *testing.T) {
	var n = NewFieldNode(NewLfNode(), "foo")
	assert.Equal(t, "foo", n.NodeKey())
}

func TestValueNode(t *testing.T) {
	var n = NewValueNode("12", NewValueType(mapping.KEYWORD_FIELD_TYPE, true))
	assert.Equal(t, "12", n.toPrintValue())
	assert.Equal(t, true, n.IsArrayType())
	WithArrayType(false)(n)
	assert.Equal(t, false, n.IsArrayType())

	n = NewValueNode(net.IP([]byte{1, 1, 1, 1}), NewValueType(mapping.IP_FIELD_TYPE, true))
	assert.Equal(t, "1.1.1.1", n.toPrintValue())

	var v, _ = version.NewVersion("v1.1.1")
	n = NewValueNode(v, NewValueType(mapping.VERSION_FIELD_TYPE, true))
	assert.Equal(t, "1.1.1", n.toPrintValue())
}

func TestKvNode(t *testing.T) {
	var n = NewKVNode(NewFieldNode(NewLfNode(), "foo"), NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, true)))
	assert.Equal(t, &kvNode{
		fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
		valueNode: valueNode{valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true}, value: "bar"},
	}, n)
}

func TestRgNode(t *testing.T) {
	var n = NewRgNode(NewFieldNode(NewLfNode(), "foo"), NewValueType(mapping.KEYWORD_FIELD_TYPE, true),
		"bar1", "bar2", GTE, LTE,
	)
	assert.Equal(t, &rgNode{
		fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
		valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true},
		lValue:    "bar1", rValue: "bar2", lCmpSym: GTE, rCmpSym: LTE,
	}, n)
}
