package convert

import (
	"reflect"
	"testing"
)

func Test_convertToInt64(t *testing.T) {
	type args struct {
		intValue string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Test_convertToInt64_01",
			args:    args{intValue: "8908"},
			want:    int64(8908),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_02",
			args:    args{intValue: "-8908"},
			want:    int64(-8908),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_03",
			args:    args{intValue: "+8908"},
			want:    int64(8908),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_04",
			args:    args{intValue: "3.45"},
			want:    int64(0),
			wantErr: true,
		},
		{
			name:    "Test_convertToInt64_05",
			args:    args{intValue: "xxdasda3.45"},
			want:    int64(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToInt64(tt.args.intValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToBool(t *testing.T) {
	type args struct {
		boolValue string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Test_convertToBool_01",
			args:    args{boolValue: "true"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test_convertToBool_02",
			args:    args{boolValue: "false"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Test_convertToBool_03",
			args:    args{boolValue: "dasdad"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Test_convertToBool_04",
			args:    args{boolValue: "0"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Test_convertToBool_05",
			args:    args{boolValue: "1"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test_convertToBool_06",
			args:    args{boolValue: "True"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test_convertToBool_07",
			args:    args{boolValue: "False"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToBool(tt.args.boolValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
