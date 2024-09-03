package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mapping "github.com/zhuliquan/es-mapping"
)

func TestBoolNodeUnionJoinLeafNode(t *testing.T) {
	t.Run("test and node union join different leaf node", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		term2 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		left := NewBoolNode(term1, AND)
		right := term2
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {term2},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test or node union join different leaf node", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		term2 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		left := NewBoolNode(term1, OR)
		right := term2
		expect := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo":  {term1},
				"foo1": {term2},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test and node union join same leaf node", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		left := NewBoolNode(term1, AND)
		right := term1
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo": {term1},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test or node union join same leaf node", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		left := NewBoolNode(term1, OR)
		right := term1
		expect := term1
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test and node union join different leaf and union join range node", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		term2 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueNode("bar1", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		range1 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar2", "bar5",
				GTE, LT,
			),
		)

		left := NewBoolNode(term1, AND)
		left, _ = left.UnionJoin(term2)
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {term2, range1},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(range1)
		assert.Equal(t, expect, actual)
	})

	t.Run("test node union join same leaf node and reduce lt op", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		term2 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueNode("bar5", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		range1 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar2", "bar5",
				GTE, LT,
			),
		)
		range2 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar2", "bar5",
				GTE, LTE,
			),
		)

		left := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {range1},
			},
			MinimumShouldMatch: 1,
		}
		right := term2
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {range2},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test node union join same leaf node and reduce gt op", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		term2 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueNode("bar2", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		range1 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar2", "bar5",
				GT, LT,
			),
		)
		range2 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar2", "bar5",
				GTE, LT,
			),
		)

		left := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {range1},
			},
			MinimumShouldMatch: 1,
		}
		right := term2
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {range2},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test node union join different range node and range compact", func(t *testing.T) {
		term1 := NewTermNode(
			NewKVNode(
				NewFieldNode(NewLfNode(), "foo"),
				NewValueNode("bar", NewValueType(mapping.KEYWORD_FIELD_TYPE, false)),
			),
		)
		range1 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar2", "bar5",
				GTE, LTE,
			),
		)
		range2 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar0", "bar3",
				GTE, LTE,
			),
		)
		range3 := NewRangeNode(
			NewRgNode(
				NewFieldNode(NewLfNode(), "foo1"),
				NewValueType(mapping.KEYWORD_FIELD_TYPE, false),
				"bar0", "bar5",
				GTE, LTE,
			),
		)
		left := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {range1},
			},
			MinimumShouldMatch: 1,
		}
		right := range2
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Must: map[string][]AstNode{
				"foo": {term1},
			},
			Should: map[string][]AstNode{
				"foo1": {range3},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.UnionJoin(right)
		assert.Equal(t, expect, actual)
	})
}

