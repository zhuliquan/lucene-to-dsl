package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
)

func convertToInt64(intValue string) (interface{}, error) {
	if i, err := strconv.ParseInt(intValue, 10, 64); err != nil {
		return 0, fmt.Errorf("int_value: '%s' is invalid, err: %s", intValue, err.Error())
	} else {
		return i, nil
	}
}

func convertToUInt64(intValue string) (interface{}, error) {
	if i, err := strconv.ParseUint(intValue, 10, 64); err != nil {
		return 0, fmt.Errorf("int_value: '%s' is invalid, err: %s", intValue, err.Error())
	} else {
		return i, nil
	}
}

func ToUpper(x string) (string, error) {
	return strings.ToUpper(x), nil
}

func ToLower(x string) (string, error) {
	return strings.ToLower(x), nil
}

func ParseDate(x string) (string, error) {
	if t, err := dateparse.ParseAny(x); err != nil {
		return "", err
	} else {
		return t.Format("2006-01-02 15:04:05"), nil
	}

}
