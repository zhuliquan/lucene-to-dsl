package convert

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/shopspring/decimal"
	"github.com/x448/float16"
	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
	"github.com/zhuliquan/lucene-to-dsl/mapping"
)

// convert func tools
type convertFunc func(string) (interface{}, error)

// convert to bool value
func convertToBool(boolValue string) (interface{}, error) {
	return strconv.ParseBool(boolValue)
}

// convert to int value
func convertToInt(bits int) convertFunc {
	return func(intValue string) (interface{}, error) {
		if v, err := strconv.ParseInt(intValue, 10, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// convert to uint value
func convertToUInt(bits int) convertFunc {
	return func(intValue string) (interface{}, error) {
		if v, err := strconv.ParseUint(intValue, 10, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// convert to float value
func convertToFloat(bits int) convertFunc {
	return func(floatValue string) (interface{}, error) {
		if bits == 16 {
			if f, err := strconv.ParseFloat(floatValue, 32); err != nil {
				return nil, err
			} else if f32 := float32(f); dsl.MinFloat16.Float32() > f32 || dsl.MaxFloat16.Float32() < f32 {
				return nil, strconv.ErrRange
			} else {
				return float16.Fromfloat32(f32), nil
			}
		} else if bits == 128 {
			if f, err := decimal.NewFromString(floatValue); err != nil {
				return nil, err
			} else {
				return f, nil
			}
		} else {
			if v, err := strconv.ParseFloat(floatValue, bits); err != nil {
				return nil, err
			} else {
				return v, nil
			}
		}
	}
}

// parse date math expr
func convertToDate(parser *datemath_parser.DateMathParser) convertFunc {
	return func(s string) (interface{}, error) {
		return parser.Parse(s)
	}
}

// parse version
func convertToVersion(versionValue string) (interface{}, error) {
	if v, err := version.NewVersion(versionValue); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}

// convert to ip value， example: {"term": {"ip_field": "172.168.1.1"}}
func convertToIp(ipValue string) (interface{}, error) {
	if ip := net.ParseIP(ipValue); ip == nil {
		return nil, fmt.Errorf("ip value: %s is invalid", ipValue)
	} else {
		return ip, nil
	}
}

// convert to ip value， example: {"term": {"ip_field": "172.168.1.0/24"}}
func convertToCidr(ipValue string) (interface{}, error) {
	if _, cidr, err := net.ParseCIDR(ipValue); err != nil {
		return nil, err
	} else {
		return cidr, nil
	}
}

var convertToString = func(s string) (interface{}, error) { return s, nil }

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

// get date range for prefix date, i.g. given 2021-01-01, we can get [2021-01-01 00:00:00 2021-01-01 23:59:59]
// TODO: 需要考虑如何解决如何处理 日缺失想查一个月的锁有天的情况，月缺失想查整年的情况, 即：2019-02 / 2019。
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
	if dateArr[2] != 0 {
		return t, time.Date(dateArr[0], month, dateArr[2], maxHour, maxMinute, maxSecond, maxNano, location)
	}

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
	mapping.BYTE_FIELD_TYPE:           8,
	mapping.SHORT_FIELD_TYPE:          16,
	mapping.INTEGER_FIELD_TYPE:        32,
	mapping.INTERGER_RANGE_FIELD_TYPE: 32,
	mapping.LONG_FIELD_TYPE:           64,
	mapping.LONG_RANGE_FIELD_TYPE:     64,
	mapping.UNSIGNED_LONG_FIELD_TYPE:  64,
	mapping.HALF_FLOAT_FIELD_TYPE:     16,
	mapping.FLOAT_FIELD_TYPE:          32,
	mapping.FLOAT_RANGE_FIELD_TYPE:    32,
	mapping.DOUBLE_FIELD_TYPE:         64,
	mapping.DOUBLE_RANGE_FIELD_TYPE:   64,
	mapping.SCALED_FLOAT_FIELD_TYPE:   128,
}

func toUpper(x string) (string, error) {
	return strings.ToUpper(x), nil
}

func toLower(x string) (string, error) {
	return strings.ToLower(x), nil
}