// TODO: 明儿继续处理此case
func TestBoolNodeIntersectMustLeafNode(t *testing.T) {
	t.Run("test or bool node intersect new leaf node", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		right := &TermNode{
			kvNode: kvNode{
				fieldNode: fieldNode{field: "foo1"},
				valueNode: valueNode{
					valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
					value:     "bar1",
				},
			},
		}
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			Must: map[string][]AstNode{
				"foo1": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo1"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar1",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.InterSect(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test or bool node intersect leaf node and error return", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		right := &RangeNode{
			rgNode: rgNode{
				fieldNode: fieldNode{field: "foo1"},
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
				rValue:    "bar4",
				lValue:    "bar2",
				rCmpSym:   LTE,
				lCmpSym:   GTE,
			},
		}
		n2, err := left.InterSect(right)
		assert.NotNil(t, err)
		assert.Nil(t, n2)
	})

	t.Run("test or bool node intersect range node and compact", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		right := &RangeNode{
			rgNode: rgNode{
				fieldNode: fieldNode{field: "foo1"},
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
				rValue:    "bar4",
				lValue:    "bar1",
				rCmpSym:   LTE,
				lCmpSym:   GTE,
			},
		}
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			Must: map[string][]AstNode{
				"foo1": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo1"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar1",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.InterSect(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test or bool node intersect leaf node and compact", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}

		right := &TermNode{
			kvNode: kvNode{
				fieldNode: fieldNode{field: "foo1"},
				valueNode: valueNode{
					valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
					value:     "bar1",
				},
			},
		}
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			Must: map[string][]AstNode{
				"foo1": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo1"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
								value:     "bar1",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.InterSect(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("test or bool node intersect range node and compact", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		right := &RangeNode{
			rgNode: rgNode{
				fieldNode: fieldNode{field: "foo1"},
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				rValue:    "bar4",
				lValue:    "bar2",
				rCmpSym:   LTE,
				lCmpSym:   GTE,
			},
		}
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			Must: map[string][]AstNode{
				"foo1": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo1"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
								value:     "bar1",
							},
						},
					},
					&RangeNode{
						rgNode: rgNode{
							fieldNode: fieldNode{field: "foo1"},
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							rValue:    "bar4",
							lValue:    "bar2",
							rCmpSym:   LTE,
							lCmpSym:   GTE,
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.InterSect(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}

		right := &RangeNode{
			rgNode: rgNode{
				fieldNode: fieldNode{field: "foo1"},
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				rValue:    "bar4",
				lValue:    "bar1",
				rCmpSym:   LTE,
				lCmpSym:   GTE,
			},
		}
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			Must: map[string][]AstNode{
				"foo1": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo1"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
								value:     "bar1",
							},
						},
					},
					&RangeNode{
						rgNode: rgNode{
							fieldNode: fieldNode{field: "foo1"},
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							rValue:    "bar4",
							lValue:    "bar2",
							rCmpSym:   LTE,
							lCmpSym:   GTE,
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.InterSect(right)
		assert.Equal(t, expect, actual)
	})

	t.Run("", func(t *testing.T) {
		left := &BoolNode{
			opNode: opNode{opType: OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		right := &RangeNode{
			rgNode: rgNode{
				fieldNode: fieldNode{field: "foo1"},
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				rValue:    "bar5",
				lValue:    "bar3",
				rCmpSym:   LTE,
				lCmpSym:   GTE,
			},
		}
		expect := &BoolNode{
			opNode: opNode{opType: AND | OR},
			Should: map[string][]AstNode{
				"foo": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
								value:     "bar",
							},
						},
					},
				},
			},
			Must: map[string][]AstNode{
				"foo1": {
					&TermNode{
						kvNode: kvNode{
							fieldNode: fieldNode{field: "foo1"},
							valueNode: valueNode{
								valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
								value:     "bar1",
							},
						},
					},
					&RangeNode{
						rgNode: rgNode{
							fieldNode: fieldNode{field: "foo1"},
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							rValue:    "bar4",
							lValue:    "bar3",
							rCmpSym:   LTE,
							lCmpSym:   GTE,
						},
					},
				},
			},
			MinimumShouldMatch: 1,
		}
		actual, _ := left.InterSect(right)
		assert.Equal(t, expect, actual)
	})
}

func TestBoolNodeIntersectFilterLeafNode(t *testing.T) {
	n := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}

	x1 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
				value:     "bar1",
			},
		},
	}
	n1, err := n.InterSect(x1)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		Filter: map[string][]AstNode{
			"foo1": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar1",
						},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, n1)

	x2 := &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
			rValue:    "bar4",
			lValue:    "bar2",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
	}
	n2, err := n.InterSect(x2)
	assert.NotNil(t, err)
	assert.Nil(t, n2)

	x3 := &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
			rValue:    "bar4",
			lValue:    "bar1",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
	}

	n3, err := n.InterSect(x3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		Filter: map[string][]AstNode{
			"foo1": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar1",
						},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, n3)

	n = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}

	x1 = &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar1",
			},
		},
	}
	n1, err = n.InterSect(x1)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		Filter: map[string][]AstNode{
			"foo1": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							value:     "bar1",
						},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, n1)

	x2 = &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
			rValue:    "bar4",
			lValue:    "bar2",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
	}
	n2, err = n.InterSect(x2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		Filter: map[string][]AstNode{
			"foo1": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							value:     "bar1",
						},
					},
				},
				&RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
						rValue:    "bar4",
						lValue:    "bar2",
						rCmpSym:   LTE,
						lCmpSym:   GTE,
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, n2)

	x3 = &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
			rValue:    "bar4",
			lValue:    "bar1",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
	}
	n3, err = n.InterSect(x3)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		Filter: map[string][]AstNode{
			"foo1": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							value:     "bar1",
						},
					},
				},
				&RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
						rValue:    "bar4",
						lValue:    "bar2",
						rCmpSym:   LTE,
						lCmpSym:   GTE,
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, n3)

	x4 := &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
			rValue:    "bar5",
			lValue:    "bar3",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
	}
	n4, _ := n.InterSect(x4)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Should: map[string][]AstNode{
			"foo": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo"},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
							value:     "bar",
						},
					},
				},
			},
		},
		Filter: map[string][]AstNode{
			"foo1": {
				&TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueNode: valueNode{
							valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
							value:     "bar1",
						},
					},
				},
				&RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo1", lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}},
						valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
						rValue:    "bar4",
						lValue:    "bar3",
						rCmpSym:   LTE,
						lCmpSym:   GTE,
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, n4)
}

