package convert

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/zhuliquan/datemath_parser"
)

// convert func tools
type convertFunc func(string) (interface{}, error)

// convert to bool value
func convertToBool(boolValue string) (interface{}, error) {
	return strconv.ParseBool(boolValue)
}

// convert to int value
func convertToInt(bits int) func(string) (interface{}, error) {
	return func(intValue string) (interface{}, error) {
		if v, err := strconv.ParseInt(intValue, 10, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// convert to uint value
func convertToUInt(bits int) func(string) (interface{}, error) {
	return func(intValue string) (interface{}, error) {
		if v, err := strconv.ParseUint(intValue, 10, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// convert to float value
func convertToFloat(bits int) func(string) (interface{}, error) {
	return func(floatValue string) (interface{}, error) {
		if v, err := strconv.ParseFloat(floatValue, bits); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

// parse date math expr
func convertToDate(parser *datemath_parser.DateMathParser) func(string) (interface{}, error) {
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

func toUpper(x string) (string, error) {
	return strings.ToUpper(x), nil
}

func toLower(x string) (string, error) {
	return strings.ToLower(x), nil
}
