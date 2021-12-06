package convert

import "testing"

func TestParseDate(t *testing.T) {
	type args struct {
		x string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test_now",
			args:    args{x: "2021-01-02 06:07:00+8y"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
