package dsl

// import (
// 	"net"
// 	"testing"
// 	"time"

// 	"github.com/hashicorp/go-version"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/zhuliquan/lucene-to-dsl/mapping"
// )

// func TestRangeNodeString(t *testing.T) {
// 	assert.Equal(t, "(1, 2)", (&RangeNode{LeftValue: 1, RightValue: 2, LeftCmpSym: GT, RightCmpSym: LT}).String())
// 	assert.Equal(t, "[1, 2)", (&RangeNode{LeftValue: 1, RightValue: 2, LeftCmpSym: GTE, RightCmpSym: LT}).String())
// 	assert.Equal(t, "(1, 2]", (&RangeNode{LeftValue: 1, RightValue: 2, LeftCmpSym: GT, RightCmpSym: LTE}).String())
// 	assert.Equal(t, "[1, 2]", (&RangeNode{LeftValue: 1, RightValue: 2, LeftCmpSym: GTE, RightCmpSym: LTE}).String())
// 	assert.Equal(t, "(\"1\", \"2\")", (&RangeNode{LeftValue: "1", RightValue: "2", LeftCmpSym: GT, RightCmpSym: LT}).String())
// 	assert.Equal(t, "[\"1\", \"2\")", (&RangeNode{LeftValue: "1", RightValue: "2", LeftCmpSym: GTE, RightCmpSym: LT}).String())
// 	assert.Equal(t, "(\"1\", \"2\"]", (&RangeNode{LeftValue: "1", RightValue: "2", LeftCmpSym: GT, RightCmpSym: LTE}).String())
// 	assert.Equal(t, "[\"1\", \"2\"]", (&RangeNode{LeftValue: "1", RightValue: "2", LeftCmpSym: GTE, RightCmpSym: LTE}).String())
// }

