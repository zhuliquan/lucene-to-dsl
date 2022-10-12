package dsl

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/zhuliquan/lucene-to-dsl/mapping"
// )

// func TestFuzzyNode(t *testing.T) {
// 	var node1 = &FuzzyNode{
// 		KvNode: KvNode{
// 			Field: "foo",
// 			Type:  mapping.TEXT_FIELD_TYPE,
// 			Value: "bar",
// 		},
// 		Fuzziness: "1",
// 	}

// 	var node2 = &ExistsNode{
// 		KvNode: KvNode{Field: "foo"},
// 	}

// 	var node3, _ = node1.UnionJoin(node2)
// 	assert.Equal(t, node2, node3)

// 	var node4, _ = node1.InterSect(node2)
// 	assert.Equal(t, node1, node4)

// 	var node5, _ = node1.Inverse()
// 	assert.Equal(t, &NotNode{
// 		Nodes: map[string][]AstNode{
// 			node4.NodeKey(): {node1},
// 		},
// 	}, node5)

// 	assert.Equal(t, FUZZY_DSL_TYPE, node1.DslType())
// 	assert.Equal(t, DSL{"fuzzy": DSL{"foo": DSL{
// 		"value":          "bar",
// 		"fuzziness":      "1",
// 		"prefix_length":  0,
// 		"max_expansions": 0,
// 		"transpositions": false,
// 	}}}, node1.ToDSL())

// }
