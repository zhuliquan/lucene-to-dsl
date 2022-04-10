package convert

import (
	"math"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
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

func Test_convertToUInt64(t *testing.T) {
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
			name:    "test_error_case_1",
			args:    args{"xxxx"},
			want:    uint64(0),
			wantErr: true,
		},
		{
			name:    "test_error_case_2",
			args:    args{"-1"},
			want:    uint64(0),
			wantErr: true,
		},
		{
			name:    "test_ok_case",
			args:    args{"12"},
			want:    uint64(12),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToUInt64(tt.args.intValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToUInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToUInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToFloat64(t *testing.T) {
	type args struct {
		floatValue string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test_error_case",
			args:    args{"1.3x"},
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "test_ok_case_1",
			args:    args{"1.2"},
			want:    1.2,
			wantErr: false,
		},
		{
			name:    "test_ok_case_2",
			args:    args{"-1.2"},
			want:    -1.2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToFloat64(tt.args.floatValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToDate(t *testing.T) {
	type args struct {
		floatValue string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test_error_case",
			args:    args{"1.3x"},
			want:    time.Unix(0, 0),
			wantErr: true,
		},
		{
			name:    "test_ok_case_1",
			args:    args{"2001-01-02T09:09:09Z"},
			want:    time.Unix(978397749+8*3600, 0).UTC(),
			wantErr: false,
		},
		{
			name:    "test_ok_case_2",
			args:    args{"2001-01-02T09:09:09Z||+7d"},
			want:    time.Unix(978397749+7*24*3600+8*3600, 0).UTC(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f = convertToDate(&datemath_parser.DateMathParser{})
			got, err := f(tt.args.floatValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToIp(t *testing.T) {
	type args struct {
		ipValue string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test_ok_case",
			args:    args{"1.1.1.2"},
			want:    net.ParseIP("1.1.1.2"),
			wantErr: false,
		},
		{
			name:    "test_error_case",
			args:    args{"1.x.4.5"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToIp(tt.args.ipValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToCidr(t *testing.T) {
	type args struct {
		ipValue string
	}
	_, cidr, _ := net.ParseCIDR("1.1.1.45/20")
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test_ok_case",
			args:    args{"1.1.1.45/20"},
			want:    cidr,
			wantErr: false,
		},
		{
			name:    "test_error_case_1",
			args:    args{"1.2.3.4/78"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_error_case_2",
			args:    args{"1.2.3.x/7"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToCidr(tt.args.ipValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToCidr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToCidr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toUpper(t *testing.T) {
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
			name:    "test_ok",
			args:    args{"okOKiu"},
			want:    "OKOKIU",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toUpper(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("toUpper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toLower(t *testing.T) {
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
			name:    "test_ok",
			args:    args{"okOKiu"},
			want:    "okokiu",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toLower(tt.args.x)
			if (err != nil) != tt.wantErr {
				t.Errorf("toLower() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareInt64(t *testing.T) {
	type args struct {
		a int64
		b int64
		c dsl.CompareType
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test_lt",
			args: args{a: 1, b: 2, c: dsl.LT},
			want: 1,
		},
		{
			name: "test_gt",
			args: args{a: 1, b: 2, c: dsl.GT},
			want: 2,
		},
		{
			name: "test_lte",
			args: args{a: 1, b: 2, c: dsl.LTE},
			want: 1,
		},
		{
			name: "test_gte",
			args: args{a: 1, b: 2, c: dsl.GTE},
			want: 2,
		},
		{
			name: "test_eq_1",
			args: args{a: 3, b: 2, c: dsl.EQ},
			want: 3,
		},
		{
			name: "test_eq_2",
			args: args{a: 2, b: 3, c: dsl.EQ},
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
		c dsl.CompareType
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "test_lt",
			args: args{a: 1, b: 2, c: dsl.LT},
			want: 1,
		},
		{
			name: "test_gt",
			args: args{a: 1, b: 2, c: dsl.GT},
			want: 2,
		},
		{
			name: "test_lte",
			args: args{a: 1, b: 2, c: dsl.LTE},
			want: 1,
		},
		{
			name: "test_gte",
			args: args{a: 1, b: 2, c: dsl.GTE},
			want: 2,
		},
		{
			name: "test_eq_1",
			args: args{a: 3, b: 2, c: dsl.EQ},
			want: 3,
		},
		{
			name: "test_eq_2",
			args: args{a: 2, b: 3, c: dsl.EQ},
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
		c dsl.CompareType
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test_lt",
			args: args{a: 1, b: 2, c: dsl.LT},
			want: 1,
		},
		{
			name: "test_gt",
			args: args{a: 1, b: 2, c: dsl.GT},
			want: 2,
		},
		{
			name: "test_lte",
			args: args{a: 1, b: 2, c: dsl.LTE},
			want: 1,
		},
		{
			name: "test_gte",
			args: args{a: 1, b: 2, c: dsl.GTE},
			want: 2,
		},
		{
			name: "test_eq_1",
			args: args{a: 3, b: 2, c: dsl.EQ},
			want: 3,
		},
		{
			name: "test_eq_2",
			args: args{a: 2, b: 3, c: dsl.EQ},
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
		c dsl.CompareType
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		{
			name: "test_lt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: dsl.LT},
			want: net.ParseIP("1.2.3.4"),
		},
		{
			name: "test_gt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: dsl.GT},
			want: net.ParseIP("1.2.4.3"),
		},
		{
			name: "test_lt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: dsl.LTE},
			want: net.ParseIP("1.2.3.4"),
		},
		{
			name: "test_gt",
			args: args{a: net.ParseIP("1.2.3.4"), b: net.ParseIP("1.2.4.3"), c: dsl.GTE},
			want: net.ParseIP("1.2.4.3"),
		},
		{
			name: "test_eq_1",
			args: args{a: net.ParseIP("1.5.3.4"), b: net.ParseIP("1.2.4.3"), c: dsl.EQ},
			want: net.ParseIP("1.5.3.4"),
		},
		{
			name: "test_eq_2",
			args: args{a: net.ParseIP("1.2.4.3"), b: net.ParseIP("1.5.3.4"), c: dsl.EQ},
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
		c dsl.CompareType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_lt",
			args: args{a: "1", b: "2", c: dsl.LT},
			want: "1",
		},
		{
			name: "test_gt",
			args: args{a: "1", b: "2", c: dsl.GT},
			want: "2",
		},
		{
			name: "test_lte",
			args: args{a: "1", b: "2", c: dsl.LTE},
			want: "1",
		},
		{
			name: "test_gte",
			args: args{a: "1", b: "2", c: dsl.GTE},
			want: "2",
		},
		{
			name: "test_eq_1",
			args: args{a: "3", b: "2", c: dsl.EQ},
			want: "3",
		},
		{
			name: "test_eq_2",
			args: args{a: "2", b: "3", c: dsl.EQ},
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
		c dsl.CompareType
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "test_lt",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: dsl.LT},
			want: time.Unix(1, 0),
		},
		{
			name: "test_gt",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: dsl.GT},
			want: time.Unix(2, 0),
		},
		{
			name: "test_lte",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: dsl.LTE},
			want: time.Unix(1, 0),
		},
		{
			name: "test_gte",
			args: args{a: time.Unix(1, 0), b: time.Unix(2, 0), c: dsl.GTE},
			want: time.Unix(2, 0),
		},
		{
			name: "test_eq_1",
			args: args{a: time.Unix(3, 0), b: time.Unix(2, 0), c: dsl.EQ},
			want: time.Unix(3, 0),
		},
		{
			name: "test_eq_2",
			args: args{a: time.Unix(2, 0), b: time.Unix(3, 0), c: dsl.EQ},
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
