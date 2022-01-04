package dsl

import (
	"reflect"
	"testing"
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
