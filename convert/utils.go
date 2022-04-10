package convert

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/zhuliquan/datemath_parser"
	"github.com/zhuliquan/lucene-to-dsl/dsl"
)

// convert func tools
type convertFunc func(string) (interface{}, error)

// convert to bool value
func convertToBool(boolValue string) (interface{}, error) {
	return strconv.ParseBool(boolValue)
}

// convert to int64 value
func convertToInt64(intValue string) (interface{}, error) {
	return strconv.ParseInt(intValue, 10, 64)
}

// convert to uint64 value
func convertToUInt64(intValue string) (interface{}, error) {
	return strconv.ParseUint(intValue, 10, 64)
}

// convert to float value
func convertToFloat64(floatValue string) (interface{}, error) {
	return strconv.ParseFloat(floatValue, 64)
}

// parse date math expr
func convertToDate(parser *datemath_parser.DateMathParser) func(string) (interface{}, error) {
	return func(s string) (interface{}, error) {
		return parser.Parse(s)
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

func toUpper(x string) (string, error) {
	return strings.ToUpper(x), nil
}

func toLower(x string) (string, error) {
	return strings.ToLower(x), nil
}

func compareInt64(a, b int64, c dsl.CompareType) int64 {
	switch c {
	case dsl.LT:
		return ltInt64(a, b)
	case dsl.GT:
		return gtInt64(a, b)
	case dsl.LTE:
		return lteInt64(a, b)
	case dsl.GTE:
		return gteInt64(a, b)
	default:
		return a
	}
}

func ltInt64(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteInt64(a, b int64) int64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtInt64(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteInt64(a, b int64) int64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func compareUInt64(a, b uint64, c dsl.CompareType) uint64 {
	switch c {
	case dsl.LT:
		return ltUInt64(a, b)
	case dsl.GT:
		return gtUInt64(a, b)
	case dsl.LTE:
		return lteUInt64(a, b)
	case dsl.GTE:
		return gteUInt64(a, b)
	default:
		return a
	}
}

func ltUInt64(a, b uint64) uint64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteUInt64(a, b uint64) uint64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtUInt64(a, b uint64) uint64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteUInt64(a, b uint64) uint64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func compareFloat64(a, b float64, c dsl.CompareType) float64 {
	switch c {
	case dsl.LT:
		return ltFloat64(a, b)
	case dsl.GT:
		return gtFloat64(a, b)
	case dsl.LTE:
		return lteFloat64(a, b)
	case dsl.GTE:
		return gteFloat64(a, b)
	default:
		return a
	}
}

func ltFloat64(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteFloat64(a, b float64) float64 {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtFloat64(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteFloat64(a, b float64) float64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

func compareIp(a, b net.IP, c dsl.CompareType) net.IP {
	switch c {
	case dsl.LT:
		return ltIP(a, b)
	case dsl.GT:
		return gtIP(a, b)
	case dsl.LTE:
		return lteIP(a, b)
	case dsl.GTE:
		return gteIP(a, b)
	default:
		return a
	}
}

func ltIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res < 0 {
		return a
	} else {
		return b
	}
}

func lteIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res <= 0 {
		return a
	} else {
		return b
	}
}

func gtIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res > 0 {
		return a
	} else {
		return b
	}
}

func gteIP(a, b net.IP) net.IP {
	var res = bytes.Compare([]byte(a), []byte(b))
	if res >= 0 {
		return a
	} else {
		return b
	}
}

func compareDate(a, b time.Time, c dsl.CompareType) time.Time {
	switch c {
	case dsl.LT:
		return ltDate(a, b)
	case dsl.GT:
		return gtDate(a, b)
	case dsl.LTE:
		return lteDate(a, b)
	case dsl.GTE:
		return gteDate(a, b)
	default:
		return a
	}
}

func ltDate(a, b time.Time) time.Time {
	if a.UnixNano() < b.UnixNano() {
		return a
	} else {
		return b
	}
}

func lteDate(a, b time.Time) time.Time {
	if a.UnixNano() <= b.UnixNano() {
		return a
	} else {
		return b
	}
}

func gtDate(a, b time.Time) time.Time {
	if a.UnixNano() > b.UnixNano() {
		return a
	} else {
		return b
	}
}

func gteDate(a, b time.Time) time.Time {
	if a.UnixNano() >= b.UnixNano() {
		return a
	} else {
		return b
	}
}

func compareString(a, b string, c dsl.CompareType) string {
	switch c {
	case dsl.LT:
		return ltString(a, b)
	case dsl.GT:
		return gtString(a, b)
	case dsl.LTE:
		return lteString(a, b)
	case dsl.GTE:
		return gteString(a, b)
	default:
		return a
	}
}

func ltString(a, b string) string {
	if a < b {
		return a
	} else {
		return b
	}
}

func lteString(a, b string) string {
	if a <= b {
		return a
	} else {
		return b
	}
}

func gtString(a, b string) string {
	if a > b {
		return a
	} else {
		return b
	}
}

func gteString(a, b string) string {
	if a >= b {
		return a
	} else {
		return b
	}
}
