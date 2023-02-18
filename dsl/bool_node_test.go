package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestBoolNodeUnionJoinLeafNode(t *testing.T) {
	n := &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
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
	}
	x1 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{field: "foo1"},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
				value:     "bar1",
			},
		},
	}
	n1, _ := n.UnionJoin(x1)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Must: map[string][]AstNode{
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
		Should: map[string][]AstNode{
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
	}, n1)

	x2 := &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1"},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
			rValue:    "bar3",
			lValue:    "bar0",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
		boostNode: boostNode{
			boost: 1.6,
		},
	}
	n2, err := n.UnionJoin(x2)
	assert.NotNil(t, err)
	assert.Nil(t, n2)

	x3 := &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1"},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
			rValue:    "bar5",
			lValue:    "bar2",
			rCmpSym:   LT,
			lCmpSym:   GTE,
		},
	}
	n3, _ := n.UnionJoin(x3)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Must: map[string][]AstNode{
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
		Should: map[string][]AstNode{
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
				&RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo1"},
						valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
						rValue:    "bar5",
						lValue:    "bar2",
						rCmpSym:   LT,
						lCmpSym:   GTE,
					},
				},
			},
		},
	}, n3)

	x4 := &TermNode{
		kvNode: kvNode{
			fieldNode: fieldNode{field: "foo1"},
			valueNode: valueNode{
				valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
				value:     "bar5",
			},
		},
	}
	n4, _ := n3.UnionJoin(x4)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Must: map[string][]AstNode{
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
		Should: map[string][]AstNode{
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
				&RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo1"},
						valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
						rValue:    "bar5",
						lValue:    "bar2",
						rCmpSym:   LTE,
						lCmpSym:   GTE,
					},
				},
			},
		},
	}, n4)

	x5 := &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{field: "foo1"},
			valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
			rValue:    "bar3",
			lValue:    "bar0",
			rCmpSym:   LTE,
			lCmpSym:   GTE,
		},
	}
	n5, _ := n4.UnionJoin(x5)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND | OR},
		Must: map[string][]AstNode{
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
		Should: map[string][]AstNode{
			"foo1": {
				&RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo1"},
						valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE, aType: false},
						rValue:    "bar5",
						lValue:    "bar0",
						rCmpSym:   LTE,
						lCmpSym:   GTE,
					},
				},
			},
		},
	}, n5)

}
