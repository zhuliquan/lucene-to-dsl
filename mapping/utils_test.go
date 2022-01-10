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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strLstHasPrefix(tt.args.va, tt.args.vb); got != tt.want {
				t.Errorf("strLstHasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