func TestBoolNodeInverse(t *testing.T) {
	child1 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo1",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar1",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child2 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo2",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar2",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child3 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo3",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar3",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child4 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo4",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar4",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	// inverse `and` node
	node1 := &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	tmp, err := node1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			OP_KEY: {
				&BoolNode{
					opNode: opNode{opType: AND},
					Must: map[string][]AstNode{
						"foo1": {child1},
						"foo2": {child2},
					},
				},
			},
		},
	}, tmp)

	// inverse `or` node
	node2 := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	tmp, err = node2.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}, tmp)

	// inverse `not` node
	node3 := &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	tmp, err = node3.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}, tmp)

	node4 := &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo1": {child1},
		},
	}
	tmp, err = node4.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, child1, tmp)

	// inverse `and or` node
	node5 := &BoolNode{
		opNode: opNode{opType: AND | OR},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}
	tmp, err = node5.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			OP_KEY: {&BoolNode{
				opNode: opNode{opType: AND | OR},
				Must: map[string][]AstNode{
					"foo1": {child1},
					"foo2": {child2},
				},
				Should: map[string][]AstNode{
					"foo3": {child3},
					"foo4": {child4},
				},
				MinimumShouldMatch: 1,
			}},
		},
	}, tmp)

	// inverse `or not` node
	node6 := &BoolNode{
		opNode: opNode{opType: NOT | OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MustNot: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}
	tmp, err = node6.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
			OP_KEY: {
				&BoolNode{
					opNode: opNode{opType: NOT},
					MustNot: map[string][]AstNode{
						"foo1": {child1},
						"foo2": {child2},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, tmp)

	// inverse `and not` node
	node7 := &BoolNode{
		opNode: opNode{opType: NOT | AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MustNot: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	tmp, err = node7.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
			OP_KEY: {
				&BoolNode{
					opNode: opNode{opType: NOT},
					MustNot: map[string][]AstNode{
						OP_KEY: {
							&BoolNode{
								opNode: opNode{opType: AND},
								Must: map[string][]AstNode{
									"foo1": {child1},
									"foo2": {child2},
								},
							},
						},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, tmp)

	// inverse `and or not` node
	node8 := &BoolNode{
		opNode: opNode{opType: NOT | AND | OR},
		Must: map[string][]AstNode{
			"foo1": {child1},
		},
		MustNot: map[string][]AstNode{
			"foo2": {child2},
		},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	tmp, err = node8.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			OP_KEY: {node8},
		},
	}, tmp)
}

func TestBoolNodeToDSL(t *testing.T) {
	child1 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo1",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar1",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child2 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo2",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar2",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child3 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo3",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar3",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child4 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo4",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar4",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	x1 := &BoolNode{
		opNode: opNode{opType: AND | OR | NOT},
		Must: map[string][]AstNode{
			"foo1": {child1},
		},
		Filter: map[string][]AstNode{
			"foo2": {child2},
		},
		Should: map[string][]AstNode{
			"foo3": {child3},
		},
		MustNot: map[string][]AstNode{
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"must":                 child1.ToDSL(),
			"filter":               child2.ToDSL(),
			"should":               child3.ToDSL(),
			"must_not":             child4.ToDSL(),
			"minimum_should_match": 1,
		},
	}, x1.ToDSL())

	x2 := &BoolNode{}
	assert.Equal(t, EmptyDSL, x2.ToDSL())
}

func TestBoolNodeUnionJoinBoolNode(t *testing.T) {
	child1 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo1",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar1",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child2 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo2",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar2",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child3 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo3",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar3",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child4 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo4",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar4",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}

	// test or union join or
	node1 := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	node2 := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}

	node3, err := node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}, node3)

	// test or union join and
	node1 = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	node2 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}

	node3, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
			OP_KEY: {&BoolNode{
				opNode: opNode{opType: AND},
				Must: map[string][]AstNode{
					"foo3": {child3},
					"foo4": {child4},
				},
			}},
		},
		MinimumShouldMatch: 1,
	}, node3)

	// test and union join or
	node1 = &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	node2 = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}

	node3, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
			OP_KEY: {&BoolNode{
				opNode: opNode{opType: AND},
				Filter: map[string][]AstNode{
					"foo1": {child1},
					"foo2": {child2},
				},
			}},
		},
		MinimumShouldMatch: 1,
	}, node3)

	// test and union join or
	node1 = &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	node2 = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}

	node3, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
			OP_KEY: {&BoolNode{
				opNode: opNode{opType: AND},
				Filter: map[string][]AstNode{
					"foo1": {child1},
					"foo2": {child2},
				},
			}},
		},
		MinimumShouldMatch: 1,
	}, node3)

	// test not node union join not node
	node1 = &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	node2 = &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}

	node3, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			OP_KEY: {
				&BoolNode{
					opNode: opNode{opType: AND},
					Must: map[string][]AstNode{
						OP_KEY: {
							&BoolNode{
								opNode: opNode{opType: OR},
								Should: map[string][]AstNode{
									"foo1": {child1},
									"foo2": {child2},
								},
								MinimumShouldMatch: 1,
							},
							&BoolNode{
								opNode: opNode{opType: OR},
								Should: map[string][]AstNode{
									"foo3": {child3},
									"foo4": {child4},
								},
								MinimumShouldMatch: 1,
							},
						},
					},
				},
			},
		},
	}, node3)

	// test and node union join and node
	node1 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	node2 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}

	node3, err = node1.UnionJoin(node2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			OP_KEY: {
				&BoolNode{
					opNode: opNode{opType: AND},
					Must: map[string][]AstNode{
						"foo1": {child1},
						"foo2": {child2},
					},
				},
				&BoolNode{
					opNode: opNode{opType: AND},
					Must: map[string][]AstNode{
						"foo3": {child3},
						"foo4": {child4},
					},
				},
			},
		},
		MinimumShouldMatch: 1,
	}, node3)
}

