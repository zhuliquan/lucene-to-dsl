package dsl

import (
	"net"
	"reflect"
	"testing"
	"time"
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
		a   *DSLTermValue
		b   *DSLTermValue
		typ DSLTermType
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "CompareInt01",
			args: args{a: &DSLTermValue{IntTerm: 1}, b: &DSLTermValue{IntTerm: 2}, typ: INT_VALUE},
			want: -1,
		},
		{
			name: "CompareInt02",
			args: args{a: &DSLTermValue{IntTerm: 2}, b: &DSLTermValue{IntTerm: 1}, typ: INT_VALUE},
			want: 1,
		},
		{
			name: "CompareInt03",
			args: args{a: &DSLTermValue{IntTerm: 1}, b: &DSLTermValue{IntTerm: 1}, typ: INT_VALUE},
			want: 0,
		},
		{
			name: "CompareFloat01",
			args: args{a: &DSLTermValue{FloatTerm: 1.13}, b: &DSLTermValue{FloatTerm: 2.1}, typ: FLOAT_VALUE},
			want: -1,
		},
		{
			name: "CompareFloat02",
			args: args{a: &DSLTermValue{FloatTerm: 2.1}, b: &DSLTermValue{FloatTerm: 1.33}, typ: FLOAT_VALUE},
			want: 1,
		},
		{
			name: "CompareFloat03",
			args: args{a: &DSLTermValue{FloatTerm: 1.3}, b: &DSLTermValue{FloatTerm: 1.3}, typ: FLOAT_VALUE},
			want: 0,
		},
		{
			name: "CompareDate01",
			args: args{a: &DSLTermValue{DateTerm: time.Unix(1, 0)}, b: &DSLTermValue{DateTerm: time.Unix(10, 0)}, typ: DATE_VALUE},
			want: -1,
		},
		{
			name: "CompareDate02",
			args: args{a: &DSLTermValue{DateTerm: time.Unix(10, 0)}, b: &DSLTermValue{DateTerm: time.Unix(1, 0)}, typ: DATE_VALUE},
			want: 1,
		},
		{
			name: "CompareDate03",
			args: args{a: &DSLTermValue{DateTerm: time.Unix(1, 0)}, b: &DSLTermValue{DateTerm: time.Unix(1, 0)}, typ: DATE_VALUE},
			want: 0,
		},
		{
			name: "CompareIp01",
			args: args{a: &DSLTermValue{IpTerm: net.ParseIP("12.23.1.1")}, b: &DSLTermValue{IpTerm: net.ParseIP("12.200.1.1")}, typ: IP_VALUE},
			want: -1,
		},
		{
			name: "CompareIp02",
			args: args{a: &DSLTermValue{IpTerm: net.ParseIP("12.200.1.1")}, b: &DSLTermValue{IpTerm: net.ParseIP("12.23.1.1")}, typ: IP_VALUE},
			want: 1,
		},
		{
			name: "CompareIp03",
			args: args{a: &DSLTermValue{IpTerm: net.ParseIP("127.0.0.1")}, b: &DSLTermValue{IpTerm: net.ParseIP("127.0.0.001")}, typ: IP_VALUE},
			want: 0,
		},
		{
			name: "CompareString01",
			args: args{a: &DSLTermValue{StringTerm: "12.23.1.1"}, b: &DSLTermValue{StringTerm: "12.200.1.1"}, typ: KEYWORD_VALUE},
			want: 1,
		},
		{
			name: "CompareString02",
			args: args{a: &DSLTermValue{StringTerm: "12.200.1.1"}, b: &DSLTermValue{StringTerm: "12.23.1.1"}, typ: KEYWORD_VALUE},
			want: -1,
		},
		{
			name: "CompareString03",
			args: args{a: &DSLTermValue{StringTerm: "127.0.0.1"}, b: &DSLTermValue{StringTerm: "127.0.0.1"}, typ: KEYWORD_VALUE},
			want: 0,
		},
		{
			name: "CompareInf01",
			args: args{a: &DSLTermValue{StringTerm: "12.23.1.1"}, b: InfValue, typ: KEYWORD_VALUE},
			want: -1,
		},
		{
			name: "CompareInf02",
			args: args{a: InfValue, b: &DSLTermValue{IntTerm: 1}, typ: KEYWORD_VALUE},
			want: 1,
		},
		{
			name: "CompareInf03",
			args: args{a: InfValue, b: InfValue, typ: KEYWORD_VALUE},
			want: 0,
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
