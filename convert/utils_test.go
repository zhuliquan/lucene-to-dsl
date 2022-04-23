package convert

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/shopspring/decimal"
	"github.com/x448/float16"
	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
)

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
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_error_case_2",
			args:    args{"-1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_over_uint64",
			args:    args{"18446744073709551616"},
			want:    nil,
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
			got, err := convertToUInt(64)(tt.args.intValue)
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

func Test_convertToFloat(t *testing.T) {
	type args struct {
		bits       int
		floatValue string
	}
	d1, _ := decimal.NewFromString("2.797693134862315708145274237317043567981e+308")
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test_error_case",
			args:    args{64, "1.3x"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_ok_case_1",
			args:    args{64, "1.2"},
			want:    1.2,
			wantErr: false,
		},
		{
			name:    "test_over_case_1",
			args:    args{64, "2.797693134862315708145274237317043567981e+308"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_ok_case_2",
			args:    args{64, "-1.2"},
			want:    -1.2,
			wantErr: false,
		},
		{
			name:    "test_over_float16_01",
			args:    args{16, fmt.Sprintf("%f", dsl.MaxFloat16.Float32()+1)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_over_float16_02",
			args:    args{16, "3.41282346638528859811704183484516925440e+38"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_over_float16_02",
			args:    args{16, fmt.Sprintf("%f", dsl.MinFloat16.Float32()-1)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_parse_ok_float16",
			args:    args{16, "-1.2"},
			want:    float16.Fromfloat32(-1.2),
			wantErr: false,
		},
		{
			name:    "test_decimal_number_err",
			args:    args{128, "eer"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_ok_decimal",
			args:    args{128, "2.797693134862315708145274237317043567981e+308"},
			want:    d1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToFloat(tt.args.bits)(tt.args.floatValue)
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

func Test_convertToVersion(t *testing.T) {
	type args struct {
		versionValue string
	}
	version1, _ := version.NewVersion("v1.1.0")
	version2, _ := version.NewVersion("1.2.0")
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test_version_v1",
			args: args{
				versionValue: "v1.1.0",
			},
			want:    version1,
			wantErr: false,
		},
		{
			name: "test_version_v2",
			args: args{
				versionValue: "1.2.0",
			},
			want:    version2,
			wantErr: false,
		},
		{
			name: "test_wrong_version",
			args: args{
				versionValue: "dsadad",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToVersion(tt.args.versionValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToInt(t *testing.T) {
	type args struct {
		bits     int
		intValue string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test_over_int32",
			args:    args{bits: 32, intValue: "214748364790"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_correct_int32",
			args:    args{bits: 32, intValue: "2147483647"},
			want:    int64(2147483647),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_01",
			args:    args{intValue: "8908"},
			want:    int64(8908),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_02",
			args:    args{bits: 32, intValue: "-8908"},
			want:    int64(-8908),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_03",
			args:    args{bits: 32, intValue: "+8908"},
			want:    int64(8908),
			wantErr: false,
		},
		{
			name:    "Test_convertToInt64_04",
			args:    args{bits: 32, intValue: "3.45"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test_convertToInt64_05",
			args:    args{bits: 32, intValue: "xxdasda3.45"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test_convert_wrong_int64",
			args:    args{bits: 64, intValue: "9223372036854775808"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToInt(tt.args.bits)(tt.args.intValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("error: %v", err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}