// func TestCheckRangeOverlap(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *RangeNode
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "test_overlap_01",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "test_overlap_02",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "test_no_overlap_01",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "test_no_overlap_02",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "test_no_overlap_03",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 0,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "test_no_overlap_04",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 0,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := checkRangeOverlap(tt.args.n, tt.args.t)
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeUnionJoinRangeNode(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *RangeNode
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    AstNode
// 		wantErr bool
// 	}{
// 		{
// 			name: "test_no_overlap",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 0,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want: &OrNode{
// 				MinimumShouldMatch: 1,
// 				Nodes: map[string][]AstNode{
// 					"LEAF:foo": {
// 						&RangeNode{
// 							KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							LeftValue: -1, RightValue: 0,
// 							LeftCmpSym: GT, RightCmpSym: LTE,
// 						},
// 						&RangeNode{
// 							KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							LeftValue: 1, RightValue: 2,
// 							LeftCmpSym: GTE, RightCmpSym: LT,
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "test_overlap_01",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 0, RightValue: 2,
// 				LeftCmpSym: GT, RightCmpSym: LT,
// 			},
// 		},
// 		{
// 			name: "test_overlap_02",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 0, RightValue: 2,
// 				LeftCmpSym: GT, RightCmpSym: LT,
// 			},
// 		},
// 		{
// 			name: "test_overlap_03",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 2,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 1,
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: -1, RightValue: 2,
// 				LeftCmpSym: GTE, RightCmpSym: LT,
// 			},
// 		},
// 		{
// 			name: "test_overlap_04",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 2,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: -1, RightValue: 2,
// 				LeftCmpSym: GT, RightCmpSym: LTE,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rangeNodeUnionJoinRangeNode(tt.args.n, tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("rangeNodeUnionJoinRangeNode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeIntersectRangeNode(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *RangeNode
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    AstNode
// 		wantErr bool
// 	}{
// 		{
// 			name: "test_no_overlap",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 0,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "test_overlap_01",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 1,
// 				LeftCmpSym: GTE, RightCmpSym: LTE,
// 			},
// 		},
// 		{
// 			name: "test_overlap_02",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 0, RightValue: 1,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 1,
// 				LeftCmpSym: GTE, RightCmpSym: LTE,
// 			},
// 		},
// 		{
// 			name: "test_overlap_03",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 2,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 1,
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: -1, RightValue: 1,
// 				LeftCmpSym: GT, RightCmpSym: LTE,
// 			},
// 		},
// 		{
// 			name: "test_overlap_04",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: -1, RightValue: 2,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 2,
// 				LeftCmpSym: GTE, RightCmpSym: LT,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rangeNodeIntersectRangeNode(tt.args.n, tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("rangeNodeIntersectRangeNode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestCheckRangeInclude(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t LeafValue
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "test_check_include_01",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: 2,
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "test_check_include_02",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GTE, RightCmpSym: LT,
// 				},
// 				t: 1,
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "test_check_include_03",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LTE,
// 				},
// 				t: 3,
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "test_check_no_include_01",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: 4,
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "test_check_no_include_02",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: 0,
// 			},
// 			want: false,
// 		},
// 		{

// 			name: "test_check_no_include_03",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: 1,
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "test_check_no_include_04",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_RANGE_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: 3,
// 			},
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := checkRangeInclude(tt.args.n, tt.args.t)
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeIntersectTermNode(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *TermNode
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    AstNode
// 		wantErr bool
// 	}{
// 		{
// 			name: "test_not_include",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 2,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermNode{
// 					KvNode: KvNode{Type: mapping.INTEGER_FIELD_TYPE, Value: 4},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "test_include",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermNode{
// 					KvNode: KvNode{Type: mapping.INTEGER_FIELD_TYPE, Value: 2},
// 				},
// 			},
// 			want: &TermNode{
// 				KvNode: KvNode{Type: mapping.INTEGER_FIELD_TYPE, Value: 2},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rangeNodeIntersectTermNode(tt.args.n, tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("rangeNodeIntersectTermNode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeIntersectTermsNode(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *TermsNode
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    AstNode
// 		wantErr bool
// 	}{
// 		{
// 			name: "test_intersect_partial",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.TEXT_FIELD_TYPE},
// 					LeftValue: "1", RightValue: "5",
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 				t: &TermsNode{
// 					KvNode: KvNode{Type: mapping.TEXT_FIELD_TYPE},
// 					Values: []LeafValue{"1", "3", "8"},
// 				},
// 			},
// 			want: &TermsNode{
// 				KvNode: KvNode{Type: mapping.TEXT_FIELD_TYPE},
// 				Values: []LeafValue{"1", "3"},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test_intersect_partial",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.TEXT_FIELD_TYPE},
// 					LeftValue: "1", RightValue: "5",
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 				t: &TermsNode{
// 					KvNode: KvNode{Type: mapping.TEXT_FIELD_TYPE},
// 					Values: []LeafValue{"0", "7", "8"},
// 				},
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rangeNodeIntersectTermsNode(tt.args.n, tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("rangeNodeIntersectTermsNode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeUnionJoinTermNode(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *TermNode
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    AstNode
// 		wantErr bool
// 	}{
// 		{
// 			name: "test_include",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermNode{
// 					KvNode: KvNode{Type: mapping.INTEGER_FIELD_TYPE, Value: 2},
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 3,
// 				LeftCmpSym: GT, RightCmpSym: LT,
// 			},
// 		},
// 		{
// 			name: "test_not_include",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermNode{
// 					KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE, Value: 4},
// 				},
// 			},
// 			want: &OrNode{
// 				MinimumShouldMatch: 1,
// 				Nodes: map[string][]AstNode{
// 					"LEAF:foo": {
// 						&RangeNode{
// 							KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							LeftValue: 1, RightValue: 3,
// 							LeftCmpSym: GT, RightCmpSym: LT,
// 						},
// 						&TermNode{
// 							KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE, Value: 4},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "test_not_include_left",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermNode{
// 					KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE, Value: 1},
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 3,
// 				LeftCmpSym: GTE, RightCmpSym: LT,
// 			},
// 		},
// 		{
// 			name: "test_not_include_right",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermNode{
// 					KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE, Value: 3},
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 3,
// 				LeftCmpSym: GT, RightCmpSym: LTE,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rangeNodeUnionJoinTermNode(tt.args.n, tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("rangeNodeUnionJoinTermNode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeUnionJoinTermsNode(t *testing.T) {
// 	type args struct {
// 		n *RangeNode
// 		t *TermsNode
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    AstNode
// 		wantErr bool
// 	}{
// 		{
// 			name: "test_include",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GTE, RightCmpSym: LTE,
// 				},
// 				t: &TermsNode{
// 					KvNode: KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 					Values: []LeafValue{1, 2, 3},
// 				},
// 			},
// 			want: &RangeNode{
// 				KvNode:    KvNode{Type: mapping.INTEGER_FIELD_TYPE},
// 				LeftValue: 1, RightValue: 3,
// 				LeftCmpSym: GTE, RightCmpSym: LTE,
// 			},
// 		},
// 		{
// 			name: "test_partial_include",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermsNode{
// 					KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					Values: []LeafValue{2, 4},
// 				},
// 			},
// 			want: &OrNode{
// 				MinimumShouldMatch: 1,
// 				Nodes: map[string][]AstNode{
// 					"LEAF:foo": {
// 						&RangeNode{
// 							KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							LeftValue: 1, RightValue: 3,
// 							LeftCmpSym: GT, RightCmpSym: LT,
// 						},
// 						&TermsNode{
// 							KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							Values: []LeafValue{4},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "test_partial_include_left",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermsNode{
// 					KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					Values: []LeafValue{1, 2, 4},
// 				},
// 			},
// 			want: &OrNode{
// 				MinimumShouldMatch: 1,
// 				Nodes: map[string][]AstNode{
// 					"LEAF:foo": {
// 						&RangeNode{
// 							KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							LeftValue: 1, RightValue: 3,
// 							LeftCmpSym: GTE, RightCmpSym: LT,
// 						},
// 						&TermsNode{
// 							KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							Values: []LeafValue{4},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "test_partial_include_right",
// 			args: args{
// 				n: &RangeNode{
// 					KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					LeftValue: 1, RightValue: 3,
// 					LeftCmpSym: GT, RightCmpSym: LT,
// 				},
// 				t: &TermsNode{
// 					KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 					Values: []LeafValue{2, 3, 4},
// 				},
// 			},
// 			want: &OrNode{
// 				MinimumShouldMatch: 1,
// 				Nodes: map[string][]AstNode{
// 					"LEAF:foo": {
// 						&RangeNode{
// 							KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							LeftValue: 1, RightValue: 3,
// 							LeftCmpSym: GT, RightCmpSym: LTE,
// 						},
// 						&TermsNode{
// 							KvNode: KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 							Values: []LeafValue{4},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := rangeNodeUnionJoinTermsNode(tt.args.n, tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("rangeNodeUnionJoinTermsNode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestRangeNodeInverse(t *testing.T) {
// 	var node1 = &RangeNode{
// 		KvNode:      KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 		LeftValue:   MinInt[32],
// 		RightValue:  MaxInt[32],
// 		LeftCmpSym:  GT,
// 		RightCmpSym: LT,
// 	}

// 	var node2, _ = node1.Inverse()
// 	assert.Equal(t, &NotNode{
// 		Nodes: map[string][]AstNode{
// 			"LEAF:foo": {&ExistsNode{KvNode: node1.KvNode}},
// 		},
// 	}, node2)

// 	node1 = &RangeNode{
// 		KvNode:      KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 		LeftValue:   MinInt[32],
// 		RightValue:  4,
// 		LeftCmpSym:  GT,
// 		RightCmpSym: LT,
// 	}

// 	node2, _ = node1.Inverse()
// 	assert.Equal(t, &RangeNode{
// 		KvNode:      node1.KvNode,
// 		LeftValue:   4,
// 		RightValue:  MaxInt[32],
// 		LeftCmpSym:  GTE,
// 		RightCmpSym: LT,
// 	}, node2)

// 	node1 = &RangeNode{
// 		KvNode:      KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 		LeftValue:   4,
// 		RightValue:  MaxInt[32],
// 		LeftCmpSym:  GT,
// 		RightCmpSym: LT,
// 	}

// 	node2, _ = node1.Inverse()
// 	assert.Equal(t, &RangeNode{
// 		KvNode:      node1.KvNode,
// 		LeftValue:   MinInt[32],
// 		RightValue:  4,
// 		LeftCmpSym:  GT,
// 		RightCmpSym: LTE,
// 	}, node2)

// 	node1 = &RangeNode{
// 		KvNode:      KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 		LeftValue:   4,
// 		RightValue:  7,
// 		LeftCmpSym:  GT,
// 		RightCmpSym: LT,
// 	}

// 	node2, _ = node1.Inverse()
// 	assert.Equal(t, &OrNode{
// 		MinimumShouldMatch: 1,
// 		Nodes: map[string][]AstNode{
// 			node1.NodeKey(): {
// 				&RangeNode{
// 					KvNode:      node1.KvNode,
// 					LeftValue:   MinInt[32],
// 					RightValue:  4,
// 					LeftCmpSym:  GT,
// 					RightCmpSym: LTE,
// 				},
// 				&RangeNode{
// 					KvNode:      node1.KvNode,
// 					LeftValue:   7,
// 					RightValue:  MaxInt[32],
// 					LeftCmpSym:  GTE,
// 					RightCmpSym: LT,
// 				},
// 			},
// 		},
// 	}, node2)

// }

// func TestRangeNodeToDsl(t *testing.T) {
// 	var node1 = &RangeNode{
// 		KvNode:    KvNode{Field: "foo", Type: mapping.INTEGER_FIELD_TYPE},
// 		LeftValue: 1, LeftCmpSym: GT,
// 		RightValue: 7, RightCmpSym: LT,
// 		Boost: 1.0,
// 	}

// 	assert.Equal(t, DSL{
// 		"range": DSL{
// 			"foo": DSL{
// 				GT.String(): 1,
// 				LT.String(): 7,
// 				"relation":  "WITHIN",
// 				"boost":     1.0,
// 			},
// 		},
// 	}, node1.ToDSL())

// 	var node2 = &RangeNode{
// 		KvNode:    KvNode{Field: "foo", Type: mapping.DATE_FIELD_TYPE},
// 		LeftValue: time.Date(2022, 01, 02, 0, 0, 0, 0, time.UTC), LeftCmpSym: GT,
// 		RightValue: time.Date(2022, 01, 03, 0, 0, 0, 0, time.UTC), RightCmpSym: LT,
// 		Boost: 1.0,
// 	}

// 	assert.Equal(t, DSL{
// 		"range": DSL{
// 			"foo": DSL{
// 				GT.String(): leafValueToPrintValue(node2.LeftValue, node2.Type),
// 				LT.String(): leafValueToPrintValue(node2.RightValue, node2.Type),
// 				"format":    "epoch_millis",
// 				"relation":  "WITHIN",
// 				"boost":     1.0,
// 			},
// 		},
// 	}, node2.ToDSL())

// 	var node3 = &RangeNode{
// 		KvNode:    KvNode{Field: "foo", Type: mapping.IP_FIELD_TYPE},
// 		LeftValue: net.IP([]byte{1, 2, 3, 4}), LeftCmpSym: GT,
// 		RightValue: net.IP([]byte{1, 2, 4, 5}), RightCmpSym: LT,
// 		Boost: 1.0,
// 	}

// 	assert.Equal(t, DSL{
// 		"range": DSL{
// 			"foo": DSL{
// 				GT.String(): leafValueToPrintValue(node3.LeftValue, node3.Type),
// 				LT.String(): leafValueToPrintValue(node3.RightValue, node3.Type),
// 				"relation":  "WITHIN",
// 				"boost":     1.0,
// 			},
// 		},
// 	}, node3.ToDSL())

// 	var v1, _ = version.NewVersion("1.2.3")
// 	var v2, _ = version.NewVersion("1.2.10")
// 	var node4 = &RangeNode{
// 		KvNode:    KvNode{Field: "foo", Type: mapping.VERSION_FIELD_TYPE},
// 		LeftValue: v1, LeftCmpSym: GT,
// 		RightValue: v2, RightCmpSym: LT,
// 		Boost: 1.0,
// 	}

// 	assert.Equal(t, DSL{
// 		"range": DSL{
// 			"foo": DSL{
// 				GT.String(): leafValueToPrintValue(node4.LeftValue, node4.Type),
// 				LT.String(): leafValueToPrintValue(node4.RightValue, node4.Type),
// 				"relation":  "WITHIN",
// 				"boost":     1.0,
// 			},
// 		},
// 	}, node4.ToDSL())
// }