func TestBoolNodeIntersectBoolNode(t *testing.T) {
	child1 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo1",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar1",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child2 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo2",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar2",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child3 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo3",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar3",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child4 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{}, field: "foo4",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar4",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child5 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}, field: "foo5",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar5",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}
	child6 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{
				lfNode: lfNode{filterCtxNode: filterCtxNode{filterCtx: true}}, field: "foo6",
			},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: true},
				value:     "bar6",
			},
		},
		boostNode: boostNode{boost: 1.2},
	}

	// test or node intersect must node
	orNode := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	mustNode := &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err := orNode.InterSect(mustNode)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR | AND},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}, res)

	// or node intersect filter node
	orNode = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	filterNode := &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err = orNode.InterSect(filterNode)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR | AND},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Filter: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}, res)

	// must node intersect must node
	mustNode1 := &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	mustNode2 := &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	res, err = mustNode1.InterSect(mustNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}, res)

	mustNode1 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	mustNode2 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err = mustNode1.InterSect(mustNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
			"foo3": {child3},
			"foo4": {child4},
		},
	}, res)

	// must node intersect filter node
	mustNode1 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	mustNode2 = &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	res, err = mustNode1.InterSect(mustNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}, res)

	mustNode1 = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	mustNode2 = &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err = mustNode1.InterSect(mustNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Filter: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}, res)

	// filter node intersect filter node
	filterNode1 := &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	filterNode2 := &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo5": {child5},
			"foo6": {child6},
		},
	}
	res, err = filterNode1.InterSect(filterNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
			"foo5": {child5},
			"foo6": {child6},
		},
	}, res)

	// filter node intersect must node
	filterNode = &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	mustNode = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	res, err = filterNode.InterSect(mustNode)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}, res)

	filterNode = &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	mustNode = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err = filterNode.InterSect(mustNode)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Filter: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Must: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}, res)

	// must node intersect or node
	mustNode = &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
	}
	orNode = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}
	res, err = mustNode.InterSect(orNode)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR | AND},
		Must: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}, res)

	// or node intersect or node
	orNode1 := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	orNode2 := &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}
	res, err = orNode1.InterSect(orNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			OP_KEY: {
				&BoolNode{
					opNode: opNode{opType: OR},
					Should: map[string][]AstNode{
						"foo1": {child1},
						"foo2": {child2},
					},
					MinimumShouldMatch: 1,
				},
				&BoolNode{
					opNode: opNode{opType: OR},
					Should: map[string][]AstNode{
						"foo3": {child3},
						"foo4": {child4},
					},
					MinimumShouldMatch: 1,
				},
			},
		},
	}, res)

	// or node intersect not node
	orNode = &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	notNode := &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err = orNode.InterSect(notNode)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR | NOT},
		Should: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MustNot: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
		MinimumShouldMatch: 1,
	}, res)

	// not node intersect not node
	notNode1 := &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
		},
		MinimumShouldMatch: 1,
	}
	notNode2 := &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo3": {child3},
			"foo4": {child4},
		},
	}
	res, err = notNode1.InterSect(notNode2)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo1": {child1},
			"foo2": {child2},
			"foo3": {child3},
			"foo4": {child4},
		},
	}, res)
}
