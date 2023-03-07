package dsl

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/x448/float16"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestFindAny(t *testing.T) {
	assert.Equal(t, FindAny([]LeafValue{int64(1), int64(2)}, int64(2), mapping.INTEGER_FIELD_TYPE), 1)
	assert.Equal(t, FindAny([]LeafValue{int64(1), int64(2)}, int64(3), mapping.INTEGER_FIELD_TYPE), -1)
}

func TestBinaryFindAny(t *testing.T) {
	assert.Equal(t, BinaryFindAny([]LeafValue{int64(1), int64(2), int64(2), int64(3)}, int64(2), mapping.INTEGER_FIELD_TYPE), 1)
	assert.Equal(t, BinaryFindAny([]LeafValue{int64(1), int64(2)}, int64(3), mapping.INTEGER_FIELD_TYPE), -1)
	assert.Equal(t, BinaryFindAny([]LeafValue{int64(1), int64(2), int64(4)}, int64(3), mapping.INTEGER_FIELD_TYPE), -1)
}

func TestUnionJoinStrLst(t *testing.T) {
	type args struct {
		al  []LeafValue
		bl  []LeafValue
		typ mapping.FieldType
	}
	tests := []struct {
		name string
		args args
		want []LeafValue
	}{
		{
			name: "TestUnionJoinStrLst01",
			args: args{
				al:  []LeafValue{"1", "2"},
				bl:  []LeafValue{"3", "2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "2", "3"},
		},
		{
			name: "TestUnionJoinStrLst02",
			args: args{
				al:  []LeafValue{"2"},
				bl:  []LeafValue{"3", "2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2", "3"},
		},
		{
			name: "TestUnionJoinStrLst03",
			args: args{
				al:  []LeafValue{"2", "3"},
				bl:  []LeafValue{"2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2", "3"},
		},
		{
			name: "TestUnionJoinStrLst04",
			args: args{
				al:  []LeafValue{"2", "3"},
				bl:  []LeafValue{},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2", "3"},
		},
		{
			name: "TestUnionJoinStrLst05",
			args: args{
				al:  []LeafValue{},
				bl:  []LeafValue{"2", "3"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnionJoinValueLst(tt.args.al, tt.args.bl, tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIntersectStrLst(t *testing.T) {
	type args struct {
		al  []LeafValue
		bl  []LeafValue
		typ mapping.FieldType
	}
	tests := []struct {
		name string
		args args
		want []LeafValue
	}{
		{
			name: "TestIntersectStrLst01",
			args: args{
				al:  []LeafValue{"1", "2"},
				bl:  []LeafValue{"1", "3", "2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "2"},
		},
		{
			name: "TestIntersectStrLst02",
			args: args{
				al:  []LeafValue{"2"},
				bl:  []LeafValue{"3", "2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2"},
		},
		{
			name: "TestIntersectStrLst03",
			args: args{
				al:  []LeafValue{"2", "2", "3"},
				bl:  []LeafValue{"2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2"},
		},
		{
			name: "TestIntersectStrLst04",
			args: args{
				al:  []LeafValue{"2", "2", "3", "1"},
				bl:  []LeafValue{"2", "2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"2"},
		},
		{
			name: "TestIntersectStrLst05",
			args: args{
				al:  []LeafValue{"2", "2", "3", "1"},
				bl:  []LeafValue{},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{},
		},
		{
			name: "TestIntersectStrLst06",
			args: args{
				al:  []LeafValue{},
				bl:  []LeafValue{"2", "2", "3", "1"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntersectValueLst(tt.args.al, tt.args.bl, tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUniqValueLst(t *testing.T) {
	type args struct {
		a   []LeafValue
		typ mapping.FieldType
	}
	tests := []struct {
		name string
		args args
		want []LeafValue
	}{
		{
			name: "TestUniqStrLst01",
			args: args{
				a:   []LeafValue{"1", "1", "2", "2", "3"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "2", "3"},
		},
		{
			name: "TestUniqStrLst02",
			args: args{
				a:   []LeafValue{"1", "2", "2", "2", "3"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "2", "3"},
		},
		{
			name: "TestUniqStrLst03",
			args: args{
				a:   []LeafValue{"1", "2", "2", "2", "2"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "2"},
		},
		{
			name: "TestUniqStrLst04",
			args: args{
				a:   []LeafValue{"1", "1", "3", "3", "3"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "3"},
		},
		{
			name: "TestUniqStrLst05",
			args: args{
				a:   []LeafValue{"1", "1", "2", "3", "3"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1", "2", "3"},
		},
		{
			name: "TestUniqStrLst06",
			args: args{
				a:   []LeafValue{"1", "1"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1"},
		},
		{
			name: "TestUniqStrLst07",
			args: args{
				a:   []LeafValue{"1"},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{"1"},
		},
		{
			name: "TestUniqStrLst08",
			args: args{
				a:   []LeafValue{},
				typ: mapping.KEYWORD_FIELD_TYPE,
			},
			want: []LeafValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UniqValueLst(tt.args.a, tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCompareAny(t *testing.T) {
	type args struct {
		a   LeafValue
		b   LeafValue
		typ mapping.FieldType
	}
	version1, _ := version.NewVersion("v1.1.0")
	version2, _ := version.NewVersion("1.1.0")
	version3, _ := version.NewVersion("v1.1.0-rc")
	version4, _ := version.NewVersion("v0-A.0-A.0-A")
	version5, _ := version.NewVersion("0")
	version6, _ := version.NewVersion("v0.0.0")
	version7, _ := version.NewVersion("v0-a.0-a.0-a")
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "CompareInt01",
			args: args{a: int64(1), b: int64(2), typ: mapping.SHORT_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareInt02",
			args: args{a: int64(2), b: int64(1), typ: mapping.SHORT_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareInt03",
			args: args{a: int64(1), b: int64(1), typ: mapping.SHORT_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareFloat01",
			args: args{a: float16.Fromfloat32(1.13), b: float16.Fromfloat32(2.1), typ: mapping.HALF_FLOAT_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareFloat02",
			args: args{a: float16.Fromfloat32(2.1), b: float16.Fromfloat32(1.33), typ: mapping.HALF_FLOAT_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareFloat03",
			args: args{a: float16.Fromfloat32(1.3), b: float16.Fromfloat32(1.3), typ: mapping.HALF_FLOAT_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareDate01",
			args: args{a: time.Unix(1, 0), b: time.Unix(10, 0), typ: mapping.DATE_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareDate02",
			args: args{a: time.Unix(10, 0), b: time.Unix(1, 0), typ: mapping.DATE_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareDate03",
			args: args{a: time.Unix(1, 0), b: time.Unix(1, 0), typ: mapping.DATE_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareIp01",
			args: args{a: net.ParseIP("12.23.1.1"), b: net.ParseIP("12.200.1.1"), typ: mapping.IP_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareIp02",
			args: args{a: net.ParseIP("12.200.1.1"), b: net.ParseIP("12.23.1.1"), typ: mapping.IP_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareIp03",
			args: args{a: net.ParseIP("127.0.0.1"), b: net.ParseIP("127.0.0.1"), typ: mapping.IP_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareString01",
			args: args{a: "12.23.1.1", b: "12.200.1.1", typ: mapping.KEYWORD_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareString02",
			args: args{a: "12.200.1.1", b: "12.23.1.1", typ: mapping.KEYWORD_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareString03",
			args: args{a: "127.0.0.1", b: "127.0.0.1", typ: mapping.KEYWORD_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareVersion01",
			args: args{a: version1, b: version2, typ: mapping.VERSION_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareVersion02",
			args: args{a: version1, b: version3, typ: mapping.VERSION_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareVersion03",
			args: args{a: version3, b: version2, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareVersion04",
			args: args{a: version4, b: version6, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareVersion05",
			args: args{a: version5, b: version6, typ: mapping.VERSION_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareVersion06",
			args: args{a: version4, b: version5, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareVersion07",
			args: args{a: version4, b: version7, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareAny(tt.args.a, tt.args.b, tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckValidRangeNode(t *testing.T) {
	type args struct {
		node *RangeNode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_error_01",
			args: args{
				node: &RangeNode{
					rgNode: rgNode{
						valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true},
						rValue:    "1",
						lValue:    "1",
						rCmpSym:   LTE,
						lCmpSym:   GT,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test_error_02",
			args: args{
				node: &RangeNode{
					rgNode: rgNode{
						valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true},
						rValue:    "1",
						lValue:    "1",
						rCmpSym:   LT,
						lCmpSym:   GTE,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test_error_03",
			args: args{
				node: &RangeNode{
					rgNode: rgNode{
						valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true},
						rValue:    "1",
						lValue:    "1",
						rCmpSym:   LT,
						lCmpSym:   GT,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test_error_04",
			args: args{
				node: &RangeNode{
					rgNode: rgNode{
						valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true},
						rValue:    "1",
						lValue:    "2",
						rCmpSym:   LTE,
						lCmpSym:   GT,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test_ok_01",
			args: args{
				node: &RangeNode{
					rgNode: rgNode{
						valueType: valueType{mapping.KEYWORD_FIELD_TYPE, true},
						rValue:    "2",
						lValue:    "1",
						rCmpSym:   LTE,
						lCmpSym:   GT,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckValidRangeNode(tt.args.node)
			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}

func TestCastInt(t *testing.T) {
	type args struct {
		x LeafValue
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test_cast_int",
			args: args{x: 1},
			want: int64(1),
		},
		{
			name: "test_cast_uint",
			args: args{x: uint(1)},
			want: int64(1),
		},
		{
			name: "test_cast_int",
			args: args{x: 1},
			want: int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := castInt(tt.args.x); got != tt.want {
				t.Errorf("castInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastUInt(t *testing.T) {
	type args struct {
		x LeafValue
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "test_int",
			args: args{x: 1},
			want: uint64(1),
		},
		{
			name: "test_uint",
			args: args{x: uint(1)},
			want: uint64(1),
		},
		{
			name: "test_uint64",
			args: args{x: uint64(1)},
			want: uint64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := castUInt(tt.args.x); got != tt.want {
				t.Errorf("castUInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMinInf(t *testing.T) {
	type args struct {
		a LeafValue
		t mapping.FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "tes_inf_01",
			args: args{a: MinInt[8], t: mapping.BYTE_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_02",
			args: args{a: MinInt[8] + 1, t: mapping.BYTE_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_03",
			args: args{a: MinInt[16], t: mapping.SHORT_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_04",
			args: args{a: MinInt[16] + 1, t: mapping.SHORT_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_05",
			args: args{a: MinInt[32], t: mapping.INTEGER_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_06",
			args: args{a: MinInt[32] + 1, t: mapping.INTEGER_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_07",
			args: args{a: MinInt[64], t: mapping.LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_08",
			args: args{a: MinInt[64] + 1, t: mapping.LONG_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_09",
			args: args{a: MinUint, t: mapping.UNSIGNED_LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_10",
			args: args{a: MinUint + 1, t: mapping.UNSIGNED_LONG_FIELD_TYPE},
			want: false,
		},
		{
			name: "test_inf_11",
			args: args{a: MinIP, t: mapping.IP_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_12",
			args: args{a: MinVersion, t: mapping.VERSION_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_13",
			args: args{a: MinFloat16, t: mapping.HALF_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_14",
			args: args{a: MinFloat[32], t: mapping.FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_15",
			args: args{a: MinFloat[64], t: mapping.DOUBLE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_16",
			args: args{a: MinFloat[128], t: mapping.SCALED_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_17",
			args: args{a: MinTime, t: mapping.DATE_FIELD_TYPE},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMinInf(tt.args.a, tt.args.t); got != tt.want {
				t.Errorf("isMinInf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMaxInf(t *testing.T) {
	type args struct {
		a LeafValue
		t mapping.FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "tes_inf_01",
			args: args{a: MaxInt[8], t: mapping.BYTE_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_02",
			args: args{a: MaxInt[8] + 1, t: mapping.BYTE_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_03",
			args: args{a: MaxInt[16], t: mapping.SHORT_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_04",
			args: args{a: MaxInt[16] - 1, t: mapping.SHORT_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_05",
			args: args{a: MaxInt[32], t: mapping.INTEGER_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_06",
			args: args{a: MaxInt[32] - 1, t: mapping.INTEGER_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_07",
			args: args{a: MaxInt[64], t: mapping.LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_08",
			args: args{a: MaxInt[64] - 1, t: mapping.LONG_FIELD_TYPE},
			want: false,
		},
		{
			name: "tes_inf_09",
			args: args{a: MaxUint[64], t: mapping.UNSIGNED_LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "tes_inf_10",
			args: args{a: MaxUint[64] - 1, t: mapping.UNSIGNED_LONG_FIELD_TYPE},
			want: false,
		},
		{
			name: "test_inf_11",
			args: args{a: MaxIP, t: mapping.IP_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_12",
			args: args{a: MaxVersion, t: mapping.VERSION_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_13",
			args: args{a: MaxFloat16, t: mapping.HALF_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_14",
			args: args{a: MaxFloat[32], t: mapping.FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_15",
			args: args{a: MaxFloat[64], t: mapping.DOUBLE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_16",
			args: args{a: MaxFloat[128], t: mapping.SCALED_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_inf_17",
			args: args{a: MaxTime, t: mapping.DATE_FIELD_TYPE},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMaxInf(tt.args.a, tt.args.t); got != tt.want {
				t.Errorf("isMaxInf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeOrMergeNode(t *testing.T) {
	type args struct {
		a AstNode
		b AstNode
	}
	tests := []struct {
		name    string
		args    args
		want    AstNode
		wantErr bool
	}{
		{
			name: "test_or_node",
			args: args{
				a: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo"},
						valueType: valueType{mapping.INTEGER_FIELD_TYPE, true},
						lValue:    MinInt[32],
						rValue:    3,
						lCmpSym:   GT,
						rCmpSym:   LTE,
					},
				},
				b: &RangeNode{
					rgNode: rgNode{
						fieldNode: fieldNode{field: "foo"},
						valueType: valueType{mapping.INTEGER_FIELD_TYPE, true},
						lValue:    4,
						rValue:    MaxInt[32],
						lCmpSym:   GTE,
						rCmpSym:   LT,
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: OR},
				Should: map[string][]AstNode{
					"foo": {
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{field: "foo"},
								valueType: valueType{mapping.INTEGER_FIELD_TYPE, true},
								lValue:    MinInt[32],
								rValue:    3,
								lCmpSym:   GT,
								rCmpSym:   LTE,
							},
						},
						&RangeNode{
							rgNode: rgNode{
								fieldNode: fieldNode{field: "foo"},
								valueType: valueType{mapping.INTEGER_FIELD_TYPE, true},
								lValue:    4,
								rValue:    MaxInt[32],
								lCmpSym:   GTE,
								rCmpSym:   LT,
							},
						},
					},
				},
				MinimumShouldMatch: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lfNodeUnionJoinLfNode(tt.args.a.NodeKey(), tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeOrMergeNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeOrMergeNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeAndMergeNode(t *testing.T) {
	type args struct {
		a AstNode
		b AstNode
	}
	tests := []struct {
		name    string
		args    args
		want    AstNode
		wantErr bool
	}{
		{
			name: "test_and_node_01",
			args: args{
				a: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
						valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
				b: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
						valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: AND},
				Filter: map[string][]AstNode{
					"foo": {
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
								valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
								valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_and_node_02",
			args: args{
				a: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
						valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
				b: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
						valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: AND},
				Must: map[string][]AstNode{
					"foo": {
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
								valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
								valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_and_node_03",
			args: args{
				a: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
						valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
				b: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
						valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: AND},
				Must: map[string][]AstNode{
					"foo": {
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
								valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
					},
				},
				Filter: map[string][]AstNode{
					"foo": {
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
								valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test_and_node_04",
			args: args{
				a: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
						valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
				b: &TermNode{
					kvNode: kvNode{
						fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
						valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
					},
				},
			},
			want: &BoolNode{
				opNode: opNode{opType: AND},
				Must: map[string][]AstNode{
					"foo": {
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{false}}, field: "foo"},
								valueNode: valueNode{value: "bar2", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
					},
				},
				Filter: map[string][]AstNode{
					"foo": {
						&TermNode{
							kvNode: kvNode{
								fieldNode: fieldNode{lfNode: lfNode{filterCtxNode: filterCtxNode{true}}, field: "foo"},
								valueNode: valueNode{value: "bar1", valueType: valueType{mapping.TEXT_FIELD_TYPE, false}},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lfNodeIntersectLfNode(tt.args.a.NodeKey(), tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("nodeAndMergeNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeAndMergeNode() = %v, want %v", got.ToDSL(), tt.want.ToDSL())
			}
		})
	}
}

func TestMinEditDistance(t *testing.T) {
	type args struct {
		termWord1 string
		termWord2 string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_change",
			args: args{termWord1: "abc", termWord2: "abd"},
			want: 1,
		},
		{
			name: "test_swap",
			args: args{termWord1: "abc", termWord2: "acb"},
			want: 2,
		},
		{
			name: "test_insert",
			args: args{termWord1: "ab", termWord2: "abc"},
			want: 1,
		},
		{
			name: "test_remove",
			args: args{termWord1: "bc", termWord2: "abc"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minEditDistance(tt.args.termWord1, tt.args.termWord2); got != tt.want {
				t.Errorf("minEditDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestAstNodes(t *testing.T) {
	var n1 = &ExistsNode{fieldNode: fieldNode{field: "foo1"}}
	var n2 = &ExistsNode{fieldNode: fieldNode{field: "foo2"}}
	var n3 = &ExistsNode{fieldNode: fieldNode{field: "foo3"}}
	type args struct {
		nodes []AstNode
		index int
	}
	tests := []struct {
		name  string
		args  args
		want1 AstNode
		want2 []AstNode
	}{
		{
			name:  "test_empty",
			args:  args{nil, 1},
			want1: nil,
			want2: nil,
		},
		{
			name:  "test_one",
			args:  args{nodes: []AstNode{&EmptyNode{}}, index: 1},
			want1: &EmptyNode{},
			want2: []AstNode{},
		},
		{
			name:  "test_three_nodes_get_first",
			args:  args{nodes: []AstNode{n1, n2, n3}, index: 0},
			want1: n1,
			want2: []AstNode{n2, n3},
		},
		{
			name:  "test_three_nodes_get_last",
			args:  args{nodes: []AstNode{n1, n2, n3}, index: 2},
			want1: n3,
			want2: []AstNode{n1, n2},
		},
		{
			name:  "test_three_nodes_get_mid",
			args:  args{nodes: []AstNode{n1, n2, n3}, index: 1},
			want1: n2,
			want2: []AstNode{n1, n3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := restAstNodes(tt.args.nodes, tt.args.index)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
		})
	}
}
