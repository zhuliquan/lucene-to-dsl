package convert

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/x448/float16"
	"github.com/zhuliquan/datemath_parser"
	mapping "github.com/zhuliquan/es-mapping"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/scaled_float"
)

// convertFunc is trait of convert func which convert string value
// to specific type (i.g. int, float, bool)
type convertFunc func(string) (interface{}, error)

// convertToInt parse bool string value to bool value
func convertToBool(boolValue string) (interface{}, error) {
	return strconv.ParseBool(boolValue)
}

// convertToInt parse int string value to int value
// bits is the bit size of int value, which allows 8,16,32,64
func convertToInt(bits int) convertFunc {
	return func(intValue string) (interface{}, error) {
		if v, err := strconv.ParseInt(intValue, 10, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// convertToUInt parse uint string value to uint value
// bits is the bit size of uint value, which allows 8,16,32,64
func convertToUInt(bits int) convertFunc {
	return func(intValue string) (interface{}, error) {
		if v, err := strconv.ParseUint(intValue, 10, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// convertToFloat parse float string value to float16/float32/float64/float128.
// bits is the bit size of float value, which allows 16,32,64,128
// scalingFactor is the scaling factor for float128.
func convertToFloat(bits int, scalingFactor float64) convertFunc {
	return func(floatValue string) (interface{}, error) {
		switch bits {
		case 16:
			return convertToFloat16(floatValue)
		case 128:
			return convertToFloat128(floatValue, scalingFactor)
		default:
			if v, err := strconv.ParseFloat(floatValue, bits); err != nil {
				return nil, err
			} else {
				return v, nil
			}
		}
	}
}

// convertToFloat16 parse float string value to float16.
func convertToFloat16(floatValue string) (interface{}, error) {
	if f, err := strconv.ParseFloat(floatValue, 32); err != nil {
		return nil, err
	} else if f32 := float32(f); dsl.MinFloat16.Float32() > f32 || dsl.MaxFloat16.Float32() < f32 {
		return nil, strconv.ErrRange
	} else {
		return float16.Fromfloat32(f32), nil
	}
}

// convertToFloat128 parse float string value to float128.
func convertToFloat128(floatValue string, scalingFactor float64) (interface{}, error) {
	if f, err := scaled_float.NewFromString(floatValue, scalingFactor); err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

// convertToDate parse date math expr.
func convertToDate(property *mapping.Property) convertFunc {
	return func(s string) (interface{}, error) {
		var parser = getDateParserFromMapping(property)
		return parser.Parse(s)
	}
}

type dateRange struct {
	from time.Time
	to   time.Time
}

// convertToDateRange parse date math expr to `dateRange` object.
// TODO: 需要考虑如何解决如何处理 日缺失想查年月中所有天，月缺失想查整年的情况,
// 例如：field:2019-02 对应查询时间区间从2019-02-01 到 2019-02-28
// 例如：field:2019 对应查询时间区间从2019-01-01 到 2019-12-31
func convertToDateRange(property *mapping.Property) convertFunc {
	return func(s string) (interface{}, error) {
		var parser = getDateParserFromMapping(property)
		d, err := parser.Parse(s)
		if err != nil {
			return nil, err
		}
		var from, to = getDateRange(d)
		return &dateRange{
			from: from, to: to,
		}, nil
	}
}

// convertToVersion parse string version value to `Version` object.
func convertToVersion(versionValue string) (interface{}, error) {
	if v, err := version.NewVersion(versionValue); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}

// convertToIp parse string ip value to ip object.
// example: {"term": {"ip_field": "172.168.1.1"}}
func convertToIp(ipValue string) (interface{}, error) {
	if ip := net.ParseIP(ipValue); ip == nil {
		return nil, fmt.Errorf("ip value: %s is invalid", ipValue)
	} else {
		return ip, nil
	}
}

// convertToCidr parse string IpCidr value to `IpCidr` object.
// example: {"term": {"ip_field": "172.168.1.0/24"}}
func convertToCidr(ipValue string) (interface{}, error) {
	if _, cidr, err := net.ParseCIDR(ipValue); err != nil {
		return nil, err
	} else {
		return cidr, nil
	}
}

// convertToString return input string value.
var convertToString = func(s string) (interface{}, error) {
	return s, nil
}

var monthDay = map[time.Month]int{
	time.January:  31,
	time.March:    31,
	time.May:      31,
	time.July:     31,
	time.August:   31,
	time.October:  31,
	time.December: 31,

	time.April:     30,
	time.June:      30,
	time.September: 30,
	time.November:  30,
}

// getMonthDay get number of day in this month
// example: January has 31 days
// this func consider leap year
// example: 2000 year is leap year and February has 29 days
func getMonthDay(year int, month time.Month) int {
	if month == time.February {
		// check year is leap year
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			return 29
		} else {
			return 28
		}
	} else {
		return monthDay[month]
	}

}

const maxHour = 23
const maxMinute = 59
const maxSecond = 59
const maxNano = 999999999

// getDateRange get date range for prefix date.
// example: given 2021-01-01, we can get [2021-01-01 00:00:00 2021-01-01 23:59:59]
func getDateRange(t time.Time) (time.Time, time.Time) {
	var month = t.Month()
	var dateArr = []int{t.Year(), int(month), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()}
	var location = t.Location()
	if dateArr[6] != 0 {
		return t, t
	}
	if dateArr[5] != 0 {
		return t, time.Date(dateArr[0], month, dateArr[2], dateArr[3], dateArr[4], dateArr[5], maxNano, location)
	}
	if dateArr[4] != 0 {
		return t, time.Date(dateArr[0], month, dateArr[2], dateArr[3], dateArr[4], maxSecond, maxNano, location)
	}
	if dateArr[3] != 0 {
		return t, time.Date(dateArr[0], month, dateArr[2], dateArr[3], maxMinute, maxSecond, maxNano, location)
	}
	if dateArr[2] != 1 {
		return t, time.Date(dateArr[0], month, dateArr[2], maxHour, maxMinute, maxSecond, maxNano, location)
	}
	// if dateArr[1] != 1 {
	// 	return t, time.Date(dateArr[0], time.December, getMonthDay(dateArr[0], time.December), maxHour, maxMinute, maxSecond, maxNano, location)
	// }
	return t, t
}

func getDateParserFromMapping(property *mapping.Property) *datemath_parser.DateMathParser {
	var opts = []datemath_parser.DateMathParserOption{}
	if property.Format != "" {
		opts = append(opts, datemath_parser.WithFormat(strings.Split(property.Format, "||")))
	} else {
		if property.Type == mapping.DATE_NANOS_FIELD_TYPE {
			opts = append(opts, datemath_parser.WithFormat(
				[]string{"strict_date_optional_time_nanos", "epoch_millis"},
			))
		} else {
			opts = append(opts, datemath_parser.WithFormat(
				[]string{"strict_date_optional_time", "epoch_millis"},
			))
		}
	}
	if dp, err := datemath_parser.NewDateMathParser(opts...); err != nil {
		panic(err)
	} else {
		return dp
	}
}

type termValue interface {
	Value(func(string) (interface{}, error)) (interface{}, error)
}

type rangeValue interface {
	termValue
	IsInf(int) bool
}

var fieldTypeBits = map[mapping.FieldType]int{
	mapping.BYTE_FIELD_TYPE:          8,
	mapping.SHORT_FIELD_TYPE:         16,
	mapping.INTEGER_FIELD_TYPE:       32,
	mapping.INTEGER_RANGE_FIELD_TYPE: 32,
	mapping.LONG_FIELD_TYPE:          64,
	mapping.LONG_RANGE_FIELD_TYPE:    64,
	mapping.UNSIGNED_LONG_FIELD_TYPE: 64,
	mapping.HALF_FLOAT_FIELD_TYPE:    16,
	mapping.FLOAT_FIELD_TYPE:         32,
	mapping.FLOAT_RANGE_FIELD_TYPE:   32,
	mapping.DOUBLE_FIELD_TYPE:        64,
	mapping.DOUBLE_RANGE_FIELD_TYPE:  64,
	mapping.SCALED_FLOAT_FIELD_TYPE:  128,
}

func toUpper(x string) (string, error) {
	return strings.ToUpper(x), nil
}

func toLower(x string) (string, error) {
	return strings.ToLower(x), nil
}

func toStrLst(x string) (interface{}, error) {
	return strings.Split(x, ","), nil
}
