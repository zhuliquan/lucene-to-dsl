package dsl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestAndNodeToDSL(t *testing.T) {
	n0 := &AndNode{
		FilterNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1", LfNode: LfNode{Filter: true},
					},
				},
			},
		},
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"filter": DSL{"term": DSL{"foo": DSL{"value": "bar1", "boost": 0.0}}},
		},
	}, n0.ToDSL())

	n1 := &AndNode{
		MustNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar2", LfNode: LfNode{Filter: false},
					},
				},
			},
		},
		FilterNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1", LfNode: LfNode{Filter: true},
					},
				},
			},
		},
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"filter": DSL{"term": DSL{"foo": DSL{"value": "bar1", "boost": 0.0}}},
			"must":   DSL{"term": DSL{"foo": DSL{"value": "bar2", "boost": 0.0}}},
		},
	}, n1.ToDSL())

	n2 := &AndNode{
		MustNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar2", LfNode: LfNode{Filter: false},
					},
				},
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar3", LfNode: LfNode{Filter: false},
					},
				},
			},
		},
		FilterNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1", LfNode: LfNode{Filter: true},
					},
				},
			},
		},
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"filter": DSL{"term": DSL{"foo": DSL{"value": "bar1", "boost": 0.0}}},
			"must": []DSL{
				{"term": DSL{"foo": DSL{"value": "bar2", "boost": 0.0}}},
				{"term": DSL{"foo": DSL{"value": "bar3", "boost": 0.0}}},
			},
		},
	}, n2.ToDSL())

	n3 := &AndNode{
		MustNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1", LfNode: LfNode{Filter: true},
					},
				},
			},
		},
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"must": DSL{"term": DSL{"foo": DSL{"value": "bar1", "boost": 0.0}}},
		},
	}, n3.ToDSL())

	n4 := &AndNode{
		FilterNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar2", LfNode: LfNode{Filter: false},
					},
				},
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar3", LfNode: LfNode{Filter: false},
					},
				},
			},
		},
		MustNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1", LfNode: LfNode{Filter: true},
					},
				},
			},
		},
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"must": DSL{"term": DSL{"foo": DSL{"value": "bar1", "boost": 0.0}}},
			"filter": []DSL{
				{"term": DSL{"foo": DSL{"value": "bar2", "boost": 0.0}}},
				{"term": DSL{"foo": DSL{"value": "bar3", "boost": 0.0}}},
			},
		},
	}, n4.ToDSL())

	n5 := &AndNode{
		FilterNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar2", LfNode: LfNode{Filter: false},
					},
				},
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar3", LfNode: LfNode{Filter: false},
					},
				},
			},
		},
		MustNodes: map[string][]AstNode{
			"LEAF:foo": {
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar1", LfNode: LfNode{Filter: true},
					},
				},
				&TermNode{
					KvNode: KvNode{
						Field: "foo", Type: mapping.TEXT_FIELD_TYPE, Value: "bar4", LfNode: LfNode{Filter: true},
					},
				},
			},
		},
	}
	assert.Equal(t, DSL{
		"bool": DSL{
			"must": []DSL{
				{"term": DSL{"foo": DSL{"value": "bar1", "boost": 0.0}}},
				{"term": DSL{"foo": DSL{"value": "bar4", "boost": 0.0}}},
			},
			"filter": []DSL{
				{"term": DSL{"foo": DSL{"value": "bar2", "boost": 0.0}}},
				{"term": DSL{"foo": DSL{"value": "bar3", "boost": 0.0}}},
			},
		},
	}, n5.ToDSL())

	n6 := &AndNode{
		FilterNodes: map[string][]AstNode{},
		MustNodes:   map[string][]AstNode{},
	}

	assert.Equal(t, EmptyDSL, n6.ToDSL())

}
