package dsl

import (
	"math"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/x448/float16"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

func TestUnionJoinStrLst(t *testing.T) {
	type args struct {
		al []string
		bl []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestUnionJoinStrLst01",
			args: args{
				al: []string{"1", "2"},
				bl: []string{"3", "2"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "TestUnionJoinStrLst02",
			args: args{
				al: []string{"2"},
				bl: []string{"3", "2"},
			},
			want: []string{"2", "3"},
		},
		{
			name: "TestUnionJoinStrLst03",
			args: args{
				al: []string{"2", "3"},
				bl: []string{"2"},
			},
			want: []string{"2", "3"},
		},
		{
			name: "TestUnionJoinStrLst04",
			args: args{
				al: []string{"2", "3"},
				bl: []string{},
			},
			want: []string{"2", "3"},
		},
		{
			name: "TestUnionJoinStrLst05",
			args: args{
				al: []string{},
				bl: []string{"2", "3"},
			},
			want: []string{"2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnionJoinStrLst(tt.args.al, tt.args.bl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnionJoinStrLst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectStrLst(t *testing.T) {
	type args struct {
		al []string
		bl []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestIntersectStrLst01",
			args: args{
				al: []string{"1", "2"},
				bl: []string{"1", "3", "2"},
			},
			want: []string{"1", "2"},
		},
		{
			name: "TestIntersectStrLst02",
			args: args{
				al: []string{"2"},
				bl: []string{"3", "2"},
			},
			want: []string{"2"},
		},
		{
			name: "TestIntersectStrLst03",
			args: args{
				al: []string{"2", "2", "3"},
				bl: []string{"2"},
			},
			want: []string{"2"},
		},
		{
			name: "TestIntersectStrLst04",
			args: args{
				al: []string{"2", "2", "3", "1"},
				bl: []string{"2", "2"},
			},
			want: []string{"2"},
		},
		{
			name: "TestIntersectStrLst05",
			args: args{
				al: []string{"2", "2", "3", "1"},
				bl: []string{},
			},
			want: []string{},
		},
		{
			name: "TestIntersectStrLst06",
			args: args{
				al: []string{},
				bl: []string{"2", "2", "3", "1"},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntersectStrLst(tt.args.al, tt.args.bl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntersectStrLst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqStrLst(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestUniqStrLst01",
			args: args{
				a: []string{"1", "1", "2", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "TestUniqStrLst02",
			args: args{
				a: []string{"1", "2", "2", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "TestUniqStrLst03",
			args: args{
				a: []string{"1", "2", "2", "2", "2"},
			},
			want: []string{"1", "2"},
		},
		{
			name: "TestUniqStrLst04",
			args: args{
				a: []string{"1", "1", "3", "3", "3"},
			},
			want: []string{"1", "3"},
		},
		{
			name: "TestUniqStrLst05",
			args: args{
				a: []string{"1", "1", "2", "3", "3"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "TestUniqStrLst06",
			args: args{
				a: []string{"1", "1"},
			},
			want: []string{"1"},
		},
		{
			name: "TestUniqStrLst07",
			args: args{
				a: []string{"1"},
			},
			want: []string{"1"},
		},
		{
			name: "TestUniqStrLst08",
			args: args{
				a: []string{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqStrLst(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqStrLst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareAny(t *testing.T) {
	type args struct {
		a   *LeafValue
		b   *LeafValue
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
			args: args{a: &LeafValue{TinyInt: 1}, b: &LeafValue{TinyInt: 2}, typ: mapping.SHORT_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareInt02",
			args: args{a: &LeafValue{TinyInt: 2}, b: &LeafValue{TinyInt: 1}, typ: mapping.SHORT_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareInt03",
			args: args{a: &LeafValue{TinyInt: 1}, b: &LeafValue{TinyInt: 1}, typ: mapping.SHORT_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareFloat01",
			args: args{a: &LeafValue{Float16: float16.Fromfloat32(1.13)}, b: &LeafValue{Float16: float16.Fromfloat32(2.1)}, typ: mapping.HALF_FLOAT_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareFloat02",
			args: args{a: &LeafValue{Float16: float16.Fromfloat32(2.1)}, b: &LeafValue{Float16: float16.Fromfloat32(1.33)}, typ: mapping.HALF_FLOAT_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareFloat03",
			args: args{a: &LeafValue{Float16: float16.Fromfloat32(1.3)}, b: &LeafValue{Float16: float16.Fromfloat32(1.3)}, typ: mapping.HALF_FLOAT_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareDate01",
			args: args{a: &LeafValue{DateTime: time.Unix(1, 0)}, b: &LeafValue{DateTime: time.Unix(10, 0)}, typ: mapping.DATE_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareDate02",
			args: args{a: &LeafValue{DateTime: time.Unix(10, 0)}, b: &LeafValue{DateTime: time.Unix(1, 0)}, typ: mapping.DATE_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareDate03",
			args: args{a: &LeafValue{DateTime: time.Unix(1, 0)}, b: &LeafValue{DateTime: time.Unix(1, 0)}, typ: mapping.DATE_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareIp01",
			args: args{a: &LeafValue{IpValue: net.ParseIP("12.23.1.1")}, b: &LeafValue{IpValue: net.ParseIP("12.200.1.1")}, typ: mapping.IP_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareIp02",
			args: args{a: &LeafValue{IpValue: net.ParseIP("12.200.1.1")}, b: &LeafValue{IpValue: net.ParseIP("12.23.1.1")}, typ: mapping.IP_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareIp03",
			args: args{a: &LeafValue{IpValue: net.ParseIP("127.0.0.1")}, b: &LeafValue{IpValue: net.ParseIP("127.0.0.001")}, typ: mapping.IP_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareString01",
			args: args{a: &LeafValue{String: "12.23.1.1"}, b: &LeafValue{String: "12.200.1.1"}, typ: mapping.KEYWORD_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareString02",
			args: args{a: &LeafValue{String: "12.200.1.1"}, b: &LeafValue{String: "12.23.1.1"}, typ: mapping.KEYWORD_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareString03",
			args: args{a: &LeafValue{String: "127.0.0.1"}, b: &LeafValue{String: "127.0.0.1"}, typ: mapping.KEYWORD_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareVersion01",
			args: args{a: &LeafValue{Version: version1}, b: &LeafValue{Version: version2}, typ: mapping.VERSION_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareVersion02",
			args: args{a: &LeafValue{Version: version1}, b: &LeafValue{Version: version3}, typ: mapping.VERSION_FIELD_TYPE},
			want: 1,
		},
		{
			name: "CompareVersion03",
			args: args{a: &LeafValue{Version: version3}, b: &LeafValue{Version: version2}, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareVersion04",
			args: args{a: &LeafValue{Version: version4}, b: &LeafValue{Version: version6}, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareVersion05",
			args: args{a: &LeafValue{Version: version5}, b: &LeafValue{Version: version6}, typ: mapping.VERSION_FIELD_TYPE},
			want: 0,
		},
		{
			name: "CompareVersion06",
			args: args{a: &LeafValue{Version: version4}, b: &LeafValue{Version: version5}, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
		{
			name: "CompareVersion07",
			args: args{a: &LeafValue{Version: version4}, b: &LeafValue{Version: version7}, typ: mapping.VERSION_FIELD_TYPE},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareAny(tt.args.a, tt.args.b, tt.args.typ); got != tt.want {
				t.Errorf("CompareAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareInt64(t *testing.T) {
	type args struct {
		a int64
		b int64
		c CompareType
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test_lt",
			args: args{a: 1, b: 2, c: LT},
			want: 1,
		},
		{
			name: "test_gt",
			args: args{a: 1, b: 2, c: GT},
			want: 2,
		},
		{
			name: "test_lte",
			args: args{a: 1, b: 2, c: LTE},
			want: 1,
		},
		{
			name: "test_gte",
			args: args{a: 1, b: 2, c: GTE},
			want: 2,
		},
		{
			name: "test_eq_1",
			args: args{a: 3, b: 2, c: EQ},
			want: 3,
		},
		{
			name: "test_eq_2",
			args: args{a: 2, b: 3, c: EQ},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareInt64(tt.args.a, tt.args.b, tt.args.c); got != tt.want {
				t.Errorf("compareInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareUInt64(t *testing.T) {
	type args struct {
		a uint64
		b uint64
		c CompareType
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "test_lt",
			args: args{a: 1, b: 2, c: LT},
			want: 1,
		},
		{
			name: "test_gt",
			args: args{a: 1, b: 2, c: GT},
			want: 2,
		},
		{
			name: "test_lte",
			args: args{a: 1, b: 2, c: LTE},
			want: 1,
		},
		{
			name: "test_gte",
			args: args{a: 1, b: 2, c: GTE},
			want: 2,
		},
		{
			name: "test_eq_1",
			args: args{a: 3, b: 2, c: EQ},
			want: 3,
		},
		{
			name: "test_eq_2",
			args: args{a: 2, b: 3, c: EQ},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareUInt64(tt.args.a, tt.args.b, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareUInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
		c CompareType
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test_lt",
			args: args{a: 1, b: 2, c: LT},
			want: 1,
		},
		{
			name: "test_gt",
			args: args{a: 1, b: 2, c: GT},
			want: 2,
		},
		{
			name: "test_lte",
			args: args{a: 1, b: 2, c: LTE},
			want: 1,
		},
		{
			name: "test_gte",
			args: args{a: 1, b: 2, c: GTE},
			want: 2,
		},
		{
			name: "test_eq_1",
			args: args{a: 3, b: 2, c: EQ},
			want: 3,
		},
		{
			name: "test_eq_2",
			args: args{a: 2, b: 3, c: EQ},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareFloat64(tt.args.a, tt.args.b, tt.args.c); math.Abs(got-tt.want) > 1E-6 {
				t.Errorf("compareFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareIp(t *testing.T) {
	type args struct {
		a net.IP
		b net.IP
		c CompareType
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		{
			name: "test_lt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: LT},
			want: net.ParseIP("1.2.3.4"),
		},
		{
			name: "test_gt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: GT},
			want: net.ParseIP("1.2.4.3"),
		},
		{
			name: "test_lt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: LTE},
			want: net.ParseIP("1.2.3.4"),
		},
		{
			name: "test_gt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: GTE},
			want: net.ParseIP("1.2.4.3"),
		},
		{
			name: "test_eq_1",
			args: args{a: net.ParseIP("1.5.3.4"), b: net.ParseIP("1.2.4.3"), c: EQ},
			want: net.ParseIP("1.5.3.4"),
		},
		{
			name: "test_eq_2",
			args: args{a: net.ParseIP("1.2.4.3"), b: net.ParseIP("1.5.3.4"), c: EQ},
			want: net.ParseIP("1.2.4.3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareIp(tt.args.a, tt.args.b, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareString(t *testing.T) {
	type args struct {
		a string
		b string
		c CompareType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_lt",
			args: args{a: "1", b: "2", c: LT},
			want: "1",
		},
		{
			name: "test_gt",
			args: args{a: "1", b: "2", c: GT},
			want: "2",
		},
		{
			name: "test_lte",
			args: args{a: "1", b: "2", c: LTE},
			want: "1",
		},
		{
			name: "test_gte",
			args: args{a: "1", b: "2", c: GTE},
			want: "2",
		},
		{
			name: "test_eq_1",
			args: args{a: "3", b: "2", c: EQ},
			want: "3",
		},
		{
			name: "test_eq_2",
			args: args{a: "2", b: "3", c: EQ},
			want: "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareString(tt.args.a, tt.args.b, tt.args.c); got != tt.want {
				t.Errorf("compareString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareDate(t *testing.T) {
	type args struct {
		a time.Time
		b time.Time
		c CompareType
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "test_lt",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: LT},
			want: time.Unix(1, 0),
		},
		{
			name: "test_gt",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: GT},
			want: time.Unix(2, 0),
		},
		{
			name: "test_lte",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: LTE},
			want: time.Unix(1, 0),
		},
		{
			name: "test_gte",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: GTE},
			want: time.Unix(2, 0),
		},
		{
			name: "test_eq_1",
			args: args{a: time.Unix(3, 0), b: time.Unix(2, 0), c: EQ},
			want: time.Unix(3, 0),
		},
		{
			name: "test_eq_2",
			args: args{a: time.Unix(2, 0), b: time.Unix(3, 0), c: EQ},
			want: time.Unix(2, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareDate(tt.args.a, tt.args.b, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compareDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
