package dsl

import (
	"net"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestRangeNode(t *testing.T) {
	var n1 = NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
			GTE,
			LTE,
		),
		WithTimeZone(time.Local.String()),
		WithRelation(INTERSECTS),
		WithBoost(1.2),
	)
	assert.Equal(t, RANGE_DSL_TYPE, n1.DslType())
	assert.Equal(t, DSL{
		"range": DSL{
			"foo": DSL{
				"gte":       time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local).UnixNano() / 1e6,
				"lte":       time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local).UnixNano() / 1e6,
				"format":    "epoch_millis",
				"time_zone": time.Local.String(),
				"relation":  INTERSECTS,
				"boost":     1.2,
			},
		},
	}, n1.ToDSL())

	n3 := NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
			GTE,
			LTE,
		),
		WithBoost(1.3),
	)

	n4 := NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
			GTE,
			LTE,
		),
		WithBoost(1.2),
	)
	n5, err := n3.UnionJoin(n4)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {n3, n4},
		},
		minimumShouldMatch: 1,
	}, n5)
	n6, err := n3.InterSect(n4)
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: AND},
		Must: map[string][]AstNode{
			"foo": {n3, n4},
		},
	}, n6)
}

func TestRangeInverse(t *testing.T) {
	var n1 = NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
			GTE,
			LTE,
		),
	)

	n2, err := n1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: OR},
		Should: map[string][]AstNode{
			"foo": {
				NewRangeNode(
					NewRgNode(
						NewFieldNode(NewLfNode(), "foo"),
						NewValueType(mapping.DATE_FIELD_TYPE, false),
						minInf[mapping.DATE_FIELD_TYPE],
						time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
						GT,
						LT,
					),
				),
				NewRangeNode(
					NewRgNode(
						NewFieldNode(NewLfNode(), "foo"),
						NewValueType(mapping.DATE_FIELD_TYPE, false),
						time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
						maxInf[mapping.DATE_FIELD_TYPE],
						GT,
						LT,
					),
				),
			},
		},
		minimumShouldMatch: 1,
	}, n2)

	n1 = NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			minInf[mapping.DATE_FIELD_TYPE],
			time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
			GT,
			LTE,
		),
	)
	n2, err = n1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			time.Date(2006, 1, 2, 0, 0, 0, 0, time.Local),
			maxInf[mapping.DATE_FIELD_TYPE],
			GT,
			LT,
		),
	), n2)

	n1 = NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
			maxInf[mapping.DATE_FIELD_TYPE],
			GTE,
			LT,
		),
	)
	n2, err = n1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			minInf[mapping.DATE_FIELD_TYPE],
			time.Date(2006, 1, 3, 0, 0, 0, 0, time.Local),
			GT,
			LT,
		),
	), n2)

	n1 = NewRangeNode(
		NewRgNode(
			NewFieldNode(NewLfNode(), "foo"),
			NewValueType(mapping.DATE_FIELD_TYPE, false),
			minInf[mapping.DATE_FIELD_TYPE],
			maxInf[mapping.DATE_FIELD_TYPE],
			GT,
			LT,
		),
	)
	n2, err = n1.Inverse()
	assert.Nil(t, err)
	assert.Equal(t, &BoolNode{
		opNode: opNode{opType: NOT},
		MustNot: map[string][]AstNode{
			"foo": {
				&ExistsNode{
					fieldNode: n1.fieldNode,
				},
			},
		},
	}, n2)

}

func TestRangeNodeString(t *testing.T) {
	assert.Equal(t, "(1, 2)", (&RangeNode{rgNode: rgNode{lValue: 1, rValue: 2, lCmpSym: GT, rCmpSym: LT, valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE}}}).String())
	assert.Equal(t, "[1, 2)", (&RangeNode{rgNode: rgNode{lValue: 1, rValue: 2, lCmpSym: GTE, rCmpSym: LT, valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE}}}).String())
	assert.Equal(t, "(1, 2]", (&RangeNode{rgNode: rgNode{lValue: 1, rValue: 2, lCmpSym: GT, rCmpSym: LTE, valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE}}}).String())
	assert.Equal(t, "[1, 2]", (&RangeNode{rgNode: rgNode{lValue: 1, rValue: 2, lCmpSym: GTE, rCmpSym: LTE, valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE}}}).String())
	assert.Equal(t, "(\"1\", \"2\")", (&RangeNode{rgNode: rgNode{lValue: "1", rValue: "2", lCmpSym: GT, rCmpSym: LT, valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE}}}).String())
	assert.Equal(t, "[\"1\", \"2\")", (&RangeNode{rgNode: rgNode{lValue: "1", rValue: "2", lCmpSym: GTE, rCmpSym: LT, valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE}}}).String())
	assert.Equal(t, "(\"1\", \"2\"]", (&RangeNode{rgNode: rgNode{lValue: "1", rValue: "2", lCmpSym: GT, rCmpSym: LTE, valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE}}}).String())
	assert.Equal(t, "[\"1\", \"2\"]", (&RangeNode{rgNode: rgNode{lValue: "1", rValue: "2", lCmpSym: GTE, rCmpSym: LTE, valueType: valueType{mType: mapping.KEYWORD_FIELD_TYPE}}}).String())
}

