package mapping

import "testing"

func Test_strLstHasPrefix(t *testing.T) {
	type args struct {
		va []string
		vb []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "strLstHasPrefix01",
			args: args{va: nil, vb: nil},
			want: true,
		},
		{
			name: "strLstHasPrefix02",
			args: args{va: []string{}, vb: nil},
			want: true,
		},
		{
			name: "strLstHasPrefix03",
			args: args{va: nil, vb: []string{}},
			want: true,
		},
		{
			name: "strLstHasPrefix04",
			args: args{va: []string{"1"}, vb: []string{"1", "2"}},
			want: false,
		},
		{
			name: "strLstHashPrefix05",
			args: args{va: []string{"1", "2"}, vb: []string{"1"}},
			want: true,
		},
		{
			name: "strLstHashPrefix06",
			args: args{va: []string{"1", "2", "3"}, vb: []string{"1", "3"}},
			want: false,
		},
		{
			name: "strLstHashPrefix07",
			args: args{va: []string{"1", "3", "3"}, vb: []string{"1", "3"}},
			want: true,
		},
		{
			name: "strLstHashPrefix08",
			args: args{va: []string{"1", "3", "3"}, vb: []string{}},
			want: true,
		},
		{
			name: "strLstHashPrefix09",
			args: args{va: []string{"1", "3", "3"}, vb: nil},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strLstHasPrefix(tt.args.va, tt.args.vb); got != tt.want {
				t.Errorf("strLstHasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
