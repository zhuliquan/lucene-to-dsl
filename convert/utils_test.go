package convert

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/x448/float16"
	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
	"github.com/zhuliquan/scaled_float"
)

func TestConvertToBool(t *testing.T) {
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
			args:    args{boolValue: "Test_convertToBool_03"},
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
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToUInt64(t *testing.T) {
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
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToFloat(t *testing.T) {
	type args struct {
		bits       int
		floatValue string
	}
	d1, _ := scaled_float.NewFromString("2.78", 100)
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
			name:    "test_decimal_number_err_01",
			args:    args{128, "eer"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test_ok_decimal",
			args:    args{128, "2.78"},
			want:    d1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToFloat(tt.args.bits, 100)(tt.args.floatValue)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToDate(t *testing.T) {
	type args struct {
		timeValue string
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
			property := &mapping.Property{
				Format: "yyyy-MM-dd'T'HH:mm:ss'Z'",
			}
			var f = convertToDate(property)
			got, err := f(tt.args.timeValue)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToIp(t *testing.T) {
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
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToCidr(t *testing.T) {
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
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToUpper(t *testing.T) {
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
			args:    args{"who are You"},
			want:    "WHO ARE YOU",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toUpper(tt.args.x)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToLower(t *testing.T) {
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
			args:    args{"How are You"},
			want:    "how are you",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toLower(tt.args.x)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToVersion(t *testing.T) {
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
				versionValue: "wrong_version",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToVersion(tt.args.versionValue)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertToInt(t *testing.T) {
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
			args:    args{bits: 32, intValue: "unknown_prefix3.45"},
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
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetMonthDay(t *testing.T) {
	type args struct {
		year  int
		month time.Month
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_leap_Feb_01",
			args: args{
				year:  2000,
				month: time.February,
			},
			want: 29,
		},
		{
			name: "test_leap_Feb_02",
			args: args{
				year:  1600,
				month: time.February,
			},
			want: 29,
		},
		{
			name: "test_leap_Feb_03",
			args: args{
				year:  2008,
				month: time.February,
			},
			want: 29,
		},
		{
			name: "test_leap_Feb_04",
			args: args{
				year:  1900,
				month: time.February,
			},
			want: 28,
		},
		{
			name: "test_no_leap_Jan",
			args: args{
				year:  2007,
				month: time.January,
			},
			want: 31,
		},
		{
			name: "test_no_leap_Feb",
			args: args{
				year:  2007,
				month: time.February,
			},
			want: 28,
		},
		{
			name: "test_no_leap_Mar",
			args: args{
				year:  2007,
				month: time.March,
			},
			want: 31,
		},
		{
			name: "test_no_leap_Apr",
			args: args{
				year:  2007,
				month: time.April,
			},
			want: 30,
		},
		{
			name: "test_no_leap_May",
			args: args{
				year:  2007,
				month: time.May,
			},
			want: 31,
		},
		{
			name: "test_no_leap_Jun",
			args: args{
				year:  2007,
				month: time.June,
			},
			want: 30,
		},
		{
			name: "test_no_leap_Jul",
			args: args{
				year:  2007,
				month: time.July,
			},
			want: 31,
		},
		{
			name: "test_no_leap_Aug",
			args: args{
				year:  2007,
				month: time.July,
			},
			want: 31,
		},
		{
			name: "test_no_leap_Sep",
			args: args{
				year:  2007,
				month: time.September,
			},
			want: 30,
		},
		{
			name: "test_no_leap_Oct",
			args: args{
				year:  2007,
				month: time.October,
			},
			want: 31,
		},
		{
			name: "test_no_leap_Nov",
			args: args{
				year:  2007,
				month: time.November,
			},
			want: 30,
		},
		{
			name: "test_no_leap_Dec",
			args: args{
				year:  2007,
				month: time.December,
			},
			want: 31,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMonthDay(tt.args.year, tt.args.month)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDateRange(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name  string
		args  args
		want1 time.Time
		want2 time.Time
	}{
		{
			name: "split_date_01",
			args: args{
				t: time.Date(1900, time.February, 2, 2, 2, 2, 2, time.UTC),
			},
			want1: time.Date(1900, time.February, 2, 2, 2, 2, 2, time.UTC),
			want2: time.Date(1900, time.February, 2, 2, 2, 2, 2, time.UTC),
		},
		{
			name: "split_date_02",
			args: args{
				t: time.Date(1900, time.February, 2, 2, 2, 2, 0, time.UTC),
			},
			want1: time.Date(1900, time.February, 2, 2, 2, 2, 0, time.UTC),
			want2: time.Date(1900, time.February, 2, 2, 2, 2, 999999999, time.UTC),
		},
		{
			name: "split_date_03",
			args: args{
				t: time.Date(1900, time.February, 2, 2, 2, 0, 0, time.UTC),
			},
			want1: time.Date(1900, time.February, 2, 2, 2, 0, 0, time.UTC),
			want2: time.Date(1900, time.February, 2, 2, 2, 59, 999999999, time.UTC),
		},
		{
			name: "split_date_04",
			args: args{
				t: time.Date(1900, time.February, 2, 2, 0, 0, 0, time.UTC),
			},
			want1: time.Date(1900, time.February, 2, 2, 0, 0, 0, time.UTC),
			want2: time.Date(1900, time.February, 2, 2, 59, 59, 999999999, time.UTC),
		},
		{
			name: "split_date_05",
			args: args{
				t: time.Date(1900, time.February, 2, 0, 0, 0, 0, time.UTC),
			},
			want1: time.Date(1900, time.February, 2, 0, 0, 0, 0, time.UTC),
			want2: time.Date(1900, time.February, 2, 23, 59, 59, 999999999, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getDateRange(tt.args.t)
			assert.Equal(t, tt.want1, got)
			assert.Equal(t, tt.want2, got1)
		})
	}

	k, _ := time.Parse("2006", "2009")
	t.Logf("parse time: %+v", k)
}

func TestGetDateParserFromMapping(t *testing.T) {
	type args struct {
		property *mapping.Property
	}
	tests := []struct {
		name string
		args args
		want *datemath_parser.DateMathParser
	}{
		{
			name: "test_date_nanos_with_no_format",
			args: args{
				property: &mapping.Property{
					Type: mapping.DATE_NANOS_FIELD_TYPE,
				},
			},
			want: &datemath_parser.DateMathParser{
				Formats: []string{
					"yyyy-MM-ddTHH:mm:ss.SSSSSSZ", "yyyy-MM-dd", "epoch_millis",
				},
				TimeZone: time.UTC,
			},
		},
		{
			name: "test_date_with_no_format",
			args: args{
				property: &mapping.Property{
					Type: mapping.DATE_RANGE_FIELD_TYPE,
				},
			},
			want: &datemath_parser.DateMathParser{
				Formats: []string{
					"yyyy-MM-ddTHH:mm:ss.SSSZ", "yyyy-MM-dd", "epoch_millis",
				},
				TimeZone: time.UTC,
			},
		},
		{
			name: "test_date_with_format",
			args: args{
				property: &mapping.Property{
					Type:   mapping.DATE_RANGE_FIELD_TYPE,
					Format: "yyyy-MM-dd||xxxx-DDD||epoch_millis",
				},
			},
			want: &datemath_parser.DateMathParser{
				Formats: []string{
					"yyyy-MM-dd", "xxxx-DDD", "epoch_millis",
				},
				TimeZone: time.UTC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDateParserFromMapping(tt.args.property)
			assert.Equal(t, tt.want, got)
		})
	}
}