func TestCheckRangeOverlap(t *testing.T) {
	type args struct {
		n *RangeNode
		t *RangeNode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_overlap_01",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2,
						lCmpSym: GTE, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1,
						lCmpSym: GT, rCmpSym: LTE,
					},
				},
			},
			want: true,
		},
		{
			name: "test_overlap_02",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 2,
						lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2,
						lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want: true,
		},
		{
			name: "test_no_overlap_01",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2,
						lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1,
						lCmpSym: GT, rCmpSym: LTE,
					},
				},
			},
			want: false,
		},
		{
			name: "test_no_overlap_02",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1,
						lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2,
						lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want: false,
		},
		{
			name: "test_no_overlap_03",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2,
						lCmpSym: GTE, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 0,
						lCmpSym: GT, rCmpSym: LTE,
					},
				},
			},
			want: false,
		},
		{
			name: "test_no_overlap_04",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 0,
						lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2,
						lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkRangeOverlap(tt.args.n, tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRangeNodeUnionJoinRangeNode(t *testing.T) {
	type args struct {
		n *RangeNode
		t *RangeNode
	}
	tests := []struct {
		name    string
		args    args
		want    AstNode
		wantErr bool
	}{
		{
			name: "test_no_overlap",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 0, lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LTE,
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: OR},
				Should: map[string][]AstNode{
					"foo": {
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
								lValue:    -1, rValue: 0, lCmpSym: GT, rCmpSym: LTE,
							},
						},
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
								lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LTE,
							},
						},
					},
				},
				minimumShouldMatch: 1,
			},
		},
		{
			name: "test_overlap_01",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1, lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    0, rValue: 2, lCmpSym: GT, rCmpSym: LT,
				},
			},
		},
		{
			name: "test_overlap_02",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1, lCmpSym: GT, rCmpSym: LTE,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    0, rValue: 2, lCmpSym: GT, rCmpSym: LT,
				},
			},
		},
		{
			name: "test_overlap_03",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 1, lCmpSym: GTE, rCmpSym: LTE,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    -1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
				},
			},
		},
		{
			name: "test_overlap_04",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LTE,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    -1, rValue: 2, lCmpSym: GT, rCmpSym: LTE,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rangeNodeUnionJoinRangeNode(tt.args.n, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("rangeNodeUnionJoinRangeNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRangeNodeIntersectRangeNode(t *testing.T) {
	type args struct {
		n *RangeNode
		t *RangeNode
	}
	tests := []struct {
		name    string
		args    args
		want    AstNode
		wantErr bool
	}{
		{
			name: "test_no_overlap",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 0, lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_no_overlap_in_array_type",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true},
						lValue:    -1, rValue: 0, lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: AND},
				Must: map[string][]AstNode{
					"foo": {
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true},
								lValue:    -1, rValue: 0, lCmpSym: GT, rCmpSym: LTE,
							},
						},
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true},
								lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_overlap_01",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1, lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    1, rValue: 1, lCmpSym: GTE, rCmpSym: LTE,
				},
			},
		},
		{
			name: "test_overlap_02",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LTE,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    0, rValue: 1, lCmpSym: GT, rCmpSym: LTE,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    1, rValue: 1, lCmpSym: GTE, rCmpSym: LTE,
				},
			},
		},
		{
			name: "test_overlap_03",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 1, lCmpSym: GTE, rCmpSym: LTE,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    -1, rValue: 1, lCmpSym: GT, rCmpSym: LTE,
				},
			},
		},
		{
			name: "test_overlap_04",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    -1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LTE,
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    1, rValue: 2, lCmpSym: GTE, rCmpSym: LT,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rangeNodeIntersectRangeNode(tt.args.n, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("rangeNodeIntersectRangeNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckRangeInclude(t *testing.T) {
	type args struct {
		n *RangeNode
		t LeafValue
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_check_include_01",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: 2,
			},
			want: true,
		},
		{
			name: "test_check_include_02",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GTE, rCmpSym: LT,
					},
				},
				t: 1,
			},
			want: true,
		},
		{
			name: "test_check_include_03",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LTE,
					},
				},
				t: 3,
			},
			want: true,
		},
		{
			name: "test_check_no_include_01",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: 4,
			},
			want: false,
		},
		{
			name: "test_check_no_include_02",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GTE, rCmpSym: LT,
					},
				},
				t: 0,
			},
			want: false,
		},
		{

			name: "test_check_no_include_03",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: 1,
			},
			want: false,
		},
		{
			name: "test_check_no_include_04",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GTE, rCmpSym: LT,
					},
				},
				t: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkRangeInclude(tt.args.n, tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRangeNodeIntersectTermNode(t *testing.T) {
	type args struct {
		n *RangeNode
		t *TermNode
	}
	tests := []struct {
		name    string
		args    args
		want    AstNode
		wantErr bool
	}{
		{
			name: "test_not_include",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 4},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_not_include_with_array_type",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true},
						lValue:    1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true}, value: 4},
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: AND},
				Must: map[string][]AstNode{
					"foo": {
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true},
								lValue:    1, rValue: 2, lCmpSym: GT, rCmpSym: LT,
							},
						},
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: true}, value: 4},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_include",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 2},
					},
				},
			},
			want: &TermNode{
				kvNode: kvNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 2},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rangeNodeIntersectTermNode(tt.args.n, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("rangeNodeIntersectTermNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRangeNodeUnionJoinTermNode(t *testing.T) {
	type args struct {
		n *RangeNode
		t *TermNode
	}
	tests := []struct {
		name    string
		args    args
		want    AstNode
		wantErr bool
	}{
		{
			name: "test_include",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 2},
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
				},
			},
		},
		{
			name: "test_not_include",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 4},
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: OR},
				Should: map[string][]AstNode{
					"foo": {
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
								lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
							},
						},
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
								valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 4},
							},
						},
					},
				},
				minimumShouldMatch: 1,
			},
		},
		{
			name: "test_not_include_left",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 1},
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    1, rValue: 3, lCmpSym: GTE, rCmpSym: LT,
				},
			},
		},
		{
			name: "test_not_include_right",
			args: args{
				n: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
						lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LT,
					},
				},
				t: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
						valueNode: valueNode{valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false}, value: 3},
					},
				},
			},
			want: &RangeNode{
				rgNode: rgNode{
					fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
					valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
					lValue:    1, rValue: 3, lCmpSym: GT, rCmpSym: LTE,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rangeNodeUnionJoinTermNode(tt.args.n, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("rangeNodeUnionJoinTermNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRangeNodeToDsl(t *testing.T) {
	var node1 = &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
			valueType: valueType{mType: mapping.INTEGER_FIELD_TYPE, aType: false},
			lValue:    1, rValue: 7, lCmpSym: GT, rCmpSym: LTE,
		},
		boostNode: boostNode{boost: 1.0},
	}

	assert.Equal(t, DSL{
		"range": DSL{
			"foo": DSL{
				GT.String():  1,
				LTE.String(): 7,
				"relation":   RelationType(""),
				"boost":      1.0,
			},
		},
	}, node1.ToDSL())
	t.Log(node1.ToDSL().String())

	var node2 = &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
			valueType: valueType{mType: mapping.DATE_FIELD_TYPE, aType: false},
			lValue:    time.Date(2022, 01, 02, 0, 0, 0, 0, time.UTC),
			rValue:    time.Date(2022, 01, 03, 0, 0, 0, 0, time.UTC),
			lCmpSym:   GT, rCmpSym: LTE,
		},
		boostNode: boostNode{boost: 1.0},
	}

	assert.Equal(t, DSL{
		"range": DSL{
			"foo": DSL{
				GT.String():  leafValueToPrintValue(node2.lValue, node2.mType),
				LTE.String(): leafValueToPrintValue(node2.rValue, node2.mType),
				"format":     "epoch_millis",
				"relation":   RelationType(""),
				"boost":      1.0,
			},
		},
	}, node2.ToDSL())
	t.Log(node2.ToDSL().String())

	var node3 = &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
			valueType: valueType{mType: mapping.IP_FIELD_TYPE, aType: false},
			lValue:    net.IP([]byte{1, 2, 3, 4}),
			rValue:    net.IP([]byte{1, 2, 4, 5}),
			lCmpSym:   GT, rCmpSym: LTE,
		},
		boostNode: boostNode{boost: 1.0},
	}

	assert.Equal(t, DSL{
		"range": DSL{
			"foo": DSL{
				GT.String():  leafValueToPrintValue(node3.lValue, node3.mType),
				LTE.String(): leafValueToPrintValue(node3.rValue, node3.mType),
				"relation":   RelationType(""),
				"boost":      1.0,
			},
		},
	}, node3.ToDSL())
	t.Log(node3.ToDSL().String())

	var v1, _ = version.NewVersion("1.2.3")
	var v2, _ = version.NewVersion("1.2.10")
	var node4 = &RangeNode{
		rgNode: rgNode{
			fieldNode: fieldNode{lfNode: lfNode{}, field: "foo"},
			valueType: valueType{mType: mapping.VERSION_FIELD_TYPE, aType: false},
			lValue:    v1,
			rValue:    v2,
			lCmpSym:   GT, rCmpSym: LTE,
		},
		boostNode: boostNode{boost: 1.0},
	}

	assert.Equal(t, DSL{
		"range": DSL{
			"foo": DSL{
				GT.String():  leafValueToPrintValue(node4.lValue, node4.mType),
				LTE.String(): leafValueToPrintValue(node4.rValue, node4.mType),
				"relation":   RelationType(""),
				"boost":      1.0,
			},
		},
	}, node4.ToDSL())
	t.Log(node4.ToDSL().String())
}
